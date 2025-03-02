package persistence

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"github.com/sh1ro/todo-api/pkg/config"
	"github.com/sh1ro/todo-api/pkg/logger"
)

// PostgresDB represents a PostgreSQL database connection
type PostgresDB struct {
	DB     *sql.DB
	logger *logger.Logger
}

// NewPostgresDB creates a new PostgreSQL database connection
func NewPostgresDB(cfg config.DatabaseConfig) (*PostgresDB, error) {
	// Construct connection string
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode)

	// Open connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(cfg.MaxConnections)
	db.SetMaxIdleConns(cfg.MaxIdleConnections)
	db.SetConnMaxLifetime(cfg.ConnectionMaxLifetime)

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &PostgresDB{
		DB: db,
	}, nil
}

// Close closes the database connection
func (db *PostgresDB) Close() error {
	return db.DB.Close()
}

// SetLogger sets the logger for the database
func (db *PostgresDB) SetLogger(logger *logger.Logger) {
	db.logger = logger
}

// Begin starts a new transaction
func (db *PostgresDB) Begin() (*sql.Tx, error) {
	return db.DB.Begin()
}

// Query executes a query that returns rows
func (db *PostgresDB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	if db.logger != nil {
		db.logger.Debug("Executing query", "query", query, "args", args)
	}
	return db.DB.Query(query, args...)
}

// QueryRow executes a query that is expected to return at most one row
func (db *PostgresDB) QueryRow(query string, args ...interface{}) *sql.Row {
	if db.logger != nil {
		db.logger.Debug("Executing query row", "query", query, "args", args)
	}
	return db.DB.QueryRow(query, args...)
}

// Exec executes a query without returning any rows
func (db *PostgresDB) Exec(query string, args ...interface{}) (sql.Result, error) {
	if db.logger != nil {
		db.logger.Debug("Executing statement", "query", query, "args", args)
	}
	return db.DB.Exec(query, args...)
}

// ExecWithTimeout executes a query with a timeout
func (db *PostgresDB) ExecWithTimeout(timeout time.Duration, query string, args ...interface{}) (sql.Result, error) {
	ctx, cancel := db.createTimeoutContext(timeout)
	defer cancel()

	if db.logger != nil {
		db.logger.Debug("Executing statement with timeout", "query", query, "args", args, "timeout", timeout)
	}
	return db.DB.ExecContext(ctx, query, args...)
}

// QueryWithTimeout executes a query with a timeout
func (db *PostgresDB) QueryWithTimeout(timeout time.Duration, query string, args ...interface{}) (*sql.Rows, error) {
	ctx, cancel := db.createTimeoutContext(timeout)
	defer cancel()

	if db.logger != nil {
		db.logger.Debug("Executing query with timeout", "query", query, "args", args, "timeout", timeout)
	}
	return db.DB.QueryContext(ctx, query, args...)
}

// QueryRowWithTimeout executes a query with a timeout
func (db *PostgresDB) QueryRowWithTimeout(timeout time.Duration, query string, args ...interface{}) *sql.Row {
	ctx, cancel := db.createTimeoutContext(timeout)
	defer cancel()

	if db.logger != nil {
		db.logger.Debug("Executing query row with timeout", "query", query, "args", args, "timeout", timeout)
	}
	return db.DB.QueryRowContext(ctx, query, args...)
}

// createTimeoutContext creates a context with a timeout
func (db *PostgresDB) createTimeoutContext(timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), timeout)
}
