// internal/app/application/query/get_overdue_todos_query.go
package query

import (
	"context"

	"github.com/google/uuid"
	"github.com/sh1ro/todo-api/internal/app/domain/model"
	"github.com/sh1ro/todo-api/internal/app/domain/service"
	"github.com/sh1ro/todo-api/pkg/logger"
)

// GetOverdueTodosQuery represents a query to get overdue todos
type GetOverdueTodosQuery struct {
	UserID uuid.UUID `json:"-"`
}

// GetOverdueTodosHandler handles the GetOverdueTodosQuery
type GetOverdueTodosHandler struct {
	todoService *service.TodoService
	logger      *logger.Logger
}

// NewGetOverdueTodosHandler creates a new GetOverdueTodosHandler
func NewGetOverdueTodosHandler(todoService *service.TodoService, logger *logger.Logger) *GetOverdueTodosHandler {
	return &GetOverdueTodosHandler{
		todoService: todoService,
		logger:      logger,
	}
}

// Handle handles the GetOverdueTodosQuery
func (h *GetOverdueTodosHandler) Handle(ctx context.Context, query GetOverdueTodosQuery) ([]*model.Todo, error) {
	h.logger.Info("Getting overdue todos", "userID", query.UserID)

	todos, err := h.todoService.GetOverdueTodos(ctx, query.UserID)
	if err != nil {
		h.logger.Error("Failed to get overdue todos", "error", err)
		return nil, err
	}

	return todos, nil
}
