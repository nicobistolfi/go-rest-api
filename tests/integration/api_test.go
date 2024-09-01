package integration

import (
	"encoding/json"
	"go-boilerplate/internal/api"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAPIEndpoints(t *testing.T) {
	// Setup the router
	router := api.SetupRouter()

	// Create a test HTTP server
	server := httptest.NewServer(router)
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
			endpoint:       "/ping",
			expectedStatus: http.StatusOK,
			expectedBody:   map[string]string{"message": "pong"},
		},
		// Add more test cases for other endpoints here
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Make a request to the test server
			resp, err := http.Get(server.URL + tc.endpoint)
			assert.NoError(t, err, "Failed to make request")
			defer resp.Body.Close()

			// Check status code
			assert.Equal(t, tc.expectedStatus, resp.StatusCode, "Unexpected status code")

			// Read and parse the response body
			body, err := ioutil.ReadAll(resp.Body)
			assert.NoError(t, err, "Failed to read response body")

			var responseBody map[string]string
			err = json.Unmarshal(body, &responseBody)
			assert.NoError(t, err, "Failed to parse response body")

			// Check response body
			assert.Equal(t, tc.expectedBody, responseBody, "Unexpected response body")
		})
	}
}
