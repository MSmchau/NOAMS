package handlers

import (
	"time"

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

type topItem struct {
	DeviceID   uint    `json:"device_id"`
	DeviceName string  `json:"device_name"`
	CPUUsage   float64 `json:"cpu_usage"`
	MemUsage   float64 `json:"mem_usage"`
}

type inspectionRow struct {
	ID          uint      `json:"id"`
	DeviceID    uint      `json:"device_id"`
	DeviceName  string    `json:"device_name"`
	DeviceIP    string    `json:"device_ip"`
	CPUUsage    *float64  `json:"cpu_usage"`
	MemoryUsage *float64  `json:"memory_usage"`
	Uptime      string    `json:"uptime"`
	Status      string    `json:"status"`
	IsAnomaly   int       `json:"is_anomaly"`
	InspectedAt time.Time `json:"inspected_at"`
}

func (h *MonitorHandler) Dashboard(c *gin.Context) {
	var deviceCount, onlineCount, offlineCount int64
	var alertTriggered, alertResolved int64

	h.db.Model(&models.Device{}).Count(&deviceCount)
	h.db.Model(&models.Device{}).Where("status = 1").Count(&onlineCount)
	h.db.Model(&models.Device{}).Where("status = 0").Count(&offlineCount)
	h.db.Model(&models.Alert{}).Where("status = ?", "triggered").Count(&alertTriggered)
	h.db.Model(&models.Alert{}).Where("status = ?", "resolved").Count(&alertResolved)

	// 最近巡检记录（JOIN 设备表获取名称）
	var recentInspections []inspectionRow
	h.db.Raw(`
		SELECT r.id, r.device_id, d.name AS device_name,
			d.management_ip AS device_ip, r.cpu_usage, r.memory_usage,
			r.uptime, r.status, r.is_anomaly, r.inspected_at
		FROM inspection_results r
		LEFT JOIN devices d ON d.id = r.device_id
		WHERE r.status = 'success'
		ORDER BY r.inspected_at DESC
		LIMIT 10
	`).Scan(&recentInspections)

	// CPU TOP10
	var cpuTop []topItem
	h.db.Raw(`
		SELECT r.device_id, d.name AS device_name, r.cpu_usage
		FROM inspection_results r
		LEFT JOIN devices d ON d.id = r.device_id
		WHERE r.cpu_usage IS NOT NULL
		AND r.id IN (
			SELECT MAX(id) FROM inspection_results
			WHERE cpu_usage IS NOT NULL GROUP BY device_id
		)
		ORDER BY r.cpu_usage DESC
		LIMIT 10
	`).Scan(&cpuTop)

	// 内存 TOP10
	var memTop []topItem
	h.db.Raw(`
		SELECT r.device_id, d.name AS device_name, r.memory_usage AS mem_usage
		FROM inspection_results r
		LEFT JOIN devices d ON d.id = r.device_id
		WHERE r.memory_usage IS NOT NULL
		AND r.id IN (
			SELECT MAX(id) FROM inspection_results
			WHERE memory_usage IS NOT NULL GROUP BY device_id
		)
		ORDER BY r.memory_usage DESC
		LIMIT 10
	`).Scan(&memTop)

	utils.Success(c, gin.H{
		"device_stats": gin.H{
			"total":   deviceCount,
			"online":  onlineCount,
			"offline": offlineCount,
		},
		"alert_stats": gin.H{
			"triggered": alertTriggered,
			"resolved":  alertResolved,
		},
		"cpu_top":       cpuTop,
		"mem_top":       memTop,
		"recent_checks": recentInspections,
	})
}
