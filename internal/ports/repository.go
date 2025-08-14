// Package ports defines the interfaces that connect the core business logic
// with external adapters, following the Hexagonal Architecture pattern.
// These interfaces represent the "ports" through which the application interacts with the outside world.
package ports

import (
	"context"

	"github.com/asp3cto/task-manager/internal/domain"
)

// TaskRepository defines the contract for task data persistence operations.
// Implementations of this interface handle the storage and retrieval of tasks
// from various data sources (memory, database, etc.).
type TaskRepository interface {
	// Create stores a new task in the repository.
	// Returns domain.ErrTaskExists if a task with the same ID already exists.
	Create(ctx context.Context, task *domain.Task) error

	// GetByID retrieves a task by its unique identifier.
	// Returns domain.ErrTaskNotFound if no task exists with the given ID.
	GetByID(ctx context.Context, id string) (*domain.Task, error)

	// GetAll retrieves all tasks, optionally filtered by status.
	// If status is empty, returns all tasks regardless of their status.
	// The status parameter should match one of the domain.TaskStatus values.
	GetAll(ctx context.Context, status string) ([]*domain.Task, error)

	// Update modifies an existing task in the repository.
	// Returns domain.ErrTaskNotFound if no task exists with the given ID.
	Update(ctx context.Context, task *domain.Task) error

	// Delete removes a task from the repository by its ID.
	// Returns domain.ErrTaskNotFound if no task exists with the given ID.
	Delete(ctx context.Context, id string) error
}
