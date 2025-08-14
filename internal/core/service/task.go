// Package service implements the core business logic for task management.
// It provides the concrete implementation of the TaskService interface,
// orchestrating domain objects and repository operations.
package service

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"log/slog"

	"github.com/asp3cto/task-manager/internal/domain"
	"github.com/asp3cto/task-manager/internal/logger"
	"github.com/asp3cto/task-manager/internal/ports"
)

var _ ports.TaskService = (*TaskService)(nil)

// TaskService implements the core business logic for task operations.
// It orchestrates domain entities and repository interactions while
// enforcing business rules and validation.
type TaskService struct {
	repo   ports.TaskRepository
	logger logger.Logger
}

// NewTaskService creates a new instance of TaskService with the provided repository.
// The repository is used for all data persistence operations.
func NewTaskService(repo ports.TaskRepository, logger logger.Logger) *TaskService {
	return &TaskService{
		repo:   repo,
		logger: logger,
	}
}

// CreateTask creates a new task with the given title and description.
// It validates the input, generates a unique ID, and stores the task.
// Returns domain.ErrEmptyTitle if the title is empty.
func (s *TaskService) CreateTask(ctx context.Context, title, description string) (*domain.Task, error) {
	s.logger.Debug(ctx, "creating task", slog.String("title", title))

	if title == "" {
		s.logger.Warn(ctx, "task creation failed: empty title")
		return nil, domain.ErrEmptyTitle
	}

	id, err := generateID()
	if err != nil {
		s.logger.Error(ctx, "failed to generate ID", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to generate ID: %w", err)
	}

	task := domain.NewTask(id, title, description)

	if err := s.repo.Create(ctx, task); err != nil {
		s.logger.Error(
			ctx,
			"failed to create task in repository",
			slog.String("task_id", id), slog.String("error", err.Error()),
		)
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	s.logger.Info(
		ctx,
		"task created successfully",
		slog.String("task_id", id), slog.String("title", title),
	)

	return task, nil
}

// GetTaskByID retrieves a task by its unique identifier.
// Returns domain.ErrTaskNotFound if no task exists with the given ID.
func (s *TaskService) GetTaskByID(ctx context.Context, id string) (*domain.Task, error) {
	s.logger.Debug(ctx, "getting task by ID", slog.String("task_id", id))

	task, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrTaskNotFound) {
			s.logger.Debug(ctx, "task not found", slog.String("task_id", id))
			return nil, err
		}

		s.logger.Error(
			ctx,
			"failed to get task from repository",
			slog.String("task_id", id), slog.String("error", err.Error()),
		)

		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	s.logger.Debug(ctx, "task retrieved successfully", slog.String("task_id", id))
	return task, nil
}

// GetAllTasks retrieves all tasks, optionally filtered by status.
// If status is empty, returns all tasks regardless of their status.
func (s *TaskService) GetAllTasks(ctx context.Context, status string) ([]*domain.Task, error) {
	s.logger.Debug(ctx, "getting all tasks", slog.String("status_filter", status))

	tasks, err := s.repo.GetAll(ctx, status)
	if err != nil {
		s.logger.Error(ctx, "failed to get tasks from repository", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}

	s.logger.Debug(
		ctx,
		"tasks retrieved successfully",
		slog.Int("count", len(tasks)), slog.String("status_filter", status),
	)
	return tasks, nil
}

// UpdateTaskStatus changes the status of an existing task.
// It retrieves the task, updates its status using domain methods, and persists the change.
// Returns domain.ErrTaskNotFound if no task exists with the given ID.
func (s *TaskService) UpdateTaskStatus(ctx context.Context, id string, status domain.TaskStatus) (*domain.Task, error) {
	s.logger.Debug(
		ctx,
		"updating task status",
		slog.String("task_id", id), slog.String("new_status", string(status)),
	)

	task, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrTaskNotFound) {
			s.logger.Debug(ctx, "task not found for status update", slog.String("task_id", id))
			return nil, err
		}

		s.logger.Error(
			ctx,
			"failed to get task for status update",
			slog.String("task_id", id), slog.String("error", err.Error()),
		)
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	oldStatus := task.Status
	task.UpdateStatus(status)

	if err := s.repo.Update(ctx, task); err != nil {
		s.logger.Error(
			ctx,
			"failed to update task in repository",
			slog.String("task_id", id), slog.String("error", err.Error()),
		)

		return nil, fmt.Errorf("failed to update task: %w", err)
	}

	s.logger.Info(
		ctx, "task status updated successfully",
		slog.String("task_id", id),
		slog.String("old_status", string(oldStatus)),
		slog.String("new_status", string(status)),
	)
	return task, nil
}

// idLength defines the number of bytes used for generating task IDs.
const idLength = 16

// generateID creates a random ID for tasks.
// It generates a 16-byte random value and returns it as a hexadecimal string.
func generateID() (string, error) {
	bytes := make([]byte, idLength)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", bytes), nil
}
