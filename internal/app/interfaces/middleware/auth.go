package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/sh1ro/todo-api/internal/app/domain/service"
	"github.com/sh1ro/todo-api/pkg/logger"
)

// contextKey is an unexported type for context keys to prevent collisions
type contextKey string

// Exported context keys
const (
	// UserIDKey is the context key for user ID
	UserIDKey contextKey = "user_id"
	
	// ClaimsKey is the context key for JWT claims
	ClaimsKey contextKey = "claims"
)

// GetUserID retrieves the user ID from the context
func GetUserID(c echo.Context) (interface{}, bool) {
	userID := c.Get(string(UserIDKey))
	return userID, userID != nil
}

// GetClaims retrieves the JWT claims from the context
func GetClaims(c echo.Context) (interface{}, bool) {
	claims := c.Get(string(ClaimsKey))
	return claims, claims != nil
}

// AuthMiddleware is a middleware that checks for a valid JWT token
type AuthMiddleware struct {
	authService *service.AuthService
	logger      *logger.Logger
}

// NewAuthMiddleware creates a new AuthMiddleware
func NewAuthMiddleware(authService *service.AuthService, logger *logger.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
		logger:      logger,
	}
}

// Authenticate is a middleware that checks for a valid JWT token
func (m *AuthMiddleware) Authenticate() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get request-specific logger with request ID using FromContext
			log := logger.FromContext(c)

			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "Authorization header is required")
			}

			// Check if the Authorization header has the correct format
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return echo.NewHTTPError(http.StatusUnauthorized, "Authorization header format must be Bearer {token}")
			}

			// Get the token
			tokenString := parts[1]

			// Validate the token
			claims, err := m.authService.ValidateToken(tokenString)
			if err != nil {
				log.Error("Invalid token", "error", err)
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid or expired token")
			}

			// Get the user from the token
			user, err := m.authService.GetUserFromToken(c.Request().Context(), tokenString)
			if err != nil {
				if errors.Is(err, errors.New("user not found")) {
					return echo.NewHTTPError(http.StatusUnauthorized, "User not found")
				}
				log.Error("Failed to get user from token", "error", err)
				return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
			}

			// Set the user ID and claims in the context
			c.Set(string(UserIDKey), user.ID)
			c.Set(string(ClaimsKey), claims)

			return next(c)
		}
	}
}
