package controllers

import (
	"net/http"

	"crackers/d2delight.com/initializers"
	"crackers/d2delight.com/models"

	"github.com/gin-gonic/gin"
)

// Create Order
func CreateOrder(c *gin.Context) {
	var input struct {
		CustomerID uint `json:"customer_id"`
		AddressID  uint `json:"address_id"` // delivery address reference
		Products   []struct {
			ProductID uint    `json:"product_id"`
			Quantity  int     `json:"quantity"`
			Price     float64 `json:"price"`
		} `json:"products"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Calculate total
	var total float64
	for _, p := range input.Products {
		total += float64(p.Quantity) * p.Price
	}

	// Create order
	order := models.Order{
		CustomerID:        input.CustomerID,
		DeliveryAddressId: input.AddressID, // save delivery address
		Total:             total,
		Status:            "pending",
	}

	if err := initializers.DB.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Create order items
	for _, p := range input.Products {
		orderItem := models.OrderItem{
			OrderID:   order.ID,
			ProductID: p.ProductID,
			Quantity:  p.Quantity,
			Price:     p.Price,
		}
		initializers.DB.Create(&orderItem)
	}

	c.JSON(http.StatusCreated, order)
}

// Get All Orders
func GetOrders(c *gin.Context) {
	var orders []models.Order
	initializers.DB.
		Preload("Customer").
		Preload("Products").
		Preload("Address"). // preload delivery address
		Find(&orders)

	c.JSON(http.StatusOK, orders)
}

// Get Order by ID
func GetOrderByID(c *gin.Context) {
	id := c.Param("id")
	var order models.Order

	if err := initializers.DB.
		Preload("Customer").
		Preload("Products").
		Preload("Address"). // preload delivery address
		First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	c.JSON(http.StatusOK, order)
}

// Update Order
func UpdateOrder(c *gin.Context) {
	id := c.Param("id")
	var order models.Order

	if err := initializers.DB.First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	var input struct {
		Status    string `json:"status"`
		AddressID uint   `json:"address_id"` // allow updating delivery address
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order.Status = input.Status
	if input.AddressID != 0 {
		order.DeliveryAddressId = input.AddressID
	}

	initializers.DB.Save(&order)

	c.JSON(http.StatusOK, order)
}

// Delete Order
func DeleteOrder(c *gin.Context) {
	id := c.Param("id")
	if err := initializers.DB.Delete(&models.Order{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete order"})
		return
	}

	// Also delete associated order_items
	initializers.DB.Where("order_id = ?", id).Delete(&models.OrderItem{})

	c.JSON(http.StatusOK, gin.H{"message": "Order deleted successfully"})
}
