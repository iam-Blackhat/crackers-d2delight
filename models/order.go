package models

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	CustomerID uint           `gorm:"not null" json:"customer_id"`
	Total      float64        `gorm:"not null" json:"total"`
	Status     string         `gorm:"type:varchar(50);default:'pending'" json:"status"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
	Customer   Customer       `json:"customer"`
	Products   []Product      `gorm:"many2many:order_items;" json:"products"`
}

type OrderItem struct {
	OrderID   uint    `gorm:"primaryKey"`
	ProductID uint    `gorm:"primaryKey"`
	Quantity  int     `gorm:"not null"`
	Price     float64 `gorm:"not null"`
}
