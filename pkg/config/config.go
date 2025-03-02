package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config holds all configuration for the application
type Config struct {
	Env      string
	Port     int
	Database DatabaseConfig
	JWT      JWTConfig
	CORS     CORSConfig
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host                  string
	Port                  int
	User                  string
	Password              string
	Name                  string
	SSLMode               string
	MaxConnections        int
	MaxIdleConnections    int
	ConnectionMaxLifetime time.Duration
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret     string
	Expiration time.Duration
}

// CORSConfig holds CORS configuration
type CORSConfig struct {
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
	MaxAge         int
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	port, err := strconv.Atoi(getEnv("PORT", "8080"))
	if err != nil {
		return nil, fmt.Errorf("invalid PORT: %w", err)
	}

	dbPort, err := strconv.Atoi(getEnv("DB_PORT", "5432"))
	if err != nil {
		return nil, fmt.Errorf("invalid DB_PORT: %w", err)
	}

	dbMaxConn, err := strconv.Atoi(getEnv("DB_MAX_CONNECTIONS", "100"))
	if err != nil {
		return nil, fmt.Errorf("invalid DB_MAX_CONNECTIONS: %w", err)
	}

	dbMaxIdleConn, err := strconv.Atoi(getEnv("DB_MAX_IDLE_CONNECTIONS", "10"))
	if err != nil {
		return nil, fmt.Errorf("invalid DB_MAX_IDLE_CONNECTIONS: %w", err)
	}

	dbConnMaxLifetime, err := time.ParseDuration(getEnv("DB_CONNECTION_MAX_LIFETIME", "1h"))
	if err != nil {
		return nil, fmt.Errorf("invalid DB_CONNECTION_MAX_LIFETIME: %w", err)
	}

	jwtExpiration, err := time.ParseDuration(getEnv("JWT_EXPIRATION", "24h"))
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_EXPIRATION: %w", err)
	}

	corsMaxAge, err := strconv.Atoi(getEnv("CORS_MAX_AGE", "300"))
	if err != nil {
		return nil, fmt.Errorf("invalid CORS_MAX_AGE: %w", err)
	}

	return &Config{
		Env:  getEnv("ENV", "development"),
		Port: port,
		Database: DatabaseConfig{
			Host:                  getEnv("DB_HOST", "localhost"),
			Port:                  dbPort,
			User:                  getEnv("DB_USER", "postgres"),
			Password:              getEnv("DB_PASSWORD", "postgres"),
			Name:                  getEnv("DB_NAME", "todo_db"),
			SSLMode:               getEnv("DB_SSL_MODE", "disable"),
			MaxConnections:        dbMaxConn,
			MaxIdleConnections:    dbMaxIdleConn,
			ConnectionMaxLifetime: dbConnMaxLifetime,
		},
		JWT: JWTConfig{
			Secret:     getEnv("JWT_SECRET", "your_jwt_secret_key_here"),
			Expiration: jwtExpiration,
		},
		CORS: CORSConfig{
			AllowedOrigins: strings.Split(getEnv("CORS_ALLOWED_ORIGINS", "*"), ","),
			AllowedMethods: strings.Split(getEnv("CORS_ALLOWED_METHODS", "GET,POST,PUT,DELETE,OPTIONS"), ","),
			AllowedHeaders: strings.Split(getEnv("CORS_ALLOWED_HEADERS", "Authorization,Content-Type"), ","),
			MaxAge:         corsMaxAge,
		},
	}, nil
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
