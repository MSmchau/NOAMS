package services

import (
	"noams/models"

	"gorm.io/gorm"
)

type DeviceService struct {
	db *gorm.DB
}

func NewDeviceService(db *gorm.DB) *DeviceService {
	return &DeviceService{db: db}
}

func (s *DeviceService) List(page, pageSize int, filters map[string]interface{}) ([]models.Device, int64, error) {
	var devices []models.Device
	var total int64

	query := s.db.Model(&models.Device{}).Preload("Group").Preload("Credential")

	if name, ok := filters["name"]; ok {
		query = query.Where("name LIKE ?", "%"+name.(string)+"%")
	}
	if ip, ok := filters["management_ip"]; ok {
		query = query.Where("management_ip LIKE ?", "%"+ip.(string)+"%")
	}
	if vendor, ok := filters["vendor"]; ok {
		query = query.Where("vendor = ?", vendor)
	}
	if role, ok := filters["role"]; ok {
		query = query.Where("role = ?", role)
	}
	if groupID, ok := filters["group_id"]; ok {
		query = query.Where("group_id = ?", groupID)
	}
	if status, ok := filters["status"]; ok {
		query = query.Where("status = ?", status)
	}
	if building, ok := filters["building"]; ok {
		query = query.Where("building LIKE ?", "%"+building.(string)+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&devices).Error; err != nil {
		return nil, 0, err
	}

	return devices, total, nil
}

func (s *DeviceService) GetByID(id uint) (*models.Device, error) {
	var device models.Device
	if err := s.db.Preload("Group").Preload("Credential").First(&device, id).Error; err != nil {
		return nil, err
	}
	return &device, nil
}

func (s *DeviceService) Create(device *models.Device) error {
	return s.db.Create(device).Error
}

func (s *DeviceService) Update(id uint, updates map[string]interface{}) error {
	return s.db.Model(&models.Device{}).Where("id = ?", id).Updates(updates).Error
}

func (s *DeviceService) Delete(id uint) error {
	return s.db.Delete(&models.Device{}, id).Error
}

func (s *DeviceService) CountByStatus() (online int64, offline int64, err error) {
	if err = s.db.Model(&models.Device{}).Where("status = 1").Count(&online).Error; err != nil {
		return
	}
	err = s.db.Model(&models.Device{}).Where("status = 0").Count(&offline).Error
	return
}

func (s *DeviceService) CountByRole() (map[string]int64, error) {
	var results []struct {
		Role  string
		Count int64
	}
	if err := s.db.Model(&models.Device{}).Select("role, count(*) as count").Group("role").Find(&results).Error; err != nil {
		return nil, err
	}
	roleMap := make(map[string]int64)
	for _, r := range results {
		roleMap[r.Role] = r.Count
	}
	return roleMap, nil
}

func (s *DeviceService) ListAll() ([]models.Device, error) {
	var devices []models.Device
	if err := s.db.Find(&devices).Error; err != nil {
		return nil, err
	}
	return devices, nil
}

func (s *DeviceService) GetByIDs(ids []uint) ([]models.Device, error) {
	var devices []models.Device
	if err := s.db.Where("id IN ?", ids).Find(&devices).Error; err != nil {
		return nil, err
	}
	return devices, nil
}
