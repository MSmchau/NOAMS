package handlers

import (
	"fmt"
	"strconv"
	"time"

	"noams/middleware"
	"noams/models"
	"noams/services"
	"noams/utils"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

type DeviceHandler struct {
	deviceService *services.DeviceService
	db            *gorm.DB
	pinger        *services.Pinger
}

func NewDeviceHandler(deviceService *services.DeviceService, db *gorm.DB, pinger *services.Pinger) *DeviceHandler {
	return &DeviceHandler{
		deviceService: deviceService,
		db:            db,
		pinger:        pinger,
	}
}

func (h *DeviceHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	filters := make(map[string]interface{})
	if name := c.Query("name"); name != "" {
		filters["name"] = name
	}
	if ip := c.Query("management_ip"); ip != "" {
		filters["management_ip"] = ip
	}
	if vendor := c.Query("vendor"); vendor != "" {
		filters["vendor"] = vendor
	}
	if role := c.Query("role"); role != "" {
		filters["role"] = role
	}
	if groupID := c.Query("group_id"); groupID != "" {
		if gid, err := strconv.Atoi(groupID); err == nil {
			filters["group_id"] = gid
		}
	}
	if status := c.Query("status"); status != "" {
		if s, err := strconv.Atoi(status); err == nil {
			filters["status"] = s
		}
	}
	if building := c.Query("building"); building != "" {
		filters["building"] = building
	}

	devices, total, err := h.deviceService.List(page, pageSize, filters)
	if err != nil {
		utils.ServerError(c, "failed to query devices")
		return
	}

	utils.SuccessPage(c, devices, total, page, pageSize)
}

func (h *DeviceHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid device id")
		return
	}

	device, err := h.deviceService.GetByID(uint(id))
	if err != nil {
		utils.BadRequest(c, "device not found")
		return
	}

	utils.Success(c, device)
}

func (h *DeviceHandler) Create(c *gin.Context) {
	var device models.Device
	if err := c.ShouldBindJSON(&device); err != nil {
		utils.BadRequest(c, "invalid request: "+err.Error())
		return
	}

	if err := h.deviceService.Create(&device); err != nil {
		utils.ServerError(c, "failed to create device: "+err.Error())
		return
	}

	// 创建后立即检测一次设备状态，无需等待后台 pinger
	if h.pinger != nil {
		h.pinger.PingDevice(&device)
	}

	utils.Success(c, device)
}

func (h *DeviceHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid device id")
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		utils.BadRequest(c, "invalid request: "+err.Error())
		return
	}

	if err := h.deviceService.Update(uint(id), updates); err != nil {
		utils.ServerError(c, "failed to update device: "+err.Error())
		return
	}

	utils.Success(c, gin.H{"id": id})
}

func (h *DeviceHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid device id")
		return
	}

	if err := h.deviceService.Delete(uint(id)); err != nil {
		utils.ServerError(c, "failed to delete device: "+err.Error())
		return
	}

	utils.Success(c, gin.H{"id": id})
}

// Ping 手动检测单个设备的在线状态
func (h *DeviceHandler) Ping(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid device id")
		return
	}

	device, err := h.deviceService.GetByID(uint(id))
	if err != nil {
		utils.BadRequest(c, "device not found")
		return
	}

	online := h.pinger.PingDevice(device)

	status := "online"
	statusCode := 1
	if !online {
		status = "offline"
		statusCode = 0
	}

	utils.Success(c, gin.H{
		"id":        device.ID,
		"name":      device.Name,
		"ip":        device.ManagementIP,
		"status":    statusCode,
		"status_txt": status,
		"checked_at": time.Now().Format(time.RFC3339),
	})
}

func (h *DeviceHandler) Stats(c *gin.Context) {
	online, offline, err := h.deviceService.CountByStatus()
	if err != nil {
		utils.ServerError(c, "failed to get device stats")
		return
	}

	roleCount, err := h.deviceService.CountByRole()
	if err != nil {
		utils.ServerError(c, "failed to get role stats")
		return
	}

	utils.Success(c, gin.H{
		"online":      online,
		"offline":     offline,
		"total":       online + offline,
		"role_counts": roleCount,
	})
}

// ListAll returns all devices without pagination (for dropdowns, etc.)
func (h *DeviceHandler) ListAll(c *gin.Context) {
	_ = middleware.GetUserID(c)
	devices, err := h.deviceService.ListAll()
	if err != nil {
		utils.ServerError(c, "failed to query devices")
		return
	}
	utils.Success(c, devices)
}

// ExportDevices 导出所有设备为 JSON 文件
func (h *DeviceHandler) ExportDevices(c *gin.Context) {
	devices, err := h.deviceService.ListAll()
	if err != nil {
		utils.ServerError(c, "failed to export devices")
		return
	}

	// 直接返回原始 JSON 数组，导出的文件可直接用于导入
	c.Header("Content-Type", "application/json; charset=utf-8")
	c.Header("Content-Disposition", "attachment; filename=devices_export.json")
	c.JSON(200, devices)
}

// ImportDevices 批量导入设备
func (h *DeviceHandler) ImportDevices(c *gin.Context) {
	var req struct {
		Devices []struct {
			Name         string `json:"name"`
			ManagementIP string `json:"management_ip"`
			DeviceType   string `json:"device_type"`
			Vendor       string `json:"vendor"`
			Role         string `json:"role"`
			Model        string `json:"model"`
			SSHPort      int    `json:"ssh_port"`
			CredentialID *uint  `json:"credential_id"`
			SSHUsername  string `json:"ssh_username"`
			SSHPassword  string `json:"ssh_password"`
			Building     string `json:"building"`
			Floor        int    `json:"floor"`
			APName       string `json:"ap_name"`
			SNMPCommunity string `json:"snmp_community"`
			Description  string `json:"description"`
		} `json:"devices"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "无效的请求格式: "+err.Error())
		return
	}

	if len(req.Devices) == 0 {
		utils.BadRequest(c, "导入设备列表为空")
		return
	}

	var success, failed int
	var errors []string

	for i, d := range req.Devices {
		if d.Name == "" || d.ManagementIP == "" {
			errors = append(errors, fmt.Sprintf("第%d行: 设备名称和管理IP为必填项", i+1))
			failed++
			continue
		}

		// 如果提供了 ssh_username/ssh_password 且未指定 credential_id，自动创建凭据
		credID := d.CredentialID
		if credID == nil && d.SSHUsername != "" && d.SSHPassword != "" {
			cred := models.Credential{
				Name:       fmt.Sprintf("导入凭据-%s", d.Name),
				Username:   d.SSHUsername,
				Password:   d.SSHPassword,
				AuthMethod: "password",
			}
			if err := h.db.Create(&cred).Error; err == nil {
				credID = &cred.ID
			}
		}

		device := models.Device{
			Name:          d.Name,
			ManagementIP:  d.ManagementIP,
			DeviceType:    d.DeviceType,
			Vendor:        d.Vendor,
			Role:          d.Role,
			Model:         d.Model,
			SSHPort:       d.SSHPort,
			CredentialID:  credID,
			Building:      d.Building,
			Floor:         d.Floor,
			APName:        d.APName,
			SNMPCommunity: d.SNMPCommunity,
			Description:   d.Description,
		}

		if device.DeviceType == "" {
			device.DeviceType = "hp_comware"
		}
		if device.Role == "" {
			device.Role = "access"
		}
		if device.SSHPort == 0 {
			device.SSHPort = 22
		}

		if err := h.deviceService.Create(&device); err != nil {
			errors = append(errors, fmt.Sprintf("第%d行(%s): %s", i+1, d.Name, err.Error()))
			failed++
			continue
		}

		// 创建后检测在线状态
		if h.pinger != nil {
			h.pinger.PingDevice(&device)
		}
		success++
	}

	result := gin.H{
		"success": success,
		"failed":  failed,
		"total":   len(req.Devices),
	}
	if len(errors) > 0 {
		result["errors"] = errors
	}

	utils.Success(c, result)
}
