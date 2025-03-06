package command

import (
	"github.com/labstack/echo/v4"
	"github.com/sh1ro/todo-api/internal/app/domain/model"
	"github.com/sh1ro/todo-api/internal/app/domain/service"
	"github.com/sh1ro/todo-api/pkg/logger"
)

// GetUserCommand represents a command to get the current user
type GetUserCommand struct {
	UserID string `json:"user_id"`
}

// GetUserHandler handles the GetUserCommand
type GetUserHandler struct {
	authService *service.AuthService
	logger      *logger.Logger
}

// NewGetUserHandler creates a new GetUserHandler
func NewGetUserHandler(authService *service.AuthService, logger *logger.Logger) *GetUserHandler {
	return &GetUserHandler{
		authService: authService,
		logger:      logger,
	}
}

// Handle handles the GetUserCommand
func (h *GetUserHandler) Handle(c echo.Context, cmd GetUserCommand) (*model.User, error) {
	// Get request-specific logger with request ID
	log := logger.FromContext(c)
	log.Info("Getting current user details")

	user, err := h.authService.GetUserFromId(c.Request().Context(), cmd.UserID)
	if err != nil {
		log.Error("Failed to get user", "user_id", cmd.UserID, "error", err)
		return nil, err
	}

	return user, nil
}
