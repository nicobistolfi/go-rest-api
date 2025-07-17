package api

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	logger "github.com/nicobistolfi/go-rest-api/pkg"
	"github.com/nicobistolfi/go-rest-api/pkg/auth" // Adjust this import path as needed
	"go.uber.org/zap"
)

func GetToken(c *gin.Context) {
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "JWT_SECRET is not set"})

		return
	}

	token, err := auth.GenerateJWT([]byte(secretKey))
	if err != nil {
		logger.Error("Error generating JWT", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

// GetProfile handles the /profile endpoint.
func GetProfile(c *gin.Context) {
	// Retrieve the user from the context
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found in context"})

		return
	}

	c.JSON(http.StatusOK, user)
}

// HealthCheck handles the /health endpoint.
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "OK",
	})
}

// Ping handles the /ping endpoint.
func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

// Register handles the /register endpoint.
func Register(c *gin.Context) {
	// This is a placeholder implementation for user registration
	c.JSON(http.StatusOK, gin.H{
		"message": "User registered successfully",
	})
}
