package handlers

import (
	"crypto/sha256"
	"fmt"
	"strconv"
	"time"

	"noams/models"
	"noams/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ConfigHandler struct {
	db *gorm.DB
}

func NewConfigHandler(db *gorm.DB) *ConfigHandler {
	return &ConfigHandler{db: db}
}

func (h *ConfigHandler) Backup(c *gin.Context) {
	var req struct {
		DeviceIDs []uint `json:"device_ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "invalid request")
		return
	}

	if len(req.DeviceIDs) == 0 {
		utils.BadRequest(c, "device_ids is required")
		return
	}

	operator := c.GetString("username")
	now := time.Now()
	var backups []models.ConfigBackup
	for _, deviceID := range req.DeviceIDs {
		// 获取设备信息
		var device models.Device
		deviceName := fmt.Sprintf("Device#%d", deviceID)
		if err := h.db.First(&device, deviceID).Error; err == nil {
			deviceName = device.Name
		}

		// 生成模拟配置内容（实际生产环境应通过SSH采集）
		content := fmt.Sprintf(`! NOAMS Configuration Backup
! Device: %s
! Backup Time: %s
! =========================================
!
version 7.1.070, Release 6608P02
sysname %s
!
interface GigabitEthernet1/0/1
 port access vlan 10
!
interface GigabitEthernet1/0/2
 port access vlan 20
!
interface Vlan-interface10
 ip address 192.168.%d.1 255.255.255.0
!
snmp-agent community read public
!
user-interface vty 0 4
 authentication-mode password
 user privilege level 15
!
return`, deviceName, now.Format("2006-01-02 15:04:05"), deviceName, deviceID)

		hash := fmt.Sprintf("%x", sha256.Sum256([]byte(content)))

		backups = append(backups, models.ConfigBackup{
			DeviceID:    deviceID,
			ConfigHash:  hash[:16],
			Content:     content,
			TriggeredBy: "manual",
			Operator:    operator,
			CreatedAt:   now,
		})
	}

	if err := h.db.Create(&backups).Error; err != nil {
		utils.ServerError(c, "failed to create backup tasks")
		return
	}

	utils.Success(c, gin.H{
		"message": "backup tasks created",
		"count":   len(backups),
	})
}

func (h *ConfigHandler) ListAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	query := h.db.Model(&models.ConfigBackup{}).Preload("Device")
	var total int64
	query.Count(&total)

	var backups []models.ConfigBackup
	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&backups)

	utils.SuccessPage(c, backups, total, page, pageSize)
}

func (h *ConfigHandler) History(c *gin.Context) {
	deviceID, err := strconv.ParseUint(c.Param("deviceId"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid device id")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	query := h.db.Model(&models.ConfigBackup{}).Where("device_id = ?", deviceID)
	var total int64
	query.Count(&total)

	var backups []models.ConfigBackup
	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&backups)

	utils.SuccessPage(c, backups, total, page, pageSize)
}

func (h *ConfigHandler) Rollback(c *gin.Context) {
	var req struct {
		BackupID uint `json:"backup_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "invalid request")
		return
	}

	var backup models.ConfigBackup
	if err := h.db.First(&backup, req.BackupID).Error; err != nil {
		utils.BadRequest(c, "backup record not found")
		return
	}

	utils.Success(c, gin.H{
		"message":   "rollback task created",
		"backup_id": backup.ID,
		"device_id": backup.DeviceID,
	})
}

func (h *ConfigHandler) Diff(c *gin.Context) {
	backupID1 := c.Query("id1")
	backupID2 := c.Query("id2")

	if backupID1 == "" || backupID2 == "" {
		utils.BadRequest(c, "id1 and id2 are required")
		return
	}

	var b1, b2 models.ConfigBackup
	if err := h.db.First(&b1, backupID1).Error; err != nil {
		utils.BadRequest(c, "backup id1 not found")
		return
	}
	if err := h.db.First(&b2, backupID2).Error; err != nil {
		utils.BadRequest(c, "backup id2 not found")
		return
	}

	utils.Success(c, gin.H{
		"backup_1": b1,
		"backup_2": b2,
	})
}

func (h *ConfigHandler) Export(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid backup id")
		return
	}

	var backup models.ConfigBackup
	if err := h.db.Preload("Device").First(&backup, uint(id)).Error; err != nil {
		utils.BadRequest(c, "backup not found")
		return
	}

	content := backup.Content
	if content == "" {
		content = "# 暂无配置内容\n"
	}

	deviceName := fmt.Sprintf("device_%d", backup.DeviceID)
	if backup.Device.Name != "" {
		deviceName = backup.Device.Name
	}

	filename := fmt.Sprintf("config_%s_%s.txt", deviceName, backup.CreatedAt.Format("20060102_150405"))
	c.Header("Content-Type", "text/plain; charset=utf-8")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.String(200, content)
}
