package controllers

import (
	"encoding/json"
	"net/http"

	"crackers/d2delight.com/initializers"
	"crackers/d2delight.com/models"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// ---------------- CREATE ----------------
// Create or Append CustomerProfile (only for customer role)
func CreateCustomerProfile(c *gin.Context) {
	// Get logged-in user
	u, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	user := u.(models.User)

	// Get role
	var role models.Role
	if err := initializers.DB.First(&role, user.RoleID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch role"})
		return
	}
	if role.Name != "customer" {
		c.JSON(http.StatusForbidden, gin.H{"error": "only customers can create addresses"})
		return
	}

	// Parse request
	var input struct {
		Address interface{} `json:"address" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var profile models.CustomerProfile
	err := initializers.DB.Where("user_id = ?", user.ID).First(&profile).Error

	// Convert new address to JSON
	newAddrJSON, _ := json.Marshal(input.Address)

	if err == gorm.ErrRecordNotFound {
		// No profile → create new
		addresses := []json.RawMessage{newAddrJSON}
		finalJSON, _ := json.Marshal(addresses)

		profile = models.CustomerProfile{
			UserID:    user.ID,
			Addresses: datatypes.JSON(finalJSON),
		}
		if err := initializers.DB.Create(&profile).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else if err == nil {
		// Profile exists → append
		var addresses []json.RawMessage
		_ = json.Unmarshal(profile.Addresses, &addresses)

		addresses = append(addresses, newAddrJSON)
		finalJSON, _ := json.Marshal(addresses)

		profile.Addresses = datatypes.JSON(finalJSON)
		if err := initializers.DB.Save(&profile).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, profile)
}

// ---------------- READ ----------------
// Get all customer profiles (admin only)
func GetCustomerProfiles(c *gin.Context) {
	var profiles []models.CustomerProfile
	if err := initializers.DB.Find(&profiles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, profiles)
}

// Get a single customer profile by ID
func GetCustomerProfileByID(c *gin.Context) {
	id := c.Param("id")
	var profile models.CustomerProfile
	if err := initializers.DB.First(&profile, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "profile not found"})
		return
	}
	c.JSON(http.StatusOK, profile)
}

// ---------------- UPDATE ----------------
// Replace entire addresses list
func UpdateCustomerProfile(c *gin.Context) {
	id := c.Param("id")
	var profile models.CustomerProfile
	if err := initializers.DB.First(&profile, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "profile not found"})
		return
	}

	var input struct {
		Addresses interface{} `json:"addresses" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newJSON, _ := json.Marshal(input.Addresses)
	profile.Addresses = datatypes.JSON(newJSON)

	if err := initializers.DB.Save(&profile).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, profile)
}

// ---------------- DELETE ----------------
// Delete a customer profile
func DeleteCustomerProfile(c *gin.Context) {
	id := c.Param("id")
	if err := initializers.DB.Delete(&models.CustomerProfile{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "profile deleted"})
}
