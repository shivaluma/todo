package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sh1ro/todo-api/pkg/logger"
)

const (
	// RequestIDHeader is the header key for request ID
	RequestIDHeader = "X-Request-ID"
	// RequestIDContextKey is the context key for request ID
	RequestIDContextKey = "requestID"
)

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
		c.Set(RequestIDContextKey, requestID)

		// Set request ID in response header
		c.Writer.Header().Set(RequestIDHeader, requestID)

		// Add request ID to logger context
		requestLogger := log.WithField("request_id", requestID)
		c.Set("logger", requestLogger)

		c.Next()
	}
}
