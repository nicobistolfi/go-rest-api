package security

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/nicobistolfi/go-rest-api/internal/api"
	"github.com/nicobistolfi/go-rest-api/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	logger "github.com/nicobistolfi/go-rest-api/pkg"
)

func setupRouter() *gin.Engine {
	cfg, _ := config.LoadConfig()
	logger.Init()
	r := gin.New()
	api.SetupRouter(r, cfg, logger.Log)
	return r
}

func TestAPISecurityEndpoints(t *testing.T) {
	// Setup mock token server
	mockTokenServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "valid_token" {
			err := json.NewEncoder(w).Encode(map[string]string{"id": "123", "email": "test@example.com", "name": "Test User"})
			if err != nil {
				t.Errorf("Failed to encode JSON response: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else {
			w.WriteHeader(http.StatusUnauthorized)
		}
	}))
	defer mockTokenServer.Close()

	// Set TOKEN_URL environment variable
	os.Setenv("TOKEN_URL", mockTokenServer.URL)
	defer os.Unsetenv("TOKEN_URL")

	router := setupRouter()

	testCases := []struct {
		name           string
		endpoint       string
		method         string
		token          string
		expectedStatus int
	}{
		{"Public Endpoint - Health Check", "/api/v1/health", "GET", "", http.StatusOK},
		{"Public Endpoint - Ping", "/api/v1/ping", "GET", "", http.StatusOK},
		{"Protected Endpoint - No Token", "/api/v1/profile", "GET", "", http.StatusUnauthorized},
		{"Protected Endpoint - Invalid Token", "/api/v1/profile", "GET", "invalid_token", http.StatusUnauthorized},
		{"Protected Endpoint - Valid Token", "/api/v1/profile", "GET", "valid_token", http.StatusOK},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest(tc.method, tc.endpoint, nil)
			if tc.token != "" {
				req.Header.Set("Authorization", tc.token)
			}
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code)
		})
	}
}
