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

	// Create a middleware to track call count
	var callCount int
	countMiddleware := func(c *gin.Context) {
		if c.Writer.Header().Get("X-Token-Cache") == "MISS" {
			callCount++
			t.Logf("Mock server called. Count: %d", callCount)
		}
		c.Next()
	}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	r.Use(countMiddleware) // Add the count middleware
	r.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	makeRequest := func() {
		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer valid_token")
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		t.Logf("Request made. Status: %d", resp.Code)
	}

	t.Log("Making first request")
	makeRequest() // First request, should call the mock server
	if callCount != 1 {
		t.Errorf("Expected 1 call to mock server, got %d", callCount)
	}

	// time.Sleep(100 * time.Millisecond) // Small delay

	t.Log("Making second request")
	makeRequest() // Second request, should use cache
	if callCount != 1 {
		t.Errorf("Expected still 1 call to mock server, got %d", callCount)
	}

	t.Log("Waiting for cache to expire")
	time.Sleep(3 * time.Second) // Wait for cache to expire

	t.Log("Making third request")
	makeRequest() // Third request, should call the mock server again
	if callCount != 2 {
		t.Errorf("Expected 2 calls to mock server, got %d", callCount)
	}

	t.Logf("Final call count: %d", callCount)
}
