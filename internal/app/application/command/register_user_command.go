// internal/app/application/command/register_user_command.go
package command

import (
	"context"

	"github.com/sh1ro/todo-api/internal/app/domain/model"
	"github.com/sh1ro/todo-api/internal/app/domain/service"
	"github.com/sh1ro/todo-api/pkg/logger"
)

// RegisterUserCommand represents a command to register a user
type RegisterUserCommand struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// RegisterUserHandler handles the RegisterUserCommand
type RegisterUserHandler struct {
	authService *service.AuthService
	logger      *logger.Logger
}

// NewRegisterUserHandler creates a new RegisterUserHandler
func NewRegisterUserHandler(authService *service.AuthService, logger *logger.Logger) *RegisterUserHandler {
	return &RegisterUserHandler{
		authService: authService,
		logger:      logger,
	}
}

// Handle handles the RegisterUserCommand
func (h *RegisterUserHandler) Handle(ctx context.Context, cmd RegisterUserCommand) (*model.User, error) {
	h.logger.Info("Registering user", "username", cmd.Username, "email", cmd.Email)

	user, err := h.authService.Register(ctx, cmd.Username, cmd.Email, cmd.Password)
	if err != nil {
		h.logger.Error("Failed to register user", "error", err)
		return nil, err
	}

	return user, nil
}
