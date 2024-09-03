---
title: Go REST API Testing Documentation
---

# Go REST API Testing Documentation

## Overview

This document outlines the comprehensive testing strategy for our Go REST API. We employ various testing approaches to ensure the reliability, security, and performance of our API. The testing suite includes unit tests, integration tests, end-to-end tests, security tests, and performance benchmarks.

## Table of Contents

1. [Unit Tests](#unit-tests)
2. [API Security Tests](#api-security-tests)
3. [Service Contract Tests](#service-contract-tests)
4. [API Flow Tests](#api-flow-tests)
5. [API Integration Tests](#api-integration-tests)
6. [API Performance Benchmarks](#api-performance-benchmarks)

## Unit Tests

File: `service_test.go`

### Testing Strategy

Unit tests focus on testing individual components or functions of the API in isolation. These tests ensure that:

1. Individual handlers work correctly
2. The API returns expected responses for specific inputs
3. Basic functionality is maintained as the codebase evolves

### Key Components

- Setup of a test router with mock configuration
- Individual test functions for each handler
- Assertions for status codes and response bodies

### Code Example

```go
func setupTestRouter() *gin.Engine {
    gin.SetMode(gin.TestMode)
    cfg := &config.Config{
        JWTSecret:   "test_secret",
        ValidAPIKey: "test_api_key",
    }
    logger, _ := zap.NewDevelopment()

    r := gin.New()

    // Add the config and logger to the gin.Context
    r.Use(func(c *gin.Context) {
        c.Set("config", cfg)
        c.Set("logger", logger)
        c.Next()
    })

    api.SetupRouter(r, cfg, logger)
    return r
}

func TestPingHandler(t *testing.T) {
    r := setupTestRouter()

    req, err := http.NewRequest(http.MethodGet, "/api/v1/ping", nil)
    assert.NoError(t, err)

    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

    assert.Equal(t, http.StatusOK, w.Code)

    var response map[string]string
    err = json.Unmarshal(w.Body.Bytes(), &response)
    assert.NoError(t, err)

    assert.Equal(t, "pong", response["message"])
}

func TestHealthCheckHandler(t *testing.T) {
    r := setupTestRouter()

    req, err := http.NewRequest(http.MethodGet, "/api/v1/health", nil)
    assert.NoError(t, err)

    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

    assert.Equal(t, http.StatusOK, w.Code)

    var response map[string]string
    err = json.Unmarshal(w.Body.Bytes(), &response)
    assert.NoError(t, err)

    assert.Equal(t, "OK", response["status"])
}
```

### Explanation

1. `setupTestRouter()`: This function creates a test environment by setting up a Gin router with a mock configuration and logger. It's used to isolate the tests from the actual application environment.

2. `TestPingHandler()`: This test verifies the `/api/v1/ping` endpoint:
   - It creates a mock GET request to the endpoint.
   - Uses `httptest.NewRecorder()` to capture the response.
   - Asserts that the status code is 200 (OK).
   - Checks that the response body contains the expected "pong" message.

3. `TestHealthCheckHandler()`: This test verifies the `/api/v1/health` endpoint:
   - Similar to the ping test, it creates a mock GET request.
   - Asserts that the status code is 200 (OK).
   - Checks that the response body contains the expected "OK" status.

These unit tests ensure that the basic endpoints of the API are functioning correctly. They provide a quick way to verify that changes to the codebase haven't broken core functionality.

## API Security Tests

File: `api_security_test.go`

### Testing Strategy

The API security tests focus on verifying the authentication and authorization mechanisms of our API. These tests ensure that:

1. Public endpoints are accessible without authentication
2. Protected endpoints require valid authentication
3. Invalid or missing authentication tokens are properly rejected

### Key Components

- Mock token server for simulating authentication responses
- Test cases for various scenarios (public endpoints, protected endpoints with valid/invalid tokens)

### Code Example

```go
func TestAPISecurityEndpoints(t *testing.T) {
    // Setup mock token server
    mockTokenServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.Header.Get("Authorization") == "valid_token" {
            json.NewEncoder(w).Encode(map[string]string{"id": "123", "email": "test@example.com", "name": "Test User"})
        } else {
            w.WriteHeader(http.StatusUnauthorized)
        }
    }))
    defer mockTokenServer.Close()

    // Set TOKEN_URL environment variable
    os.Setenv("TOKEN_URL", mockTokenServer.URL)
    defer os.Unsetenv("TOKEN_URL")

    router := setupRouter()

    testCases := []struct {
        name           string
        endpoint       string
        method         string
        token          string
        expectedStatus int
    }{
        {"Public Endpoint - Health Check", "/api/v1/health", "GET", "", http.StatusOK},
        {"Protected Endpoint - No Token", "/api/v1/profile", "GET", "", http.StatusUnauthorized},
        {"Protected Endpoint - Valid Token", "/api/v1/profile", "GET", "valid_token", http.StatusOK},
    }

    // Test execution code...
}
```

## Service Contract Tests

File: `service_contract_test.go`

### Testing Strategy

Service contract tests verify that our API adheres to its defined contract. These tests ensure that:

1. Endpoints return the expected status codes
2. Response bodies match the expected structure and content
3. The API behaves consistently across different environments

### Key Components

- Setup of the router with test configuration
- Test cases for each endpoint, checking status codes and response bodies

### Code Example

```go
func TestServiceContract(t *testing.T) {
    // Setup the router
    cfg, err := config.LoadConfig()
    assert.NoError(t, err, "Failed to load configuration")

    logger, err := zap.NewProduction()
    assert.NoError(t, err, "Failed to initialize logger")
    defer logger.Sync()

    r := gin.New()
    api.SetupRouter(r, cfg, logger)

    // Create a test HTTP server
    server := httptest.NewServer(r)
    defer server.Close()

    // Test the /ping endpoint
    t.Run("Ping Endpoint", func(t *testing.T) {
        resp, err := http.Get(server.URL + "/api/v1/ping")
        assert.NoError(t, err, "Failed to make request to /api/v1/ping")
        defer resp.Body.Close()

        assert.Equal(t, http.StatusOK, resp.StatusCode, "Unexpected status code for /api/v1/ping")

        var response map[string]string
        err = json.NewDecoder(resp.Body).Decode(&response)
        assert.NoError(t, err, "Failed to decode response body")
        assert.Equal(t, "pong", response["message"], "Unexpected response message for /api/v1/ping")
    })

    // Additional endpoint tests...
}
```

## API Flow Tests

File: `api_flow_test.go`

### Testing Strategy

API flow tests simulate real-world usage scenarios by testing multiple endpoints in sequence. These tests ensure that:

1. The API functions correctly in typical user flows
2. Different parts of the API integrate properly
3. State changes are correctly reflected across multiple requests

### Key Components

- Setup of a test router with mock configuration
- Sequential tests that mimic user interactions
- Verification of expected outcomes at each step

### Code Example

```go
func TestAPIFlow(t *testing.T) {
    router := setupTestRouter()
    server := httptest.NewServer(router.Handler())
    defer server.Close()

    // Step 1: Ping
    t.Run("Ping", func(t *testing.T) {
        resp, err := http.Get(server.URL + "/api/v1/ping")
        assert.NoError(t, err)
        defer resp.Body.Close()

        assert.Equal(t, http.StatusOK, resp.StatusCode)

        var response map[string]string
        err = json.NewDecoder(resp.Body).Decode(&response)
        assert.NoError(t, err)
        assert.Equal(t, "pong", response["message"])
    })

    // Step 2: Health Check
    t.Run("Health Check", func(t *testing.T) {
        resp, err := http.Get(server.URL + "/api/v1/health")
        assert.NoError(t, err)
        defer resp.Body.Close()

        assert.Equal(t, http.StatusOK, resp.StatusCode)

        var response map[string]string
        err = json.NewDecoder(resp.Body).Decode(&response)
        assert.NoError(t, err)
        assert.Equal(t, "OK", response["status"])
    })

    // Additional flow steps...
}
```

## API Integration Tests

File: `api_test.go`

### Testing Strategy

API integration tests verify that all components of the API work together correctly. These tests ensure that:

1. All endpoints are accessible and return expected results
2. The API handles various input scenarios correctly
3. Error handling and edge cases are properly managed

### Key Components

- Setup of the full API router with actual configuration
- Comprehensive test cases covering all endpoints
- Verification of status codes and response bodies

### Code Example

```go
func TestAPIEndpoints(t *testing.T) {
    cfg, err := config.LoadConfig()
    assert.NoError(t, err, "Failed to load configuration")

    logger, err := zap.NewProduction()
    assert.NoError(t, err, "Failed to initialize logger")
    defer logger.Sync()

    r := gin.New()
    api.SetupRouter(r, cfg, logger)

    server := httptest.NewServer(r)
    defer server.Close()

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
        // Additional test cases...
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            resp, err := http.Get(server.URL + tc.endpoint)
            assert.NoError(t, err, "Failed to make request")
            defer resp.Body.Close()

            assert.Equal(t, tc.expectedStatus, resp.StatusCode, "Unexpected status code")

            body, err := ioutil.ReadAll(resp.Body)
            assert.NoError(t, err, "Failed to read response body")

            var responseBody map[string]string
            err = json.Unmarshal(body, &responseBody)
            assert.NoError(t, err, "Failed to parse response body")

            assert.Equal(t, tc.expectedBody, responseBody, "Unexpected response body")
        })
    }
}
```

## API Performance Benchmarks

File: `api_benchmark_test.go`

### Testing Strategy

API performance benchmarks measure the efficiency and scalability of our API. These tests ensure that:

1. The API can handle a high volume of requests
2. Response times remain within acceptable limits under load
3. Rate limiting functionality works as expected

### Key Components

- Benchmarks for individual endpoints
- Tests for rate limiting effectiveness
- Measurement of response times and request handling capacity

### Code Example

```go
func BenchmarkPingEndpoint(b *testing.B) {
    cfg, _ := config.LoadConfig()
    logger, _ := zap.NewProduction()
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
    // Setup code...

    client := &http.Client{}

    b.ResetTimer()
    for j := 0; j < 11; j++ {
        resp, err := client.Get(server.URL + "/api/v1/ping")
        if err != nil {
            b.Fatalf("Failed to make request: %v", err)
        }

        if j == 11 && resp.StatusCode != http.StatusTooManyRequests {
            b.Errorf("Expected rate limiting to be triggered (status 429), got %d", resp.StatusCode)
        }

        resp.Body.Close()
    }
}
```