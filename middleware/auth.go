package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"crackers/d2delight.com/initializers"
	"crackers/d2delight.com/models"
	"crackers/d2delight.com/utils"
)

type Role struct {
	ID   uuid.UUID `db:"id" json:"id"`
	Name string    `db:"name" json:"name"`
}

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Expect Authorization: Bearer <token>
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
			return
		}

		parts := strings.SplitN(auth, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header"})
			return
		}

		tokenStr := parts[1]
		claims, err := utils.ParseToken(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		// Load user from DB
		var user models.User
		if err := initializers.DB.First(&user, claims.UserID).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
			return
		}

		var role Role
		if err := initializers.DB.First(&role, user.RoleID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch role"})
			return
		}

		// attach user to context
		c.Set("currentUser", user)
		c.Next()
	}
}

func RoleRequired(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		u, exists := c.Get("currentUser")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user not found in context"})
			return
		}

		user := u.(models.User)

		// Fetch role using RoleID
		var role models.Role
		if err := initializers.DB.First(&role, user.RoleID).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch role"})
			return
		}

		// Check if role matches allowed roles
		for _, allowedRole := range allowedRoles {
			if role.Name == allowedRole {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "access denied"})
	}
}

// Helper to retrieve current user in handlers
func CurrentUser(c *gin.Context) (models.User, bool) {
	u, exists := c.Get("currentUser")
	if !exists {
		return models.User{}, false
	}
	user, ok := u.(models.User)
	return user, ok
}
