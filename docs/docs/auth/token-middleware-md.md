---
title: Token Middleware
---

# Token Middleware

The Token Middleware, defined in `token.go`, is responsible for validating authentication tokens and retrieving user profiles. It works in conjunction with the Auth Middleware to provide a complete authentication solution. This middleware ensures that only valid tokens are accepted and provides user profile information for authenticated requests.

:::warning

This middleware is built as an example using GitHub API. It requires to be modified to work with other authentication services.

:::

## Key Features

- Validates tokens using an external authentication service
- Retrieves and parses user profiles
- Implements token caching for improved performance
- Supports different profile structures (e.g., GitHub profiles)

## Token Validation Process

1. Retrieve the token from the request context (set by Auth Middleware)
2. Check the token cache for a valid, non-expired entry
3. If not in cache, validate the token using the external `TOKEN_URL`
4. Parse the user profile from the validation response
5. Cache the validated token and profile
6. Set the user profile in the request context

## Usage in Routes

```go
protected := r.Group("/api/v1")
protected.Use(middleware.AuthMiddleware(), middleware.VerifyToken())
```
