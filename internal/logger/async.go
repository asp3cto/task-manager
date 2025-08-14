// Package logger provides an asynchronous logging system with JSON output.
// It features a single goroutine worker and configurable buffer size for high-performance logging.
package logger

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"os"
	"sync"
	"time"
)

var (
	_ Logger = (*AsyncLogger)(nil)
)

// LogEntry represents a single log entry that will be processed asynchronously.
// It contains all the information needed to generate a JSON log line.
type LogEntry struct {
	// Level is the log level (Debug, Info, Warn, Error)
	Level slog.Level
	// Message is the main log message
	Message string
	// Time is when the log entry was created
	Time time.Time
	// Attrs contains structured attributes to be included in the log output
	Attrs []slog.Attr
}

// AsyncLogger provides asynchronous logging with JSON output format.
// It uses a single background goroutine to process log entries from a buffered channel,
// ensuring non-blocking log operations in the calling goroutines.
type AsyncLogger struct {
	// ch is the buffered channel for log entries
	ch chan LogEntry
	// output is where log entries are written (e.g., os.Stdout, file)
	output io.Writer
	// level is the minimum log level to process
	level slog.Level
	// wg ensures graceful shutdown waits for worker completion
	wg sync.WaitGroup
}

// New creates a new AsyncLogger instance with the specified configuration.
// It starts a background worker goroutine immediately upon creation.
//
// Parameters:
//   - output: Writer where log entries will be written (uses os.Stdout if nil)
//   - level: Minimum log level to process (Debug, Info, Warn, Error)
//   - bufSize: Buffer size for the log entry channel
//
// Returns a fully initialized AsyncLogger ready for use.
// Remember to call Close() when done to ensure graceful shutdown.
func New(output io.Writer, level slog.Level, bufSize int) *AsyncLogger {
	if output == nil {
		output = os.Stdout
	}

	logger := &AsyncLogger{
		ch:     make(chan LogEntry, bufSize),
		output: output,
		level:  level,
	}

	return logger
}

// Start initializes and launches the background worker goroutine.
func (l *AsyncLogger) Start(ctx context.Context) {
	l.wg.Add(1)
	go l.worker(ctx)
}

// worker is the background goroutine that processes log entries.
// It continuously reads from the log channel and writes entries to the output.
// During shutdown, it processes all remaining entries before exiting.
func (l *AsyncLogger) worker(ctx context.Context) {
	defer l.wg.Done()

	for {
		select {
		case entry := <-l.ch:
			l.writeEntry(entry)
		case <-ctx.Done():
			for len(l.ch) > 0 {
				entry := <-l.ch
				l.writeEntry(entry)
			}

			return
		}
	}
}

// writeEntry formats and writes a single log entry as JSON.
// It filters entries based on the configured log level and marshals
// the entry data into JSON format with a newline terminator.
func (l *AsyncLogger) writeEntry(entry LogEntry) {
	if entry.Level < l.level {
		return
	}

	logData := map[string]interface{}{
		"time":    entry.Time.Format(time.RFC3339),
		"level":   entry.Level.String(),
		"message": entry.Message,
	}

	for _, attr := range entry.Attrs {
		logData[attr.Key] = attr.Value.Any()
	}

	jsonData, err := json.Marshal(logData)
	if err != nil {
		return
	}

	jsonData = append(jsonData, '\n')
	_, _ = l.output.Write(jsonData)
}

// log is the internal method that creates and queues log entries.
// if the context is done, it returns immediately.
func (l *AsyncLogger) log(ctx context.Context, level slog.Level, msg string, attrs ...slog.Attr) {
	if level < l.level {
		return
	}

	entry := LogEntry{
		Level:   level,
		Message: msg,
		Time:    time.Now(),
		Attrs:   attrs,
	}

	select {
	case l.ch <- entry:
	case <-ctx.Done():
		return
	}
}

// Debug logs a debug-level message with optional structured attributes.
func (l *AsyncLogger) Debug(ctx context.Context, msg string, attrs ...slog.Attr) {
	l.log(ctx, slog.LevelDebug, msg, attrs...)
}

// Info logs an info-level message with optional structured attributes.
func (l *AsyncLogger) Info(ctx context.Context, msg string, attrs ...slog.Attr) {
	l.log(ctx, slog.LevelInfo, msg, attrs...)
}

// Warn logs a warning-level message with optional structured attributes.
func (l *AsyncLogger) Warn(ctx context.Context, msg string, attrs ...slog.Attr) {
	l.log(ctx, slog.LevelWarn, msg, attrs...)
}

// Error logs an error-level message with optional structured attributes.
func (l *AsyncLogger) Error(ctx context.Context, msg string, attrs ...slog.Attr) {
	l.log(ctx, slog.LevelError, msg, attrs...)
}

// Close performs graceful shutdown of the async logger.
func (l *AsyncLogger) Close() {
	close(l.ch)
	l.wg.Wait()
}
