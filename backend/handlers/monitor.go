package handlers

import (
	"noams/models"
	"noams/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MonitorHandler struct {
	db *gorm.DB
}

func NewMonitorHandler(db *gorm.DB) *MonitorHandler {
	return &MonitorHandler{db: db}
}

func (h *MonitorHandler) Dashboard(c *gin.Context) {
	var deviceCount, onlineCount, offlineCount int64
	var alertTriggered, alertResolved int64
	var lastInspection []models.InspectionResult

	h.db.Model(&models.Device{}).Count(&deviceCount)
	h.db.Model(&models.Device{}).Where("status = 1").Count(&onlineCount)
	h.db.Model(&models.Device{}).Where("status = 0").Count(&offlineCount)
	h.db.Model(&models.Alert{}).Where("status = ?", "triggered").Count(&alertTriggered)
	h.db.Model(&models.Alert{}).Where("status = ?", "resolved").Count(&alertResolved)

	h.db.Preload("Device").
		Where("status = ?", "success").
		Order("inspected_at DESC").
		Limit(10).
		Find(&lastInspection)

	var cpuTop []struct {
		DeviceID  uint    `json:"device_id"`
		DeviceName string `json:"device_name"`
		CPUUsage  float64 `json:"cpu_usage"`
	}
	h.db.Model(&models.InspectionResult{}).
		Select("device_id, MAX(cpu_usage) as cpu_usage").
		Where("cpu_usage IS NOT NULL").
		Group("device_id").
		Order("cpu_usage DESC").
		Limit(10).
		Scan(&cpuTop)

	var memTop []struct {
		DeviceID   uint    `json:"device_id"`
		DeviceName string `json:"device_name"`
		MemUsage   float64 `json:"mem_usage"`
	}
	h.db.Model(&models.InspectionResult{}).
		Select("device_id, MAX(memory_usage) as mem_usage").
		Where("memory_usage IS NOT NULL").
		Group("device_id").
		Order("mem_usage DESC").
		Limit(10).
		Scan(&memTop)

	utils.Success(c, gin.H{
		"device_stats": gin.H{
			"total":  deviceCount,
			"online": onlineCount,
			"offline": offlineCount,
		},
		"alert_stats": gin.H{
			"triggered": alertTriggered,
			"resolved":  alertResolved,
		},
		"cpu_top":        cpuTop,
		"mem_top":        memTop,
		"recent_checks":  lastInspection,
	})
}
