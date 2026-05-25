package services

import (
	"noams/models"

	"gorm.io/gorm"
)

type CredentialService struct {
	db *gorm.DB
}

func NewCredentialService(db *gorm.DB) *CredentialService {
	return &CredentialService{db: db}
}

func (s *CredentialService) List() ([]models.Credential, error) {
	var list []models.Credential
	if err := s.db.Order("id DESC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (s *CredentialService) GetByID(id uint) (*models.Credential, error) {
	var c models.Credential
	if err := s.db.First(&c, id).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (s *CredentialService) Create(c *models.Credential) error {
	return s.db.Create(c).Error
}

func (s *CredentialService) Update(id uint, updates map[string]interface{}) error {
	return s.db.Model(&models.Credential{}).Where("id = ?", id).Updates(updates).Error
}

func (s *CredentialService) Delete(id uint) error {
	// 解除引用该凭证的所有设备
	s.db.Model(&models.Device{}).Where("credential_id = ?", id).Update("credential_id", nil)
	return s.db.Delete(&models.Credential{}, id).Error
}
