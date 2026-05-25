package models

import (
	"time"
)

type ScheduledTask struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"type:varchar(128);not null" json:"name"`
	TaskType    string    `gorm:"type:varchar(32);not null;index" json:"task_type"`
	CronExpr    string    `gorm:"type:varchar(64);not null" json:"cron_expr"`
	Params      string    `gorm:"type:json" json:"params"`
	Status      int       `gorm:"default:1" json:"status"`
	Description string    `gorm:"type:varchar(256)" json:"description"`
	LastRunAt   *time.Time `json:"last_run_at"`
	LastResult  string    `gorm:"type:varchar(32)" json:"last_result"`
	CreatedBy   string    `gorm:"type:varchar(64)" json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (ScheduledTask) TableName() string {
	return "scheduled_tasks"
}
