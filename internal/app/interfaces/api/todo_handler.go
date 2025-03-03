package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
func (h *TodoHandler) CreateTodo(c *gin.Context) {
	var cmd command.CreateTodoCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		response.RespondWithBadRequest(c, "Invalid JSON format")
		return
	}

	// Get request-specific logger
	log := h.GetLogger(c)

	// Validate the command
	if errors := h.validator.Validate(cmd); errors != nil {
		log.Error("Validation failed for create todo", "errors", errors)
		response.RespondWithValidationError(c, "Validation failed", errors)
		return
	}

	// Get user ID from context
	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.RespondWithUnauthorized(c, "User ID not found in context")
		return
	}
	cmd.UserID = userID.(uuid.UUID)

	// Handle the command
	todo, err := h.createTodoHandler.Handle(c.Request.Context(), cmd)
	if err != nil {
		log.Error("Failed to create todo", "error", err)
		response.RespondWithInternalError(c, err.Error())
		return
	}

	// Use the generic response helper for type safety
	response.RespondWithGenericCreated(c, "Todo created successfully", todo)
}

// GetTodo handles getting a todo by ID
func (h *TodoHandler) GetTodo(c *gin.Context) {
	// Get user ID from context
	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.RespondWithUnauthorized(c, "User ID not found in context")
		return
	}

	// Get todo ID from URL
	todoID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.RespondWithBadRequest(c, "Invalid todo ID")
		return
	}

	// Get request-specific logger
	log := h.GetLogger(c)

	// Create query
	q := query.GetTodoQuery{
		UserID: userID.(uuid.UUID),
		TodoID: todoID,
	}

	// Handle the query
	todo, err := h.getTodoHandler.Handle(c.Request.Context(), q)
	if err != nil {
		log.Error("Failed to get todo", "error", err)
		response.RespondWithInternalError(c, err.Error())
		return
	}

	// Use the generic response helper for type safety
	response.RespondWithGenericOK(c, "Todo retrieved successfully", todo)
}

// ListTodos handles listing todos with filters
func (h *TodoHandler) ListTodos(c *gin.Context) {
	// Get user ID from context
	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.RespondWithUnauthorized(c, "User ID not found in context")
		return
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

	// Parse query parameters
	if page := c.Query("page"); page != "" {
		if err := h.validator.ValidateVar(page, "numeric"); err == nil {
			var pageInt int
			if _, err := fmt.Sscanf(page, "%d", &pageInt); err == nil && pageInt > 0 {
				q.Page = pageInt
			}
		}
	}

	if pageSize := c.Query("page_size"); pageSize != "" {
		if err := h.validator.ValidateVar(pageSize, "numeric"); err == nil {
			var pageSizeInt int
			if _, err := fmt.Sscanf(pageSize, "%d", &pageSizeInt); err == nil && pageSizeInt > 0 && pageSizeInt <= 100 {
				q.PageSize = pageSizeInt
			}
		}
	}

	if status := c.Query("status"); status != "" {
		todoStatus := model.TodoStatus(status)
		q.Status = &todoStatus
	}

	if priority := c.Query("priority"); priority != "" {
		todoPriority := model.TodoPriority(priority)
		q.Priority = &todoPriority
	}

	if search := c.Query("search"); search != "" {
		q.Search = &search
	}

	if sortBy := c.Query("sort_by"); sortBy != "" {
		q.SortBy = sortBy
	}

	if sortOrder := c.Query("sort_order"); sortOrder != "" {
		q.SortOrder = sortOrder
	}

	// Handle the query
	result, err := h.listTodosHandler.Handle(c.Request.Context(), q)
	if err != nil {
		log.Error("Failed to list todos", "error", err)
		response.RespondWithInternalError(c, err.Error())
		return
	}

	// Use the generic response helper for type safety
	response.RespondWithGenericPaginated(c, http.StatusOK, "Todos retrieved successfully", result.Todos, result)
}

// UpdateTodo handles updating a todo
func (h *TodoHandler) UpdateTodo(c *gin.Context) {
	var cmd command.UpdateTodoCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		response.RespondWithBadRequest(c, "Invalid JSON format")
		return
	}

	// Get request-specific logger
	log := h.GetLogger(c)

	// Validate the command
	if errors := h.validator.Validate(cmd); errors != nil {
		log.Error("Validation failed for update todo", "errors", errors)
		response.RespondWithValidationError(c, "Validation failed", errors)
		return
	}

	// Get user ID from context
	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.RespondWithUnauthorized(c, "User ID not found in context")
		return
	}
	cmd.UserID = userID.(uuid.UUID)

	// Get todo ID from URL
	todoID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.RespondWithBadRequest(c, "Invalid todo ID")
		return
	}
	cmd.TodoID = todoID

	// Handle the command
	todo, err := h.updateTodoHandler.Handle(c.Request.Context(), cmd)
	if err != nil {
		log.Error("Failed to update todo", "error", err)
		response.RespondWithInternalError(c, err.Error())
		return
	}

	// Use the generic response helper for type safety
	response.RespondWithGenericOK(c, "Todo updated successfully", todo)
}

// DeleteTodo handles deleting a todo
func (h *TodoHandler) DeleteTodo(c *gin.Context) {
	var cmd command.DeleteTodoCommand

	// Get request-specific logger
	log := h.GetLogger(c)

	// Get user ID from context
	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.RespondWithUnauthorized(c, "User ID not found in context")
		return
	}
	cmd.UserID = userID.(uuid.UUID)

	// Get todo ID from URL
	todoID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		log.Error("Invalid todo ID", "error", err)
		response.RespondWithBadRequest(c, "Invalid todo ID")
		return
	}
	cmd.ID = todoID

	// Handle the command
	err = h.deleteTodoHandler.Handle(c.Request.Context(), cmd)
	if err != nil {
		log.Error("Failed to delete todo", "error", err)
		response.RespondWithInternalError(c, err.Error())
		return
	}

	// Use the generic response helper for type safety with an empty struct
	type EmptyResponse struct{}
	response.RespondWithGenericOK(c, "Todo deleted successfully", EmptyResponse{})
}

// GetOverdueTodos handles getting overdue todos
func (h *TodoHandler) GetOverdueTodos(c *gin.Context) {
	// Get user ID from context
	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.RespondWithUnauthorized(c, "User ID not found in context")
		return
	}

	// Get request-specific logger
	log := h.GetLogger(c)

	// Create query
	q := query.GetOverdueTodosQuery{
		UserID: userID.(uuid.UUID),
	}

	// Handle the query
	todos, err := h.getOverdueTodosHandler.Handle(c.Request.Context(), q)
	if err != nil {
		log.Error("Failed to get overdue todos", "error", err)
		response.RespondWithInternalError(c, err.Error())
		return
	}

	// Use the generic response helper for type safety
	response.RespondWithGenericOK(c, "Overdue todos retrieved successfully", todos)
}

