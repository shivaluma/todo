// internal/app/application/command/login_user_command.go
package command

import (
	"github.com/labstack/echo/v4"
	"github.com/sh1ro/todo-api/internal/app/domain/model"
	"github.com/sh1ro/todo-api/internal/app/domain/service"
	"github.com/sh1ro/todo-api/pkg/logger"
)

// LoginUserCommand represents a command to login a user
type LoginUserCommand struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// LoginResult represents the result of a login
type LoginResult struct {
	Token string `json:"token"`
	User  *model.User `json:"user"`
}

// LoginUserHandler handles the LoginUserCommand
type LoginUserHandler struct {
	authService *service.AuthService
	logger      *logger.Logger
}

// NewLoginUserHandler creates a new LoginUserHandler
func NewLoginUserHandler(authService *service.AuthService, logger *logger.Logger) *LoginUserHandler {
	return &LoginUserHandler{
		authService: authService,
		logger:      logger,
	}
}

// Handle handles the LoginUserCommand
func (h *LoginUserHandler) Handle(c echo.Context, cmd LoginUserCommand) (*LoginResult, error) {
	// Get request-specific logger with request ID
	log := logger.FromContext(c)
	log.Info("Logging in user", "email", cmd.Email)

	user, token, err := h.authService.Login(c.Request().Context(), cmd.Email, cmd.Password)
	if err != nil {
		log.Error("Failed to login user", "error", err)
		return nil, err
	}

	return &LoginResult{Token: token, User: user}, nil
}
