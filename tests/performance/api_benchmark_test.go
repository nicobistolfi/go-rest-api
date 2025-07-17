package performance

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nicobistolfi/go-rest-api/internal/api"
	"github.com/nicobistolfi/go-rest-api/internal/config"
	logger "github.com/nicobistolfi/go-rest-api/pkg"
)

func BenchmarkPingEndpoint(b *testing.B) {
	cfg, err := config.LoadConfig()
	if err != nil {
		b.Fatalf("Failed to load configuration: %v", err)
	}

	logger.Init()

	router := gin.New()
	api.SetupRouter(router, cfg, logger.Log, api.WithoutRateLimiting())

	server := httptest.NewServer(router)
	defer server.Close()

	b.ResetTimer()

	for range b.N {
		resp, pingErr := http.Get(server.URL + "/api/v1/ping")
		if pingErr != nil {
			b.Fatalf("Failed to make request: %v", pingErr)
		}

		resp.Body.Close()
	}
}

func BenchmarkRateLimiting(b *testing.B) {
	cfg, err := config.LoadConfig()
	if err != nil {
		b.Fatalf("Failed to load configuration: %v", err)
	}

	logger.Init()

	router := gin.New()
	api.SetupRouter(router, cfg, logger.Log) // Use default setup with rate limiting

	server := httptest.NewServer(router)
	defer server.Close()

	client := &http.Client{}

	b.ResetTimer()
	// for i := 0; i < b.N; i++ {
	// Send multiple requests in quick succession
	for j := range 11 {
		resp, clientErr := client.Get(server.URL + "/api/v1/ping")
		if clientErr != nil {
			b.Fatalf("Failed to make request: %v", clientErr)
		}

		// Check if rate limiting is triggered (expecting 429 status code)
		if j == 11 && resp.StatusCode != http.StatusTooManyRequests {
			b.Errorf("Expected rate limiting to be triggered (status 429), got %d", resp.StatusCode)
		}

		resp.Body.Close()
	}
	// }
}
