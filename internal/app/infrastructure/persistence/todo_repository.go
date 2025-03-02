package persistence

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/sh1ro/todo-api/internal/app/domain/model"
	"github.com/sh1ro/todo-api/internal/app/domain/repository"
)

// PostgresTodoRepository implements the TodoRepository interface for PostgreSQL
type PostgresTodoRepository struct {
	db *PostgresDB
}

// NewPostgresTodoRepository creates a new PostgresTodoRepository
func NewPostgresTodoRepository(db *PostgresDB) repository.TodoRepository {
	return &PostgresTodoRepository{
		db: db,
	}
}

// Create creates a new todo
func (r *PostgresTodoRepository) Create(ctx context.Context, todo *model.Todo) error {
	query := `
		INSERT INTO todos (id, user_id, title, description, status, priority, due_date, created_at, updated_at, completed_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err := r.db.Exec(query,
		todo.ID,
		todo.UserID,
		todo.Title,
		todo.Description,
		todo.Status,
		todo.Priority,
		todo.DueDate,
		todo.CreatedAt,
		todo.UpdatedAt,
		todo.CompletedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create todo: %w", err)
	}

	return nil
}

// GetByID gets a todo by ID
func (r *PostgresTodoRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Todo, error) {
	query := `
		SELECT id, user_id, title, description, status, priority, due_date, created_at, updated_at, completed_at
		FROM todos
		WHERE id = $1
	`

	row := r.db.QueryRow(query, id)
	todo, err := r.scanTodo(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("todo not found")
		}
		return nil, fmt.Errorf("failed to get todo by ID: %w", err)
	}

	return todo, nil
}

// GetByUserIDAndID gets a todo by user ID and todo ID
func (r *PostgresTodoRepository) GetByUserIDAndID(ctx context.Context, userID, todoID uuid.UUID) (*model.Todo, error) {
	query := `
		SELECT id, user_id, title, description, status, priority, due_date, created_at, updated_at, completed_at
		FROM todos
		WHERE user_id = $1 AND id = $2
	`

	row := r.db.QueryRow(query, userID, todoID)
	todo, err := r.scanTodo(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("todo not found")
		}
		return nil, fmt.Errorf("failed to get todo by user ID and todo ID: %w", err)
	}

	return todo, nil
}

// List lists todos based on filter
func (r *PostgresTodoRepository) List(ctx context.Context, filter repository.TodoFilter) ([]*model.Todo, error) {
	query, args := r.buildListQuery(filter)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list todos: %w", err)
	}
	defer rows.Close()

	var todos []*model.Todo
	for rows.Next() {
		todo, err := r.scanTodoFromRows(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan todo: %w", err)
		}
		todos = append(todos, todo)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating todo rows: %w", err)
	}

	return todos, nil
}

// Count counts todos based on filter
func (r *PostgresTodoRepository) Count(ctx context.Context, filter repository.TodoFilter) (int, error) {
	whereClause, args := r.buildWhereClause(filter)

	query := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM todos
		%s
	`, whereClause)

	var count int
	err := r.db.QueryRow(query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count todos: %w", err)
	}

	return count, nil
}

// Update updates a todo
func (r *PostgresTodoRepository) Update(ctx context.Context, todo *model.Todo) error {
	query := `
		UPDATE todos
		SET title = $1, description = $2, status = $3, priority = $4, due_date = $5, updated_at = $6, completed_at = $7
		WHERE id = $8
	`

	_, err := r.db.Exec(query,
		todo.Title,
		todo.Description,
		todo.Status,
		todo.Priority,
		todo.DueDate,
		time.Now().UTC(),
		todo.CompletedAt,
		todo.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update todo: %w", err)
	}

	return nil
}

// Delete deletes a todo
func (r *PostgresTodoRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM todos
		WHERE id = $1
	`

	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete todo: %w", err)
	}

	return nil
}

// DeleteByUserID deletes all todos for a user
func (r *PostgresTodoRepository) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	query := `
		DELETE FROM todos
		WHERE user_id = $1
	`

	_, err := r.db.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete todos by user ID: %w", err)
	}

	return nil
}

// scanTodo scans a todo from a row
func (r *PostgresTodoRepository) scanTodo(row *sql.Row) (*model.Todo, error) {
	var todo model.Todo
	var dueDate sql.NullTime
	var completedAt sql.NullTime

	err := row.Scan(
		&todo.ID,
		&todo.UserID,
		&todo.Title,
		&todo.Description,
		&todo.Status,
		&todo.Priority,
		&dueDate,
		&todo.CreatedAt,
		&todo.UpdatedAt,
		&completedAt,
	)

	if err != nil {
		return nil, err
	}

	if dueDate.Valid {
		todo.DueDate = &dueDate.Time
	}

	if completedAt.Valid {
		todo.CompletedAt = &completedAt.Time
	}

	return &todo, nil
}

// scanTodoFromRows scans a todo from rows
func (r *PostgresTodoRepository) scanTodoFromRows(rows *sql.Rows) (*model.Todo, error) {
	var todo model.Todo
	var dueDate sql.NullTime
	var completedAt sql.NullTime

	err := rows.Scan(
		&todo.ID,
		&todo.UserID,
		&todo.Title,
		&todo.Description,
		&todo.Status,
		&todo.Priority,
		&dueDate,
		&todo.CreatedAt,
		&todo.UpdatedAt,
		&completedAt,
	)

	if err != nil {
		return nil, err
	}

	if dueDate.Valid {
		todo.DueDate = &dueDate.Time
	}

	if completedAt.Valid {
		todo.CompletedAt = &completedAt.Time
	}

	return &todo, nil
}

// buildListQuery builds a query for listing todos
func (r *PostgresTodoRepository) buildListQuery(filter repository.TodoFilter) (string, []interface{}) {
	whereClause, args := r.buildWhereClause(filter)

	// Add sorting
	orderBy := "created_at DESC"
	if filter.SortBy != "" {
		direction := "ASC"
		if strings.ToLower(filter.SortOrder) == "desc" {
			direction = "DESC"
		}
		orderBy = fmt.Sprintf("%s %s", filter.SortBy, direction)
	}

	// Add pagination
	limit := 10
	offset := 0
	if filter.Limit > 0 {
		limit = filter.Limit
	}
	if filter.Offset >= 0 {
		offset = filter.Offset
	}

	query := fmt.Sprintf(`
		SELECT id, user_id, title, description, status, priority, due_date, created_at, updated_at, completed_at
		FROM todos
		%s
		ORDER BY %s
		LIMIT %d OFFSET %d
	`, whereClause, orderBy, limit, offset)

	return query, args
}

// buildWhereClause builds a WHERE clause for filtering todos
func (r *PostgresTodoRepository) buildWhereClause(filter repository.TodoFilter) (string, []interface{}) {
	var conditions []string
	var args []interface{}
	argIndex := 1

	// Add user ID filter
	if filter.UserID != nil {
		conditions = append(conditions, fmt.Sprintf("user_id = $%d", argIndex))
		args = append(args, *filter.UserID)
		argIndex++
	}

	// Add status filter
	if filter.Status != nil {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, *filter.Status)
		argIndex++
	}

	// Add priority filter
	if filter.Priority != nil {
		conditions = append(conditions, fmt.Sprintf("priority = $%d", argIndex))
		args = append(args, *filter.Priority)
		argIndex++
	}

	// Add due date from filter
	if filter.DueDateFrom != nil {
		conditions = append(conditions, fmt.Sprintf("due_date >= $%d", argIndex))
		args = append(args, *filter.DueDateFrom)
		argIndex++
	}

	// Add due date to filter
	if filter.DueDateTo != nil {
		conditions = append(conditions, fmt.Sprintf("due_date <= $%d", argIndex))
		args = append(args, *filter.DueDateTo)
		argIndex++
	}

	// Add search filter
	if filter.Search != nil {
		conditions = append(conditions, fmt.Sprintf("(title ILIKE $%d OR description ILIKE $%d)", argIndex, argIndex))
		args = append(args, "%"+*filter.Search+"%")
		argIndex++
	}

	// Build WHERE clause
	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	return whereClause, args
}
