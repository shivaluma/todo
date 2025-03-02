package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/sh1ro/todo-api/internal/app/domain/model"
	"github.com/sh1ro/todo-api/internal/app/domain/repository"
	"github.com/sh1ro/todo-api/pkg/logger"
)

// TodoService provides todo related functionality
type TodoService struct {
	todoRepo repository.TodoRepository
	logger   *logger.Logger
}

// NewTodoService creates a new todo service
func NewTodoService(todoRepo repository.TodoRepository, logger *logger.Logger) *TodoService {
	return &TodoService{
		todoRepo: todoRepo,
		logger:   logger,
	}
}

// CreateTodo creates a new todo
func (s *TodoService) CreateTodo(ctx context.Context, userID uuid.UUID, title, description string, priority model.TodoPriority, dueDate *time.Time) (*model.Todo, error) {
	todo := model.NewTodo(userID, title, description, priority, dueDate)

	if err := s.todoRepo.Create(ctx, todo); err != nil {
		s.logger.Error("Failed to create todo", "error", err)
		return nil, err
	}

	return todo, nil
}

// GetTodo gets a todo by ID
func (s *TodoService) GetTodo(ctx context.Context, id uuid.UUID) (*model.Todo, error) {
	todo, err := s.todoRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get todo", "id", id, "error", err)
		return nil, err
	}

	return todo, nil
}

// GetUserTodo gets a todo by user ID and todo ID
func (s *TodoService) GetUserTodo(ctx context.Context, userID, todoID uuid.UUID) (*model.Todo, error) {
	todo, err := s.todoRepo.GetByUserIDAndID(ctx, userID, todoID)
	if err != nil {
		s.logger.Error("Failed to get user todo", "userID", userID, "todoID", todoID, "error", err)
		return nil, err
	}

	return todo, nil
}

// ListTodos lists todos based on filter
func (s *TodoService) ListTodos(ctx context.Context, filter repository.TodoFilter) ([]*model.Todo, int, error) {
	todos, err := s.todoRepo.List(ctx, filter)
	if err != nil {
		s.logger.Error("Failed to list todos", "error", err)
		return nil, 0, err
	}

	count, err := s.todoRepo.Count(ctx, filter)
	if err != nil {
		s.logger.Error("Failed to count todos", "error", err)
		return nil, 0, err
	}

	return todos, count, nil
}

// UpdateTodo updates a todo
func (s *TodoService) UpdateTodo(ctx context.Context, userID uuid.UUID, todoID uuid.UUID, title, description *string, status *model.TodoStatus, priority *model.TodoPriority, dueDate *time.Time) (*model.Todo, error) {
	todo, err := s.todoRepo.GetByUserIDAndID(ctx, userID, todoID)
	if err != nil {
		s.logger.Error("Failed to get todo for update", "userID", userID, "todoID", todoID, "error", err)
		return nil, err
	}

	if todo == nil {
		return nil, errors.New("todo not found")
	}

	// Update fields if provided
	if title != nil {
		todo.UpdateTitle(*title)
	}

	if description != nil {
		todo.UpdateDescription(*description)
	}

	if status != nil {
		todo.UpdateStatus(*status)
	}

	if priority != nil {
		todo.UpdatePriority(*priority)
	}

	if dueDate != nil {
		todo.UpdateDueDate(dueDate)
	}

	if err := s.todoRepo.Update(ctx, todo); err != nil {
		s.logger.Error("Failed to update todo", "todoID", todoID, "error", err)
		return nil, err
	}

	return todo, nil
}

// DeleteTodo deletes a todo
func (s *TodoService) DeleteTodo(ctx context.Context, userID, todoID uuid.UUID) error {
	todo, err := s.todoRepo.GetByUserIDAndID(ctx, userID, todoID)
	if err != nil {
		s.logger.Error("Failed to get todo for delete", "userID", userID, "todoID", todoID, "error", err)
		return err
	}

	if todo == nil {
		return errors.New("todo not found")
	}

	if err := s.todoRepo.Delete(ctx, todoID); err != nil {
		s.logger.Error("Failed to delete todo", "todoID", todoID, "error", err)
		return err
	}

	return nil
}

// MarkTodoAsCompleted marks a todo as completed
func (s *TodoService) MarkTodoAsCompleted(ctx context.Context, userID, todoID uuid.UUID) (*model.Todo, error) {
	todo, err := s.todoRepo.GetByUserIDAndID(ctx, userID, todoID)
	if err != nil {
		s.logger.Error("Failed to get todo for completion", "userID", userID, "todoID", todoID, "error", err)
		return nil, err
	}

	if todo == nil {
		return nil, errors.New("todo not found")
	}

	todo.MarkAsCompleted()

	if err := s.todoRepo.Update(ctx, todo); err != nil {
		s.logger.Error("Failed to mark todo as completed", "todoID", todoID, "error", err)
		return nil, err
	}

	return todo, nil
}

// GetOverdueTodos gets all overdue todos for a user
func (s *TodoService) GetOverdueTodos(ctx context.Context, userID uuid.UUID) ([]*model.Todo, error) {
	now := time.Now().UTC()
	filter := repository.TodoFilter{
		UserID:      &userID,
		DueDateTo:   &now,
		SortBy:      "due_date",
		SortOrder:   "asc",
	}

	todos, err := s.todoRepo.List(ctx, filter)
	if err != nil {
		s.logger.Error("Failed to get overdue todos", "userID", userID, "error", err)
		return nil, err
	}

	// Filter out completed and cancelled todos
	var overdueTodos []*model.Todo
	for _, todo := range todos {
		if todo.Status != model.TodoStatusCompleted && todo.Status != model.TodoStatusCancelled {
			overdueTodos = append(overdueTodos, todo)
		}
	}

	return overdueTodos, nil
}
