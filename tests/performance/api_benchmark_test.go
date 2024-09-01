package performance

import (
	"go-boilerplate/internal/api"
	"net/http"
	"net/http/httptest"
	"testing"
)

func BenchmarkPingEndpoint(b *testing.B) {
	router := api.SetupRouter()
	server := httptest.NewServer(router)
	defer server.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp, err := http.Get(server.URL + "/ping")
		if err != nil {
			b.Fatalf("Failed to make request: %v", err)
		}
		resp.Body.Close()
	}
}
