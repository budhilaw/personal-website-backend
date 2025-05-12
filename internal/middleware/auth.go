package middleware

import (
	"strings"

	"github.com/budhilaw/personal-website-backend/config"
	"github.com/budhilaw/personal-website-backend/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

var (
	jwtManager *JWTManager
)

// JWTClaims represents the JWT claims
type JWTClaims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	IsAdmin  bool   `json:"is_admin"`
	jwt.RegisteredClaims
}

// InitJWTManager initializes the JWT manager
func InitJWTManager(cfg config.Config) {
	jwtManager = NewJWTManager(cfg)
	logger.Info("JWT Manager initialized with secret rotation")
}

// GenerateToken generates a new JWT token
func GenerateToken(userID string, username string, isAdmin bool, cfg config.Config) (string, error) {
	if jwtManager == nil {
		InitJWTManager(cfg)
	}
	return jwtManager.GenerateToken(userID, username, isAdmin)
}

// GenerateRefreshToken generates a new refresh token
func GenerateRefreshToken(userID string, username string, isAdmin bool, cfg config.Config) (string, error) {
	if jwtManager == nil {
		InitJWTManager(cfg)
	}
	return jwtManager.GenerateRefreshToken(userID, username, isAdmin)
}

// Protected middleware for protecting routes
func Protected(cfg config.Config) fiber.Handler {
	// Ensure JWT Manager is initialized
	if jwtManager == nil {
		InitJWTManager(cfg)
	}

	return func(c *fiber.Ctx) error {
		// Get authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}

		// Check if the Authorization header starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid authorization format",
			})
		}

		// Extract the token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Verify token using JWT Manager
		claims, err := jwtManager.VerifyToken(tokenString)
		if err != nil {
			logger.Warn("Invalid JWT token", zap.Error(err), zap.String("token", tokenString[:10]+"..."))
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired token",
			})
		}

		// Set claims in context
		c.Locals("user_id", claims.UserID)
		c.Locals("username", claims.Username)
		c.Locals("is_admin", claims.IsAdmin)

		return c.Next()
	}
}

// AdminOnly middleware for admin-only routes
func AdminOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		isAdmin, ok := c.Locals("is_admin").(bool)
		if !ok || !isAdmin {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Access denied",
			})
		}

		return c.Next()
	}
}
