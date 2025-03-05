package api

import (
	"github.com/labstack/echo/v4"
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
// If a request-specific logger is available in the context, it will be used
// Otherwise, the default logger will be used
func (h *BaseHandler) GetLogger(c echo.Context) *logger.Logger {
	if l, ok := c.Get("logger").(*logger.Logger); ok {
		return l
	}
	return h.logger
}

// GetRequestID returns the request ID for the current request
func (h *BaseHandler) GetRequestID(c echo.Context) string {
	return c.Response().Header().Get(echo.HeaderXRequestID)
}

// GetUserID returns the user ID from the context
func (h *BaseHandler) GetUserID(c echo.Context) interface{} {
	userID := c.Get("user_id")
	if userID == nil {
		return nil
	}
	return userID
}
