package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	customLogger "github.com/nicobistolfi/go-rest-api/pkg"
)

func TestLoggerMiddlewareFunc(t *testing.T) {
	customLogger.Init()
	// Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	// Create a new Gin engine
	r := gin.New()

	// Create a custom logger
	logger := customLogger.Log

	// Use the LoggerMiddleware
	r.Use(LoggerMiddleware(logger))

	// Define a test route
	r.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "test")
	})

	// Create a test request
	req, _ := http.NewRequest("GET", "/test", nil)
	resp := httptest.NewRecorder()

	// Serve the request
	r.ServeHTTP(resp, req)

	// Check the status code
	if resp.Code != http.StatusOK {
		t.Errorf("Expected status %d; got %d", http.StatusOK, resp.Code)
	}

	// Add more assertions here to check logging behavior
	// For example, you could use a custom io.Writer to capture log output
	// and assert on its contents
}
