# Database Package

This package provides database connectivity and management for the Todo API using PostgreSQL and the `pgx` driver.

## Overview

The database package offers:

1. Connection management for PostgreSQL
2. Migration support using `golang-migrate`
3. Transaction management
4. Query execution helpers
5. Connection pooling
6. Metrics collection for database operations

## Usage

### Basic Usage

```go
import (
    "context"
    "github.com/sh1ro/todo-api/pkg/config"
    "github.com/sh1ro/todo-api/pkg/database"
)

func main() {
    // Load configuration
    cfg, err := config.Load()
    if err != nil {
        panic(err)
    }

    // Connect to the database
    db, err := database.Connect(cfg.Database)
    if err != nil {
        panic(err)
    }
    defer db.Close()

    // Run migrations
    if err := database.MigrateUp(cfg.Database); err != nil {
        panic(err)
    }

    // Use the database
    ctx := context.Background()
    if err := db.Ping(ctx); err != nil {
        panic(err)
    }

    // ... use the database for queries
}
```

### Query Execution

```go
import (
    "context"
    "github.com/sh1ro/todo-api/pkg/database"
)

func GetUser(ctx context.Context, db *database.DB, id string) (*User, error) {
    var user User

    query := `SELECT id, username, email, created_at, updated_at FROM users WHERE id = $1`
    err := db.QueryRow(ctx, query, id).Scan(
        &user.ID,
        &user.Username,
        &user.Email,
        &user.CreatedAt,
        &user.UpdatedAt,
    )
    if err != nil {
        return nil, err
    }

    return &user, nil
}

func GetUsers(ctx context.Context, db *database.DB, limit, offset int) ([]User, error) {
    var users []User

    query := `SELECT id, username, email, created_at, updated_at FROM users ORDER BY created_at DESC LIMIT $1 OFFSET $2`
    rows, err := db.Query(ctx, query, limit, offset)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var user User
        if err := rows.Scan(
            &user.ID,
            &user.Username,
            &user.Email,
            &user.CreatedAt,
            &user.UpdatedAt,
        ); err != nil {
            return nil, err
        }
        users = append(users, user)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }

    return users, nil
}
```

### Transaction Management

```go
import (
    "context"
    "github.com/sh1ro/todo-api/pkg/database"
)

func CreateUserWithProfile(ctx context.Context, db *database.DB, user User, profile Profile) error {
    return db.WithTransaction(ctx, func(tx database.Tx) error {
        // Insert user
        query := `INSERT INTO users (id, username, email, password_hash, created_at, updated_at)
                 VALUES ($1, $2, $3, $4, $5, $6)`
        _, err := tx.Exec(ctx, query,
            user.ID,
            user.Username,
            user.Email,
            user.PasswordHash,
            user.CreatedAt,
            user.UpdatedAt,
        )
        if err != nil {
            return err
        }

        // Insert profile
        query = `INSERT INTO profiles (user_id, first_name, last_name, bio, created_at, updated_at)
                VALUES ($1, $2, $3, $4, $5, $6)`
        _, err = tx.Exec(ctx, query,
            user.ID,
            profile.FirstName,
            profile.LastName,
            profile.Bio,
            profile.CreatedAt,
            profile.UpdatedAt,
        )
        if err != nil {
            return err
        }

        return nil
    })
}
```

## Migration Management

The package provides functions for managing database migrations:

```go
// Run migrations up to the latest version
if err := database.MigrateUp(cfg.Database); err != nil {
    panic(err)
}

// Run migrations up to a specific version
if err := database.MigrateTo(cfg.Database, 5); err != nil {
    panic(err)
}

// Roll back the last migration
if err := database.MigrateDown(cfg.Database); err != nil {
    panic(err)
}

// Roll back to a specific version
if err := database.MigrateDownTo(cfg.Database, 3); err != nil {
    panic(err)
}

// Get the current migration version
version, dirty, err := database.MigrateVersion(cfg.Database)
if err != nil {
    panic(err)
}
```

## Connection Configuration

The database connection can be configured with:

```go
type DatabaseConfig struct {
    Host     string
    Port     int
    User     string
    Password string
    Name     string
    SSLMode  string

    MaxOpenConns     int
    MaxIdleConns     int
    ConnMaxLifetime  time.Duration
    ConnMaxIdleTime  time.Duration

    MigrationsPath string
}
```

## Metrics Collection

The package automatically collects metrics for database operations when the metrics package is initialized:

-   `database_connections` - Number of open connections
-   `database_connections_max` - Maximum number of open connections
-   `database_connections_in_use` - Number of connections currently in use
-   `database_connections_idle` - Number of idle connections
-   `database_operations_total` - Total number of database operations
-   `database_operation_duration_seconds` - Duration of database operations
