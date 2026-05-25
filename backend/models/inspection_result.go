package models

import (
	"time"
)

type InspectionResult struct {
	ID              uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	TaskID          string     `gorm:"type:varchar(64);index" json:"task_id"`
	DeviceID        uint       `gorm:"index;not null" json:"device_id"`
	Device          Device     `gorm:"foreignKey:DeviceID" json:"device,omitempty"`
	CPUUsage        *float64   `gorm:"type:decimal(5,2)" json:"cpu_usage"`
	MemoryUsage     *float64   `gorm:"type:decimal(5,2)" json:"memory_usage"`
	Temperature     *float64   `gorm:"type:decimal(5,2)" json:"temperature"`
	Uptime          string     `gorm:"type:varchar(128)" json:"uptime"`
	InterfaceStatus string     `gorm:"type:json" json:"interface_status"`
	RawOutput       string     `gorm:"type:longtext" json:"raw_output"`
	IsAnomaly       int        `gorm:"default:0" json:"is_anomaly"`
	AnomalyMsg      string     `gorm:"type:varchar(512)" json:"anomaly_msg"`
	Duration        int        `gorm:"default:0" json:"duration"`
	Status          string     `gorm:"type:varchar(16);default:success" json:"status"`
	InspectedAt     time.Time  `gorm:"index" json:"inspected_at"`
	CreatedAt       time.Time  `json:"created_at"`
}

func (InspectionResult) TableName() string {
	return "inspection_results"
}
