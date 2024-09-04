package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCORSMiddlewareFunc(t *testing.T) {
	// Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	// Create a new Gin engine
	r := gin.New()

	// Use the CORSMiddleware
	r.Use(CORSMiddleware())

	// Define a test route
	r.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "test")
	})

	tests := []struct {
		name            string
		method          string
		origin          string
		expectedStatus  int
		expectedHeaders map[string]string
	}{
		{
			name:           "OPTIONS request",
			method:         "OPTIONS",
			origin:         "http://example.com",
			expectedStatus: http.StatusOK,
			expectedHeaders: map[string]string{
				"Access-Control-Allow-Origin":      "http://example.com",
				"Access-Control-Allow-Methods":     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
				"Access-Control-Allow-Headers":     "Authorization,Content-Type",
				"Access-Control-Allow-Credentials": "true",
			},
		},
		{
			name:           "GET request",
			method:         "GET",
			origin:         "http://example.com",
			expectedStatus: http.StatusOK,
			expectedHeaders: map[string]string{
				"Access-Control-Allow-Origin": "http://example.com",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("ALLOWED_ORIGINS", tt.origin)
			req, _ := http.NewRequest(tt.method, "/test", nil)
			req.Header.Set("Origin", tt.origin)
			resp := httptest.NewRecorder()

			r.ServeHTTP(resp, req)

			if resp.Code != tt.expectedStatus {
				t.Errorf("Expected status %d; got %d", tt.expectedStatus, resp.Code)
			}

			for key, value := range tt.expectedHeaders {
				if resp.Header().Get(key) != value {
					t.Errorf("Expected header %s to be %s; got %s", key, value, resp.Header().Get(key))
				}
			}
		})
	}
}
