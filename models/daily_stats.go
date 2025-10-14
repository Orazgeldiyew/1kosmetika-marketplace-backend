package models

import "time"

// DailyStats — ежедневная статистика для админа
type DailyStats struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	Date          string    `json:"date" gorm:"uniqueIndex"`
	TotalUsers    int64     `json:"total_users"`
	TotalOrders   int64     `json:"total_orders"`
	TotalProducts int64     `json:"total_products"`
	TotalRevenue  float64   `json:"total_revenue"`
	CreatedAt     time.Time `json:"created_at"`
}
