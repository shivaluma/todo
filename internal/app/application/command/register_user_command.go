// internal/app/application/command/register_user_command.go
package command

import (
	"github.com/labstack/echo/v4"
	"github.com/sh1ro/todo-api/internal/app/domain/model"
	"github.com/sh1ro/todo-api/internal/app/domain/service"
	"github.com/sh1ro/todo-api/pkg/logger"
)

// RegisterUserCommand represents a command to register a user
type RegisterUserCommand struct {
	Fullname string `json:"fullname" validate:"required,min=3,max=50"`
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
func (h *RegisterUserHandler) Handle(c echo.Context, cmd RegisterUserCommand) (*model.User, error) {
	// Get request-specific logger with request ID
	log := logger.FromContext(c)
	log.Info("Registering user", "fullname", cmd.Fullname, "email", cmd.Email)

	user, err := h.authService.Register(c.Request().Context(), cmd.Fullname, cmd.Email, cmd.Password)
	if err != nil {
		log.Error("Failed to register user", "error", err)
		return nil, err
	}

	return user, nil
}
