package auth

import (
	"errors"

	"go-boilerplate/internal/config"
	"go-boilerplate/internal/models"
)

func ValidateAPIKey(apiKey string, cfg *config.Config) (*models.User, error) {
	// In a real-world scenario, you would typically check the API key against a database
	// For this example, we'll use a simple in-memory check
	if apiKey == cfg.ValidAPIKey {
		return &models.User{
			ID:    "api_user",
			Email: "api@example.com",
			Name:  "API User",
		}, nil
	}

	return nil, errors.New("invalid API key")
}
