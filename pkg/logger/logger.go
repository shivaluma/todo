package logger

import (
	"io"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// RequestIDKey is the key used for request ID in the logger
const RequestIDKey = "request_id"

// Logger is a wrapper around zap.Logger
type Logger struct {
	*zap.Logger
	sugar *zap.SugaredLogger
}

// NewLogger creates a new logger with the given level and format
func NewLogger(level, format string) *Logger {
	// Configure log level
	var logLevel zapcore.Level
	switch strings.ToLower(level) {
	case "debug":
		logLevel = zapcore.DebugLevel
	case "info":
		logLevel = zapcore.InfoLevel
	case "warn":
		logLevel = zapcore.WarnLevel
	case "error":
		logLevel = zapcore.ErrorLevel
	default:
		logLevel = zapcore.InfoLevel
	}

	// Configure encoder based on format
	var encoder zapcore.Encoder
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	switch strings.ToLower(format) {
	case "json":
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	default:
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	// Create core
	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(os.Stdout),
		logLevel,
	)

	// Create logger
	zapLogger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return &Logger{
		Logger: zapLogger,
		sugar:  zapLogger.Sugar(),
	}
}

// WithOutput returns a new logger with the given output
func (l *Logger) WithOutput(w io.Writer) *Logger {
	// Create a new core with the same encoder and level but different output
	core := zapcore.NewCore(
		l.getEncoder(),
		zapcore.AddSync(w),
		l.getLevel(),
	)

	// Create a new logger with the new core
	zapLogger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return &Logger{
		Logger: zapLogger,
		sugar:  zapLogger.Sugar(),
	}
}

// getEncoder returns the encoder used by the logger
func (l *Logger) getEncoder() zapcore.Encoder {
	// Default to console encoder if we can't determine
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	
	// Try to determine if JSON or console encoder is being used
	// This is a simplification; in a real implementation, you might want to store this information
	return zapcore.NewConsoleEncoder(encoderConfig)
}

// getLevel returns the level used by the logger
func (l *Logger) getLevel() zapcore.Level {
	// Default to info level
	return zapcore.InfoLevel
}

// WithField returns a new logger with the given field
func (l *Logger) WithField(key string, value interface{}) *Logger {
	zapLogger := l.Logger.With(zap.Any(key, value))
	return &Logger{
		Logger: zapLogger,
		sugar:  zapLogger.Sugar(),
	}
}

// WithRequestID returns a new logger with the request ID field
func (l *Logger) WithRequestID(requestID string) *Logger {
	if requestID == "" {
		return l
	}
	return l.WithField(RequestIDKey, requestID)
}

// FromContext creates a logger from an Echo context with the request ID
func FromContext(c echo.Context) *Logger {
	requestID := c.Response().Header().Get(echo.HeaderXRequestID)
	if requestID == "" {
		// Try to get from request header if not in response
		requestID = c.Request().Header.Get(echo.HeaderXRequestID)
	}
	
	// Get the base logger - this assumes a logger has been set in the context
	// If not, create a new one
	var logger *Logger
	if l, ok := c.Get("logger").(*Logger); ok {
		logger = l
	} else {
		logger = NewLogger("info", "console")
	}
	
	return logger.WithRequestID(requestID)
}

// WithFields returns a new logger with the given fields
func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	zapFields := make([]zap.Field, 0, len(fields))
	for k, v := range fields {
		zapFields = append(zapFields, zap.Any(k, v))
	}
	
	zapLogger := l.Logger.With(zapFields...)
	return &Logger{
		Logger: zapLogger,
		sugar:  zapLogger.Sugar(),
	}
}

// Debug logs a message at debug level
func (l *Logger) Debug(msg string, args ...interface{}) {
	l.sugar.Debugw(msg, args...)
}

// Info logs a message at info level
func (l *Logger) Info(msg string, args ...interface{}) {
	l.sugar.Infow(msg, args...)
}

// Warn logs a message at warn level
func (l *Logger) Warn(msg string, args ...interface{}) {
	l.sugar.Warnw(msg, args...)
}

// Error logs a message at error level
func (l *Logger) Error(msg string, args ...interface{}) {
	l.sugar.Errorw(msg, args...)
}

// Fatal logs a message at error level and then exits the program
func (l *Logger) Fatal(msg string, args ...interface{}) {
	l.sugar.Errorw(msg, args...)
	os.Exit(1)
}
