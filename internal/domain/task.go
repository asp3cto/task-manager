// Package domain contains the core business entities and domain logic for the task management system.
// It defines the Task entity, its statuses, and domain-specific errors following Domain-Driven Design principles.
package domain

import (
	"errors"
	"time"
)

// Domain errors represent business rule violations and expected error conditions.
var (
	// ErrTaskNotFound is returned when a task with the specified ID does not exist.
	ErrTaskNotFound = errors.New("task not found")
	// ErrEmptyTitle is returned when attempting to create a task without a title.
	ErrEmptyTitle = errors.New("title in task cannot be empty")
	// ErrTaskExists is returned when attempting to create a task with an ID that already exists.
	ErrTaskExists = errors.New("task already exists")
)

// TaskStatus represents the current state of a task in its lifecycle.
type TaskStatus string

// Task status constants define the possible states a task can be in.
const (
	// StatusPending indicates a task that has been created but not yet started.
	StatusPending TaskStatus = "pending"
	// StatusInProgress indicates a task that is currently being worked on.
	StatusInProgress TaskStatus = "in_progress"
	// StatusCompleted indicates a task that has been finished successfully.
	StatusCompleted TaskStatus = "completed"
	// StatusCancelled indicates a task that was stopped before completion.
	StatusCancelled TaskStatus = "cancelled"
)

// Task represents a work item in the task management system.
// It contains all the information needed to track and manage a single task.
type Task struct {
	// ID is the unique identifier for the task.
	ID string `json:"id"`
	// Title is the short name or summary of the task.
	Title string `json:"title"`
	// Description provides detailed information about what the task involves.
	Description string `json:"description"`
	// Status indicates the current state of the task.
	Status TaskStatus `json:"status"`
	// CreatedAt is the timestamp when the task was first created.
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt is the timestamp when the task was last modified.
	UpdatedAt time.Time `json:"updated_at"`
}

// NewTask creates a new task with the provided details.
// The task is initialized with StatusPending and current timestamps.
// The id parameter should be unique across all tasks.
func NewTask(id, title, description string) *Task {
	now := time.Now()
	return &Task{
		ID:          id,
		Title:       title,
		Description: description,
		Status:      StatusPending,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// UpdateStatus changes the task's status and updates the UpdatedAt timestamp.
// This method should be used whenever the task's state changes.
func (t *Task) UpdateStatus(status TaskStatus) {
	t.Status = status
	t.UpdatedAt = time.Now()
}

// IsValidStatus checks if the provided status string is a valid TaskStatus.
// Returns true if the status is one of the defined constants, false otherwise.
func IsValidStatus(status string) bool {
	switch TaskStatus(status) {
	case StatusPending, StatusInProgress, StatusCompleted, StatusCancelled:
		return true
	default:
		return false
	}
}
