// internal/app/application/command/delete_todo_command.go
package command

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sh1ro/todo-api/internal/app/domain/service"
	"github.com/sh1ro/todo-api/pkg/logger"
)

// DeleteTodoCommand represents a command to delete a todo
type DeleteTodoCommand struct {
	UserID uuid.UUID `json:"-"`
	ID     uuid.UUID `json:"-"`
}

// DeleteTodoHandler handles the DeleteTodoCommand
type DeleteTodoHandler struct {
	todoService *service.TodoService
	logger      *logger.Logger
}

// NewDeleteTodoHandler creates a new DeleteTodoHandler
func NewDeleteTodoHandler(todoService *service.TodoService, logger *logger.Logger) *DeleteTodoHandler {
	return &DeleteTodoHandler{
		todoService: todoService,
		logger:      logger,
	}
}

// Handle handles the DeleteTodoCommand
func (h *DeleteTodoHandler) Handle(c echo.Context, cmd DeleteTodoCommand) error {
	// Get request-specific logger with request ID
	log := logger.FromContext(c)
	log.Info("Deleting todo", "userID", cmd.UserID, "todoID", cmd.ID)

	err := h.todoService.DeleteTodo(c.Request().Context(), cmd.UserID, cmd.ID)
	if err != nil {
		log.Error("Failed to delete todo", "error", err)
		return err
	}

	return nil
}
