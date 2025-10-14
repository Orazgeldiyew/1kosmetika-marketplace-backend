package models

import (
	"time"
)

type Product struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"not null" json:"name" binding:"required"`
	Description string    `json:"description"`
	Price       float64   `gorm:"not null" json:"price" binding:"required,gt=0"`
	ImageURL    string    `json:"image_url"`
	Category    string    `json:"category" binding:"required"`
	Brand       string    `json:"brand" binding:"required"`
	Stock       int       `gorm:"default:0" json:"stock"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}