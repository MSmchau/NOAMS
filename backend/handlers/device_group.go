package handlers

import (
	"strconv"

	"noams/models"
	"noams/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DeviceGroupHandler struct {
	db *gorm.DB
}

func NewDeviceGroupHandler(db *gorm.DB) *DeviceGroupHandler {
	return &DeviceGroupHandler{db: db}
}

func (h *DeviceGroupHandler) List(c *gin.Context) {
	var groups []models.DeviceGroup
	if err := h.db.Preload("Children").Order("id ASC").Find(&groups).Error; err != nil {
		utils.ServerError(c, "查询分组失败")
		return
	}
	// 只返回顶级分组（Children 由 GORM 自动预加载）
	var topLevel []models.DeviceGroup
	for _, g := range groups {
		if g.ParentID == nil {
			topLevel = append(topLevel, g)
		}
	}
	utils.Success(c, topLevel)
}

func (h *DeviceGroupHandler) Create(c *gin.Context) {
	var req struct {
		Name     string `json:"name"`
		ParentID *uint  `json:"parent_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "无效的请求参数")
		return
	}
	if req.Name == "" {
		utils.BadRequest(c, "分组名称不能为空")
		return
	}

	group := models.DeviceGroup{
		Name:     req.Name,
		ParentID: req.ParentID,
	}
	if err := h.db.Create(&group).Error; err != nil {
		utils.ServerError(c, "创建分组失败")
		return
	}
	utils.Success(c, group)
}

func (h *DeviceGroupHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "无效的分组 ID")
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		utils.BadRequest(c, "无效的请求参数")
		return
	}

	if err := h.db.Model(&models.DeviceGroup{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		utils.ServerError(c, "更新分组失败")
		return
	}
	utils.Success(c, gin.H{"id": id})
}

func (h *DeviceGroupHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "无效的分组 ID")
		return
	}

	// 将属于该分组的设备设为无分组
	h.db.Model(&models.Device{}).Where("group_id = ?", id).Update("group_id", nil)
	// 将子分组设为无父分组
	h.db.Model(&models.DeviceGroup{}).Where("parent_id = ?", id).Update("parent_id", nil)
	// 删除分组
	h.db.Delete(&models.DeviceGroup{}, id)

	utils.Success(c, gin.H{"id": id})
}
