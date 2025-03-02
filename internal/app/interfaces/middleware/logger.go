package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sh1ro/todo-api/pkg/logger"
)

// Logger returns a middleware that logs request information
func Logger(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Stop timer
		end := time.Now()
		latency := end.Sub(start)

		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()

		if raw != "" {
			path = path + "?" + raw
		}

		// Get request ID from context if available
		var requestID string
		if id, exists := c.Get(RequestIDContextKey); exists {
			requestID = id.(string)
		}

		// Use the logger from context if available (with request ID)
		var contextLogger *logger.Logger
		if l, exists := c.Get("logger"); exists {
			contextLogger = l.(*logger.Logger)
		} else if requestID != "" {
			// If we have a request ID but no logger in context, create one with the request ID
			contextLogger = log.WithField("request_id", requestID)
		} else {
			contextLogger = log
		}

		// Log request details
		contextLogger.Info("Request",
			"status", statusCode,
			"method", method,
			"path", path,
			"ip", clientIP,
			"latency", latency,
			"error", errorMessage,
		)
	}
}
