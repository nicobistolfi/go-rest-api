package api

import (
	"net/http"
	"os"

	"go-rest-api/pkg/auth" // Adjust this import path as needed

	"github.com/gin-gonic/gin"
)

// ProfileResponse represents the structure of the profile data
type ProfileResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

func GetToken(c *gin.Context) {
	secretKey := []byte(os.Getenv("JWT_SECRET")) // Replace with a secure secret key
	token, err := auth.GenerateJWT(secretKey)
	if err != nil {
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
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
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
