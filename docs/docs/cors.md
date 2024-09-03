---
title: "CORS"
sidebar_position: 5
---

# CORS Middleware Configuration and Usage

This document explains how to configure and use the CORS (Cross-Origin Resource Sharing) middleware in the Go REST API Boilerplate project. The middleware is implemented using the `gin-contrib/cors` package and is designed to be flexible and configurable.

## Overview

The CORS middleware allows you to control which origins can access your API. It's crucial for securing your API while still allowing legitimate cross-origin requests.

## Configuration

The middleware is configured in the `CORSMiddleware()` function, which returns a `gin.HandlerFunc`. Here's how it works:

1. **Default Configuration**: It starts with the default CORS configuration.

2. **Allowed Origins**: 
   - The middleware checks for an environment variable `ALLOWED_ORIGINS`.
   - If set, it uses these origins as the allowed origins.
   - If not set, it allows any origin that starts with `http://localhost` or `https://localhost`.

3. **Allowed Methods**: The middleware allows the following HTTP methods:
   - GET
   - POST
   - PUT
   - PATCH
   - DELETE
   - OPTIONS

4. **Allowed Headers**: The middleware allows the following headers:
   - Origin
   - Content-Type
   - Accept
   - Authorization

## Usage

To use this middleware in your Gin application, follow these steps:

1. Import the middleware package:

```go
import "go-rest-api/internal/api/middleware"
```

2. Add the middleware to your Gin router:

```go
func main() {
    router := gin.Default()
    
    // Apply the CORS middleware
    router.Use(middleware.CORSMiddleware())

    // Your routes go here
    // ...

    router.Run()
}
```

## Configuration via Environment Variables

To configure allowed origins using environment variables:

1. Set the `ALLOWED_ORIGINS` environment variable before running your application. Multiple origins should be comma-separated.

   Example:
   ```bash
   export ALLOWED_ORIGINS="https://example.com,https://api.example.com"
   ```

2. If `ALLOWED_ORIGINS` is not set, the middleware will default to allowing any origin that starts with `http://localhost` or `https://localhost`.

## Customization

If you need to customize the CORS settings further:

1. Modify the `CORSMiddleware()` function in the middleware package.
2. You can adjust allowed methods, headers, or add more sophisticated origin checking logic.

Example of adding a custom header:

```go
config.AllowHeaders = append(config.AllowHeaders, "X-Custom-Header")
```

## Security Considerations

1. **Restrict Origins**: In production, always set `ALLOWED_ORIGINS` to a specific list of trusted domains.
2. **Least Privilege**: Only expose the methods and headers that your API actually needs.
3. **Credentials**: If your API requires credentials (cookies, HTTP authentication), you may need to set `config.AllowCredentials = true`. Use this with caution and ensure `AllowOrigins` is not set to `*`.

## Troubleshooting

If you're experiencing CORS issues:

1. Check that the `ALLOWED_ORIGINS` environment variable is set correctly.
2. Ensure that the origin making the request matches exactly with one of the allowed origins (including the protocol, `http://` or `https://`).
3. Verify that the request is using an allowed method and only includes allowed headers.

## Example

Here's a complete example of setting up a Gin router with the CORS middleware:

```go
package main

import (
    "github.com/gin-gonic/gin"
    "go-rest-api/internal/api/middleware"
)

func main() {
    // Set up Gin
    router := gin.Default()

    // Apply CORS middleware
    router.Use(middleware.CORSMiddleware())

    // Define a route
    router.GET("/api/data", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "This is CORS-enabled data",
        })
    })

    // Run the server
    router.Run(":8080")
}
```

In this example, the CORS middleware will be applied to all routes. The allowed origins will be determined by the `ALLOWED_ORIGINS` environment variable, or default to localhost if not set.

By following these guidelines, you can effectively implement and customize CORS in your Go REST API, ensuring that your API is accessible to the intended clients while maintaining security.