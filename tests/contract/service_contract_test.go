package contract

import (
	"encoding/json"
	"go-boilerplate/internal/api"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServiceContract(t *testing.T) {
	// Setup the router
	router := api.SetupRouter()

	// Create a test HTTP server
	server := httptest.NewServer(router)
	defer server.Close()

	// Test the /ping endpoint
	t.Run("Ping Endpoint", func(t *testing.T) {
		// Make a GET request to the /ping endpoint
		resp, err := http.Get(server.URL + "/ping")
		assert.NoError(t, err, "Failed to make request to /ping")
		defer resp.Body.Close()

		// Check the status code
		assert.Equal(t, http.StatusOK, resp.StatusCode, "Unexpected status code for /ping")

		// Check the response body
		var response map[string]string
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err, "Failed to decode response body")
		assert.Equal(t, "pong", response["message"], "Unexpected response message for /ping")
	})

	// Add more contract tests for other endpoints here as they are implemented
}
