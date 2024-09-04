package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestVerifyToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup a mock token server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "Bearer valid_token" {
			profile := Profile{
				ID:    "123",
				Email: "test@example.com",
				Name:  "Test User",
			}
			json.NewEncoder(w).Encode(profile)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
		}
	}))
	defer mockServer.Close()

	os.Setenv("TOKEN_URL", mockServer.URL)
	defer os.Unsetenv("TOKEN_URL")

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("auth_token", c.GetHeader("Authorization"))
		c.Set("auth_header", "Authorization")
		c.Next()
	})
	r.Use(VerifyToken())
	r.GET("/test", func(c *gin.Context) {
		user, _ := c.Get("user")
		c.JSON(http.StatusOK, user)
	})

	tests := []struct {
		name           string
		token          string
		expectedStatus int
	}{
		{"Valid token", "Bearer valid_token", http.StatusOK},
		{"Invalid token", "Bearer invalid_token", http.StatusUnauthorized},
		{"Missing token", "", http.StatusUnauthorized},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/test", nil)
			if tt.token != "" {
				req.Header.Set("Authorization", tt.token)
			}
			resp := httptest.NewRecorder()
			r.ServeHTTP(resp, req)

			if resp.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, resp.Code)
			}

			if tt.expectedStatus == http.StatusOK {
				var profile Profile
				json.NewDecoder(resp.Body).Decode(&profile)
				if profile.ID != "123" || profile.Email != "test@example.com" || profile.Name != "Test User" {
					t.Errorf("Unexpected profile data")
				}
			}
		})
	}
}

func TestVerifyTokenCaching(t *testing.T) {
	gin.SetMode(gin.TestMode)

	callCount := 0
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		profile := Profile{
			ID:    "123",
			Email: "test@example.com",
			Name:  "Test User",
		}
		json.NewEncoder(w).Encode(profile)
	}))
	defer mockServer.Close()

	os.Setenv("TOKEN_URL", mockServer.URL)
	os.Setenv("TOKEN_CACHE_EXPIRY", "1s")

	defer os.Unsetenv("TOKEN_URL")
	defer os.Unsetenv("TOKEN_CACHE_EXPIRY")

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("auth_token", c.GetHeader("Authorization"))
		c.Set("auth_header", "Authorization")
		c.Next()
	})
	r.Use(VerifyToken(os.Getenv("TOKEN_CACHE_EXPIRY")))
	r.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	makeRequest := func() {
		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer valid_token")
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
	}

	makeRequest() // First request, should call the mock server
	if callCount != 1 {
		t.Errorf("Expected 1 call to mock server, got %d", callCount)
	}

	makeRequest() // Second request, should use cache
	if callCount != 1 {
		t.Errorf("Expected still 1 call to mock server, got %d", callCount)
	}

	time.Sleep(3 * time.Second) // Wait for cache to expire

	makeRequest() // Third request, should call the mock server again
	if callCount != 2 {
		t.Errorf("Expected 2 calls to mock server, got %d", callCount)
	}
}
