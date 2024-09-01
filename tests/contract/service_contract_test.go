package contract

import (
	"encoding/json"
	"go-boilerplate/internal/api"
	"go-boilerplate/internal/config"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestServiceContract(t *testing.T) {
	// Setup the router
	cfg, err := config.LoadConfig()
	assert.NoError(t, err, "Failed to load configuration")

	logger, err := zap.NewProduction()
	assert.NoError(t, err, "Failed to initialize logger")
	defer logger.Sync()

	r := gin.New()
	api.SetupRouter(r, cfg, logger)

	// Create a test HTTP server
	server := httptest.NewServer(r)
	defer server.Close()

	// Test the /ping endpoint
	t.Run("Ping Endpoint", func(t *testing.T) {
		// Make a GET request to the /ping endpoint
		resp, err := http.Get(server.URL + "/api/v1/ping")
		assert.NoError(t, err, "Failed to make request to /api/v1/ping")
		defer resp.Body.Close()

		// Check the status code
		assert.Equal(t, http.StatusOK, resp.StatusCode, "Unexpected status code for /api/v1/ping")

		// Check the response body
		var response map[string]string
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err, "Failed to decode response body")
		assert.Equal(t, "pong", response["message"], "Unexpected response message for/api/v1/ping")
	})

	// Add more contract tests for other endpoints here as they are implemented
}
