// controllers/productController.go
package controllers

import (
	"net/http"

	"crackers/d2delight.com/initializers"
	"crackers/d2delight.com/models"

	"github.com/gin-gonic/gin"
)

// Create Product
func CreateProduct(c *gin.Context) {
	var input models.Product
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := initializers.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusCreated, input)
}

// Get All Products
func GetProducts(c *gin.Context) {
	var products []models.Product
	if err := initializers.DB.Preload("Category").Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, products)
}

// Get Product by ID
func GetProductByID(c *gin.Context) {
	id := c.Param("id")
	var product models.Product

	if err := initializers.DB.Preload("Category").Find(&product, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

// Update Product
func UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	var product models.Product

	if err := initializers.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	var input models.Product
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	initializers.DB.Model(&product).Updates(input)
	c.JSON(http.StatusOK, product)
}

// Delete Product
func DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	if err := initializers.DB.Delete(&models.Product{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}
