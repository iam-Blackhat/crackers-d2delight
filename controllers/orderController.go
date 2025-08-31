package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"crackers/d2delight.com/initializers"
	"crackers/d2delight.com/models"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

var user models.User
var deliveryAddress string

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

	// Check if user exists and has role CUSTOMER
	if err := initializers.DB.Where("id = ? AND role_id = (SELECT id FROM roles WHERE name = ?)", input.CustomerID, "CUSTOMER").First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid customer ID"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
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

	if err := initializers.DB.Preload("Customer").
		Preload("CustomerProfile").
		Find(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var addresses []string
	if err := json.Unmarshal(order.CustomerProfile.Addresses, &addresses); err != nil {
		// handle error
	}

	// Pick the right address
	if int(order.DeliveryAddressId) < len(addresses) {
		deliveryAddress := addresses[int(order.DeliveryAddressId)]
		order.DeliveryAddress = deliveryAddress
	} else {
		fmt.Println("Invalid address index")
	}
	c.JSON(http.StatusCreated, order)
}

// Get All Orders
func GetOrders(c *gin.Context) {
	var orders []models.Order
	if err := initializers.DB.
		Preload("Customer").
		Preload("OrderItems.Product"). // load order items + product
		Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
}

// Get Order by ID
func GetOrderByID(c *gin.Context) {
	id := c.Param("id")
	var order models.Order

	if err := initializers.DB.
		Preload("Customer").
		Preload("OrderItems.Product"). // load order items + product
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
