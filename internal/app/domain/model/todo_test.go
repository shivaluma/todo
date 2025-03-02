package model

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestNewTodo(t *testing.T) {
	userID := uuid.New()
	title := "Test Todo"
	description := "This is a test todo"
	priority := TodoPriorityHigh
	now := time.Now().UTC()
	dueDate := now.Add(24 * time.Hour)

	todo := NewTodo(userID, title, description, priority, &dueDate)

	if todo.ID == uuid.Nil {
		t.Error("Expected todo ID to be set, got nil")
	}

	if todo.UserID != userID {
		t.Errorf("Expected user ID to be %v, got %v", userID, todo.UserID)
	}

	if todo.Title != title {
		t.Errorf("Expected title to be %s, got %s", title, todo.Title)
	}

	if todo.Description != description {
		t.Errorf("Expected description to be %s, got %s", description, todo.Description)
	}

	if todo.Priority != priority {
		t.Errorf("Expected priority to be %s, got %s", priority, todo.Priority)
	}

	if todo.Status != TodoStatusPending {
		t.Errorf("Expected status to be %s, got %s", TodoStatusPending, todo.Status)
	}

	if todo.DueDate == nil || !todo.DueDate.Equal(dueDate) {
		t.Errorf("Expected due date to be %v, got %v", dueDate, todo.DueDate)
	}

	if todo.CompletedAt != nil {
		t.Errorf("Expected completed at to be nil, got %v", todo.CompletedAt)
	}
}

func TestMarkAsCompleted(t *testing.T) {
	userID := uuid.New()
	todo := NewTodo(userID, "Test Todo", "This is a test todo", TodoPriorityMedium, nil)

	todo.MarkAsCompleted()

	if todo.Status != TodoStatusCompleted {
		t.Errorf("Expected status to be %s, got %s", TodoStatusCompleted, todo.Status)
	}

	if todo.CompletedAt == nil {
		t.Error("Expected completed at to be set, got nil")
	}
}

func TestIsOverdue(t *testing.T) {
	userID := uuid.New()

	// Test with no due date
	todo1 := NewTodo(userID, "Test Todo 1", "This is a test todo", TodoPriorityMedium, nil)
	if todo1.IsOverdue() {
		t.Error("Expected todo with no due date to not be overdue")
	}

	// Test with future due date
	futureDueDate := time.Now().UTC().Add(24 * time.Hour)
	todo2 := NewTodo(userID, "Test Todo 2", "This is a test todo", TodoPriorityMedium, &futureDueDate)
	if todo2.IsOverdue() {
		t.Error("Expected todo with future due date to not be overdue")
	}

	// Test with past due date
	pastDueDate := time.Now().UTC().Add(-24 * time.Hour)
	todo3 := NewTodo(userID, "Test Todo 3", "This is a test todo", TodoPriorityMedium, &pastDueDate)
	if !todo3.IsOverdue() {
		t.Error("Expected todo with past due date to be overdue")
	}

	// Test with past due date but completed
	todo4 := NewTodo(userID, "Test Todo 4", "This is a test todo", TodoPriorityMedium, &pastDueDate)
	todo4.MarkAsCompleted()
	if todo4.IsOverdue() {
		t.Error("Expected completed todo to not be overdue")
	}
}
