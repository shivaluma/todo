package api

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sh1ro/todo-api/internal/app/application/command"
	"github.com/sh1ro/todo-api/internal/app/application/query"
	"github.com/sh1ro/todo-api/internal/app/domain/service"
	"github.com/sh1ro/todo-api/internal/app/infrastructure/auth"
	"github.com/sh1ro/todo-api/internal/app/infrastructure/persistence"
	"github.com/sh1ro/todo-api/internal/app/interfaces/middleware"
	"github.com/sh1ro/todo-api/pkg/config"
	"github.com/sh1ro/todo-api/pkg/logger"
	"github.com/sh1ro/todo-api/pkg/metrics"
	"github.com/sh1ro/todo-api/pkg/response"
	"github.com/sh1ro/todo-api/pkg/validator"
)

// RegisterRoutes registers all routes for the API
func RegisterRoutes(router *gin.RouterGroup, db *persistence.PostgresDB, log *logger.Logger, cfg *config.Config) {
	// Create validator
	validator := validator.NewValidator()

	// Register metrics middleware
	router.Use(metrics.MetricsMiddleware())

	// Register metrics endpoint
	metrics.RegisterMetricsEndpoint(router)

	// Create repositories
	userRepo := persistence.NewPostgresUserRepository(db)
	todoRepo := persistence.NewPostgresTodoRepository(db)

	// Create services
	authService := auth.NewAuthService(userRepo, log, cfg.JWT.Secret, cfg.JWT.Expiration)
	todoService := service.NewTodoService(todoRepo, log)

	// Create command handlers
	registerUserHandler := command.NewRegisterUserHandler(authService, log)
	loginUserHandler := command.NewLoginUserHandler(authService, log)
	getUserHandler := command.NewGetUserHandler(authService, log)
	createTodoHandler := command.NewCreateTodoHandler(todoService, log)
	updateTodoHandler := command.NewUpdateTodoHandler(todoService, log)
	deleteTodoHandler := command.NewDeleteTodoHandler(todoService, log)

	// Create query handlers
	getTodoHandler := query.NewGetTodoHandler(todoService, log)
	listTodosHandler := query.NewListTodosHandler(todoService, log)
	getOverdueTodosHandler := query.NewGetOverdueTodosHandler(todoService, log)

	// Create API handlers
	authHandler := NewAuthHandler(registerUserHandler, loginUserHandler, getUserHandler, validator, log)
	todoHandler := NewTodoHandler(
		createTodoHandler,
		updateTodoHandler,
		deleteTodoHandler,
		getTodoHandler,
		listTodosHandler,
		getOverdueTodosHandler,
		validator,
		log,
	)

	// Create middleware
	authMiddleware := middleware.NewAuthMiddleware(authService, log)

	// Register auth routes
	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/register", authHandler.Register)
		authRoutes.POST("/login", authHandler.Login)
	}

	userRoutes := router.Group("/users")
	userRoutes.Use(authMiddleware.Authenticate())
	{
		userRoutes.GET("/me", authHandler.Me)
	}

	// Register todo routes (protected by auth middleware)
	todoRoutes := router.Group("/todos")
	todoRoutes.Use(authMiddleware.Authenticate())
	{
		todoRoutes.POST("", todoHandler.CreateTodo)
		todoRoutes.GET("", todoHandler.ListTodos)
		todoRoutes.GET("/overdue", todoHandler.GetOverdueTodos)
		todoRoutes.GET("/:id", todoHandler.GetTodo)
		todoRoutes.PUT("/:id", todoHandler.UpdateTodo)
		todoRoutes.DELETE("/:id", todoHandler.DeleteTodo)
	}

	// Register health check route
	router.GET("/health", func(c *gin.Context) {
		// Create a strongly typed health response
		type HealthResponse struct {
			Status string    `json:"status"`
			Time   time.Time `json:"time"`
		}

		healthData := HealthResponse{
			Status: "ok",
			Time:   time.Now().UTC(),
		}

		response.RespondWithGenericOK(c, "Service is healthy", healthData)
	})
}
