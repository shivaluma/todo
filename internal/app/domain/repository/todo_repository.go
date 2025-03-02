package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/sh1ro/todo-api/internal/app/domain/model"
)

// TodoFilter defines the filter options for querying todos
type TodoFilter struct {
	UserID      *uuid.UUID
	Status      *model.TodoStatus
	Priority    *model.TodoPriority
	DueDateFrom *time.Time
	DueDateTo   *time.Time
	Search      *string
	Limit       int
	Offset      int
	SortBy      string
	SortOrder   string
}

// TodoRepository defines the interface for todo repository operations
type TodoRepository interface {
	// Create creates a new todo
	Create(ctx context.Context, todo *model.Todo) error

	// GetByID gets a todo by ID
	GetByID(ctx context.Context, id uuid.UUID) (*model.Todo, error)

	// GetByUserIDAndID gets a todo by user ID and todo ID
	GetByUserIDAndID(ctx context.Context, userID, todoID uuid.UUID) (*model.Todo, error)

	// List lists todos based on filter
	List(ctx context.Context, filter TodoFilter) ([]*model.Todo, error)

	// Count counts todos based on filter
	Count(ctx context.Context, filter TodoFilter) (int, error)

	// Update updates a todo
	Update(ctx context.Context, todo *model.Todo) error

	// Delete deletes a todo
	Delete(ctx context.Context, id uuid.UUID) error

	// DeleteByUserID deletes all todos for a user
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error
}
