package security

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-boilerplate/internal/api"
	"go-boilerplate/internal/config"
	"go-boilerplate/pkg/auth"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"gopkg.in/h2non/gock.v1"
)

func setupTestRouter() (*gin.Engine, *config.Config) {
	gin.SetMode(gin.TestMode)
	cfg := &config.Config{
		JWTSecret:         "test_secret",
		ValidAPIKey:       "test_api_key",
		OIDCIssuer:        "https://test-oidc-provider.com",
		OAuthClientID:     "test_client_id",
		OAuthClientSecret: "test_client_secret",
		OAuthRedirectURL:  "http://localhost:8080/auth/callback",
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
	return r, cfg
}

func setupTestRouterOAuth() (*gin.Engine, *config.Config) {
	gin.SetMode(gin.TestMode)
	cfg := &config.Config{
		JWTSecret:         "test_secret",
		ValidAPIKey:       "test_api_key",
		OIDCIssuer:        "https://test-oidc-provider.com",
		OAuthClientID:     "test_client_id",
		OAuthClientSecret: "test_client_secret",
		OAuthRedirectURL:  "http://localhost:8080/auth/callback",
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
	token, err := auth.GenerateJWTToken("test_user", "test@example.com", cfg.JWTSecret, 6000)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/jwt/profile", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response api.ProfileResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Add more detailed assertions and error logging
	if assert.NotEmpty(t, response.ID, "ID should not be empty") {
		assert.Equal(t, "test_user", response.ID, "Unexpected ID value")
	}
	if assert.NotEmpty(t, response.Email, "Email should not be empty") {
		assert.Equal(t, "test@example.com", response.Email, "Unexpected Email value")
	}

	// Log the entire response for debugging
	t.Logf("Response body: %s", w.Body.String())
}

func TestGetProfileWithAPIKey(t *testing.T) {
	router, cfg := setupTestRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/apikey/profile", nil)
	req.Header.Set("X-API-Key", cfg.ValidAPIKey)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response api.ProfileResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "api_user", response.ID)
	assert.Equal(t, "api@example.com", response.Email)
}

// Note: Testing OAuth requires mocking the OIDC provider, which is more complex.
// For this example, we'll skip the OAuth test, but in a real-world scenario,
// you should mock the OIDC provider and test the OAuth flow as well.

func TestGetProfileWithOAuth(t *testing.T) {
	router, cfg := setupTestRouterOAuth()

	// Mock OIDC provider
	defer gock.Off()
	gock.New(cfg.OIDCIssuer).
		Get("/.well-known/openid-configuration").
		Reply(200).
		JSON(map[string]interface{}{
			"issuer":                 cfg.OIDCIssuer,
			"authorization_endpoint": cfg.OIDCIssuer + "/auth",
			"token_endpoint":         cfg.OIDCIssuer + "/token",
			"userinfo_endpoint":      cfg.OIDCIssuer + "/userinfo",
			"jwks_uri":               cfg.OIDCIssuer + "/jwks",
		})

	gock.New(cfg.OIDCIssuer).
		Get("/userinfo").
		MatchHeader("Authorization", "Bearer valid_oauth_token").
		Reply(200).
		JSON(map[string]interface{}{
			"sub":   "oauth_user_id",
			"email": "oauth_user@example.com",
			"name":  "OAuth User",
		})

	// Initialize OAuth
	err := auth.InitOAuth(cfg)
	assert.NoError(t, err, "Failed to initialize OAuth")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/oauth/profile", nil)
	req.Header.Set("Authorization", "Bearer valid_oauth_token")
	router.ServeHTTP(w, req)

	// Check response status
	assert.Equal(t, http.StatusOK, w.Code, "Unexpected status code")

	// Check response headers
	contentType := w.Header().Get("Content-Type")
	assert.Contains(t, contentType, "application/json", "Unexpected Content-Type")

	// Parse response body
	var response api.ProfileResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err, "Failed to unmarshal response body")

	// Check response fields
	assert.Equal(t, "oauth_user_id", response.ID, "Unexpected user ID")
	assert.Equal(t, "oauth_user@example.com", response.Email, "Unexpected email")
	// assert.Equal(t, "OAuth User", response.Name, "Unexpected name")

	// Additional checks
	assert.NotEmpty(t, response.ID, "User ID should not be empty")
	assert.Contains(t, response.Email, "@", "Email should contain '@'")
	// assert.True(t, len(response.Name) > 0, "Name should not be empty")

	// Log response for debugging
	t.Logf("Response body: %s", w.Body.String())
}

// ... (keep existing code)
