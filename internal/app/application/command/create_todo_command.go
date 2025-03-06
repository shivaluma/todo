// internal/app/application/command/create_todo_command.go
package command

import (
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sh1ro/todo-api/internal/app/domain/model"
	"github.com/sh1ro/todo-api/internal/app/domain/service"
	"github.com/sh1ro/todo-api/pkg/logger"
)

// CreateTodoCommand represents a command to create a todo
type CreateTodoCommand struct {
	UserID      uuid.UUID         `json:"-"`
	Title       string            `json:"title" validate:"required,min=1,max=255"`
	Description string            `json:"description"`
	Priority    model.TodoPriority `json:"priority" validate:"required,oneof=low medium high"`
	DueDate     *time.Time        `json:"due_date"`
}

// CreateTodoHandler handles the CreateTodoCommand
type CreateTodoHandler struct {
	todoService *service.TodoService
	logger      *logger.Logger
}

// NewCreateTodoHandler creates a new CreateTodoHandler
func NewCreateTodoHandler(todoService *service.TodoService, logger *logger.Logger) *CreateTodoHandler {
	return &CreateTodoHandler{
		todoService: todoService,
		logger:      logger,
	}
}

// Handle handles the CreateTodoCommand
func (h *CreateTodoHandler) Handle(c echo.Context, cmd CreateTodoCommand) (*model.Todo, error) {
	// Get request-specific logger with request ID
	log := logger.FromContext(c)
	log.Info("Creating todo", "userID", cmd.UserID, "title", cmd.Title)

	todo, err := h.todoService.CreateTodo(
		c.Request().Context(),
		cmd.UserID,
		cmd.Title,
		cmd.Description,
		cmd.Priority,
		cmd.DueDate,
	)

	if err != nil {
		log.Error("Failed to create todo", "error", err)
		return nil, err
	}

	return todo, nil
}
