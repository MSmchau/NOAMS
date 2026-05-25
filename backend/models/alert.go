package models

import (
	"time"
)

type Alert struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	DeviceID    *uint     `gorm:"index" json:"device_id"`
	Device      Device    `gorm:"foreignKey:DeviceID" json:"device,omitempty"`
	AlertType   string    `gorm:"type:varchar(64);index;not null" json:"alert_type"`
	Severity    string    `gorm:"type:varchar(16);index;not null" json:"severity"`
	Message     string    `gorm:"type:varchar(512);not null" json:"message"`
	Detail      string    `gorm:"type:text" json:"detail"`
	Status      string    `gorm:"type:varchar(16);default:triggered;index" json:"status"`
	ConfirmedBy string    `gorm:"type:varchar(64)" json:"confirmed_by"`
	ConfirmedAt *time.Time `json:"confirmed_at"`
	ResolvedBy  string    `gorm:"type:varchar(64)" json:"resolved_by"`
	TriggeredAt time.Time `gorm:"index;not null" json:"triggered_at"`
	ResolvedAt  *time.Time `json:"resolved_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (Alert) TableName() string {
	return "alerts"
}
