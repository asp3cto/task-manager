package logger

import (
	"io"
	"log/slog"
	"os"
	"strconv"
)

const defaultBufferSize = 100

// NewFromEnv creates a new AsyncLogger configured from environment variables.
// This is the recommended way to create a logger for most applications.
//
// Environment variables used:
//   - LOG_BUFFER_SIZE: Buffer size for the log channel (default: 100)
//   - LOG_LEVEL: Minimum log level - DEBUG, INFO, WARN, ERROR (default: INFO)
//
// Parameters:
//   - output: Writer where log entries will be written (uses os.Stdout if nil)
//
// Returns a configured AsyncLogger ready for use.
func NewFromEnv(output io.Writer) *AsyncLogger {
	bufSize := getLogBufferSize()
	level := getLogLevel()

	return New(output, level, bufSize)
}

// getLogBufferSize reads the LOG_BUFFER_SIZE environment variable
// and returns the buffer size for the log channel.
//
// Returns 100 if the environment variable is not set, invalid, or <= 0.
// The buffer size determines how many log entries can be queued before
// log calls become blocking or entries are dropped.
func getLogBufferSize() int {
	bufSizeStr := os.Getenv("LOG_BUFFER_SIZE")
	if bufSizeStr == "" {
		return defaultBufferSize
	}

	bufSize, err := strconv.Atoi(bufSizeStr)
	if err != nil || bufSize <= 0 {
		panic("LOG_BUFFER_SIZE must be a positive integer, got: " + bufSizeStr)
	}

	return bufSize
}

// getLogLevel reads the LOG_LEVEL environment variable
// and returns the corresponding slog.Level.
//
// Supported values (case-sensitive):
//   - DEBUG: Most verbose, includes all log levels
//   - INFO:  General information (default)
//   - WARN:  Warning conditions
//   - ERROR: Error conditions only
//
// Returns slog.LevelInfo if the environment variable is not set
// or contains an unrecognized value.
func getLogLevel() slog.Level {
	levelStr := os.Getenv("LOG_LEVEL")
	switch levelStr {
	case "DEBUG":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelInfo
	case "WARN":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
