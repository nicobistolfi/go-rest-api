---
title: "Logging"
sidebar_position: 3
---

# Logging

This document explains how to use the custom logger implemented in the Go REST API Boilerplate project. The logger is a singleton wrapper around the `zap` logging library, providing a simple and consistent way to log messages across your application.

## Overview

The logger package provides a global `Log` variable of type `*Logger`, which is a wrapper around `zap.Logger`. It's designed as a singleton to ensure that only one instance of the logger is created and used throughout the application.

## Initialization

Before using the logger, you need to initialize it. This is typically done at the start of your application:

```go
package main

import (
    logger "go-rest-api/pkg"
)

func main() {
    logger.Init()
    // Rest of your application code
}
```

The `Init()` function uses a `sync.Once` to ensure that the logger is only initialized once, even if `Init()` is called multiple times.

## Basic Usage

After initialization, you can use the logger throughout your application. The package provides several logging methods:

### Logging at Different Levels

```go
logger.Info("This is an info message")
logger.Error("This is an error message")
logger.Fatal("This is a fatal message") // This will also call os.Exit(1)
```

### Logging with Additional Fields

You can add structured fields to your log messages:

```go
logger.Info("User logged in", zap.String("username", "john_doe"), zap.Int("user_id", 12345))
```


### Creating a Logger with Preset Fields

If you want to create a logger with some fields already set (useful for adding context to all logs from a particular component):

```go
componentLogger := logger.With(zap.String("component", "auth_service"))
componentLogger.Info("Starting authentication")
```

## Advanced Usage

### Direct Access to Underlying Zap Logger

If you need access to the underlying `zap.Logger` for advanced use cases:

```go
zapLogger := logger.Log.Logger
// Use zapLogger for zap-specific functionality
```

### Custom Log Fields

You can create custom log fields using `zap.Field` constructors:

```go
customField := zap.Any("custom_data", someComplexStruct)
logger.Info("Message with custom data", customField)
```

## Best Practices

1. **Initialization**: Always call `logger.Init()` at the start of your application.

2. **Log Levels**: Use appropriate log levels:
   - `Info` for general information
   - `Error` for error conditions
   - `Fatal` for unrecoverable errors (use sparingly as it terminates the program)

3. **Structured Logging**: Prefer using structured fields over string interpolation for better searchability and analysis of logs.

4. **Context**: Use `logger.With()` to add context to logs from specific components or request handlers.

5. **Performance**: The logger is designed to be performant, but avoid excessive logging in hot paths.

## Customization

The logger is initialized with a production configuration. If you need to customize the logger (e.g., for development environments), you can modify the `Init()` function in `logger.go`.

## Examples

### In HTTP Handlers

```go
func UserHandler(w http.ResponseWriter, r *http.Request) {
    userID := getUserIDFromRequest(r)
    requestLogger := logger.With(zap.String("handler", "UserHandler"), zap.Int("user_id", userID))
    
    requestLogger.Info("Processing user request")
    
    // Handler logic...
    
    if err != nil {
        requestLogger.Error("Failed to process user request", zap.Error(err))
        // Handle error...
    }
    
    requestLogger.Info("User request processed successfully")
}
```

### In Services

```go
type AuthService struct {
    logger *logger.Logger
}

func NewAuthService() *AuthService {
    return &AuthService{
        logger: logger.With(zap.String("service", "auth")),
    }
}

func (s *AuthService) Authenticate(username, password string) error {
    s.logger.Info("Attempting authentication", zap.String("username", username))
    
    // Authentication logic...
    
    if authFailed {
        s.logger.Error("Authentication failed", zap.String("username", username))
        return ErrAuthFailed
    }
    
    s.logger.Info("Authentication successful", zap.String("username", username))
    return nil
}
```

By following these guidelines and examples, you can effectively use the logger throughout your application to produce consistent, structured logs that will aid in monitoring, debugging, and maintaining your Go REST API.