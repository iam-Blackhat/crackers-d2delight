package controllers

import (
	"net/http"

	"crackers/d2delight.com/initializers"
	"crackers/d2delight.com/models"

	"github.com/gin-gonic/gin"
)

// Create Customer Profile
func CreateCustomerProfile(c *gin.Context) {
	var input models.CustomerProfile
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := initializers.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, input)
}

// Get All Customer Profiles
func GetCustomerProfiles(c *gin.Context) {
	var customers []models.CustomerProfile
	if err := initializers.DB.Find(&customers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch customers"})
		return
	}
	c.JSON(http.StatusOK, customers)
}

// Get Customer Profile by ID
func GetCustomerProfileByID(c *gin.Context) {
	id := c.Param("id")
	var customer models.CustomerProfile

	if err := initializers.DB.First(&customer, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer profile not found"})
		return
	}

	c.JSON(http.StatusOK, customer)
}

// Update Customer Profile
func UpdateCustomerProfile(c *gin.Context) {
	id := c.Param("id")
	var customer models.CustomerProfile

	if err := initializers.DB.First(&customer, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer profile not found"})
		return
	}

	var input models.CustomerProfile
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := initializers.DB.Model(&customer).Updates(input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update customer profile"})
		return
	}

	c.JSON(http.StatusOK, customer)
}

// Delete Customer Profile
func DeleteCustomerProfile(c *gin.Context) {
	id := c.Param("id")
	if err := initializers.DB.Delete(&models.CustomerProfile{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete customer profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Customer profile deleted successfully"})
}
