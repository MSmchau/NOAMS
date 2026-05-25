package models

import (
	"time"
)

type ConfigBackup struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	DeviceID    uint      `gorm:"index;not null" json:"device_id"`
	Device      Device    `gorm:"foreignKey:DeviceID" json:"device,omitempty"`
	ConfigHash  string    `gorm:"type:varchar(64);index" json:"config_hash"`
	GitCommitID string    `gorm:"type:varchar(40)" json:"git_commit_id"`
	FilePath    string    `gorm:"type:varchar(512)" json:"file_path"`
	Content     string    `gorm:"type:longtext" json:"content,omitempty"`
	TriggeredBy string    `gorm:"type:varchar(32)" json:"triggered_by"`
	Operator    string    `gorm:"type:varchar(64)" json:"operator"`
	CreatedAt   time.Time `json:"created_at"`
}

func (ConfigBackup) TableName() string {
	return "config_backups"
}
