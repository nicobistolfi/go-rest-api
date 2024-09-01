package middleware

import (
	"net/http"

	"go-boilerplate/internal/config"
	"go-boilerplate/pkg/auth"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func APIKeyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := c.MustGet("logger").(*zap.Logger)
		cfg := c.MustGet("config").(*config.Config)

		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing API key"})
			c.Abort()
			return
		}

		user, err := auth.ValidateAPIKey(apiKey, cfg)
		if err != nil {
			logger.Error("Failed to validate API key", zap.Error(err))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key"})
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Next()
	}
}
