package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sh1ro/todo-api/internal/app/domain/service"
	"github.com/sh1ro/todo-api/pkg/logger"
)

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
func (m *AuthMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get request-specific logger if available
		var log *logger.Logger
		if l, exists := c.Get("logger"); exists {
			log = l.(*logger.Logger)
		} else {
			log = m.logger
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		}

		// Check if the Authorization header has the correct format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
			return
		}

		// Get the token
		tokenString := parts[1]

		// Validate the token
		claims, err := m.authService.ValidateToken(tokenString)
		if err != nil {
			log.Error("Invalid token", "error", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		// Get the user from the token
		user, err := m.authService.GetUserFromToken(c.Request.Context(), tokenString)
		if err != nil {
			if errors.Is(err, errors.New("user not found")) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
				return
			}
			log.Error("Failed to get user from token", "error", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		// Set the user ID and username in the context
		c.Set("userID", user.ID)
		c.Set("username", user.Username)
		c.Set("claims", claims)

		c.Next()
	}
}
