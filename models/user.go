package models

import (
	"time"
)

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	FullName  string    `json:"full_name" binding:"required"`
	Email     string    `gorm:"uniqueIndex;not null" json:"email" binding:"required,email"`
	Password  string    `json:"-" binding:"required,min=6"`
	Role      string    `gorm:"default:user" json:"role"` // user/admin
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Orders    []Order   `json:"orders,omitempty"`
}