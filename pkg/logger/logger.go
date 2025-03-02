package logger

import (
	"io"
	"log/slog"
	"os"
	"strings"
)

// Logger is a wrapper around slog.Logger
type Logger struct {
	*slog.Logger
}

// NewLogger creates a new logger with the given level and format
func NewLogger(level, format string) *Logger {
	var logLevel slog.Level
	switch strings.ToLower(level) {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	var handler slog.Handler
	switch strings.ToLower(format) {
	case "json":
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: logLevel,
		})
	default:
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: logLevel,
		})
	}

	return &Logger{
		Logger: slog.New(handler),
	}
}

// WithOutput returns a new logger with the given output
func (l *Logger) WithOutput(w io.Writer) *Logger {
	var handler slog.Handler

	// Create a new handler with the same level as the current logger
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo, // Default to info level
	}

	// Determine if the current handler is JSON or text
	if _, ok := l.Logger.Handler().(*slog.JSONHandler); ok {
		handler = slog.NewJSONHandler(w, opts)
	} else {
		handler = slog.NewTextHandler(w, opts)
	}

	return &Logger{
		Logger: slog.New(handler),
	}
}

// WithField returns a new logger with the given field
func (l *Logger) WithField(key string, value interface{}) *Logger {
	return &Logger{
		Logger: l.Logger.With(key, value),
	}
}

// WithFields returns a new logger with the given fields
func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	attrs := make([]any, 0, len(fields)*2)
	for k, v := range fields {
		attrs = append(attrs, k, v)
	}
	return &Logger{
		Logger: l.Logger.With(attrs...),
	}
}

// Fatal logs a message at error level and then exits the program
func (l *Logger) Fatal(msg string, args ...interface{}) {
	l.Error(msg, args...)
	os.Exit(1)
}
