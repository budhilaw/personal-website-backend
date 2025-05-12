package middleware

import (
	"time"

	"github.com/budhilaw/personal-website-backend/internal/logger"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ZapLogger is a middleware that logs HTTP requests using zap
func ZapLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Handle request
		err := c.Next()

		// Get latency
		latency := time.Since(start)

		// Determine status for color coding
		status := c.Response().StatusCode()

		// Prepare fields
		fields := []zapcore.Field{
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.Int("status", status),
			zap.String("ip", c.IP()),
			zap.String("user-agent", c.Get("User-Agent")),
			zap.Duration("latency", latency),
			zap.String("reqId", c.GetRespHeader("X-Request-Id")),
		}

		// Log based on status code
		switch {
		case status >= 500:
			logger.Error("Server error", fields...)
		case status >= 400:
			logger.Warn("Client error", fields...)
		case status >= 300:
			logger.Info("Redirection", fields...)
		default:
			logger.Info("Success", fields...)
		}

		return err
	}
}
