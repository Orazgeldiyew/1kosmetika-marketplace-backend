package models

import "time"

type Order struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	UserID        uint      `json:"user_id"`
	User          User      `gorm:"foreignKey:UserID" json:"user"`
	Products      []Product `gorm:"many2many:order_products;" json:"products"`
	Total         float64   `json:"total"`
	Status        string    `gorm:"default:pending" json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	TotalAmount   float64   `json:"total_amount"` 
	PaymentMethod string    `json:"payment_method"`
}

type OrderProduct struct {
    OrderID   uint    `gorm:"primaryKey"`
    ProductID uint    `gorm:"primaryKey"`
    Quantity  int     `gorm:"not null;default:1"`
    Price     float64 `gorm:"not null;default:0"`
}
