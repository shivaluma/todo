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

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sh1ro/todo-api/internal/app/infrastructure/persistence"
	"github.com/sh1ro/todo-api/internal/app/interfaces/api"
	"github.com/sh1ro/todo-api/internal/app/interfaces/middleware"
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

	// Set Gin mode
	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize database connection
	db, err := persistence.NewPostgresDB(cfg.Database)
	if err != nil {
		log.Fatal("Failed to connect to database", "error", err)
	}
	defer db.Close()

	// Create router
	router := gin.New()
	router.Use(gin.Recovery())

	router.Use(middleware.RequestID(log))
	router.Use(middleware.Logger(log))
	router.Use(middleware.CORS(cfg.CORS))

	// Setup API routes
	apiVersion := os.Getenv("API_VERSION")
	if apiVersion == "" {
		apiVersion = "v1"
	}

	apiGroup := router.Group(fmt.Sprintf("/api/%s", apiVersion))
	api.RegisterRoutes(apiGroup, db, log, cfg)

	// Start server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		log.Info("Server starting", "port", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
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

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown", "error", err)
	}

	log.Info("Server exiting")
}
