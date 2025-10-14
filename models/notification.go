package models

import (
	"time"
	
	"gorm.io/gorm"
)

type Notification struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    uint           `json:"user_id" gorm:"not null"`
	Title     string         `json:"title" gorm:"size:255;not null"`
	Message   string         `json:"message" gorm:"type:text"`
	IsRead    bool           `json:"is_read" gorm:"default:false"`
	Type      string         `json:"type" gorm:"size:100;default:info"` // info, warning, success, error
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

func (Notification) TableName() string {
	return "notifications"
}