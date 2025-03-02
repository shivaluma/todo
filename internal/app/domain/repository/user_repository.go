package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/sh1ro/todo-api/internal/app/domain/model"
)

// UserRepository defines the interface for user repository operations
type UserRepository interface {
	// Create creates a new user
	Create(ctx context.Context, user *model.User) error

	// GetByID gets a user by ID
	GetByID(ctx context.Context, id uuid.UUID) (*model.User, error)

	// GetByEmail gets a user by email
	GetByEmail(ctx context.Context, email string) (*model.User, error)

	// GetByUsername gets a user by username
	GetByUsername(ctx context.Context, username string) (*model.User, error)

	// Update updates a user
	Update(ctx context.Context, user *model.User) error

	// Delete deletes a user
	Delete(ctx context.Context, id uuid.UUID) error

	// Exists checks if a user exists by email or username
	Exists(ctx context.Context, email, username string) (bool, error)
}
