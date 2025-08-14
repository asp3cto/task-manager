package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	httpAdapter "github.com/asp3cto/task-manager/internal/adapters/http"
	"github.com/asp3cto/task-manager/internal/adapters/repository"
	"github.com/asp3cto/task-manager/internal/core/service"
)

// shutdownDelay defines the maximum time to wait for graceful shutdown.
// The server will force shutdown if active connections don't close within this time.
const shutdownDelay = 30 * time.Second

func main() {
	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":8080"
	}

	repo := repository.NewMemoryTaskRepository()
	taskService := service.NewTaskService(repo)
	server := httpAdapter.NewServer(addr, taskService)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Printf("server starting on %s", server.Addr())
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("failed to start server: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("received shutdown signal, shutting down gracefully")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownDelay)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Println("server forced to shutdown: ", err)
	}

	log.Println("server exited")
}
