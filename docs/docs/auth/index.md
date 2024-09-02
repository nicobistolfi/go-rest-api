---
title: Authentication
---
## Importance of Authentication

Authentication is a crucial aspect of any API, serving several important purposes:

1. **Security**: It ensures that only authorized users can access protected resources and perform sensitive operations.
2. **User Identity**: Authentication helps identify who is making requests, enabling personalized experiences and audit trails.
3. **Access Control**: It allows for fine-grained control over what actions different users can perform.
4. **Data Protection**: By restricting access to authenticated users, sensitive data is protected from unauthorized access.
5. **Compliance**: Many regulatory standards require proper authentication mechanisms to be in place.

In this boilerplate, we've implemented a flexible authentication system that supports multiple authentication methods, including API keys, JWT tokens, and OAuth 2.0 with OpenID Connect (OIDC).

## Modifying routes.go to Add Authenticated Endpoints

To add new endpoints that require authentication, you'll need to modify the `routes.go` file. Here's how you can do it:

1. Open the `routes.go` file:


```1:38:internal/api/routes.go
package api

import (
	"time"

	"go-boilerplate/internal/api/middleware"
	"go-boilerplate/internal/config"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

func SetupRouter(r *gin.Engine, cfg *config.Config, logger *zap.Logger) {
	// Add global middleware
	r.Use(gin.Recovery())
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.LoggerMiddleware(logger))
	r.Use(middleware.RateLimiter(rate.Every(time.Second), 10)) // 10 requests per second

	// Public routes
	public := r.Group("/api/v1")
	{
		public.GET("/health", HealthCheck)
		public.GET("/ping", Ping)
		// Add other public routes
	}

	// Protected routes
	protected := r.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware())
	protected.Use(middleware.VerifyToken())
	{
		protected.GET("/profile", GetProfile)
		// Add other protected routes
	}
}

```


2. To add a new authenticated endpoint, add it to the `protected` group. For example, to add a new `/api/v1/user-data` endpoint:

```go
protected.GET("/user-data", GetUserData)
```

3. Implement the corresponding handler function (e.g., `GetUserData`) in the `handlers.go` file.

4. If you need different authentication methods for different endpoints, you can create new route groups with specific middleware. For example:

```go
jwtProtected := r.Group("/api/v1/jwt")
jwtProtected.Use(middleware.JWTAuthMiddleware())
{
    jwtProtected.GET("/jwt-specific-data", GetJWTSpecificData)
}
```

## Setting the TOKEN_URL Environment Variable

The `TOKEN_URL` environment variable is crucial for token validation in the OAuth 2.0 flow. It should be set to the URL of your OAuth provider's token introspection endpoint. Here's how to set it:

1. In your `.env` file or deployment environment, add:

```
TOKEN_URL=https://your-oauth-provider.com/oauth/token/introspect
```

2. Ensure that your application loads this environment variable. In the `config/config.go` file, add:

```go
TokenURL: os.Getenv("TOKEN_URL"),
```

3. Update your `Config` struct to include this new field.

## How the Auth Middlewares Work

The authentication system uses two main middleware components: `auth.go` and `token.go`.

### auth.go Middleware


```1:43:internal/api/middleware/auth.go
package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var token string
		var authHeader string

		// Check for Authorization header
		token = c.GetHeader("Authorization")
		authHeader = "Authorization"

		// If not found, check for X-API-Key header
		if token == "" {
			token = c.GetHeader("X-API-Key")
			authHeader = "X-API-Key"
		}

		// If still not found, check for access_token query parameter
		if token == "" {
			token = c.Query("access_token")
			authHeader = "access_token"
		}

		// If no token found in any of the above methods
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication token is required"})
			c.Abort()
			return
		}

		// Add the token to the context
		c.Set("auth_token", token)
		c.Set("auth_header", authHeader)

		c.Next()
	}
}
```


This middleware:
1. Checks for an authentication token in the request headers or query parameters.
2. Supports multiple authentication methods (Authorization header, X-API-Key header, access_token query parameter).
3. If a token is found, it adds it to the Gin context for further processing.
4. If no token is found, it returns a 401 Unauthorized response.

### token.go Middleware


```1:120:internal/api/middleware/token.go
package middleware

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var logger *zap.Logger

func init() {
	logger, _ = zap.NewProduction()
}

type Profile struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

func VerifyToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the token from the context set by AuthMiddleware
		token, exists := c.Get("auth_token")
		authHeader, authHeaderExists := c.Get("auth_header")

		if !authHeaderExists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication header is missing"})
			c.Abort()
			return
		}

		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication token is missing"})
			c.Abort()
			return
		}

		tokenString, ok := token.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}

		tokenURL := os.Getenv("TOKEN_URL")

		if tokenURL == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "TOKEN_URL not set"})
			c.Abort()
			return
		}

		// Validate token using the TOKEN_URL
		req, err := http.NewRequest("GET", tokenURL, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
			c.Abort()
			return
		}
		req.Header.Set(authHeader.(string), tokenString)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate token"})
			c.Abort()
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			// use logger to log the error
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
			c.Abort()
			return
		}

		var profile Profile
		if err := json.Unmarshal(body, &profile); err != nil {
			// Use a custom unmarshaling to map the fields correctly
			var githubProfile struct {
				ID    uint   `json:"id"`
				Email string `json:"email"`
				Name  string `json:"name"`
				Login string `json:"login"`
			}
			if err := json.Unmarshal(body, &githubProfile); err != nil {
				logger.Error("Failed to parse GitHub profile", zap.Error(err))
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse profile"})
				c.Abort()
				return
			}
			profile = Profile{
				ID:    fmt.Sprintf("%d", githubProfile.ID),
				Email: githubProfile.Email,
				Name:  githubProfile.Name,
			}
			// If Email is empty, use Login as a fallback
			if profile.Email == "" {
				profile.Email = githubProfile.Login
			}
		}

		c.Set("user", profile)
		c.Next()
	}
}
```


This middleware:
1. Retrieves the token from the Gin context (set by the auth middleware).
2. Validates the token by making a request to the `TOKEN_URL` endpoint.
3. If the token is valid, it extracts user information from the response.
4. Sets the user information in the Gin context for use in subsequent request handling.
5. If the token is invalid, it returns a 401 Unauthorized response.

Together, these middlewares provide a robust authentication system that can be easily extended to support various authentication methods and token validation strategies.

To use these middlewares, ensure they are properly added to your routes in the `routes.go` file, as shown in the earlier section on modifying routes.

By understanding and properly configuring these authentication components, you can ensure that your API endpoints are secure and accessible only to authorized users.