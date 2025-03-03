package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sh1ro/todo-api/pkg/logger"
)

const (
	// RequestIDHeader is the header key for request ID
	RequestIDHeader = "X-Request-ID"
)

// RequestIDKey is the context key for request ID
const RequestIDKey contextKey = "requestID"

// GetRequestID retrieves the request ID from the context
func GetRequestID(c *gin.Context) (string, bool) {
	if id, exists := c.Get(string(RequestIDKey)); exists {
		if requestID, ok := id.(string); ok {
			return requestID, true
		}
	}
	return "", false
}

// RequestID returns a middleware that adds a unique request ID to each request
func RequestID(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if request already has an ID
		requestID := c.GetHeader(RequestIDHeader)

		// If no request ID is provided, generate a new one
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Set request ID in context for other handlers to use
		c.Set(string(RequestIDKey), requestID)

		// Set request ID in response header
		c.Writer.Header().Set(RequestIDHeader, requestID)

		// Add request ID to logger context
		requestLogger := log.WithField("request_id", requestID)
		c.Set("logger", requestLogger)

		c.Next()
	}
}
