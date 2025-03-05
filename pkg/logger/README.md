# Logger Package

This package provides structured logging capabilities for the Todo API using the `uber-go/zap` package.

## Overview

The logger package offers:

1. Pre-configured logger instance with common settings
2. Support for different log levels (debug, info, warn, error, fatal)
3. Structured logging with key-value pairs
4. HTTP request logging middleware for Echo
5. Context-aware logging with request IDs

## Usage

### Basic Usage

```go
import (
    "github.com/sh1ro/todo-api/pkg/logger"
)

func main() {
    // Initialize the logger
    log := logger.NewLogger("info", "console")

    // Log messages at different levels
    log.Debug("Debug message")
    log.Info("Info message")
    log.Warn("Warning message")
    log.Error("Error message")

    // Log with additional fields
    log.Info("User created",
        "user_id", "123",
        "username", "johndoe",
        "age", 30,
    )

    // Log with a single field
    log.WithField("request_id", "abc123").Info("Processing request")

    // Log with multiple fields
    log.WithFields(map[string]interface{}{
        "user_id": "123",
        "action": "login",
    }).Info("User logged in")

    // Log errors
    err := someFunction()
    if err != nil {
        log.Error("Failed to execute function",
            "error", err,
            "function", "someFunction",
        )
    }
}
```

### HTTP Request Logging

```go
import (
    "github.com/labstack/echo/v4"
    "github.com/sh1ro/todo-api/pkg/logger"
    customMiddleware "github.com/sh1ro/todo-api/internal/app/interfaces/middleware"
)

func main() {
    // Initialize the logger
    log := logger.NewLogger("info", "console")

    // Create Echo instance
    e := echo.New()

    // Register the logger middleware
    e.Use(customMiddleware.Logger(log))

    // ... rest of your application setup
}
```

### Context-Aware Logging with Request IDs

The logger provides two ways to add request IDs to your logs:

#### Method 1: Using FromContext

```go
import (
    "github.com/labstack/echo/v4"
    "github.com/sh1ro/todo-api/pkg/logger"
)

func Handler(c echo.Context) error {
    // Get logger with request ID automatically extracted from context
    log := logger.FromContext(c)

    // Log with request ID automatically included
    log.Info("Processing request")

    // ... process the request

    log.Info("Request processed successfully",
        "result", "success",
    )

    return nil
}
```

#### Method 2: Using WithRequestID

```go
import (
    "github.com/labstack/echo/v4"
    "github.com/sh1ro/todo-api/pkg/logger"
)

func Handler(c echo.Context) error {
    // Get request ID from header
    requestID := c.Response().Header().Get(echo.HeaderXRequestID)

    // Create logger with request ID
    log := logger.NewLogger("info", "console").WithRequestID(requestID)

    // Log with request ID
    log.Info("Processing request")

    // ... process the request

    log.Info("Request processed successfully",
        "result", "success",
    )

    return nil
}
```

## Configuration

The logger can be configured with:

1. Log level (debug, info, warn, error)
2. Log format (json or console)
3. Output destination (can be changed with WithOutput)

Example configuration:

```go
// Initialize with info level and console format
log := logger.NewLogger("info", "console")

// Initialize with debug level and JSON format
log := logger.NewLogger("debug", "json")

// Change output to a file
file, _ := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
log = log.WithOutput(file)
```

## Available Log Levels

-   `Debug` - Detailed information, typically of interest only when diagnosing problems
-   `Info` - Confirmation that things are working as expected
-   `Warn` - Indication that something unexpected happened, or may happen in the future
-   `Error` - Due to a more serious problem, the software has not been able to perform a function
-   `Fatal` - Very severe error events that will presumably lead the application to abort

## Implementation Details

Under the hood, this logger uses Uber's zap logging library, which provides:

1. Extremely fast, structured logging
2. Type-safe field additions
3. Leveled logging
4. Sampling capabilities
5. Hooks for extending functionality

The implementation uses both the structured Logger and the more flexible SugaredLogger from zap to provide a balance of performance and usability.

### Request ID Handling

Request IDs are handled consistently throughout the application:

1. The `WithRequestID` method adds a request ID with the standard key "request_id"
2. The `FromContext` method automatically extracts the request ID from the Echo context
3. Request IDs are preserved when using methods like `WithField` or `WithFields`
