package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"crackers/d2delight.com/initializers"
	"crackers/d2delight.com/models"
	"crackers/d2delight.com/utils"
)

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

		// attach user to context
		c.Set("currentUser", user)
		c.Next()
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
