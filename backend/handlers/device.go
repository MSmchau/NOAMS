package handlers

import (
	"strconv"
	"time"

	"noams/middleware"
	"noams/models"
	"noams/services"
	"noams/utils"

	"github.com/gin-gonic/gin"
)

type DeviceHandler struct {
	deviceService *services.DeviceService
	pinger        *services.Pinger
}

func NewDeviceHandler(deviceService *services.DeviceService, pinger *services.Pinger) *DeviceHandler {
	return &DeviceHandler{
		deviceService: deviceService,
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
