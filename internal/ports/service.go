package ports

import (
	"context"

	"github.com/asp3cto/task-manager/internal/domain"
)

// TaskService defines the contract for task business logic operations.
// This interface encapsulates all the use cases and business rules for task management,
// providing a clean API for the application's core functionality.
type TaskService interface {
	// CreateTask creates a new task with the given title and description.
	// The task is automatically assigned a unique ID and set to pending status.
	// Returns domain.ErrEmptyTitle if the title is empty or whitespace.
	CreateTask(ctx context.Context, title, description string) (*domain.Task, error)

	// GetTaskByID retrieves a task by its unique identifier.
	// Returns domain.ErrTaskNotFound if no task exists with the given ID.
	GetTaskByID(ctx context.Context, id string) (*domain.Task, error)

	// GetAllTasks retrieves all tasks, optionally filtered by status.
	// If status is empty, returns all tasks regardless of their status.
	// The status parameter should match one of the domain.TaskStatus values.
	GetAllTasks(ctx context.Context, status string) ([]*domain.Task, error)

	// UpdateTaskStatus changes the status of an existing task.
	// Returns the updated task on success.
	// Returns domain.ErrTaskNotFound if no task exists with the given ID.
	UpdateTaskStatus(ctx context.Context, id string, status domain.TaskStatus) (*domain.Task, error)
}
