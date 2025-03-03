package api

import (
	"github.com/gin-gonic/gin"
	"github.com/sh1ro/todo-api/internal/app/application/command"
	"github.com/sh1ro/todo-api/pkg/logger"
	"github.com/sh1ro/todo-api/pkg/response"
	"github.com/sh1ro/todo-api/pkg/validator"
)

// AuthHandler handles authentication related requests
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
		response.RespondWithBadRequest(c, err.Error())
		return
	}

	// Create a user response struct
	type UserResponse struct {
		ID       string `json:"id"`
		Fullname string `json:"fullname"`
		Email    string `json:"email"`
	}

	userData := UserResponse{
		ID:       user.ID.String(),
		Fullname: user.Fullname,
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

// Me returns the current user
func (h *AuthHandler) Me(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.RespondWithUnauthorized(c, "Unauthorized")
		return
	}

	// Get request-specific logger
	log := h.GetLogger(c)

	user, err := h.getUserHandler.Handle(c.Request.Context(), command.GetUserCommand{
		UserID: userID.(string),
	})
	if err != nil {
		log.Error("Failed to get user", "error", err)
		response.RespondWithUnauthorized(c, "Failed to get user")
		return
	}

	// Create a user response struct
	type UserResponse struct {
		ID       string `json:"id"`
		Fullname string `json:"fullname"`
		Email    string `json:"email"`
	}

	userData := UserResponse{
		ID:       user.ID.String(),
		Fullname: user.Fullname,
		Email:    user.Email,
	}

	response.RespondWithGenericOK(c, "User retrieved successfully", userData)
}
