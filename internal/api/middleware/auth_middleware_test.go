package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestAuthMiddlewareFunc(t *testing.T) {
	// Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	// Create a new Gin engine
	r := gin.New()

	// Use the AuthMiddleware
	protected := r.Group("/")
	protected.Use(AuthMiddleware())
	// Define a test route
	protected.GET("/protected", func(c *gin.Context) {
		c.String(http.StatusOK, "protected")
	})

	tests := []struct {
		name           string
		token          string
		expectedStatus int
	}{
		{"Valid Token", "Bearer valid_token", http.StatusOK},
		{"Invalid Token", "Bearer invalid_token", http.StatusOK}, // This is expected to pass, as the auth only checks for the token presence
		{"No Token", "", http.StatusUnauthorized},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/protected", nil)
			if tt.token != "" {
				req.Header.Set("Authorization", tt.token)
			}
			resp := httptest.NewRecorder()

			r.ServeHTTP(resp, req)

			if resp.Code != tt.expectedStatus {
				t.Errorf("Expected status %d; got %d", tt.expectedStatus, resp.Code)
			}
		})
	}
}
