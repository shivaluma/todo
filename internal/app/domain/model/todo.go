package model

import (
	"time"

	"github.com/google/uuid"
)

// TodoStatus represents the status of a todo item
type TodoStatus string

// TodoPriority represents the priority of a todo item
type TodoPriority string

const (
	// Todo statuses
	TodoStatusPending   TodoStatus = "pending"
	TodoStatusInProgress TodoStatus = "in_progress"
	TodoStatusCompleted TodoStatus = "completed"
	TodoStatusCancelled TodoStatus = "cancelled"

	// Todo priorities
	TodoPriorityLow    TodoPriority = "low"
	TodoPriorityMedium TodoPriority = "medium"
	TodoPriorityHigh   TodoPriority = "high"
)

// Todo represents a todo item
type Todo struct {
	ID          uuid.UUID    `json:"id"`
	UserID      uuid.UUID    `json:"user_id"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Status      TodoStatus   `json:"status"`
	Priority    TodoPriority `json:"priority"`
	DueDate     *time.Time   `json:"due_date"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
	CompletedAt *time.Time   `json:"completed_at"`
}

// NewTodo creates a new todo item
func NewTodo(userID uuid.UUID, title, description string, priority TodoPriority, dueDate *time.Time) *Todo {
	now := time.Now().UTC()
	return &Todo{
		ID:          uuid.New(),
		UserID:      userID,
		Title:       title,
		Description: description,
		Status:      TodoStatusPending,
		Priority:    priority,
		DueDate:     dueDate,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// UpdateTitle updates the todo's title
func (t *Todo) UpdateTitle(title string) {
	t.Title = title
	t.UpdatedAt = time.Now().UTC()
}

// UpdateDescription updates the todo's description
func (t *Todo) UpdateDescription(description string) {
	t.Description = description
	t.UpdatedAt = time.Now().UTC()
}

// UpdatePriority updates the todo's priority
func (t *Todo) UpdatePriority(priority TodoPriority) {
	t.Priority = priority
	t.UpdatedAt = time.Now().UTC()
}

// UpdateDueDate updates the todo's due date
func (t *Todo) UpdateDueDate(dueDate *time.Time) {
	t.DueDate = dueDate
	t.UpdatedAt = time.Now().UTC()
}

// UpdateStatus updates the todo's status
func (t *Todo) UpdateStatus(status TodoStatus) {
	t.Status = status
	t.UpdatedAt = time.Now().UTC()

	// If the status is completed, set the completed at time
	if status == TodoStatusCompleted {
		now := time.Now().UTC()
		t.CompletedAt = &now
	} else {
		t.CompletedAt = nil
	}
}

// MarkAsCompleted marks the todo as completed
func (t *Todo) MarkAsCompleted() {
	t.UpdateStatus(TodoStatusCompleted)
}

// MarkAsInProgress marks the todo as in progress
func (t *Todo) MarkAsInProgress() {
	t.UpdateStatus(TodoStatusInProgress)
}

// MarkAsPending marks the todo as pending
func (t *Todo) MarkAsPending() {
	t.UpdateStatus(TodoStatusPending)
}

// MarkAsCancelled marks the todo as cancelled
func (t *Todo) MarkAsCancelled() {
	t.UpdateStatus(TodoStatusCancelled)
}

// IsOverdue checks if the todo is overdue
func (t *Todo) IsOverdue() bool {
	if t.DueDate == nil {
		return false
	}
	return t.Status != TodoStatusCompleted && t.Status != TodoStatusCancelled && t.DueDate.Before(time.Now().UTC())
}

// TodoStatusPtr converts a TodoStatus to a pointer
func TodoStatusPtr(status TodoStatus) *TodoStatus {
	return &status
}

// TodoPriorityPtr converts a TodoPriority to a pointer
func TodoPriorityPtr(priority TodoPriority) *TodoPriority {
	return &priority
}
