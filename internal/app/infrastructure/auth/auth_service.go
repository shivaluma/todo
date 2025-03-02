package auth

import (
	"time"

	"github.com/sh1ro/todo-api/internal/app/domain/repository"
	"github.com/sh1ro/todo-api/internal/app/domain/service"
	"github.com/sh1ro/todo-api/pkg/logger"
)

// NewAuthService creates a new authentication service
// This is an adapter that creates the domain AuthService with infrastructure dependencies
func NewAuthService(
	userRepo repository.UserRepository,
	logger *logger.Logger,
	jwtSecret string,
	jwtExpiration time.Duration,
) *service.AuthService {
	// We directly use the domain AuthService implementation
	// The infrastructure layer is just providing the dependencies
	return service.NewAuthService(userRepo, logger, jwtSecret, jwtExpiration)
}
