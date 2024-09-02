package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	// OAuth configuration
	OIDCIssuer        string
	OAuthClientID     string
	OAuthClientSecret string
	OAuthRedirectURL  string

	// JWT configuration
	JWTSecret            string
	JWTExpirationMinutes int

	// API Key configuration
	ValidAPIKey string

	// Rate Limiting configuration
	RateLimitRequests int
	RateLimitDuration time.Duration

	// Other configuration options
	// ...
}

func LoadConfig() (*Config, error) {
	// Load .env file if it exists
	_ = godotenv.Load() // Ignore error if .env file doesn't exist

	config := &Config{
		OIDCIssuer:        getEnv("OIDC_ISSUER", "https://accounts.google.com"),
		OAuthClientID:     getEnv("OAUTH_CLIENT_ID", ""),
		OAuthClientSecret: getEnv("OAUTH_CLIENT_SECRET", ""),
		OAuthRedirectURL:  getEnv("OAUTH_REDIRECT_URL", "http://localhost:8080/auth/callback"),

		JWTSecret:            getEnv("JWT_SECRET", "default_jwt_secret"),
		JWTExpirationMinutes: getEnvAsInt("JWT_EXPIRATION_MINUTES", 60),

		ValidAPIKey: getEnv("VALID_API_KEY", "default_api_key"),

		RateLimitRequests: getEnvAsInt("RATE_LIMIT_REQUESTS", 10),
		RateLimitDuration: getEnvAsDuration("RATE_LIMIT_DURATION", time.Second),
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr := getEnv(key, "")
	if value, err := time.ParseDuration(valueStr); err == nil {
		return value
	}
	return defaultValue
}
