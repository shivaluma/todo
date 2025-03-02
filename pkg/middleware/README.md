# Middleware Package

This package provides HTTP middleware components for the Todo API using the Gin framework.

## Overview

The middleware package offers:

1. Authentication middleware for JWT validation
2. Request ID middleware for request tracking
3. Logging middleware for HTTP request logging
4. CORS middleware for cross-origin resource sharing
5. Rate limiting middleware for API protection
6. Recovery middleware for panic handling
7. Metrics middleware for Prometheus metrics collection

## Usage

### Basic Usage

```go
import (
    "github.com/gin-gonic/gin"
    "github.com/sh1ro/todo-api/pkg/middleware"
)

func main() {
    router := gin.Default()

    // Apply global middleware
    router.Use(middleware.RequestID())
    router.Use(middleware.Logger())
    router.Use(middleware.Recovery())
    router.Use(middleware.CORS())
    router.Use(middleware.Metrics())

    // Apply rate limiting to specific routes
    api := router.Group("/api")
    api.Use(middleware.RateLimit(100, 1*time.Minute))

    // Apply authentication to protected routes
    protected := api.Group("/v1")
    protected.Use(middleware.Auth())

    // ... register your routes
}
```

## Available Middleware

### Authentication Middleware

The authentication middleware validates JWT tokens and sets the authenticated user in the context:

```go
// Apply to routes that require authentication
protected := router.Group("/api/v1")
protected.Use(middleware.Auth())

// Access the authenticated user in handlers
func Handler(c *gin.Context) {
    user := middleware.GetAuthenticatedUser(c)
    // ... use the authenticated user
}
```

### Request ID Middleware

The request ID middleware generates a unique ID for each request and adds it to the response headers:

```go
// Apply globally
router.Use(middleware.RequestID())

// Access the request ID in handlers
func Handler(c *gin.Context) {
    requestID := middleware.GetRequestID(c)
    // ... use the request ID
}
```

### Logging Middleware

The logging middleware logs HTTP requests with details such as method, path, status code, and duration:

```go
// Apply globally
router.Use(middleware.Logger())
```

### CORS Middleware

The CORS middleware handles Cross-Origin Resource Sharing headers:

```go
// Apply globally with default settings
router.Use(middleware.CORS())

// Apply with custom settings
router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
    AllowOrigins:     []string{"https://example.com"},
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
    AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
    ExposeHeaders:    []string{"Content-Length"},
    AllowCredentials: true,
    MaxAge:           12 * time.Hour,
}))
```

### Rate Limiting Middleware

The rate limiting middleware protects the API from abuse by limiting the number of requests per client:

```go
// Apply globally with default settings (100 requests per minute)
router.Use(middleware.RateLimit(100, 1*time.Minute))

// Apply with custom settings
router.Use(middleware.RateLimitWithConfig(middleware.RateLimitConfig{
    Limit:      100,
    Window:     1 * time.Minute,
    KeyFunc:    func(c *gin.Context) string { return c.ClientIP() },
    ExcludeFunc: func(c *gin.Context) bool { return c.Request.URL.Path == "/health" },
}))
```

### Recovery Middleware

The recovery middleware recovers from panics and returns a 500 Internal Server Error response:

```go
// Apply globally
router.Use(middleware.Recovery())
```

### Metrics Middleware

The metrics middleware collects Prometheus metrics for HTTP requests:

```go
// Apply globally
router.Use(middleware.Metrics())
```

## Configuration

Each middleware can be configured with custom settings using the corresponding `WithConfig` function:

```go
// Example: Custom authentication configuration
router.Use(middleware.AuthWithConfig(middleware.AuthConfig{
    TokenLookup:  "header:Authorization",
    AuthScheme:   "Bearer",
    IgnoreRoutes: []string{"/health", "/metrics"},
}))
```

## Middleware Order

The recommended order for applying middleware is:

1. `RequestID` - Generates a unique ID for each request
2. `Logger` - Logs HTTP requests
3. `Recovery` - Recovers from panics
4. `CORS` - Handles Cross-Origin Resource Sharing
5. `Metrics` - Collects Prometheus metrics
6. `RateLimit` - Limits the number of requests per client
7. `Auth` - Validates JWT tokens (only for protected routes)
