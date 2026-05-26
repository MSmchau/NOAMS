package services

import (
	"strings"

	"noams/config"
	"noams/models"
	"noams/utils"

	"gorm.io/gorm"
)

type CredentialService struct {
	db *gorm.DB
}

const encryptPrefix = "AES_GCM_V1:"

func encryptionKey() string {
	key := config.AppConfig.JWT.Secret
	if key == "" {
		key = "noams-default-encryption-key"
	}
	// Ensure key is exactly 32 bytes for AES-256
	if len(key) < 32 {
		key = key + strings.Repeat("0", 32-len(key))
	}
	return key[:32]
}

func encryptField(value string) (string, error) {
	if value == "" || strings.HasPrefix(value, encryptPrefix) {
		return value, nil
	}
	enc, err := utils.EncryptAES(value, encryptionKey())
	if err != nil {
		return value, err
	}
	return encryptPrefix + enc, nil
}

func decryptField(value string) (string, error) {
	if !strings.HasPrefix(value, encryptPrefix) {
		return value, nil // plaintext (legacy data or empty)
	}
	dec, err := utils.DecryptAES(value[len(encryptPrefix):], encryptionKey())
	if err != nil {
		return value, err
	}
	return dec, nil
}

func NewCredentialService(db *gorm.DB) *CredentialService {
	return &CredentialService{db: db}
}

func (s *CredentialService) List() ([]models.Credential, error) {
	var list []models.Credential
	if err := s.db.Order("id DESC").Find(&list).Error; err != nil {
		return nil, err
	}
	// 解密
	for i := range list {
		list[i].Password, _ = decryptField(list[i].Password)
	}
	return list, nil
}

func (s *CredentialService) GetByID(id uint) (*models.Credential, error) {
	var c models.Credential
	if err := s.db.First(&c, id).Error; err != nil {
		return nil, err
	}
	c.Password, _ = decryptField(c.Password)
	return &c, nil
}

func (s *CredentialService) Create(c *models.Credential) error {
	var err error
	c.Password, err = encryptField(c.Password)
	if err != nil {
		return err
	}
	c.EnablePW, err = encryptField(c.EnablePW)
	if err != nil {
		return err
	}
	return s.db.Create(c).Error
}

func (s *CredentialService) Update(id uint, updates map[string]interface{}) error {
	if pw, ok := updates["password"]; ok {
		if v, ok := pw.(string); ok {
			enc, err := encryptField(v)
			if err == nil {
				updates["password"] = enc
			}
		}
	}
	if pw, ok := updates["enable_pw"]; ok {
		if v, ok := pw.(string); ok {
			enc, err := encryptField(v)
			if err == nil {
				updates["enable_pw"] = enc
			}
		}
	}
	return s.db.Model(&models.Credential{}).Where("id = ?", id).Updates(updates).Error
}

func (s *CredentialService) Delete(id uint) error {
	// 解除引用该凭证的所有设备
	s.db.Model(&models.Device{}).Where("credential_id = ?", id).Update("credential_id", nil)
	return s.db.Delete(&models.Credential{}, id).Error
}
