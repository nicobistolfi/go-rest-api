package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var token string
		var authHeader string

		// Check for Authorization header
		token = c.GetHeader("Authorization")
		authHeader = "Authorization"

		// If not found, check for X-API-Key header
		if token == "" {
			token = c.GetHeader("X-API-Key")
			authHeader = "X-API-Key"
		}

		// If still not found, check for access_token query parameter
		if token == "" {
			token = c.Query("access_token")
			authHeader = "access_token"
		}

		// If no token found in any of the above methods
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication token is required"})
			c.Abort()
			return
		}

		// Add the token to the context
		c.Set("auth_token", token)
		c.Set("auth_header", authHeader)

		c.Next()
	}
}
