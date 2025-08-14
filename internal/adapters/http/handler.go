package http

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/asp3cto/task-manager/internal/domain"
	"github.com/asp3cto/task-manager/internal/logger"
	"github.com/asp3cto/task-manager/internal/ports"
)

// TaskHandler handles HTTP requests for task-related operations.
// It translates HTTP requests into service calls and formats responses.
type TaskHandler struct {
	service ports.TaskService
	logger  logger.Logger
}

// NewTaskHandler creates a new HTTP handler for task operations.
func NewTaskHandler(service ports.TaskService, logger logger.Logger) *TaskHandler {
	return &TaskHandler{
		service: service,
		logger:  logger,
	}
}

// CreateTaskRequest represents the JSON payload for creating a new task.
type CreateTaskRequest struct {
	// Title is the short name or summary of the task
	Title string `json:"title"`
	// Description provides detailed information about the task
	Description string `json:"description"`
}

// UpdateTaskStatusRequest represents the JSON payload for updating a task's status.
type UpdateTaskStatusRequest struct {
	// Status is the new status to set for the task
	Status domain.TaskStatus `json:"status"`
}

// ErrorResponse represents the JSON format for error responses.
type ErrorResponse struct {
	// Error contains the error message to return to the client
	Error string `json:"error"`
}

// HTTP-specific error messages for consistent API responses.
var (
	// ErrInternalServerError is returned when an unexpected server error occurs.
	ErrInternalServerError = errors.New("internal server error")
	// ErrInvalidStatus is returned when an invalid status parameter is provided.
	ErrInvalidStatus = errors.New("invalid status parameter")
	// ErrTaskNotFound is returned when a requested task does not exist.
	ErrTaskNotFound = errors.New("task not found")
	// ErrInvalidRequestFormat is returned when the request JSON cannot be parsed.
	ErrInvalidRequestFormat = errors.New("invalid request format")
	// ErrTitleRequired is returned when attempting to create a task without a title.
	ErrTitleRequired = errors.New("title is required")
)

// GetTasks handles GET /tasks requests to retrieve all tasks.
// Supports optional status query parameter for filtering tasks by status.
// Returns a JSON array of tasks or an error response.
func (h *TaskHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	status := r.URL.Query().Get("status")
	h.logger.Info(ctx, "getting tasks", slog.String("status_filter", status))

	if status != "" && !domain.IsValidStatus(status) {
		h.logger.Warn(ctx, "invalid status parameter", slog.String("status", status))
		h.writeError(w, ErrInvalidStatus, http.StatusBadRequest)
		return
	}

	tasks, err := h.service.GetAllTasks(r.Context(), status)
	if err != nil {
		h.logger.Error(ctx, "failed to get tasks", slog.String("error", err.Error()))
		h.writeError(w, ErrInternalServerError, http.StatusInternalServerError)
		return
	}

	h.writeJSONResponse(w, http.StatusOK, tasks)
}

// GetTask handles GET /tasks/{id} requests to retrieve a specific task by ID.
// Returns the task as JSON or a 404 error if the task doesn't exist.
func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	taskID := r.PathValue("id")
	h.logger.Info(ctx, "getting task by ID", slog.String("task_id", taskID))

	if taskID == "" {
		h.logger.Warn(ctx, "empty task ID in request")
		h.writeError(w, ErrTaskNotFound, http.StatusNotFound)
		return
	}

	task, err := h.service.GetTaskByID(r.Context(), taskID)
	if err != nil {
		if errors.Is(err, domain.ErrTaskNotFound) {
			h.logger.Warn(ctx, "task not found", slog.String("task_id", taskID))
			h.writeError(w, ErrTaskNotFound, http.StatusNotFound)
		} else {
			h.logger.Error(ctx, "failed to get task", slog.String("task_id", taskID), slog.String("error", err.Error()))
			h.writeError(w, ErrInternalServerError, http.StatusInternalServerError)
		}

		return
	}

	h.writeJSONResponse(w, http.StatusOK, task)
}

// CreateTask handles POST /tasks requests to create a new task.
// Expects a JSON payload with title and description fields.
// Returns the created task with a generated ID and pending status.
func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	h.logger.Info(ctx, "creating new task")

	var req CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn(ctx, "invalid request format", slog.String("error", err.Error()))
		h.writeError(w, ErrInvalidRequestFormat, http.StatusBadRequest)
		return
	}

	h.logger.Debug(ctx, "parsed create task request", slog.String("title", req.Title))
	task, err := h.service.CreateTask(r.Context(), req.Title, req.Description)
	if err != nil {
		if errors.Is(err, domain.ErrEmptyTitle) {
			h.logger.Warn(ctx, "task creation failed: empty title")
			h.writeError(w, ErrTitleRequired, http.StatusBadRequest)
		} else {
			h.logger.Error(ctx, "failed to create task", slog.String("error", err.Error()))
			h.writeError(w, ErrInternalServerError, http.StatusInternalServerError)
		}
		return
	}

	h.writeJSONResponse(w, http.StatusCreated, task)
}

// writeError writes an error response in JSON format with the specified status code.
// The err parameter can be a string, error, or any other type (converted to string).
func (h *TaskHandler) writeError(w http.ResponseWriter, err any, statusCode int) {
	w.WriteHeader(statusCode)

	var errorMsg string
	switch v := err.(type) {
	case string:
		errorMsg = v
	case error:
		errorMsg = v.Error()
	default:
		errorMsg = ErrInternalServerError.Error()
	}

	_ = json.NewEncoder(w).Encode(ErrorResponse{Error: errorMsg})
}

// writeJSONResponse writes a successful JSON response with the specified status code.
// Sets the appropriate Content-Type header and encodes the data as JSON.
func (h *TaskHandler) writeJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.writeError(w, ErrInternalServerError, http.StatusInternalServerError)
	}
}
