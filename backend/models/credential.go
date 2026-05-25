package models

import (
	"time"
)

type Credential struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"type:varchar(128);not null" json:"name"`
	Username    string    `gorm:"type:varchar(128);not null" json:"username"`
	Password    string    `gorm:"type:varchar(512);not null" json:"-"`
	EnablePW    string    `gorm:"type:varchar(512)" json:"-"`
	SSHKey      string    `gorm:"type:text" json:"-"`
	AuthMethod  string    `gorm:"type:varchar(16);default:password" json:"auth_method"`
	Description string    `gorm:"type:varchar(256)" json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (Credential) TableName() string {
	return "credentials"
}
