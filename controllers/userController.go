package controllers

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/datatypes"

	"crackers/d2delight.com/initializers"
	"crackers/d2delight.com/models"
	"crackers/d2delight.com/utils"
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

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user with the provided input
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.User true "User Data"
// @Success 201 {object} models.User
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/register [post]
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

	// If RoleID is empty, assign default "CUSTOMER" role
	if input.RoleID == uuid.Nil {
		var customerRole models.Role
		if err := initializers.DB.Where("name = ?", "CUSTOMER").First(&customerRole).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Default customer role not found"})
			return
		}
		input.RoleID = customerRole.ID
	}

	// Generate OTP
	otp := utils.GenerateNumericOTP(6)
	otpHash := utils.HashOTP(otp)

	meta := map[string]any{
		"name":    input.Name,
		"email":   input.Email,
		"phone":   input.Phone,
		"pw_hash": string(hashedPassword),
		"role_id": input.RoleID.String(),
	}
	metaJSON, _ := json.Marshal(meta)

	// Delete old OTPs
	initializers.DB.Where("phone = ? AND purpose = ?", input.Phone, "register").
		Delete(&models.PhoneOTP{})

	otpRow := models.PhoneOTP{
		Phone:       input.Phone,
		Purpose:     "register",
		CodeHash:    otpHash,
		Meta:        datatypes.JSON(metaJSON),
		Attempts:    0,
		MaxAttempts: 5,
		ExpiresAt:   time.Now().Add(5 * time.Minute),
	}
	initializers.DB.Create(&otpRow)

	// Send SMS
	_ = utils.SendSMS(input.Phone, "Your OTP is: "+otp)

	c.JSON(http.StatusOK, gin.H{"message": "OTP sent"})

}

func VerifyRegisterOTP(c *gin.Context) {
	var input struct {
		Phone string `json:"phone" binding:"required"`
		Code  string `json:"code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var otpRow models.PhoneOTP
	if err := initializers.DB.Where("phone = ? AND purpose = ?", input.Phone, "register").
		First(&otpRow).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "OTP not found"})
		return
	}

	if time.Now().After(otpRow.ExpiresAt) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "OTP expired"})
		return
	}

	if otpRow.Attempts >= otpRow.MaxAttempts {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Too many attempts"})
		return
	}

	if !utils.CheckOTP(input.Code, otpRow.CodeHash) {
		initializers.DB.Model(&otpRow).Update("attempts", otpRow.Attempts+1)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid OTP"})
		return
	}

	// Parse metadata
	var meta map[string]any
	json.Unmarshal(otpRow.Meta, &meta)

	// Create user
	user := models.User{
		Name:     meta["name"].(string),
		Email:    meta["email"].(string),
		Phone:    meta["phone"].(string),
		Password: meta["pw_hash"].(string),
		RoleID:   uuid.MustParse(meta["role_id"].(string)),
	}

	if err := initializers.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User create failed"})
		return
	}
	// Clean OTP row
	initializers.DB.Delete(&otpRow)

	var role Role
	if err := initializers.DB.First(&role, user.RoleID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch role"})
		return
	}

	token, err := utils.CreateToken(user.ID, user.Email, user.Name, role.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create token"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"access_token": token,
		"token_type":   "bearer",
		"expires_in":   os.Getenv("JWT_EXPIRES_MINUTES"),
		"user": gin.H{
			"id":     user.ID,
			"name":   user.Name,
			"email":  user.Email,
			"phone":  user.Phone,
			"roleId": user.RoleID,
		},
		"message": "User registered successfully",
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
