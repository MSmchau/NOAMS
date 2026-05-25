package services

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"noams/models"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Pinger struct {
	db          *gorm.DB
	interval    time.Duration
	timeout     time.Duration
	retry       int
	method      string
	concurrency int
}

type PingerConfig struct {
	Interval    time.Duration
	Timeout     time.Duration
	Retry       int
	Method      string
	Concurrency int
}

func NewPinger(db *gorm.DB, cfg PingerConfig) *Pinger {
	if cfg.Concurrency < 1 {
		cfg.Concurrency = 50
	}
	return &Pinger{
		db:          db,
		interval:    cfg.Interval,
		timeout:     cfg.Timeout,
		retry:       cfg.Retry,
		method:      cfg.Method,
		concurrency: cfg.Concurrency,
	}
}

// Start 启动后台周期性检测。首次启动时立即执行一次检测。
func (p *Pinger) Start(ctx context.Context) {
	zap.L().Info("启动设备状态检测服务",
		zap.Duration("interval", p.interval),
		zap.Duration("timeout", p.timeout),
		zap.Int("retry", p.retry),
		zap.String("method", p.method),
		zap.Int("concurrency", p.concurrency))

	// 启动时立即执行一次
	p.CheckAll()

	ticker := time.NewTicker(p.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			p.CheckAll()
		case <-ctx.Done():
			zap.L().Info("设备状态检测服务已停止")
			return
		}
	}
}

// CheckAll 检测所有设备的状态
func (p *Pinger) CheckAll() {
	var devices []models.Device
	if err := p.db.Find(&devices).Error; err != nil {
		zap.L().Error("查询设备列表失败", zap.Error(err))
		return
	}

	if len(devices) == 0 {
		return
	}

	zap.L().Info("开始批量设备状态检测", zap.Int("设备数量", len(devices)))

	sem := make(chan struct{}, p.concurrency)
	var wg sync.WaitGroup

	for i := range devices {
		wg.Add(1)
		go func(d *models.Device) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			p.checkDevice(d)
		}(&devices[i])
	}

	wg.Wait()
	zap.L().Info("批量设备状态检测完成")
}

// checkDevice 检测单个设备并更新数据库
func (p *Pinger) checkDevice(device *models.Device) {
	online := p.tryPing(device.ManagementIP, device.SSHPort)

	newStatus := 0
	if online {
		newStatus = 1
	}

	updates := map[string]interface{}{
		"status": newStatus,
	}
	if online {
		now := time.Now()
		updates["last_seen"] = &now
	}

	if err := p.db.Model(&models.Device{}).Where("id = ?", device.ID).Updates(updates).Error; err != nil {
		zap.L().Error("更新设备状态失败",
			zap.String("name", device.Name),
			zap.String("ip", device.ManagementIP),
			zap.Error(err))
		return
	}

	if device.Status != newStatus {
		state := "在线"
		if newStatus == 0 {
			state = "离线"
		}
		zap.L().Info("设备状态变更",
			zap.String("name", device.Name),
			zap.String("ip", device.ManagementIP),
			zap.String("status", state))
	}
}

// CheckDevice 手动检测单个设备，返回是否在线（不写数据库，仅返回检测结果）
func (p *Pinger) CheckDevice(device *models.Device) bool {
	return p.tryPing(device.ManagementIP, device.SSHPort)
}

// PingDevice 手动检测并更新单个设备状态
func (p *Pinger) PingDevice(device *models.Device) bool {
	online := p.tryPing(device.ManagementIP, device.SSHPort)

	newStatus := 0
	if online {
		newStatus = 1
	}

	updates := map[string]interface{}{
		"status": newStatus,
	}
	if online {
		now := time.Now()
		updates["last_seen"] = &now
	}

	_ = p.db.Model(&models.Device{}).Where("id = ?", device.ID).Updates(updates).Error
	return online
}

// tryPing 尝试多次 ping 设备
func (p *Pinger) tryPing(ip string, port int) bool {
	for i := 0; i <= p.retry; i++ {
		if p.ping(ip, port) {
			return true
		}
		if i < p.retry {
			time.Sleep(500 * time.Millisecond)
		}
	}
	return false
}

// ping 执行单次检测
func (p *Pinger) ping(ip string, port int) bool {
	switch p.method {
	case "icmp":
		return p.icmpPing(ip)
	default:
		return p.tcpPing(ip, port)
	}
}

// tcpPing 通过 TCP 端口检测设备可达性
func (p *Pinger) tcpPing(ip string, port int) bool {
	addr := net.JoinHostPort(ip, fmt.Sprintf("%d", port))
	conn, err := net.DialTimeout("tcp", addr, p.timeout)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

// icmpPing 通过 ICMP 检测设备可达性
func (p *Pinger) icmpPing(ip string) bool {
	addr := net.JoinHostPort(ip, "0")
	conn, err := net.DialTimeout("ip:icmp", addr, p.timeout)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}
