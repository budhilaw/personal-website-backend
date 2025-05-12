package logger

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ctxLoggerKey struct{}

// WithContext adds logger fields to a context
func WithContext(ctx context.Context, fields ...zapcore.Field) context.Context {
	return context.WithValue(ctx, ctxLoggerKey{}, fields)
}

// WithContextFields adds an array of logger fields to a context
func WithContextFields(ctx context.Context, fields []zapcore.Field) context.Context {
	return context.WithValue(ctx, ctxLoggerKey{}, fields)
}

// FromContext retrieves logger fields from a context
func FromContext(ctx context.Context) []zapcore.Field {
	fields, ok := ctx.Value(ctxLoggerKey{}).([]zapcore.Field)
	if !ok {
		return []zapcore.Field{}
	}
	return fields
}

// InfoContext logs a message at info level with context fields
func InfoContext(ctx context.Context, message string, fields ...zapcore.Field) {
	ctxFields := FromContext(ctx)
	if len(ctxFields) > 0 {
		fields = append(fields, ctxFields...)
	}
	Info(message, fields...)
}

// DebugContext logs a message at debug level with context fields
func DebugContext(ctx context.Context, message string, fields ...zapcore.Field) {
	ctxFields := FromContext(ctx)
	if len(ctxFields) > 0 {
		fields = append(fields, ctxFields...)
	}
	Debug(message, fields...)
}

// WarnContext logs a message at warn level with context fields
func WarnContext(ctx context.Context, message string, fields ...zapcore.Field) {
	ctxFields := FromContext(ctx)
	if len(ctxFields) > 0 {
		fields = append(fields, ctxFields...)
	}
	Warn(message, fields...)
}

// ErrorContext logs a message at error level with context fields
func ErrorContext(ctx context.Context, message string, fields ...zapcore.Field) {
	ctxFields := FromContext(ctx)
	if len(ctxFields) > 0 {
		fields = append(fields, ctxFields...)
	}
	Error(message, fields...)
}

// RequestLogger creates fields for logging HTTP request information
func RequestLogger(userID, method, path string) []zapcore.Field {
	return []zapcore.Field{
		zap.String("user_id", userID),
		zap.String("method", method),
		zap.String("path", path),
	}
}

// DatabaseLogger creates fields for logging database operations
func DatabaseLogger(operation, entity string) []zapcore.Field {
	return []zapcore.Field{
		zap.String("db_operation", operation),
		zap.String("entity", entity),
	}
}
