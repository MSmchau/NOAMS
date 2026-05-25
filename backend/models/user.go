package models

import (
	"time"
)

type User struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Username  string    `gorm:"type:varchar(64);uniqueIndex;not null" json:"username"`
	Password  string    `gorm:"type:varchar(256);not null" json:"-"`
	Email     string    `gorm:"type:varchar(128)" json:"email"`
	Phone     string    `gorm:"type:varchar(32)" json:"phone"`
	Nickname  string    `gorm:"type:varchar(64)" json:"nickname"`
	Role      string    `gorm:"type:varchar(16);default:operator" json:"role"`
	Status    int       `gorm:"default:1" json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}
