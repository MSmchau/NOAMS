package models

import (
	"time"
)

type Device struct {
	ID           uint              `gorm:"primaryKey;autoIncrement" json:"id"`
	Name         string            `gorm:"type:varchar(128);not null;index" json:"name"`
	DeviceType   string            `gorm:"type:varchar(64);not null" json:"device_type"`
	Vendor       string            `gorm:"type:varchar(32)" json:"vendor"`
	Model        string            `gorm:"type:varchar(64)" json:"model"`
	Role         string            `gorm:"type:varchar(32);default:access" json:"role"`
	ManagementIP string            `gorm:"type:varchar(45);not null;uniqueIndex" json:"management_ip"`
	SSHPort      int               `gorm:"default:22" json:"ssh_port"`
	CredentialID *uint             `gorm:"index" json:"credential_id"`
	Credential   *Credential       `gorm:"foreignKey:CredentialID" json:"credential,omitempty"`
	GroupID      *uint             `gorm:"index" json:"group_id"`
	Group        *DeviceGroup      `gorm:"foreignKey:GroupID" json:"group,omitempty"`
	Building     string            `gorm:"type:varchar(64)" json:"building"`
	Floor        int               `gorm:"default:0" json:"floor"`
	APName       string            `gorm:"type:varchar(128)" json:"ap_name"`
	SNMPCommunity string           `gorm:"type:varchar(64)" json:"snmp_community"`
	Status       int               `gorm:"default:0;index" json:"status"`
	LastSeen     *time.Time        `json:"last_seen"`
	Description  string            `gorm:"type:varchar(256)" json:"description"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
}

func (Device) TableName() string {
	return "devices"
}
