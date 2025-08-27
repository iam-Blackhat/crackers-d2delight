package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"crackers/d2delight.com/initializers"
	"crackers/d2delight.com/models"
)

// Input struct for creating/updating users
type Input struct {
	Name     string    `json:"name" binding:"required,notnumeric"`
	Email    string    `json:"email" binding:"required,email"`
	Phone    string    `json:"phone" binding:"required"`
	Password string    `json:"password" binding:"required,min=6"`
	RoleID   uuid.UUID `json:"role_id"` // Reference to roles table
}

type UpdateInput struct {
	Name     string `json:"name" binding:"required,notnumeric"`
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

// Create a new user
func Create(c *gin.Context) {
	var input Input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Password hashing failed"})
		return
	}

	// If RoleID is empty, assign default "customer" role
	if input.RoleID == uuid.Nil {
		var customerRole models.Role
		if err := initializers.DB.Where("name = ?", "customer").First(&customerRole).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Default customer role not found"})
			return
		}
		input.RoleID = customerRole.ID
	}

	user := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Phone:    input.Phone,
		Password: string(hashedPassword),
		RoleID:   input.RoleID,
	}

	if err := initializers.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user": gin.H{
			"id":     user.ID,
			"name":   user.Name,
			"email":  user.Email,
			"phone":  user.Phone,
			"roleId": user.RoleID,
		},
	})
}

// Get all users
func GetUsers(c *gin.Context) {
	var users []models.User
	initializers.DB.Preload("Role").Find(&users) // preload role if you want role info
	c.JSON(http.StatusOK, gin.H{"users": users})
}

// Get user by ID
func GetUserByID(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	if err := initializers.DB.Preload("Role").First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": user})
}

// Update user (name, role, password if provided)
func Update(c *gin.Context) {
	id := c.Param("id")

	var user models.User
	if err := initializers.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var input UpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedFields := map[string]interface{}{}

	if input.Name != "" {
		updatedFields["name"] = input.Name
	}
	if input.Phone != "" {
		updatedFields["phone"] = input.Phone
	}
	if input.Password != "" {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		updatedFields["password"] = string(hashedPassword)
	}

	initializers.DB.Model(&user).Updates(updatedFields)
	c.JSON(http.StatusOK, gin.H{"user": user})
}

// Delete user
func Delete(c *gin.Context) {
	id := c.Param("id")
	if err := initializers.DB.Delete(&models.User{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}
