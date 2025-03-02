# Config Package

This package provides configuration management for the Todo API using environment variables and configuration files.

## Overview

The config package offers:

1. Environment-based configuration loading
2. Support for multiple configuration sources (env vars, .env files)
3. Type-safe configuration access
4. Default values for optional configuration
5. Validation of required configuration

## Usage

### Basic Usage

```go
import (
    "github.com/sh1ro/todo-api/pkg/config"
)

func main() {
    // Load configuration
    cfg, err := config.Load()
    if err != nil {
        panic(err)
    }

    // Access configuration values
    dbHost := cfg.Database.Host
    dbPort := cfg.Database.Port
    dbUser := cfg.Database.User
    dbPassword := cfg.Database.Password
    dbName := cfg.Database.Name

    // Use configuration values
    // ...
}
```

### Configuration Structure

The configuration is organized into logical sections:

```go
type Config struct {
    App struct {
        Name        string
        Environment string
        Port        int
        LogLevel    string
        Debug       bool
    }

    Database struct {
        Host     string
        Port     int
        User     string
        Password string
        Name     string
        SSLMode  string
    }

    Auth struct {
        JWTSecret     string
        TokenDuration int
    }

    // Other sections...
}
```

### Environment Variables

Configuration values can be set using environment variables:

```
# App configuration
APP_NAME=todo-api
APP_ENV=development
APP_PORT=8080
APP_LOG_LEVEL=debug
APP_DEBUG=true

# Database configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=todo_db
DB_SSL_MODE=disable

# Auth configuration
AUTH_JWT_SECRET=your-secret-key
AUTH_TOKEN_DURATION=24
```

### Environment Files

The package supports loading configuration from `.env` files:

```
# .env file
APP_NAME=todo-api
APP_ENV=development
APP_PORT=8080
APP_LOG_LEVEL=debug
APP_DEBUG=true

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=todo_db
DB_SSL_MODE=disable

AUTH_JWT_SECRET=your-secret-key
AUTH_TOKEN_DURATION=24
```

### Environment-Specific Configuration

The package supports loading environment-specific configuration files:

-   `.env` - Base configuration for all environments
-   `.env.development` - Development environment configuration
-   `.env.test` - Test environment configuration
-   `.env.production` - Production environment configuration

Environment-specific files override values from the base configuration.

## Configuration Sections

### App Configuration

-   `APP_NAME` - Application name
-   `APP_ENV` - Environment (development, test, production)
-   `APP_PORT` - HTTP server port
-   `APP_LOG_LEVEL` - Log level (debug, info, warn, error)
-   `APP_DEBUG` - Debug mode flag

### Database Configuration

-   `DB_HOST` - Database host
-   `DB_PORT` - Database port
-   `DB_USER` - Database user
-   `DB_PASSWORD` - Database password
-   `DB_NAME` - Database name
-   `DB_SSL_MODE` - SSL mode (disable, require, verify-ca, verify-full)

### Auth Configuration

-   `AUTH_JWT_SECRET` - JWT signing secret
-   `AUTH_TOKEN_DURATION` - Token duration in hours

### Metrics Configuration

-   `METRICS_ENABLED` - Enable metrics collection
-   `METRICS_PORT` - Metrics server port

### Tracing Configuration

-   `TRACING_ENABLED` - Enable distributed tracing
-   `TRACING_PROVIDER` - Tracing provider (jaeger, zipkin)
-   `TRACING_ENDPOINT` - Tracing collector endpoint
