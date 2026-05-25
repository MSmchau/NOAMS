package handlers

import (
	"strconv"
	"time"

	"noams/models"
	"noams/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AlertHandler struct {
	db *gorm.DB
}

func NewAlertHandler(db *gorm.DB) *AlertHandler {
	return &AlertHandler{db: db}
}

func (h *AlertHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	severity := c.Query("severity")
	status := c.Query("status")
	alertType := c.Query("alert_type")

	query := h.db.Model(&models.Alert{}).Preload("Device")

	if severity != "" {
		query = query.Where("severity = ?", severity)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if alertType != "" {
		query = query.Where("alert_type = ?", alertType)
	}

	var total int64
	query.Count(&total)

	var alerts []models.Alert
	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&alerts)

	utils.SuccessPage(c, alerts, total, page, pageSize)
}

func (h *AlertHandler) Resolve(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid alert id")
		return
	}

	operator := c.GetString("username")
	now := time.Now()

	if err := h.db.Model(&models.Alert{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":       "resolved",
		"resolved_by":  operator,
		"resolved_at":  now,
	}).Error; err != nil {
		utils.ServerError(c, "failed to resolve alert")
		return
	}

	utils.Success(c, gin.H{"id": id, "status": "resolved"})
}

func (h *AlertHandler) Confirm(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid alert id")
		return
	}

	operator := c.GetString("username")
	now := time.Now()

	if err := h.db.Model(&models.Alert{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":       "confirmed",
		"confirmed_by": operator,
		"confirmed_at": now,
	}).Error; err != nil {
		utils.ServerError(c, "failed to confirm alert")
		return
	}

	utils.Success(c, gin.H{"id": id, "status": "confirmed"})
}

func (h *AlertHandler) Stats(c *gin.Context) {
	var critical, warning, info, triggered, resolved int64

	h.db.Model(&models.Alert{}).Where("severity = ? AND status != ?", "critical", "resolved").Count(&critical)
	h.db.Model(&models.Alert{}).Where("severity = ? AND status != ?", "warning", "resolved").Count(&warning)
	h.db.Model(&models.Alert{}).Where("severity = ? AND status != ?", "info", "resolved").Count(&info)
	h.db.Model(&models.Alert{}).Where("status = ?", "triggered").Count(&triggered)
	h.db.Model(&models.Alert{}).Where("status = ?", "resolved").Count(&resolved)

	utils.Success(c, gin.H{
		"critical":  critical,
		"warning":   warning,
		"info":      info,
		"triggered": triggered,
		"resolved":  resolved,
	})
}
