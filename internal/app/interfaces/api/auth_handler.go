package api

import (
	"github.com/labstack/echo/v4"
	"github.com/sh1ro/todo-api/internal/app/application/command"
	"github.com/sh1ro/todo-api/internal/app/interfaces/middleware"
	"github.com/sh1ro/todo-api/pkg/logger"
	"github.com/sh1ro/todo-api/pkg/response"
	"github.com/sh1ro/todo-api/pkg/validator"
)

// AuthHandler handles authentication requests
type AuthHandler struct {
	BaseHandler
	registerUserHandler *command.RegisterUserHandler
	loginUserHandler    *command.LoginUserHandler
	getUserHandler      *command.GetUserHandler
	validator           *validator.Validator
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(
	registerUserHandler *command.RegisterUserHandler,
	loginUserHandler *command.LoginUserHandler,
	getUserHandler *command.GetUserHandler,
	validator *validator.Validator,
	logger *logger.Logger,
) *AuthHandler {
	return &AuthHandler{
		BaseHandler:         NewBaseHandler(logger),
		registerUserHandler: registerUserHandler,
		loginUserHandler:    loginUserHandler,
		getUserHandler:      getUserHandler,
		validator:           validator,
	}
}

// Register handles user registration
func (h *AuthHandler) Register(c echo.Context) error {
	var cmd command.RegisterUserCommand
	if err := c.Bind(&cmd); err != nil {
		return response.RespondWithBadRequest(c, "Invalid JSON format")
	}

	// Get request-specific logger
	log := h.GetLogger(c)

	// Validate the command
	if errors := h.validator.Validate(cmd); errors != nil {
		log.Error("Validation failed for user registration", "errors", errors)
		return response.RespondWithValidationError(c, "Validation failed", errors)
	}

	// Handle the command
	user, err := h.registerUserHandler.Handle(c, cmd)
	if err != nil {
		log.Error("Failed to register user", "error", err)
		return response.RespondWithInternalError(c, err.Error())
	}

	// Return the user with a JWT token
	return response.RespondWithCreated(c, "User registered successfully", user)
}

// Login handles user login
func (h *AuthHandler) Login(c echo.Context) error {
	var cmd command.LoginUserCommand
	if err := c.Bind(&cmd); err != nil {
		return response.RespondWithBadRequest(c, "Invalid JSON format")
	}

	// Get request-specific logger
	log := h.GetLogger(c)
    

	// Validate the command
	if errors := h.validator.Validate(cmd); errors != nil {
		log.Error("Validation failed for user login", "errors", errors)
		return response.RespondWithValidationError(c, "Validation failed", errors)
	}

	// Handle the command
	user, err := h.loginUserHandler.Handle(c, cmd)
	if err != nil {
		log.Error("Failed to login user", "error", err)
		if err.Error() == "invalid credentials" {
			return response.RespondWithUnauthorized(c, "Invalid email or password")
		}
		return response.RespondWithInternalError(c, err.Error())
	}

	// Return the user with a JWT token
	return response.RespondWithOK(c, "User logged in successfully", user)
}

// Me handles getting the current user
func (h *AuthHandler) Me(c echo.Context) error {
	// Get user ID from context
	userID, exists := middleware.GetUserID(c)
	if !exists {
		return response.RespondWithUnauthorized(c, "User ID not found in context")
	}

	// Get request-specific logger
	log := h.GetLogger(c)

	// Create command
	cmd := command.GetUserCommand{
		UserID: userID.(string),
	}

	// Handle the command
	user, err := h.getUserHandler.Handle(c, cmd)
	if err != nil {
		log.Error("Failed to get user", "error", err)
		if err.Error() == "user not found" {
			return response.RespondWithNotFound(c, "User not found")
		}
		return response.RespondWithInternalError(c, err.Error())
	}

	// Return the user
	return response.RespondWithOK(c, "User retrieved successfully", user)
}
