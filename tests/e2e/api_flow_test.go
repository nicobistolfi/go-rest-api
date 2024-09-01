package e2e

import (
	"encoding/json"
	"go-boilerplate/internal/api"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAPIFlow(t *testing.T) {
	// Setup the router
	router := api.SetupRouter()

	// Create a test HTTP server
	server := httptest.NewServer(router)
	defer server.Close()

	// Step 1: Ping
	t.Run("Ping", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/ping")
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]string
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "pong", response["message"])
	})

	/*
		// Step 2: Create Resource
		t.Run("Create Resource", func(t *testing.T) {
			payload := []byte(`{"name": "Test Resource"}`)
			resp, err := http.Post(server.URL+"/resources", "application/json", bytes.NewBuffer(payload))
			assert.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusCreated, resp.StatusCode)

			var response map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&response)
			assert.NoError(t, err)
			assert.Equal(t, "Test Resource", response["name"])
			assert.NotEmpty(t, response["id"])
		})

		// Step 3: Retrieve Resource
		var resourceID string
		t.Run("Retrieve Resource", func(t *testing.T) {
			resp, err := http.Get(server.URL + "/resources/1")
			assert.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusOK, resp.StatusCode)

			var response map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&response)
			assert.NoError(t, err)
			assert.Equal(t, "Test Resource", response["name"])
			resourceID = response["id"].(string)
		})

		// Step 4: Update Resource
		t.Run("Update Resource", func(t *testing.T) {
			payload := []byte(`{"name": "Updated Test Resource"}`)
			req, _ := http.NewRequest(http.MethodPut, server.URL+"/resources/"+resourceID, bytes.NewBuffer(payload))
			req.Header.Set("Content-Type", "application/json")
			resp, err := http.DefaultClient.Do(req)
			assert.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusOK, resp.StatusCode)

			var response map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&response)
			assert.NoError(t, err)
			assert.Equal(t, "Updated Test Resource", response["name"])
		})

		// Step 5: Delete Resource
		t.Run("Delete Resource", func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodDelete, server.URL+"/resources/"+resourceID, nil)
			resp, err := http.DefaultClient.Do(req)
			assert.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusNoContent, resp.StatusCode)
		})
	*/
}
