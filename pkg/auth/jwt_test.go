package auth

import (
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestGenerateJWT(t *testing.T) {
	secretKey := []byte("test-secret-key")

	token, err := GenerateJWT(secretKey)
	if err != nil {
		t.Fatalf("Failed to generate JWT: %v", err)
	}

	// Parse and validate the token
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		t.Fatalf("Failed to parse JWT: %v", err)
	}

	if !parsedToken.Valid {
		t.Fatalf("Generated token is not valid")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatalf("Failed to get claims from token")
	}

	// Check if the claims contain the expected user data
	if claims["user_id"] != "123456" {
		t.Errorf("Expected user_id '123456', got '%v'", claims["user_id"])
	}
	if claims["username"] != "user" {
		t.Errorf("Expected username 'user', got '%v'", claims["username"])
	}
	if claims["email"] != "@example.com" {
		t.Errorf("Expected email '@example.com', got '%v'", claims["email"])
	}

	// Check if the expiration time is set correctly
	exp, ok := claims["exp"].(float64)
	if !ok {
		t.Fatalf("Failed to get expiration time from token")
	}
	expectedExp := time.Now().Add(time.Minute * time.Duration(getJWTExpirationMinutes())).Unix()
	if int64(exp) != expectedExp {
		t.Errorf("Expected expiration time %v, got %v", expectedExp, int64(exp))
	}
}

func TestGetJWTExpirationMinutes(t *testing.T) {
	// Test default value
	os.Unsetenv("JWT_EXPIRATION_MINUTES")
	if exp := getJWTExpirationMinutes(); exp != 1440 {
		t.Errorf("Expected default expiration 1440, got %d", exp)
	}

	// Test custom value
	os.Setenv("JWT_EXPIRATION_MINUTES", "60")
	if exp := getJWTExpirationMinutes(); exp != 60 {
		t.Errorf("Expected expiration 60, got %d", exp)
	}

	// Test invalid value
	os.Setenv("JWT_EXPIRATION_MINUTES", "invalid")
	if exp := getJWTExpirationMinutes(); exp != 1440 {
		t.Errorf("Expected default expiration 1440 for invalid input, got %d", exp)
	}

	// Clean up
	os.Unsetenv("JWT_EXPIRATION_MINUTES")
}
