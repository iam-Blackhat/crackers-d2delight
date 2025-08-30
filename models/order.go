package models

type Order struct {
	ID                uint    `gorm:"primaryKey" json:"id"`
	CustomerID        uint    `json:"customer_id"`
	DeliveryAddressId uint    `json:"delivery_address_id"`
	Total             float64 `json:"total"`
	Status            string  `json:"status"`

	// Relationships
	Customer   User        `gorm:"foreignKey:CustomerID" json:"customer"`
	OrderItems []OrderItem `json:"order_items"`
}

type OrderItem struct {
	ID        uint    `gorm:"primaryKey" json:"id"`
	OrderID   uint    `json:"order_id"`
	ProductID uint    `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`

	Product Product `gorm:"foreignKey:ProductID" json:"product"`
}
