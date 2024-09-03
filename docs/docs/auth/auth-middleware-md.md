---
title: Auth Middleware
---

# Auth Middleware

The Auth Middleware, defined in `auth.go`, is responsible for extracting authentication tokens from incoming requests. It supports multiple authentication methods to provide flexibility in how clients can authenticate.

## Key Features

- Checks multiple sources for authentication tokens
- Supports different authentication methods
- Sets extracted token information in the request context

## Token Extraction Process

1. Check the `Authorization` header
2. If not found, check the `X-API-Key` header
3. If still not found, check the `access_token` query parameter

## Usage in Routes

```go
protected := r.Group("/api/v1")
protected.Use(middleware.AuthMiddleware())
```

This middleware sets the stage for token validation, which is handled by the Token Middleware.
