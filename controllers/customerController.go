// controllers/customerController.go
package controllers

import (
	"net/http"

	"crackers/d2delight.com/initializers"
	"crackers/d2delight.com/models"

	"github.com/gin-gonic/gin"
)

// Create Customer
func CreateCustomer(c *gin.Context) {
	var input models.Customer
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := initializers.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create customer"})
		return
	}

	c.JSON(http.StatusCreated, input)
}

// Get All Customers
func GetCustomers(c *gin.Context) {
	var customers []models.Customer
	initializers.DB.Find(&customers)
	c.JSON(http.StatusOK, customers)
}

// Get Customer by ID
func GetCustomerByID(c *gin.Context) {
	id := c.Param("id")
	var customer models.Customer

	if err := initializers.DB.First(&customer, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}

	c.JSON(http.StatusOK, customer)
}

// Update Customer
func UpdateCustomer(c *gin.Context) {
	id := c.Param("id")
	var customer models.Customer

	if err := initializers.DB.First(&customer, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}

	var input models.Customer
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	initializers.DB.Model(&customer).Updates(input)
	c.JSON(http.StatusOK, customer)
}

// Delete Customer
func DeleteCustomer(c *gin.Context) {
	id := c.Param("id")
	if err := initializers.DB.Delete(&models.Customer{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete customer"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Customer deleted successfully"})
}
