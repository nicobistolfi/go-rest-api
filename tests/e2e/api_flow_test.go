package e2e

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

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	cfg := &config.Config{
		JWTSecret:   "test_secret",
		ValidAPIKey: "test_api_key",
	}
	logger, _ := zap.NewDevelopment()

	r := gin.New()

	// Add the config and logger to the gin.Context
	r.Use(func(c *gin.Context) {
		c.Set("config", cfg)
		c.Set("logger", logger)
		c.Next()
	})

	api.SetupRouter(r, cfg, logger)
	return r
}

func TestAPIFlow(t *testing.T) {
	// Setup the router
	router := setupTestRouter()

	// Create a test HTTP server
	server := httptest.NewServer(router.Handler())
	defer server.Close()

	// Step 1: Ping
	t.Run("Ping", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/api/v1/ping")
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]string
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "pong", response["message"])
	})

	// Step 2: Health Check
	t.Run("Health Check", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/api/v1/health")
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]string
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "OK", response["status"])
	})

	// Additional tests for authenticated routes can be added here
	// You'll need to generate valid tokens/keys for each auth method
	// and include them in the request headers

	/*
		// Example: Test JWT authenticated route
		t.Run("JWT Authenticated Route", func(t *testing.T) {
			// Generate a valid JWT token
			token := generateValidJWTToken()

			req, _ := http.NewRequest("GET", server.URL+"/api/v1/jwt/profile", nil)
			req.Header.Set("Authorization", "Bearer "+token)
			resp, err := http.DefaultClient.Do(req)
			assert.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusOK, resp.StatusCode)

			var profile api.ProfileResponse
			err = json.NewDecoder(resp.Body).Decode(&profile)
			assert.NoError(t, err)
			assert.NotEmpty(t, profile.ID)
			assert.NotEmpty(t, profile.Email)
		})
	*/
}
