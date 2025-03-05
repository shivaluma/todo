package middleware

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sh1ro/todo-api/pkg/logger"
)

// Logger returns a middleware that logs request information
func Logger(log *logger.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			// Get request ID from context
			requestID := c.Response().Header().Get(echo.HeaderXRequestID)

			// Create a request-specific logger with request ID
			reqLogger := log.WithRequestID(requestID)

			// Store logger in context for handlers to use
			c.Set("logger", reqLogger)

			// Process request
			err := next(c)
			if err != nil {
				c.Error(err)
			}

			// Log request details after completion
			latency := time.Since(start)
			status := c.Response().Status
			method := c.Request().Method
			path := c.Request().URL.Path
			ip := c.RealIP()

			// Log at appropriate level based on status code
			switch {
			case status >= 500:
				reqLogger.Error("Request completed",
					"status", status,
					"method", method,
					"path", path,
					"ip", ip,
					"latency", latency,
				)
			case status >= 400:
				reqLogger.Warn("Request completed",
					"status", status,
					"method", method,
					"path", path,
					"ip", ip,
					"latency", latency,
				)
			default:
				reqLogger.Info("Request completed",
					"status", status,
					"method", method,
					"path", path,
					"ip", ip,
					"latency", latency,
				)
			}

			return err
		}
	}
}
