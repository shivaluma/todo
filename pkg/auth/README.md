# Auth Package

This package provides authentication and authorization functionality for the Todo API.

## Overview

The auth package offers:

1. JWT token generation and validation
2. Password hashing and verification
3. User authentication
4. Role-based authorization
5. Token management (refresh, revocation)

## Usage

### JWT Token Management

```go
import (
    "time"
    "github.com/sh1ro/todo-api/pkg/auth"
)

func main() {
    // Initialize the JWT manager
    jwtManager := auth.NewJWTManager("your-secret-key", 24*time.Hour)

    // Generate a token
    userID := "123"
    username := "johndoe"
    role := "user"

    token, err := jwtManager.Generate(userID, username, role)
    if err != nil {
        panic(err)
    }

    // Validate a token
    claims, err := jwtManager.Verify(token)
    if err != nil {
        // Handle invalid token
        panic(err)
    }

    // Access claims
    userID = claims.UserID
    username = claims.Username
    role = claims.Role
}
```

### Password Management

```go
import (
    "github.com/sh1ro/todo-api/pkg/auth"
)

func RegisterUser(username, password string) error {
    // Hash the password
    hashedPassword, err := auth.HashPassword(password)
    if err != nil {
        return err
    }

    // Store the user with the hashed password
    // ...

    return nil
}

func LoginUser(username, password string) error {
    // Retrieve the user's hashed password from the database
    hashedPassword := getUserHashedPassword(username)

    // Verify the password
    if !auth.CheckPasswordHash(password, hashedPassword) {
        return errors.New("invalid credentials")
    }

    // Password is correct, proceed with login
    // ...

    return nil
}
```

### User Authentication

```go
import (
    "github.com/gin-gonic/gin"
    "github.com/sh1ro/todo-api/pkg/auth"
    "github.com/sh1ro/todo-api/pkg/response"
)

func LoginHandler(c *gin.Context) {
    var req LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.RespondWithBadRequest(c, "Invalid request", nil)
        return
    }

    // Authenticate the user
    user, err := auth.AuthenticateUser(req.Username, req.Password)
    if err != nil {
        response.RespondWithUnauthorized(c, "Invalid credentials", nil)
        return
    }

    // Generate a token
    token, err := auth.GenerateToken(user.ID, user.Username, user.Role)
    if err != nil {
        response.RespondWithInternalError(c, "Failed to generate token", nil)
        return
    }

    // Return the token
    response.RespondWithSuccess(c, 200, "Login successful", gin.H{
        "token": token,
        "user": gin.H{
            "id": user.ID,
            "username": user.Username,
            "role": user.Role,
        },
    })
}
```

### Role-Based Authorization

```go
import (
    "github.com/gin-gonic/gin"
    "github.com/sh1ro/todo-api/pkg/auth"
    "github.com/sh1ro/todo-api/pkg/response"
)

// Middleware to require admin role
func RequireAdmin() gin.HandlerFunc {
    return func(c *gin.Context) {
        user := auth.GetUserFromContext(c)
        if user == nil || user.Role != "admin" {
            response.RespondWithForbidden(c, "Admin access required", nil)
            c.Abort()
            return
        }
        c.Next()
    }
}

// Apply the middleware to routes that require admin access
func SetupRoutes(router *gin.Engine) {
    admin := router.Group("/api/v1/admin")
    admin.Use(auth.RequireAuth()) // First check if user is authenticated
    admin.Use(RequireAdmin())     // Then check if user is an admin

    admin.GET("/users", ListUsersHandler)
    admin.POST("/users", CreateUserHandler)
    // ... other admin routes
}
```

### Token Refresh

```go
import (
    "github.com/gin-gonic/gin"
    "github.com/sh1ro/todo-api/pkg/auth"
    "github.com/sh1ro/todo-api/pkg/response"
)

func RefreshTokenHandler(c *gin.Context) {
    var req RefreshTokenRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.RespondWithBadRequest(c, "Invalid request", nil)
        return
    }

    // Verify the refresh token
    claims, err := auth.VerifyRefreshToken(req.RefreshToken)
    if err != nil {
        response.RespondWithUnauthorized(c, "Invalid refresh token", nil)
        return
    }

    // Generate a new access token
    accessToken, err := auth.GenerateToken(claims.UserID, claims.Username, claims.Role)
    if err != nil {
        response.RespondWithInternalError(c, "Failed to generate token", nil)
        return
    }

    // Return the new access token
    response.RespondWithSuccess(c, 200, "Token refreshed", gin.H{
        "token": accessToken,
    })
}
```

## Configuration

The auth package can be configured with:

```go
type AuthConfig struct {
    JWTSecret            string
    AccessTokenDuration  time.Duration
    RefreshTokenDuration time.Duration
    PasswordHashCost     int
}
```

## Available Functions

### JWT Management

-   `NewJWTManager(secret string, duration time.Duration) *JWTManager`
-   `Generate(userID, username, role string) (string, error)`
-   `Verify(token string) (*Claims, error)`
-   `GenerateRefreshToken(userID, username, role string) (string, error)`
-   `VerifyRefreshToken(token string) (*Claims, error)`

### Password Management

-   `HashPassword(password string) (string, error)`
-   `CheckPasswordHash(password, hash string) bool`

### User Authentication

-   `AuthenticateUser(username, password string) (*User, error)`
-   `GetUserFromContext(c *gin.Context) *User`

### Middleware

-   `RequireAuth() gin.HandlerFunc`
-   `RequireRole(role string) gin.HandlerFunc`
