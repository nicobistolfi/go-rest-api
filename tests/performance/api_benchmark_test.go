package performance

import (
	"go-boilerplate/internal/api"
	"go-boilerplate/internal/config"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func BenchmarkPingEndpoint(b *testing.B) {
	cfg, err := config.LoadConfig()
	if err != nil {
		b.Fatalf("Failed to load configuration: %v", err)
	}

	logger, err := zap.NewProduction()
	if err != nil {
		b.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	router := gin.New()
	api.SetupRouter(router, cfg, logger, api.WithoutRateLimiting())
	server := httptest.NewServer(router)
	defer server.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp, err := http.Get(server.URL + "/api/v1/ping")
		if err != nil {
			b.Fatalf("Failed to make request: %v", err)
		}
		resp.Body.Close()
	}
}

func BenchmarkRateLimiting(b *testing.B) {
	cfg, err := config.LoadConfig()
	if err != nil {
		b.Fatalf("Failed to load configuration: %v", err)
	}

	logger, err := zap.NewProduction()
	if err != nil {
		b.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	router := gin.New()
	api.SetupRouter(router, cfg, logger) // Use default setup with rate limiting

	server := httptest.NewServer(router)
	defer server.Close()

	client := &http.Client{}

	b.ResetTimer()
	// for i := 0; i < b.N; i++ {
	// Send multiple requests in quick succession
	for j := 0; j < 11; j++ {
		resp, err := client.Get(server.URL + "/api/v1/ping")
		if err != nil {
			b.Fatalf("Failed to make request: %v", err)
		}

		// Check if rate limiting is triggered (expecting 429 status code)
		if j == 11 && resp.StatusCode != http.StatusTooManyRequests {
			b.Errorf("Expected rate limiting to be triggered (status 429), got %d", resp.StatusCode)
		}

		resp.Body.Close()
	}
	// }
}
