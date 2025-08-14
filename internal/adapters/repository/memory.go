// Package repository provides concrete implementations of data persistence interfaces.
// This package contains adapters that implement the repository ports defined in the ports package.
package repository

import (
	"context"
	"sync"

	"github.com/asp3cto/task-manager/internal/domain"
	"github.com/asp3cto/task-manager/internal/ports"
)

var _ ports.TaskRepository = (*MemoryTaskRepository)(nil)

// MemoryTaskRepository provides an in-memory implementation of the TaskRepository interface.
// It stores tasks in a map with thread-safe access using read-write mutexes.
// Data is lost when the application restarts since it's stored only in memory.
type MemoryTaskRepository struct {
	// tasks stores the task data indexed by task ID
	tasks map[string]*domain.Task
	// mu provides thread-safe access to the tasks map
	mu sync.RWMutex
}

// NewMemoryTaskRepository creates a new instance of the in-memory task repository.
func NewMemoryTaskRepository() *MemoryTaskRepository {
	return &MemoryTaskRepository{
		tasks: make(map[string]*domain.Task),
	}
}

// Create stores a new task in the in-memory repository.
// Returns domain.ErrTaskExists if a task with the same ID already exists.
func (r *MemoryTaskRepository) Create(_ context.Context, task *domain.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.tasks[task.ID]; exists {
		return domain.ErrTaskExists
	}

	r.tasks[task.ID] = task
	return nil
}

// GetByID retrieves a task by its unique identifier from the in-memory repository.
// Returns a copy of the task to prevent external modifications to the stored data.
// Returns domain.ErrTaskNotFound if no task exists with the given ID.
func (r *MemoryTaskRepository) GetByID(_ context.Context, id string) (*domain.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	task, exists := r.tasks[id]
	if !exists {
		return nil, domain.ErrTaskNotFound
	}

	// Return a copy to prevent external modifications
	taskCopy := *task
	return &taskCopy, nil
}

// GetAll retrieves all tasks from the in-memory repository, optionally filtered by status.
// If status is empty, returns all tasks regardless of their status.
// Returns copies of tasks to prevent external modifications to the stored data.
func (r *MemoryTaskRepository) GetAll(_ context.Context, status string) ([]*domain.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tasks := make([]*domain.Task, 0)
	for _, task := range r.tasks {
		if status == "" || string(task.Status) == status {
			// Create a copy to prevent external modifications
			taskCopy := *task
			tasks = append(tasks, &taskCopy)
		}
	}

	return tasks, nil
}

// Update modifies an existing task in the in-memory repository.
// Returns domain.ErrTaskNotFound if no task exists with the given ID.
func (r *MemoryTaskRepository) Update(_ context.Context, task *domain.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.tasks[task.ID]; !exists {
		return domain.ErrTaskNotFound
	}

	r.tasks[task.ID] = task
	return nil
}

// Delete removes a task from the in-memory repository by its ID.
// Returns domain.ErrTaskNotFound if no task exists with the given ID.
func (r *MemoryTaskRepository) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.tasks[id]; !exists {
		return domain.ErrTaskNotFound
	}

	delete(r.tasks, id)
	return nil
}
