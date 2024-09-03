package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/nicobistolfi/go-rest-api/internal/api/middleware"
)

func TestGetToken(t *testing.T) {
	r := gin.Default()

	r.GET("/token", GetToken)

	os.Setenv("JWT_SECRET", "test_secret")
	os.Setenv("JWT_EXPIRATION_MINUTES", "1")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/token", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "token")
}

func TestGetProfile(t *testing.T) {
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Set("user", middleware.Profile{ID: "1", Email: "test@example.com", Name: "Test User"})
		c.Next()
	})
	r.GET("/profile", GetProfile)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/profile", nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response middleware.Profile
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "1", response.ID)
	assert.Equal(t, "test@example.com", response.Email)
	assert.Equal(t, "Test User", response.Name)
}

func TestHealthCheck(t *testing.T) {
	r := gin.Default()

	r.GET("/health", HealthCheck)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "OK", response["status"])
}

func TestPing(t *testing.T) {
	r := gin.Default()

	r.GET("/ping", Ping)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "pong", response["message"])
}

func TestRegister(t *testing.T) {
	r := gin.Default()

	r.POST("/register", Register)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/register", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "User registered successfully", response["message"])
}
