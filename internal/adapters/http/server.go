// Package http provides HTTP transport layer for the task management API.
// It implements the REST API endpoints and HTTP server functionality,
// serving as an adapter between HTTP requests and the core business logic.
package http

import (
	"context"
	"net/http"
	"time"

	"github.com/asp3cto/task-manager/internal/logger"
	"github.com/asp3cto/task-manager/internal/ports"
)

// Server wraps an HTTP server with task management capabilities.
type Server struct {
	http *http.Server
	// handler contains the HTTP request handlers for task operations
	handler *TaskHandler
}

// readHeaderTimeout defines the maximum time allowed to read request headers.
// This helps prevent Slowloris attacks by limiting the time spent reading headers.
const readHeaderTimeout = 2 * time.Second

// NewServer creates a new HTTP server instance with task management endpoints.
func NewServer(addr string, service ports.TaskService, logger logger.Logger) *Server {
	handler := NewTaskHandler(service, logger)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /tasks", handler.GetTasks)
	mux.HandleFunc("GET /tasks/{id}", handler.GetTask)
	mux.HandleFunc("POST /tasks", handler.CreateTask)

	httpServer := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	return &Server{
		http:    httpServer,
		handler: handler,
	}
}

// ListenAndServe starts the HTTP server and begins accepting connections.
// This method blocks until the server is shut down or an error occurs.
func (s *Server) ListenAndServe() error {
	return s.http.ListenAndServe()
}

// Shutdown gracefully shuts down the HTTP server without interrupting active connections.
// It waits for active connections to finish or for the context to be cancelled.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.http.Shutdown(ctx)
}

// Addr returns the network address the server is configured to listen on.
func (s *Server) Addr() string {
	return s.http.Addr
}
