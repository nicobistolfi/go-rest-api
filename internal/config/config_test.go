package config

import (
	"os"
	"testing"
	"time"
)

func TestLoadConfig(t *testing.T) {
	// Set up test environment variables
	os.Setenv("OIDC_ISSUER", "https://test.issuer.com")
	os.Setenv("OAUTH_CLIENT_ID", "test_client_id")
	os.Setenv("OAUTH_CLIENT_SECRET", "test_client_secret")
	os.Setenv("OAUTH_REDIRECT_URL", "http://test.redirect.url")
	os.Setenv("JWT_SECRET", "test_jwt_secret")
	os.Setenv("JWT_EXPIRATION_MINUTES", "120")
	os.Setenv("VALID_API_KEY", "test_api_key")
	os.Setenv("RATE_LIMIT_REQUESTS", "20")
	os.Setenv("RATE_LIMIT_DURATION", "1m")

	// Load the configuration
	config, err := LoadConfig()

	// Check for errors
	if err != nil {
		t.Fatalf("LoadConfig() returned an error: %v", err)
	}

	// Test each configuration value
	tests := []struct {
		name     string
		got      interface{}
		expected interface{}
	}{
		{"OIDCIssuer", config.OIDCIssuer, "https://test.issuer.com"},
		{"OAuthClientID", config.OAuthClientID, "test_client_id"},
		{"OAuthClientSecret", config.OAuthClientSecret, "test_client_secret"},
		{"OAuthRedirectURL", config.OAuthRedirectURL, "http://test.redirect.url"},
		{"JWTSecret", config.JWTSecret, "test_jwt_secret"},
		{"JWTExpirationMinutes", config.JWTExpirationMinutes, 120},
		{"ValidAPIKey", config.ValidAPIKey, "test_api_key"},
		{"RateLimitRequests", config.RateLimitRequests, 20},
		{"RateLimitDuration", config.RateLimitDuration, time.Minute},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.expected {
				t.Errorf("%s = %v, want %v", tt.name, tt.got, tt.expected)
			}
		})
	}
}

func TestGetEnv(t *testing.T) {
	// Test with existing environment variable
	os.Setenv("TEST_ENV_VAR", "test_value")
	if got := getEnv("TEST_ENV_VAR", "default"); got != "test_value" {
		t.Errorf("getEnv() = %v, want %v", got, "test_value")
	}

	// Test with non-existing environment variable
	if got := getEnv("NON_EXISTING_VAR", "default"); got != "default" {
		t.Errorf("getEnv() = %v, want %v", got, "default")
	}
}

func TestGetEnvAsInt(t *testing.T) {
	// Test with valid integer
	os.Setenv("TEST_INT_VAR", "42")
	if got := getEnvAsInt("TEST_INT_VAR", 0); got != 42 {
		t.Errorf("getEnvAsInt() = %v, want %v", got, 42)
	}

	// Test with invalid integer
	os.Setenv("TEST_INVALID_INT", "not_an_int")
	if got := getEnvAsInt("TEST_INVALID_INT", 10); got != 10 {
		t.Errorf("getEnvAsInt() = %v, want %v", got, 10)
	}
}

func TestGetEnvAsDuration(t *testing.T) {
	// Test with valid duration
	os.Setenv("TEST_DURATION_VAR", "5m")
	if got := getEnvAsDuration("TEST_DURATION_VAR", time.Second); got != 5*time.Minute {
		t.Errorf("getEnvAsDuration() = %v, want %v", got, 5*time.Minute)
	}

	// Test with invalid duration
	os.Setenv("TEST_INVALID_DURATION", "not_a_duration")
	if got := getEnvAsDuration("TEST_INVALID_DURATION", time.Hour); got != time.Hour {
		t.Errorf("getEnvAsDuration() = %v, want %v", got, time.Hour)
	}
}
