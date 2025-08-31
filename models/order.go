package models

import "time"

type Order struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	CustomerID        uint      `json:"customer_id"`
	DeliveryAddressId uint      `json:"delivery_address_id"`
	DeliveryAddress   string    `json:"delivery_address" gorm:"-"`
	Total             float64   `json:"total"`
	Status            string    `json:"status"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`

	// Relations
	Customer        User            `gorm:"foreignKey:CustomerID;references:ID" json:"customer"`
	CustomerProfile CustomerProfile `gorm:"foreignKey:UserID;references:CustomerID" json:"customer_profile"`
}

type OrderItem struct {
	ID        uint    `gorm:"primaryKey" json:"id"`
	OrderID   uint    `json:"order_id"`
	ProductID uint    `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`

	Product Product `gorm:"foreignKey:ProductID" json:"product"`
}
