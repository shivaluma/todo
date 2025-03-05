package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sh1ro/todo-api/internal/app/infrastructure/persistence"
	"github.com/sh1ro/todo-api/internal/app/interfaces/api"
	customMiddleware "github.com/sh1ro/todo-api/internal/app/interfaces/middleware"
	"github.com/sh1ro/todo-api/pkg/config"
	"github.com/sh1ro/todo-api/pkg/logger"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	// Initialize logger
	log := logger.NewLogger(os.Getenv("LOG_LEVEL"), os.Getenv("LOG_FORMAT"))
	log.Info("Starting Todo API service")

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load configuration", "error", err)
	}

	// Initialize database connection
	db, err := persistence.NewPostgresDB(cfg.Database)
	if err != nil {
		log.Fatal("Failed to connect to database", "error", err)
	}
	defer db.Close()

	// Create Echo instance
	e := echo.New()
	
	// Configure Echo
	e.HideBanner = true
	e.HidePort = true
	
	// Add middleware
	e.Use(middleware.Recover())
	e.Use(customMiddleware.RequestID(log))
	e.Use(customMiddleware.Logger(log))
	e.Use(customMiddleware.CORS(cfg.CORS))

	// Setup API routes
	apiVersion := os.Getenv("API_VERSION")
	if apiVersion == "" {
		apiVersion = "v1"
	}

	apiGroup := e.Group(fmt.Sprintf("/api/%s", apiVersion))
	api.RegisterRoutes(apiGroup, db, log, cfg)

	// Start server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: e,
	}

	// Graceful shutdown
	go func() {
		log.Info("Server starting", "port", cfg.Port)
		if err := e.StartServer(srv); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server", "error", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown", "error", err)
	}

	log.Info("Server exiting")
}
