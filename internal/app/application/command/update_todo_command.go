// internal/app/application/command/update_todo_command.go
package command

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/sh1ro/todo-api/internal/app/domain/model"
	"github.com/sh1ro/todo-api/internal/app/domain/service"
	"github.com/sh1ro/todo-api/pkg/logger"
)

// UpdateTodoCommand represents a command to update a todo
type UpdateTodoCommand struct {
	UserID      uuid.UUID          `json:"-"`
	TodoID      uuid.UUID          `json:"-"`
	Title       *string            `json:"title" validate:"omitempty,min=1,max=255"`
	Description *string            `json:"description"`
	Status      *model.TodoStatus  `json:"status" validate:"omitempty,oneof=pending in_progress completed cancelled"`
	Priority    *model.TodoPriority `json:"priority" validate:"omitempty,oneof=low medium high"`
	DueDate     *time.Time         `json:"due_date"`
}

// UpdateTodoHandler handles the UpdateTodoCommand
type UpdateTodoHandler struct {
	todoService *service.TodoService
	logger      *logger.Logger
}

// NewUpdateTodoHandler creates a new UpdateTodoHandler
func NewUpdateTodoHandler(todoService *service.TodoService, logger *logger.Logger) *UpdateTodoHandler {
	return &UpdateTodoHandler{
		todoService: todoService,
		logger:      logger,
	}
}

// Handle handles the UpdateTodoCommand
func (h *UpdateTodoHandler) Handle(ctx context.Context, cmd UpdateTodoCommand) (*model.Todo, error) {
	h.logger.Info("Updating todo", "userID", cmd.UserID, "todoID", cmd.TodoID)

	todo, err := h.todoService.UpdateTodo(
		ctx,
		cmd.UserID,
		cmd.TodoID,
		cmd.Title,
		cmd.Description,
		cmd.Status,
		cmd.Priority,
		cmd.DueDate,
	)

	if err != nil {
		h.logger.Error("Failed to update todo", "error", err)
		return nil, err
	}

	return todo, nil
}
