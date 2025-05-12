package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

// Security middleware for adding security headers and protections
func Security(frontendURL string) fiber.Handler {
	// Use cors middleware
	return cors.New(cors.Config{
		AllowOrigins:     frontendURL,
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
		MaxAge:           86400, // 24 hours
	})
}

// Helmet middleware for adding secure headers
func Helmet() fiber.Handler {
	return helmet.New(helmet.Config{
		ContentSecurityPolicy: "default-src 'self'; img-src 'self' data:; style-src 'self' 'unsafe-inline'; script-src 'self' 'unsafe-inline'",
		XSSProtection:         "1; mode=block",
		ContentTypeNosniff:    "nosniff",
		ReferrerPolicy:        "no-referrer-when-downgrade",
	})
}

// RateLimiter middleware for rate limiting
func RateLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        100,              // max 100 requests
		Expiration: 60 * 1000000000, // 1 minute
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP() // use IP as key
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Too many requests",
			})
		},
	})
} 