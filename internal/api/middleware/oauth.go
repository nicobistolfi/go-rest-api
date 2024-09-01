package middleware

import (
	"net/http"

	"go-boilerplate/internal/config"
	"go-boilerplate/pkg/auth"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func OAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := c.MustGet("logger").(*zap.Logger)
		cfg := c.MustGet("config").(*config.Config)

		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing authorization token"})
			c.Abort()
			return
		}

		userInfo, err := auth.ValidateOAuthToken(c.Request.Context(), token, cfg)
		if err != nil {
			logger.Error("Failed to validate OAuth token", zap.Error(err))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("user", userInfo)
		c.Next()
	}
}
