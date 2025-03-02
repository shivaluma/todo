// internal/app/application/query/get_todo_query.go
package query

import (
	"context"

	"github.com/google/uuid"
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
func (h *GetTodoHandler) Handle(ctx context.Context, query GetTodoQuery) (*model.Todo, error) {
	h.logger.Info("Getting todo", "userID", query.UserID, "todoID", query.TodoID)

	todo, err := h.todoService.GetUserTodo(ctx, query.UserID, query.TodoID)
	if err != nil {
		h.logger.Error("Failed to get todo", "error", err)
		return nil, err
	}

	return todo, nil
}
