package logger

import (
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	// Log is the global logger instance
	Log  *zap.Logger
	once sync.Once
)

// InitLogger initializes the global logger
func InitLogger(isProduction bool) *zap.Logger {
	once.Do(func() {
		var config zap.Config

		if isProduction {
			config = zap.NewProductionConfig()
			config.EncoderConfig.TimeKey = "timestamp"
			config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		} else {
			config = zap.NewDevelopmentConfig()
			config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
			config.EncoderConfig.TimeKey = "timestamp"
			config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		}

		var err error
		Log, err = config.Build(zap.AddCaller(), zap.AddCallerSkip(1))
		if err != nil {
			panic(err)
		}
	})

	return Log
}

// Info logs a message at info level
func Info(message string, fields ...zapcore.Field) {
	Log.Info(message, fields...)
}

// Debug logs a message at debug level
func Debug(message string, fields ...zapcore.Field) {
	Log.Debug(message, fields...)
}

// Warn logs a message at warn level
func Warn(message string, fields ...zapcore.Field) {
	Log.Warn(message, fields...)
}

// Error logs a message at error level
func Error(message string, fields ...zapcore.Field) {
	Log.Error(message, fields...)
}

// Fatal logs a message at fatal level and then calls os.Exit(1)
func Fatal(message string, fields ...zapcore.Field) {
	Log.Fatal(message, fields...)
	os.Exit(1)
}

// With returns a logger with the provided fields attached
func With(fields ...zapcore.Field) *zap.Logger {
	return Log.With(fields...)
}

// Sync flushes any buffered log entries
func Sync() error {
	return Log.Sync()
}
