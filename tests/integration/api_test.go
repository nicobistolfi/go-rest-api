package integration

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nicobistolfi/go-rest-api/internal/api"
	"github.com/nicobistolfi/go-rest-api/internal/config"
	logger "github.com/nicobistolfi/go-rest-api/pkg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAPIEndpoints(t *testing.T) {
	// Setup the router
	cfg, err := config.LoadConfig()
	require.NoError(t, err, "Failed to load configuration")

	logger.Init()

	r := gin.New()
	api.SetupRouter(r, cfg, logger.Log)

	// Create a test HTTP server
	server := httptest.NewServer(r)
	defer server.Close()

	// Test cases
	testCases := []struct {
		name           string
		endpoint       string
		expectedStatus int
		expectedBody   map[string]string
	}{
		{
			name:           "Ping Endpoint",
			endpoint:       "/api/v1/ping",
			expectedStatus: http.StatusOK,
			expectedBody:   map[string]string{"message": "pong"},
		},
		// Add more test cases for other endpoints here
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Make a request to the test server
			resp, reqErr := http.Get(server.URL + tc.endpoint)
			require.NoError(t, reqErr, "Failed to make request")

			defer func() {
				if closeErr := resp.Body.Close(); closeErr != nil {
					t.Logf("Failed to close response body: %v", closeErr)
				}
			}()

			// Check status code
			assert.Equal(t, tc.expectedStatus, resp.StatusCode, "Unexpected status code")

			// Read and parse the response body
			body, readErr := io.ReadAll(resp.Body)
			require.NoError(t, readErr, "Failed to read response body")

			var responseBody map[string]string
			parseErr := json.Unmarshal(body, &responseBody)
			require.NoError(t, parseErr, "Failed to parse response body")

			// Check response body
			assert.Equal(t, tc.expectedBody, responseBody, "Unexpected response body")
		})
	}
}
