package api

import (
	"fmt"
	"net/http"
	"os"

	"github.com/nicobistolfi/go-rest-api/pkg/auth" // Adjust this import path as needed

	"github.com/gin-gonic/gin"
)

func GetToken(c *gin.Context) {
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "JWT_SECRET is not set"})
		return
	}

	token, err := auth.GenerateJWT([]byte(secretKey))
	if err != nil {
		fmt.Printf("Error generating JWT: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

// GetProfile handles the /profile endpoint
func GetProfile(c *gin.Context) {
	// Retrieve the user from the context
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found in context"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// HealthCheck handles the /health endpoint
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "OK",
	})
}

// Ping handles the /ping endpoint
func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

// Register handles the /register endpoint
func Register(c *gin.Context) {
	// TODO: Implement user registration logic
	c.JSON(http.StatusOK, gin.H{
		"message": "User registered successfully",
	})
}
