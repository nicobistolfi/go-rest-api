package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

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
			if err := json.NewEncoder(w).Encode(profile); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		} else {
			w.WriteHeader(http.StatusUnauthorized)
		}
	}))
	defer mockServer.Close()

	os.Setenv("TOKEN_URL", mockServer.URL)

	defer func() {
		if err := os.Unsetenv("TOKEN_URL"); err != nil { //nolint:staticcheck // Empty test cleanup is acceptable
			// Ignore error in test cleanup
		}
	}()

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
			req, _ := http.NewRequest(http.MethodGet, "/test", nil)
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
				if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
					t.Errorf("Failed to decode response: %v", err)
				}

				if profile.ID != "123" || profile.Email != "test@example.com" || profile.Name != "Test User" {
					t.Errorf("Unexpected profile data")
				}
			}
		})
	}
}
