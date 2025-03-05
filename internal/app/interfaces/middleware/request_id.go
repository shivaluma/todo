package middleware

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sh1ro/todo-api/pkg/logger"
)

// RequestIDKey is the key used to store the request ID in the context
const RequestIDKey = "request_id"

// RequestID returns a middleware that adds a request ID to the context
func RequestID(log *logger.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Check if request ID is already set in header
			requestID := c.Request().Header.Get(echo.HeaderXRequestID)
			if requestID == "" {
				// Generate a new request ID
				requestID = uuid.New().String()
			}

			// Set request ID in response header
			c.Response().Header().Set(echo.HeaderXRequestID, requestID)

			// Create a request-specific logger with request ID
			reqLogger := log.WithField("request_id", requestID)

			// Store logger in context for handlers to use
			c.Set("logger", reqLogger)

			// Continue processing
			return next(c)
		}
	}
}
