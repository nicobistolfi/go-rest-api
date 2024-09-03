package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-rest-api/internal/api"
	"go-rest-api/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	logger "go-rest-api/pkg"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	cfg := &config.Config{
		JWTSecret:   "test_secret",
		ValidAPIKey: "test_api_key",
	}

	logger.Init()

	r := gin.New()

	// Add the config and logger to the gin.Context
	r.Use(func(c *gin.Context) {
		c.Set("config", cfg)
		c.Set("logger", logger.Log)
		c.Next()
	})

	api.SetupRouter(r, cfg, logger.Log)
	return r
}

func TestPingHandler(t *testing.T) {
	// Setup the router
	r := setupTestRouter()

	// Create a mock request to the /ping endpoint
	req, err := http.NewRequest(http.MethodGet, "/api/v1/ping", nil)
	assert.NoError(t, err)

	// Create a response recorder
	w := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(w, req)

	// Assert the status code
	assert.Equal(t, http.StatusOK, w.Code)

	// Parse the response body
	var response map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Assert the response body
	assert.Equal(t, "pong", response["message"])
}

func TestHealthCheckHandler(t *testing.T) {
	// Setup the router
	r := setupTestRouter()

	// Create a mock request to the /health endpoint
	req, err := http.NewRequest(http.MethodGet, "/api/v1/health", nil)
	assert.NoError(t, err)

	// Create a response recorder
	w := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(w, req)

	// Assert the status code
	assert.Equal(t, http.StatusOK, w.Code)

	// Parse the response body
	var response map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Assert the response body
	assert.Equal(t, "OK", response["status"])
}
