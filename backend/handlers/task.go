package handlers

import (
	"strconv"
	"time"

	"noams/models"
	"noams/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TaskHandler struct {
	db *gorm.DB
}

func NewTaskHandler(db *gorm.DB) *TaskHandler {
	return &TaskHandler{db: db}
}

func (h *TaskHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	var total int64
	var tasks []models.ScheduledTask

	h.db.Model(&models.ScheduledTask{}).Count(&total)
	offset := (page - 1) * pageSize
	h.db.Offset(offset).Limit(pageSize).Order("id DESC").Find(&tasks)

	utils.SuccessPage(c, tasks, total, page, pageSize)
}

func (h *TaskHandler) Create(c *gin.Context) {
	var task models.ScheduledTask
	if err := c.ShouldBindJSON(&task); err != nil {
		utils.BadRequest(c, "invalid request: "+err.Error())
		return
	}

	operator := c.GetString("username")
	task.CreatedBy = operator
	task.Status = 1

	if err := h.db.Create(&task).Error; err != nil {
		utils.ServerError(c, "failed to create task")
		return
	}

	utils.Success(c, task)
}

func (h *TaskHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid task id")
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		utils.BadRequest(c, "invalid request")
		return
	}

	if err := h.db.Model(&models.ScheduledTask{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		utils.ServerError(c, "failed to update task")
		return
	}

	utils.Success(c, gin.H{"id": id})
}

func (h *TaskHandler) Toggle(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid task id")
		return
	}

	var task models.ScheduledTask
	if err := h.db.First(&task, id).Error; err != nil {
		utils.BadRequest(c, "task not found")
		return
	}

	newStatus := 0
	if task.Status == 0 {
		newStatus = 1
	}
	h.db.Model(&task).Update("status", newStatus)

	utils.Success(c, gin.H{"id": id, "status": newStatus})
}

func (h *TaskHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid task id")
		return
	}

	h.db.Delete(&models.ScheduledTask{}, id)
	utils.Success(c, gin.H{"id": id})
}

func (h *TaskHandler) Logs(c *gin.Context) {
	taskID := c.Query("task_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	query := h.db.Model(&models.InspectionResult{})
	if taskID != "" {
		query = query.Where("task_id LIKE ?", "%"+taskID+"%")
	}

	var total int64
	query.Count(&total)

	var logs []models.InspectionResult
	offset := (page - 1) * pageSize
	query.Preload("Device").Offset(offset).Limit(pageSize).Order("id DESC").Find(&logs)

	utils.SuccessPage(c, logs, total, page, pageSize)
}

var _ = time.Now
