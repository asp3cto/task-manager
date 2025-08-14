package logger

import (
	"context"
	"log/slog"
)

// Logger defines the interface for structured logging operations.
type Logger interface {
	// Debug logs a debug-level message with optional structured attributes.
	Debug(ctx context.Context, msg string, attrs ...slog.Attr)

	// Info logs an info-level message with optional structured attributes.
	Info(ctx context.Context, msg string, attrs ...slog.Attr)

	// Warn logs a warning-level message with optional structured attributes.
	Warn(ctx context.Context, msg string, attrs ...slog.Attr)

	// Error logs an error-level message with optional structured attributes.
	Error(ctx context.Context, msg string, attrs ...slog.Attr)
}
