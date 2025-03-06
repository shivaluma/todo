// internal/app/application/query/list_todos_query.go
package query

import (
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sh1ro/todo-api/internal/app/domain/model"
	"github.com/sh1ro/todo-api/internal/app/domain/repository"
	"github.com/sh1ro/todo-api/internal/app/domain/service"
	"github.com/sh1ro/todo-api/pkg/logger"
)

// ListTodosQuery represents a query to list todos
type ListTodosQuery struct {
	UserID      uuid.UUID          `json:"-"`
	Status      *model.TodoStatus  `json:"status"`
	Priority    *model.TodoPriority `json:"priority"`
	DueDateFrom *time.Time         `json:"due_date_from"`
	DueDateTo   *time.Time         `json:"due_date_to"`
	Search      *string            `json:"search"`
	Page        int                `json:"page" validate:"min=1"`
	PageSize    int                `json:"page_size" validate:"min=1,max=100"`
	SortBy      string             `json:"sort_by"`
	SortOrder   string             `json:"sort_order" validate:"omitempty,oneof=asc desc"`
}

// TodosResult represents the result of listing todos
type TodosResult struct {
	Todos      []*model.Todo `json:"todos"`
	TotalCount int           `json:"total_count"`
	Page       int           `json:"page"`
	PageSize   int           `json:"page_size"`
	TotalPages int           `json:"total_pages"`
}

// ListTodosHandler handles the ListTodosQuery
type ListTodosHandler struct {
	todoService *service.TodoService
	logger      *logger.Logger
}

// NewListTodosHandler creates a new ListTodosHandler
func NewListTodosHandler(todoService *service.TodoService, logger *logger.Logger) *ListTodosHandler {
	return &ListTodosHandler{
		todoService: todoService,
		logger:      logger,
	}
}

// Handle handles the ListTodosQuery
func (h *ListTodosHandler) Handle(c echo.Context, query ListTodosQuery) (*TodosResult, error) {
	// Get request-specific logger with request ID
	log := logger.FromContext(c)
	log.Info("Listing todos", "userID", query.UserID)

	// Set default values
	page := 1
	if query.Page > 0 {
		page = query.Page
	}

	pageSize := 10
	if query.PageSize > 0 {
		pageSize = query.PageSize
	}

	// Calculate offset
	offset := (page - 1) * pageSize

	// Create filter
	filter := repository.TodoFilter{
		UserID:      &query.UserID,
		Status:      query.Status,
		Priority:    query.Priority,
		DueDateFrom: query.DueDateFrom,
		DueDateTo:   query.DueDateTo,
		Search:      query.Search,
		Limit:       pageSize,
		Offset:      offset,
		SortBy:      query.SortBy,
		SortOrder:   query.SortOrder,
	}

	// Get todos and count
	todos, count, err := h.todoService.ListTodos(c.Request().Context(), filter)
	if err != nil {
		log.Error("Failed to list todos", "error", err)
		return nil, err
	}

	// Calculate total pages
	totalPages := count / pageSize
	if count%pageSize > 0 {
		totalPages++
	}

	return &TodosResult{
		Todos:      todos,
		TotalCount: count,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}
