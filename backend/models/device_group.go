package models

import (
	"time"
)

type DeviceGroup struct {
	ID        uint          `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string        `gorm:"type:varchar(128);not null" json:"name"`
	ParentID  *uint         `gorm:"index" json:"parent_id"`
	Children  []DeviceGroup `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	Devices   []Device      `gorm:"foreignKey:GroupID" json:"devices,omitempty"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

func (DeviceGroup) TableName() string {
	return "device_groups"
}
