package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func TestRateLimiter(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		limit          rate.Limit
		burst          int
		requests       int
		expectedStatus int
	}{
		{"Under limit", rate.Limit(10), 5, 5, http.StatusOK},
		{"At limit", rate.Limit(10), 10, 10, http.StatusOK},
		{"Over limit", rate.Limit(10), 10, 11, http.StatusTooManyRequests},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			r.Use(RateLimiter(tt.limit, tt.burst, tt.name))
			r.GET("/test", func(c *gin.Context) {
				c.String(http.StatusOK, "test")
			})

			for i := 0; i < tt.requests; i++ {
				req, _ := http.NewRequest("GET", "/test", nil)
				resp := httptest.NewRecorder()
				r.ServeHTTP(resp, req)

				t.Logf("Request %d: Status Code %d", i+1, resp.Code)

				if i == tt.requests-1 {
					if resp.Code != tt.expectedStatus {
						t.Errorf("Expected status %d, got %d", tt.expectedStatus, resp.Code)
					}
				} else if resp.Code != http.StatusOK {
					t.Errorf("Request %d: Expected status %d, got %d", i+1, http.StatusOK, resp.Code)
				}
			}
		})
	}
}

func TestRateLimiterWithDifferentClients(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(RateLimiter(rate.Limit(1), 1))
	r.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "test")
	})

	clients := []string{"1.1.1.1", "2.2.2.2"}

	for _, client := range clients {
		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("X-Forwarded-For", client)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		if resp.Code != http.StatusOK {
			t.Errorf("Expected status %d for client %s, got %d", http.StatusOK, client, resp.Code)
		}
	}
}

func TestRateLimiterWithAuthHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(RateLimiter(rate.Limit(1), 1))
	r.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "test")
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer token1")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.Code)
	}

	// Same IP, different token
	req, _ = http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer token2")
	resp = httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.Code)
	}
}

func TestRateLimiterBurst(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(RateLimiter(rate.Limit(1), 3)) // 1 request per second, burst of 3
	r.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "test")
	})

	for i := 0; i < 4; i++ {
		req, _ := http.NewRequest("GET", "/test", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		if i < 3 {
			if resp.Code != http.StatusOK {
				t.Errorf("Expected status %d for request %d, got %d", http.StatusOK, i+1, resp.Code)
			}
		} else {
			if resp.Code != http.StatusTooManyRequests {
				t.Errorf("Expected status %d for request %d, got %d", http.StatusTooManyRequests, i+1, resp.Code)
			}
		}
	}
}

func TestRateLimiterRecovery(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(RateLimiter(rate.Limit(1), 1)) // 1 request per second
	r.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "test")
	})

	// First request should succeed
	req, _ := http.NewRequest("GET", "/test", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	if resp.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.Code)
	}

	// Second request should fail
	resp = httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	if resp.Code != http.StatusTooManyRequests {
		t.Errorf("Expected status %d, got %d", http.StatusTooManyRequests, resp.Code)
	}

	// Wait for rate limit to reset
	time.Sleep(1100 * time.Millisecond)

	// Third request should succeed
	resp = httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	if resp.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.Code)
	}
}
