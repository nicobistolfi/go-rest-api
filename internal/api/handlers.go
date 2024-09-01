package api

import (
	"net/http"

	"go-boilerplate/internal/models"
	"go-boilerplate/pkg/auth" // Adjust this import path as needed

	"github.com/coreos/go-oidc"
	"github.com/gin-gonic/gin"
)

// ProfileResponse represents the structure of the profile data
type ProfileResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

// GetProfile handles the /profile endpoint
func GetProfile(c *gin.Context) {
	// Retrieve the user from the context
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	var profile ProfileResponse

	// Determine the type of user and populate the profile accordingly
	switch u := user.(type) {
	case *oidc.UserInfo: // OAuth user
		profile = ProfileResponse{
			ID:    u.Subject,
			Email: u.Email,
			Name:  u.Profile,
		}
	case *auth.Claims: // JWT user
		profile = ProfileResponse{
			ID:    u.UserID,
			Email: u.Email,
			Name:  "JWT User", // You might want to store and retrieve the name in your JWT claims
		}
	case *models.User: // API Key user
		profile = ProfileResponse{
			ID:    u.ID,
			Email: u.Email,
			Name:  u.Name,
		}
	default: // Unknown type
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unknown user type"})
		return
	}

	c.JSON(http.StatusOK, profile)
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
