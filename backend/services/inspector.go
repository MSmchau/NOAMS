package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"noams/models"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// 设备类型对应的巡检命令
type deviceSpec struct {
	CPUCommand       string
	MemoryCommand    string
	UptimeCommand    string
	InterfaceCommand string
}

var deviceCommands = map[string]deviceSpec{
	"hp_comware": {
		CPUCommand:       "display cpu-usage",
		MemoryCommand:    "display memory",
		UptimeCommand:    "display version",
		InterfaceCommand: "display interface brief",
	},
	"huawei": {
		CPUCommand:       "display cpu-usage",
		MemoryCommand:    "display memory-usage",
		UptimeCommand:    "display version",
		InterfaceCommand: "display interface brief",
	},
	"cisco_ios": {
		CPUCommand:       "show processes cpu sorted | include CPU utilization",
		MemoryCommand:    "show memory statistics",
		UptimeCommand:    "show version",
		InterfaceCommand: "show ip interface brief",
	},
	"ruijie_os": {
		CPUCommand:       "show cpu",
		MemoryCommand:    "show memory",
		UptimeCommand:    "show version",
		InterfaceCommand: "show interface brief",
	},
}

type Inspector struct {
	db         *gorm.DB
	workerURL  string
	interval   time.Duration
	httpClient *http.Client
}

type InspectorConfig struct {
	Interval  time.Duration
	WorkerURL string
	Timeout   time.Duration
}

func NewInspector(db *gorm.DB, cfg InspectorConfig) *Inspector {
	if cfg.Timeout <= 0 {
		cfg.Timeout = 30 * time.Second
	}
	return &Inspector{
		db:        db,
		workerURL: strings.TrimRight(cfg.WorkerURL, "/"),
		interval:  cfg.Interval,
		httpClient: &http.Client{
			Timeout: cfg.Timeout,
		},
	}
}

// Start 启动后台巡检任务处理
func (ins *Inspector) Start(ctx context.Context) {
	zap.L().Info("启动巡检任务处理服务",
		zap.Duration("interval", ins.interval),
		zap.String("worker_url", ins.workerURL))

	ticker := time.NewTicker(ins.interval)
	defer ticker.Stop()

	// 启动时立即处理一次
	ins.processPending()

	for {
		select {
		case <-ticker.C:
			ins.processPending()
		case <-ctx.Done():
			zap.L().Info("巡检任务处理服务已停止")
			return
		}
	}
}

// processPending 处理所有 pending 状态的巡检任务
func (ins *Inspector) processPending() {
	var tasks []models.InspectionResult
	if err := ins.db.Where("status = ?", "pending").
		Preload("Device").
		Preload("Device.Credential").
		Find(&tasks).Error; err != nil {
		zap.L().Error("查询待处理巡检任务失败", zap.Error(err))
		return
	}

	if len(tasks) == 0 {
		return
	}

	zap.L().Info("处理待处理巡检任务", zap.Int("count", len(tasks)))

	for i := range tasks {
		ins.executeTask(&tasks[i])

		// 任务间间隔，避免并发 SSH 压力
		time.Sleep(500 * time.Millisecond)
	}
}

// executeTask 执行单个巡检任务
func (ins *Inspector) executeTask(task *models.InspectionResult) {
	device := task.Device

	// 检查设备是否有凭证
	if device.CredentialID == nil || device.Credential == nil {
		ins.failTask(task, "设备未配置凭证")
		return
	}

	spec, ok := deviceCommands[device.DeviceType]
	if !ok {
		ins.failTask(task, fmt.Sprintf("不支持的设备类型: %s", device.DeviceType))
		return
	}

	commands := []string{spec.CPUCommand}
	if spec.MemoryCommand != "" {
		commands = append(commands, spec.MemoryCommand)
	}
	if spec.UptimeCommand != "" {
		commands = append(commands, spec.UptimeCommand)
	}
	if spec.InterfaceCommand != "" {
		commands = append(commands, spec.InterfaceCommand)
	}

	startTime := time.Now()

	// 调用 netmiko-worker
	reqBody := executeRequest{
		IP:         device.ManagementIP,
		Port:       device.SSHPort,
		Username:   device.Credential.Username,
		Password:   device.Credential.Password,
		DeviceType: device.DeviceType,
		Commands:   commands,
	}

	resp, err := ins.callWorker(reqBody)
	duration := int(time.Since(startTime).Seconds())

	if err != nil {
		ins.failTask(task, fmt.Sprintf("调用 worker 失败: %s", err.Error()))
		return
	}

	if !resp.Success {
		errMsg := "执行命令失败"
		if len(resp.Results) > 0 && !resp.Results[0].Success {
			errMsg = resp.Results[0].Error
		}
		ins.failTask(task, errMsg)
		return
	}

	// 解析巡检数据
	var cpuUsage, memoryUsage *float64
	var uptime string
	var rawOutputs []string

	for _, r := range resp.Results {
		rawOutputs = append(rawOutputs, fmt.Sprintf("=== %s ===\n%s", r.Command, r.Output))

		if r.Command == spec.CPUCommand {
			cpuUsage = parseCPU(r.Output)
		}
		if spec.MemoryCommand != "" && r.Command == spec.MemoryCommand {
			memoryUsage = parseMemory(r.Output)
		}
		if spec.UptimeCommand != "" && r.Command == spec.UptimeCommand {
			uptime = parseUptime(r.Output)
		}
	}

	// 更新数据库
	updates := map[string]interface{}{
		"status":    "success",
		"cpu_usage": cpuUsage,
		"uptime":    uptime,
		"duration":  duration,
		"raw_output": strings.Join(rawOutputs, "\n\n"),
	}
	if memoryUsage != nil {
		updates["memory_usage"] = memoryUsage
	}

	if err := ins.db.Model(&models.InspectionResult{}).
		Where("id = ?", task.ID).Updates(updates).Error; err != nil {
		zap.L().Error("更新巡检结果失败", zap.Uint("task_id", task.ID), zap.Error(err))
		return
	}

	// 同时更新设备的 last_seen
	now := time.Now()
	_ = ins.db.Model(&models.Device{}).Where("id = ?", device.ID).
		Update("last_seen", &now).Error

	zap.L().Info("巡检完成",
		zap.String("device", device.Name),
		zap.String("ip", device.ManagementIP),
		zap.Int("duration", duration),
		zap.Any("cpu", cpuUsage),
		zap.Any("memory", memoryUsage))
}

// failTask 将任务标记为失败
func (ins *Inspector) failTask(task *models.InspectionResult, errMsg string) {
	zap.L().Warn("巡检失败",
		zap.String("device", task.Device.Name),
		zap.String("ip", task.Device.ManagementIP),
		zap.String("error", errMsg))

	_ = ins.db.Model(&models.InspectionResult{}).
		Where("id = ?", task.ID).Updates(map[string]interface{}{
		"status":  "failed",
		"anomaly_msg": errMsg,
	}).Error
}

// callWorker 调用 netmiko-worker 执行命令
func (ins *Inspector) callWorker(req executeRequest) (*executeResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %w", err)
	}

	url := ins.workerURL + "/api/v1/execute"
	httpResp, err := ins.httpClient.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("HTTP 请求失败: %w", err)
	}
	defer httpResp.Body.Close()

	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	if httpResp.StatusCode != 200 {
		return nil, fmt.Errorf("worker 返回错误状态 %d: %s", httpResp.StatusCode, string(respBody))
	}

	var result executeResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return &result, nil
}

// parseCPU 从命令输出中解析 CPU 使用率
func parseCPU(output string) *float64 {
	// 通用匹配：第一个百分比数字
	re := regexp.MustCompile(`(\d+(?:\.\d+)?)%`)
	matches := re.FindStringSubmatch(output)
	if len(matches) > 1 {
		v, err := strconv.ParseFloat(matches[1], 64)
		if err == nil && v >= 0 && v <= 100 {
			return &v
		}
	}
	return nil
}

// parseMemory 从命令输出中解析内存使用率
func parseMemory(output string) *float64 {
	// 1. 直接百分比： "Current memory usage: 18%"
	re1 := regexp.MustCompile(`(?i)(?:memory|mem)\s*(?:usage|used|utilization)?[\s:\-]+(\d+(?:\.\d+)?)%`)
	if matches := re1.FindStringSubmatch(output); len(matches) > 1 {
		if v, err := strconv.ParseFloat(matches[1], 64); err == nil && v >= 0 && v <= 100 {
			return &v
		}
	}

	// 2. "used/total" 同行: "Used: 1234, Total: 6912"
	re2 := regexp.MustCompile(`(?i)used[\s:]+(\d+)[\s\S]{0,30}total[\s:]+(\d+)`)
	if matches := re2.FindStringSubmatch(output); len(matches) > 2 {
		used, _ := strconv.ParseFloat(matches[1], 64)
		total, _ := strconv.ParseFloat(matches[2], 64)
		if total > 0 && used > 0 {
			v := used / total * 100
			return &v
		}
	}

	// 3. "Total/Used" 分行
	re3 := regexp.MustCompile(`(?i)Total[\s:]+(\d+)`)
	re4 := regexp.MustCompile(`(?i)Used[\s:]+(\d+)`)
	totalMatch := re3.FindStringSubmatch(output)
	usedMatch := re4.FindStringSubmatch(output)
	if len(totalMatch) > 1 && len(usedMatch) > 1 {
		total, _ := strconv.ParseFloat(totalMatch[1], 64)
		used, _ := strconv.ParseFloat(usedMatch[1], 64)
		if total > 0 && used > 0 {
			v := used / total * 100
			return &v
		}
	}

	return nil
}

// parseUptime 从命令输出中解析系统运行时长
func parseUptime(output string) string {
	// 常见格式: "uptime is 45 weeks, 5 days, 3 hours, 18 minutes"
	// "System uptime is 2 days, 1 hour, 30 minutes"
	re := regexp.MustCompile(`(?i)(?:uptime\s+is|up\s+for)\s+(.+?)(?:\n|$)`)
	if matches := re.FindStringSubmatch(output); len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}
	return ""
}

// executeRequest 发送给 netmiko-worker 的请求结构
type executeRequest struct {
	IP         string   `json:"ip"`
	Port       int      `json:"port"`
	Username   string   `json:"username"`
	Password   string   `json:"password"`
	DeviceType string   `json:"device_type"`
	Commands   []string `json:"commands"`
}

// executeResponse netmiko-worker 返回的响应结构
type executeResponse struct {
	Success bool `json:"success"`
	Results []struct {
		Command string `json:"command"`
		Output  string `json:"output"`
		Success bool   `json:"success"`
		Error   string `json:"error,omitempty"`
	} `json:"results"`
	Elapsed float64 `json:"elapsed"`
}
