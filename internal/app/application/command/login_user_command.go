// internal/app/application/command/login_user_command.go
package command

import (
	"context"

	"github.com/sh1ro/todo-api/internal/app/domain/service"
	"github.com/sh1ro/todo-api/pkg/logger"
)

// LoginUserCommand represents a command to login a user
type LoginUserCommand struct {
	UsernameOrEmail string `json:"username_or_email" validate:"required"`
	Password        string `json:"password" validate:"required"`
}

// LoginResult represents the result of a login
type LoginResult struct {
	Token string `json:"token"`
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
func (h *LoginUserHandler) Handle(ctx context.Context, cmd LoginUserCommand) (*LoginResult, error) {
	h.logger.Info("Logging in user", "usernameOrEmail", cmd.UsernameOrEmail)

	token, err := h.authService.Login(ctx, cmd.UsernameOrEmail, cmd.Password)
	if err != nil {
		h.logger.Error("Failed to login user", "error", err)
		return nil, err
	}

	return &LoginResult{Token: token}, nil
}
