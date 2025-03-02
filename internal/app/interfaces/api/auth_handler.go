package api

import (
	"github.com/gin-gonic/gin"
	"github.com/sh1ro/todo-api/internal/app/application/command"
	"github.com/sh1ro/todo-api/pkg/logger"
	"github.com/sh1ro/todo-api/pkg/response"
	"github.com/sh1ro/todo-api/pkg/validator"
)

// AuthHandler handles authentication requests
type AuthHandler struct {
	BaseHandler
	registerUserHandler *command.RegisterUserHandler
	loginUserHandler    *command.LoginUserHandler
	validator           *validator.Validator
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(
	registerUserHandler *command.RegisterUserHandler,
	loginUserHandler *command.LoginUserHandler,
	validator *validator.Validator,
	logger *logger.Logger,
) *AuthHandler {
	return &AuthHandler{
		BaseHandler:         NewBaseHandler(logger),
		registerUserHandler: registerUserHandler,
		loginUserHandler:    loginUserHandler,
		validator:           validator,
	}
}

// Register handles user registration
func (h *AuthHandler) Register(c *gin.Context) {
	var cmd command.RegisterUserCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		response.RespondWithBadRequest(c, "Invalid JSON format")
		return
	}

	// Get request-specific logger
	log := h.GetLogger(c)

	// Validate the command
	if errors := h.validator.Validate(cmd); errors != nil {
		log.Error("Validation failed for user registration", "errors", errors)
		response.RespondWithValidationError(c, "Validation failed", errors)
		return
	}

	// Handle the command
	user, err := h.registerUserHandler.Handle(c.Request.Context(), cmd)
	if err != nil {
		log.Error("Failed to register user", "error", err)
		response.RespondWithInternalError(c, err.Error())
		return
	}

	// Create a user response struct
	type UserResponse struct {
		ID       string `json:"id"`
		Username string `json:"username"`
		Email    string `json:"email"`
	}

	userData := UserResponse{
		ID:       user.ID.String(),
		Username: user.Username,
		Email:    user.Email,
	}

	// Use the generic response helper
	response.RespondWithGenericCreated(c, "User registered successfully", userData)
}

// LoginResponse represents the response for a successful login
type LoginResponse struct {
	Token string `json:"token"`
}

// Login handles user login
func (h *AuthHandler) Login(c *gin.Context) {
	var cmd command.LoginUserCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		response.RespondWithBadRequest(c, "Invalid JSON format")
		return
	}

	// Get request-specific logger
	log := h.GetLogger(c)

	// Validate the command
	if errors := h.validator.Validate(cmd); errors != nil {
		log.Error("Validation failed for user login", "errors", errors)
		response.RespondWithValidationError(c, "Validation failed", errors)
		return
	}

	// Handle the command
	result, err := h.loginUserHandler.Handle(c.Request.Context(), cmd)
	if err != nil {
		log.Error("Failed to login user", "error", err)
		response.RespondWithUnauthorized(c, err.Error())
		return
	}

	// Create a strongly typed response
	loginResponse := LoginResponse{
		Token: result.Token,
	}

	// Use the generic response helper for type safety
	response.RespondWithGenericOK(c, "Login successful", loginResponse)
}
