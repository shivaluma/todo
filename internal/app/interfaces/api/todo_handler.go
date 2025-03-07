package api

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sh1ro/todo-api/internal/app/application/command"
	"github.com/sh1ro/todo-api/internal/app/application/query"
	"github.com/sh1ro/todo-api/internal/app/domain/model"
	"github.com/sh1ro/todo-api/internal/app/interfaces/middleware"
	"github.com/sh1ro/todo-api/pkg/logger"
	"github.com/sh1ro/todo-api/pkg/response"
	"github.com/sh1ro/todo-api/pkg/validator"
)

// TodoHandler handles todo requests
type TodoHandler struct {
	BaseHandler
	createTodoHandler       *command.CreateTodoHandler
	updateTodoHandler       *command.UpdateTodoHandler
	deleteTodoHandler       *command.DeleteTodoHandler
	getTodoHandler          *query.GetTodoHandler
	listTodosHandler        *query.ListTodosHandler
	getOverdueTodosHandler  *query.GetOverdueTodosHandler
	validator               *validator.Validator
}

// NewTodoHandler creates a new TodoHandler
func NewTodoHandler(
	createTodoHandler *command.CreateTodoHandler,
	updateTodoHandler *command.UpdateTodoHandler,
	deleteTodoHandler *command.DeleteTodoHandler,
	getTodoHandler *query.GetTodoHandler,
	listTodosHandler *query.ListTodosHandler,
	getOverdueTodosHandler *query.GetOverdueTodosHandler,
	validator *validator.Validator,
	logger *logger.Logger,
) *TodoHandler {
	return &TodoHandler{
		BaseHandler:            NewBaseHandler(logger),
		createTodoHandler:      createTodoHandler,
		updateTodoHandler:      updateTodoHandler,
		deleteTodoHandler:      deleteTodoHandler,
		getTodoHandler:         getTodoHandler,
		listTodosHandler:       listTodosHandler,
		getOverdueTodosHandler:  getOverdueTodosHandler,
		validator:              validator,
	}
}

// CreateTodo handles creating a new todo
func (h *TodoHandler) CreateTodo(c echo.Context) error {
	var cmd command.CreateTodoCommand
	if err := c.Bind(&cmd); err != nil {
		return response.RespondWithBadRequest(c, "Invalid JSON format")
	}

	// Get request-specific logger
	log := h.GetLogger(c)

	// Validate the command
	if errors := h.validator.Validate(cmd); errors != nil {
		log.Error("Validation failed for create todo", "errors", errors)
		return response.RespondWithValidationError(c, "Validation failed", errors)
	}

	// Get user ID from context
	userID, exists := middleware.GetUserID(c)
	if !exists {
		return response.RespondWithUnauthorized(c, "User ID not found in context")
	}
	cmd.UserID = userID.(uuid.UUID)

	// Handle the command
	todo, err := h.createTodoHandler.Handle(c, cmd)
	if err != nil {
		log.Error("Failed to create todo", "error", err)
		return response.RespondWithInternalError(c, err.Error())
	}

	// Use the generic response helper for type safety
	return response.RespondWithGenericCreated(c, "Todo created successfully", todo)
}

// GetTodo handles getting a todo by ID
func (h *TodoHandler) GetTodo(c echo.Context) error {
	// Get user ID from context
	userID, exists := middleware.GetUserID(c)
	if !exists {
		return response.RespondWithUnauthorized(c, "User ID not found in context")
	}

	// Get todo ID from URL
	todoIDStr := c.Param("id")
	if todoIDStr == "" {
		return response.RespondWithBadRequest(c, "Todo ID is required")
	}
	
	// Parse todo ID
	todoID, err := uuid.Parse(todoIDStr)
	if err != nil {
		return response.RespondWithBadRequest(c, "Invalid todo ID format")
	}

	// Get request-specific logger
	log := h.GetLogger(c)

	// Create query
	q := query.GetTodoQuery{
		UserID: userID.(uuid.UUID),
		TodoID: todoID,
	}

	// Handle the query
	todo, err := h.getTodoHandler.Handle(c, q)
	if err != nil {
		log.Error("Failed to get todo", "error", err)
		if err.Error() == "todo not found" {
			return response.RespondWithNotFound(c, "Todo not found")
		}
		return response.RespondWithInternalError(c, err.Error())
	}

	// Return the todo
	return response.RespondWithOK(c, "Todo retrieved successfully", todo)
}

// ListTodos handles listing todos with pagination and filtering
func (h *TodoHandler) ListTodos(c echo.Context) error {
	// Get user ID from context
	userID, exists := middleware.GetUserID(c)
	if !exists {
		return response.RespondWithUnauthorized(c, "User ID not found in context")
	}

	// Get request-specific logger
	log := h.GetLogger(c)

	// Create query with default values
	q := query.ListTodosQuery{
		UserID:    userID.(uuid.UUID),
		Page:      1,
		PageSize:  10,
		SortBy:    "created_at",
		SortOrder: "desc",
	}

	// Bind query parameters
	if err := c.Bind(&q); err != nil {
		return response.RespondWithBadRequest(c, "Invalid query parameters")
	}

	// Parse status filter
	if statusStr := c.QueryParam("status"); statusStr != "" {
		status := model.TodoStatus(statusStr)
		q.Status = &status
	}

	// Parse priority filter
	if priorityStr := c.QueryParam("priority"); priorityStr != "" {
		priority := model.TodoPriority(priorityStr)
		q.Priority = &priority
	}

	// Parse search filter
	if search := c.QueryParam("search"); search != "" {
		q.Search = &search
	}

	// Handle the query
	result, err := h.listTodosHandler.Handle(c, q)
	if err != nil {
		log.Error("Failed to list todos", "error", err)
		return response.RespondWithInternalError(c, err.Error())
	}

	// Return the todos
	return response.RespondWithOK(c, "Todos retrieved successfully", result)
}

// GetOverdueTodos handles getting overdue todos
func (h *TodoHandler) GetOverdueTodos(c echo.Context) error {
	// Get user ID from context
	userID, exists := middleware.GetUserID(c)
	if !exists {
		return response.RespondWithUnauthorized(c, "User ID not found in context")
	}

	// Get request-specific logger
	log := h.GetLogger(c)

	// Create query
	q := query.GetOverdueTodosQuery{
		UserID: userID.(uuid.UUID),
	}

	// Handle the query
	todos, err := h.getOverdueTodosHandler.Handle(c, q)
	if err != nil {
		log.Error("Failed to get overdue todos", "error", err)
		return response.RespondWithInternalError(c, err.Error())
	}

	// Return the todos
	return response.RespondWithOK(c, "Overdue todos retrieved successfully", todos)
}

// UpdateTodo handles updating a todo
func (h *TodoHandler) UpdateTodo(c echo.Context) error {
	// Get user ID from context
	userID, exists := middleware.GetUserID(c)
	if !exists {
		return response.RespondWithUnauthorized(c, "User ID not found in context")
	}

	// Get todo ID from URL
	todoIDStr := c.Param("id")
	if todoIDStr == "" {
		return response.RespondWithBadRequest(c, "Todo ID is required")
	}
	
	// Parse todo ID
	todoID, err := uuid.Parse(todoIDStr)
	if err != nil {
		return response.RespondWithBadRequest(c, "Invalid todo ID format")
	}

	// Parse request body
	var cmd command.UpdateTodoCommand
	if err := c.Bind(&cmd); err != nil {
		return response.RespondWithBadRequest(c, "Invalid JSON format")
	}

	// Get request-specific logger
	log := h.GetLogger(c)

	// Set the todo ID and user ID
	cmd.TodoID = todoID
	cmd.UserID = userID.(uuid.UUID)

	// Validate the command
	if errors := h.validator.Validate(cmd); errors != nil {
		log.Error("Validation failed for update todo", "errors", errors)
		return response.RespondWithValidationError(c, "Validation failed", errors)
	}

	// Handle the command
	todo, err := h.updateTodoHandler.Handle(c, cmd)
	if err != nil {
		log.Error("Failed to update todo", "error", err)
		if err.Error() == "todo not found" {
			return response.RespondWithNotFound(c, "Todo not found")
		}
		return response.RespondWithInternalError(c, err.Error())
	}

	// Return the updated todo
	return response.RespondWithOK(c, "Todo updated successfully", todo)
}

// DeleteTodo handles deleting a todo
func (h *TodoHandler) DeleteTodo(c echo.Context) error {
	// Get user ID from context
	userID, exists := middleware.GetUserID(c)
	if !exists {
		return response.RespondWithUnauthorized(c, "User ID not found in context")
	}

	// Get todo ID from URL
	todoIDStr := c.Param("id")
	if todoIDStr == "" {
		return response.RespondWithBadRequest(c, "Todo ID is required")
	}
	
	// Parse todo ID
	todoID, err := uuid.Parse(todoIDStr)
	if err != nil {
		return response.RespondWithBadRequest(c, "Invalid todo ID format")
	}

	// Get request-specific logger
	log := h.GetLogger(c)

	// Create command
	cmd := command.DeleteTodoCommand{
		ID:     todoID,
		UserID: userID.(uuid.UUID),
	}

	// Handle the command
	err = h.deleteTodoHandler.Handle(c, cmd)
	if err != nil {
		log.Error("Failed to delete todo", "error", err)
		if err.Error() == "todo not found" {
			return response.RespondWithNotFound(c, "Todo not found")
		}
		return response.RespondWithInternalError(c, err.Error())
	}

	// Return success with no content
	return response.RespondWithNoContent(c)
}

