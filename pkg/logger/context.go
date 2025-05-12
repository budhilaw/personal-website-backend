package logger

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ContextKey is the key type for context values
type ContextKey string

// Context keys
const (
	ContextLoggerKey ContextKey = "logger"
)

// RequestLogger returns zap fields for a request
func RequestLogger(userID, action, resource string) []zapcore.Field {
	requestID := fmt.Sprintf("%d", time.Now().UnixNano())
	fields := []zapcore.Field{
		zap.String("request_id", requestID),
	}
	if userID != "" {
		fields = append(fields, zap.String("user_id", userID))
	}
	if action != "" {
		fields = append(fields, zap.String("action", action))
	}
	if resource != "" {
		fields = append(fields, zap.String("resource", resource))
	}
	return fields
}

// WithContextFields adds fields to the logger
func WithContextFields(ctx context.Context, fields []zapcore.Field) context.Context {
	// Get existing logger if it exists
	existingLogger, ok := ctx.Value(ContextLoggerKey).(*zap.Logger)
	var newLogger *zap.Logger

	if ok && existingLogger != nil {
		// Add fields to existing logger
		newLogger = existingLogger.With(fields...)
	} else {
		// Create new logger with fields
		newLogger = Log.With(fields...)
	}

	// Create new context with logger
	return context.WithValue(ctx, ContextLoggerKey, newLogger)
}

// DebugContext logs a debug message with context
func DebugContext(ctx context.Context, message string, fields ...zapcore.Field) {
	logger := getLoggerFromContext(ctx)
	logger.Debug(message, fields...)
}

// InfoContext logs an info message with context
func InfoContext(ctx context.Context, message string, fields ...zapcore.Field) {
	logger := getLoggerFromContext(ctx)
	logger.Info(message, fields...)
}

// WarnContext logs a warning message with context
func WarnContext(ctx context.Context, message string, fields ...zapcore.Field) {
	logger := getLoggerFromContext(ctx)
	logger.Warn(message, fields...)
}

// ErrorContext logs an error message with context
func ErrorContext(ctx context.Context, message string, fields ...zapcore.Field) {
	logger := getLoggerFromContext(ctx)
	logger.Error(message, fields...)
}

// FatalContext logs a fatal message with context and exits
func FatalContext(ctx context.Context, message string, fields ...zapcore.Field) {
	logger := getLoggerFromContext(ctx)
	logger.Fatal(message, fields...)
}

// getLoggerFromContext gets a logger from context
func getLoggerFromContext(ctx context.Context) *zap.Logger {
	logger, ok := ctx.Value(ContextLoggerKey).(*zap.Logger)
	if !ok || logger == nil {
		return Log
	}
	return logger
}
