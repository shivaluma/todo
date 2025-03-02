package api

import (
	"github.com/gin-gonic/gin"
	"github.com/sh1ro/todo-api/internal/app/interfaces/middleware"
	"github.com/sh1ro/todo-api/pkg/logger"
)

// BaseHandler provides common functionality for all handlers
type BaseHandler struct {
	logger *logger.Logger
}

// NewBaseHandler creates a new BaseHandler
func NewBaseHandler(logger *logger.Logger) BaseHandler {
	return BaseHandler{
		logger: logger,
	}
}

// GetLogger returns the logger for the current request
func (h *BaseHandler) GetLogger(c *gin.Context) *logger.Logger {
	// Get request-specific logger from context if available
	if l, exists := c.Get("logger"); exists {
		return l.(*logger.Logger)
	}

	// Fall back to default logger
	return h.logger
}

// GetRequestID returns the request ID from the current request
func (h *BaseHandler) GetRequestID(c *gin.Context) string {
	if id, exists := c.Get(middleware.RequestIDContextKey); exists {
		return id.(string)
	}
	return ""
}
