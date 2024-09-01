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
	api.SetupRouter(router, cfg, logger)
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
