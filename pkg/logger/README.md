# Logger Package

This package provides structured logging capabilities for the Todo API using the `uber-go/zap` package.

## Overview

The logger package offers:

1. Pre-configured logger instance with common settings
2. Support for different log levels (debug, info, warn, error, fatal)
3. Structured logging with key-value pairs
4. HTTP request logging middleware for Gin
5. Context-aware logging with request IDs

## Usage

### Basic Usage

```go
import (
    "github.com/sh1ro/todo-api/pkg/logger"
)

func main() {
    // Initialize the logger
    logger.Init("info", false)

    // Log messages at different levels
    logger.Debug("Debug message")
    logger.Info("Info message")
    logger.Warn("Warning message")
    logger.Error("Error message")

    // Log with additional fields
    logger.Info("User created",
        logger.String("user_id", "123"),
        logger.String("username", "johndoe"),
        logger.Int("age", 30),
    )

    // Log errors with stack traces
    err := someFunction()
    if err != nil {
        logger.Error("Failed to execute function",
            logger.Error(err),
            logger.String("function", "someFunction"),
        )
    }
}
```

### HTTP Request Logging

```go
import (
    "github.com/gin-gonic/gin"
    "github.com/sh1ro/todo-api/pkg/logger"
)

func main() {
    router := gin.Default()

    // Register the logger middleware
    router.Use(logger.GinMiddleware())

    // ... rest of your application setup
}
```

### Context-Aware Logging

```go
import (
    "github.com/gin-gonic/gin"
    "github.com/sh1ro/todo-api/pkg/logger"
)

func Handler(c *gin.Context) {
    // Get logger with request context
    log := logger.FromContext(c)

    // Log with request ID and other context information
    log.Info("Processing request")

    // ... process the request

    log.Info("Request processed successfully",
        logger.String("result", "success"),
    )
}
```

## Configuration

The logger can be configured with:

1. Log level (debug, info, warn, error, fatal)
2. Development mode (enables more verbose logging)
3. Output destination (stdout, file, or both)
4. Log format (JSON or console)

Example configuration:

```go
// Initialize with info level in production mode
logger.Init("info", false)

// Initialize with debug level in development mode
logger.Init("debug", true)
```

## Available Log Levels

-   `Debug` - Detailed information, typically of interest only when diagnosing problems
-   `Info` - Confirmation that things are working as expected
-   `Warn` - Indication that something unexpected happened, or may happen in the future
-   `Error` - Due to a more serious problem, the software has not been able to perform a function
-   `Fatal` - Very severe error events that will presumably lead the application to abort

## Field Helpers

The package provides helper functions for adding structured fields to log messages:

-   `String(key, value)` - Add a string field
-   `Int(key, value)` - Add an integer field
-   `Bool(key, value)` - Add a boolean field
-   `Float64(key, value)` - Add a float field
-   `Error(err)` - Add an error field with stack trace
-   `Any(key, value)` - Add a field of any type

## Middleware Features

The Gin middleware automatically logs:

1. Request method and path
2. Response status code
3. Request duration
4. Client IP address
5. User agent
6. Request ID (generated or from X-Request-ID header)
7. Request body size
8. Response body size
