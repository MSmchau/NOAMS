package handlers

import (
	"strconv"
	"time"

	"noams/middleware"
	"noams/models"
	"noams/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type InspectionHandler struct {
	db *gorm.DB
}

func NewInspectionHandler(db *gorm.DB) *InspectionHandler {
	return &InspectionHandler{db: db}
}

func (h *InspectionHandler) InspectDevice(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid device id")
		return
	}

	var device models.Device
	if err := h.db.First(&device, uint(id)).Error; err != nil {
		utils.BadRequest(c, "device not found")
		return
	}

	taskID := "insp_" + strconv.FormatInt(time.Now().UnixNano(), 36)
	result := models.InspectionResult{
		TaskID:          taskID,
		DeviceID:        uint(id),
		Status:          "pending",
		InterfaceStatus: "{}",
		InspectedAt:     time.Now(),
	}

	if err := h.db.Omit("Device").Create(&result).Error; err != nil {
		utils.ServerError(c, "failed to create inspection task")
		return
	}

	utils.Success(c, gin.H{
		"message": "inspection task created",
		"id":      result.ID,
	})
}

func (h *InspectionHandler) BatchInspect(c *gin.Context) {
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

	taskID := "insp_" + strconv.FormatInt(time.Now().UnixNano(), 36)
	now := time.Now()

	var results []models.InspectionResult
	for _, deviceID := range req.DeviceIDs {
		results = append(results, models.InspectionResult{
			TaskID:      taskID,
			InterfaceStatus: "{}",
			DeviceID:    deviceID,
			Status:      "pending",
			InspectedAt: now,
		})
	}

	if err := h.db.Create(&results).Error; err != nil {
		utils.ServerError(c, "failed to create batch inspection")
		return
	}

	utils.Success(c, gin.H{
		"task_id": taskID,
		"count":   len(results),
	})
}

func (h *InspectionHandler) Report(c *gin.Context) {
	taskID := c.Query("task_id")
	deviceID := c.Query("device_id")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	query := h.db.Model(&models.InspectionResult{}).Preload("Device")

	if taskID != "" {
		query = query.Where("task_id = ?", taskID)
	}
	if deviceID != "" {
		query = query.Where("device_id = ?", deviceID)
	}
	if startDate != "" {
		query = query.Where("inspected_at >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("inspected_at <= ?", endDate)
	}

	var total int64
	query.Count(&total)

	var results []models.InspectionResult
	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&results)

	utils.SuccessPage(c, results, total, page, pageSize)
}

func (h *InspectionHandler) LatestReport(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	var results []models.InspectionResult
	h.db.Preload("Device").
		Where("status = ?", "success").
		Order("inspected_at DESC").
		Limit(limit).
		Find(&results)

	// Get summary counts
	var total, normal, anomaly int64
	h.db.Model(&models.InspectionResult{}).Count(&total)
	h.db.Model(&models.InspectionResult{}).Where("is_anomaly = 1").Count(&anomaly)
	normal = total - anomaly

	utils.Success(c, gin.H{
		"results": results,
		"summary": gin.H{
			"total":   total,
			"normal":  normal,
			"anomaly": anomaly,
		},
	})
}

func (h *InspectionHandler) History(c *gin.Context) {
	deviceID, err := strconv.ParseUint(c.Param("deviceId"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid device id")
		return
	}

	var results []models.InspectionResult
	h.db.Where("device_id = ?", deviceID).
		Order("inspected_at DESC").
		Limit(30).
		Find(&results)

	utils.Success(c, results)
}

var _ = middleware.GetUsername
