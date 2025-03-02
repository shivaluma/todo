# Validator Package

This package provides enhanced validation capabilities for the Todo API using the `go-playground/validator/v10` package.

## Overview

The validator package offers:

1. Pre-configured validator instance with common settings
2. Custom validation rules specific to the application
3. Internationalization (i18n) support for validation error messages
4. Helper functions for translating validation errors to user-friendly messages

## Usage

### Basic Usage

```go
import (
    "github.com/gin-gonic/gin"
    "github.com/sh1ro/todo-api/pkg/validator"
)

type CreateUserRequest struct {
    Username string `json:"username" validate:"required,min=3,max=50"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
}

func RegisterHandler(c *gin.Context) {
    var req CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        // Handle JSON binding error
        return
    }

    // Validate the request
    if err := validator.Validate(req); err != nil {
        // Get translated validation errors
        validationErrors := validator.TranslateError(err)

        // Respond with validation errors
        // ...
        return
    }

    // Process valid request
    // ...
}
```

### Custom Validation Rules

The package includes several custom validation rules:

1. `priority` - Validates that a priority value is one of the allowed values (low, medium, high)
2. `status` - Validates that a status value is one of the allowed values (pending, in_progress, completed)
3. `future_date` - Validates that a date is in the future

Example:

```go
type CreateTodoRequest struct {
    Title       string    `json:"title" validate:"required,min=3,max=100"`
    Description string    `json:"description" validate:"max=500"`
    DueDate     time.Time `json:"due_date" validate:"required,future_date"`
    Priority    string    `json:"priority" validate:"required,priority"`
    Status      string    `json:"status" validate:"required,status"`
}
```

### Validation Error Translation

The package provides functions to translate validation errors into user-friendly messages:

```go
// Translate a validation error to a map of field names to error messages
errorMap := validator.TranslateError(err)

// Example output:
// {
//   "username": "Username is required",
//   "email": "Email must be a valid email address",
//   "password": "Password must be at least 8 characters long"
// }
```

## Configuration

The validator is configured with:

1. JSON tag name support (uses JSON field names in error messages)
2. English translations for common validation errors
3. Custom validation rules specific to the application
4. Support for struct-level validation

## Available Validation Tags

In addition to the standard validation tags provided by `go-playground/validator/v10`, this package supports:

-   `priority` - Validates priority values (low, medium, high)
-   `status` - Validates status values (pending, in_progress, completed)
-   `future_date` - Validates that a date is in the future

## Error Handling

Validation errors are returned as a map of field names to error messages, making it easy to include them in API responses:

```go
if err := validator.Validate(req); err != nil {
    validationErrors := validator.TranslateError(err)

    response.RespondWithValidationError(c, "Validation failed", validationErrors)
    return
}
```
