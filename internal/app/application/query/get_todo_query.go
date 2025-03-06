// internal/app/application/query/get_todo_query.go
package query

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sh1ro/todo-api/internal/app/domain/model"
	"github.com/sh1ro/todo-api/internal/app/domain/service"
	"github.com/sh1ro/todo-api/pkg/logger"
)

// GetTodoQuery represents a query to get a todo
type GetTodoQuery struct {
	UserID uuid.UUID `json:"-"`
	TodoID uuid.UUID `json:"-"`
}

// GetTodoHandler handles the GetTodoQuery
type GetTodoHandler struct {
	todoService *service.TodoService
	logger      *logger.Logger
}

// NewGetTodoHandler creates a new GetTodoHandler
func NewGetTodoHandler(todoService *service.TodoService, logger *logger.Logger) *GetTodoHandler {
	return &GetTodoHandler{
		todoService: todoService,
		logger:      logger,
	}
}

// Handle handles the GetTodoQuery
func (h *GetTodoHandler) Handle(c echo.Context, query GetTodoQuery) (*model.Todo, error) {
	// Get request-specific logger with request ID
	log := logger.FromContext(c)
	log.Info("Getting todo", "userID", query.UserID, "todoID", query.TodoID)

	todo, err := h.todoService.GetUserTodo(c.Request().Context(), query.UserID, query.TodoID)
	if err != nil {
		log.Error("Failed to get todo", "error", err)
		return nil, err
	}

	return todo, nil
}
