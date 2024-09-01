package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-boilerplate/internal/config"
	"go-boilerplate/pkg/auth"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func setupTestRouter() (*gin.Engine, *config.Config) {
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

	SetupRouter(r, cfg, logger)
	return r, cfg
}

func TestGetProfileUnauthorized(t *testing.T) {
	router, _ := setupTestRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/oauth/profile", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestGetProfileWithJWT(t *testing.T) {
	router, cfg := setupTestRouter()

	// Generate a valid JWT token
	token, err := auth.GenerateJWTToken("test_user", "test@example.com", cfg.JWTSecret, 60)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/jwt/profile", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response ProfileResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "test_user", response.ID)
	assert.Equal(t, "test@example.com", response.Email)
}

func TestGetProfileWithAPIKey(t *testing.T) {
	router, cfg := setupTestRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/apikey/profile", nil)
	req.Header.Set("X-API-Key", cfg.ValidAPIKey)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response ProfileResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "api_user", response.ID)
	assert.Equal(t, "api@example.com", response.Email)
}

// Note: Testing OAuth requires mocking the OIDC provider, which is more complex.
// For this example, we'll skip the OAuth test, but in a real-world scenario,
// you should mock the OIDC provider and test the OAuth flow as well.
