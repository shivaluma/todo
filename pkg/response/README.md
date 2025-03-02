# Response Package

This package provides standardized response structures and helper functions for the Todo API.

## Overview

The response package offers:

1. Standard response templates for success and error responses
2. Generic response templates for type-safe responses
3. Helper functions for common response patterns
4. Consistent error handling

## Response Structure

All API responses follow a consistent structure:

```json
{
  "status": "success|error",
  "message": "Human-readable message",
  "data": { ... },  // For success responses
  "error": { ... }, // For error responses
  "timestamp": "2023-01-01T12:00:00Z",
  "request_id": "unique-request-id"
}
```

## Usage

### Standard Responses

```go
import (
    "github.com/gin-gonic/gin"
    "github.com/sh1ro/todo-api/pkg/response"
)

func Handler(c *gin.Context) {
    // Success response
    response.RespondWithSuccess(c, 200, "Item retrieved successfully", map[string]interface{}{
        "item": item,
    })

    // Error response
    response.RespondWithError(c, 400, "Invalid request", map[string]interface{}{
        "details": "Field 'name' is required",
    })
}
```

### Generic Responses

```go
import (
    "github.com/gin-gonic/gin"
    "github.com/sh1ro/todo-api/pkg/response"
)

type User struct {
    ID       string `json:"id"`
    Username string `json:"username"`
    Email    string `json:"email"`
}

func GetUser(c *gin.Context) {
    user := User{
        ID:       "123",
        Username: "johndoe",
        Email:    "john@example.com",
    }

    // Type-safe success response
    response.RespondWithGenericOK(c, "User retrieved successfully", user)
}

type PaginatedUsers struct {
    Users      []User `json:"users"`
    TotalCount int    `json:"total_count"`
    Page       int    `json:"page"`
    PageSize   int    `json:"page_size"`
}

func ListUsers(c *gin.Context) {
    users := []User{
        {ID: "123", Username: "johndoe", Email: "john@example.com"},
        {ID: "456", Username: "janedoe", Email: "jane@example.com"},
    }

    // Type-safe paginated response
    response.RespondWithGenericPaginated(c, "Users retrieved successfully", PaginatedUsers{
        Users:      users,
        TotalCount: 2,
        Page:       1,
        PageSize:   10,
    })
}
```

## Available Helper Functions

### Success Responses

-   `RespondWithSuccess(c *gin.Context, code int, message string, data interface{})`
-   `RespondWithGenericOK[T any](c *gin.Context, message string, data T)`
-   `RespondWithGenericCreated[T any](c *gin.Context, message string, data T)`
-   `RespondWithGenericPaginated[T any](c *gin.Context, message string, data T)`

### Error Responses

-   `RespondWithError(c *gin.Context, code int, message string, details interface{})`
-   `RespondWithBadRequest(c *gin.Context, message string, details interface{})`
-   `RespondWithUnauthorized(c *gin.Context, message string, details interface{})`
-   `RespondWithForbidden(c *gin.Context, message string, details interface{})`
-   `RespondWithNotFound(c *gin.Context, message string, details interface{})`
-   `RespondWithInternalError(c *gin.Context, message string, details interface{})`
-   `RespondWithValidationError(c *gin.Context, message string, details interface{})`

### Generic Error Responses

-   `RespondWithGenericError[T any](c *gin.Context, code int, message string, details T)`
-   `RespondWithGenericBadRequest[T any](c *gin.Context, message string, details T)`
-   `RespondWithGenericValidationError[T any](c *gin.Context, message string, details T)`
