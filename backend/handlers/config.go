package handlers

import (
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

	now := time.Now()
	var backups []models.ConfigBackup
	for _, deviceID := range req.DeviceIDs {
		backups = append(backups, models.ConfigBackup{
			DeviceID:    deviceID,
			TriggeredBy: "manual",
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
