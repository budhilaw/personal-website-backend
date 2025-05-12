package middleware

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"sync"
	"time"

	"github.com/budhilaw/personal-website-backend/config"
	"github.com/budhilaw/personal-website-backend/internal/logger"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

// JWTManager handles JWT token operations with secret rotation
type JWTManager struct {
	currentSecret    []byte
	previousSecret   []byte
	secretCreatedAt  time.Time
	rotationInterval time.Duration
	mutex            sync.RWMutex
	config           config.Config
}

// NewJWTManager creates a new JWT manager with secret rotation
func NewJWTManager(cfg config.Config) *JWTManager {
	manager := &JWTManager{
		currentSecret:    []byte(cfg.JWTSecret),
		previousSecret:   nil, // Initially no previous secret
		secretCreatedAt:  time.Now(),
		rotationInterval: time.Hour * 24 * 7, // Default 7 days, adjust as needed
		config:           cfg,
	}

	// Start secret rotation in background
	go manager.rotateSecretsPeriodically()

	return manager
}

// rotateSecretsPeriodically rotates JWT secrets at the specified interval
func (m *JWTManager) rotateSecretsPeriodically() {
	ticker := time.NewTicker(m.rotationInterval / 2) // Check at half the interval
	defer ticker.Stop()

	for range ticker.C {
		if time.Since(m.secretCreatedAt) >= m.rotationInterval {
			if err := m.rotateSecrets(); err != nil {
				logger.Error("Failed to rotate JWT secrets", zap.Error(err))
			}
		}
	}
}

// rotateSecrets generates a new secret and rotates the existing one
func (m *JWTManager) rotateSecrets() error {
	newSecret := make([]byte, 32) // 256-bit secret
	_, err := rand.Read(newSecret)
	if err != nil {
		return err
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.previousSecret = m.currentSecret
	m.currentSecret = newSecret
	m.secretCreatedAt = time.Now()

	// Log secret rotation (without exposing the secrets)
	logger.Info("JWT secrets rotated successfully",
		zap.Time("rotation_time", m.secretCreatedAt),
		zap.Time("next_rotation", m.secretCreatedAt.Add(m.rotationInterval)))

	return nil
}

// GenerateToken generates a new JWT token using the current secret
func (m *JWTManager) GenerateToken(userID string, username string, isAdmin bool) (string, error) {
	// Create token claims
	claims := JWTClaims{
		UserID:   userID,
		Username: username,
		IsAdmin:  isAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.config.JWTExpiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	// Create token with current secret
	m.mutex.RLock()
	secret := m.currentSecret
	m.mutex.RUnlock()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// GenerateRefreshToken generates a new refresh token
func (m *JWTManager) GenerateRefreshToken(userID string, username string, isAdmin bool) (string, error) {
	// Create token claims
	claims := JWTClaims{
		UserID:   userID,
		Username: username,
		IsAdmin:  isAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.config.JWTRefreshExpiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	// Use dedicated refresh secret from config
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(m.config.JWTRefreshSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// VerifyToken verifies a JWT token against current and previous secrets
func (m *JWTManager) VerifyToken(tokenString string) (*JWTClaims, error) {
	var lastError error

	// Try with current secret first
	m.mutex.RLock()
	currentSecret := m.currentSecret
	previousSecret := m.previousSecret
	m.mutex.RUnlock()

	// Try with current secret
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return currentSecret, nil
	})

	if err == nil && token.Valid {
		if claims, ok := token.Claims.(*JWTClaims); ok {
			return claims, nil
		}
		return nil, errors.New("invalid token claims")
	}

	lastError = err

	// If verification with current secret fails and we have a previous secret, try with that
	if previousSecret != nil {
		token, err = jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return previousSecret, nil
		})

		if err == nil && token.Valid {
			if claims, ok := token.Claims.(*JWTClaims); ok {
				// Log that we used previous secret (for monitoring purposes)
				logger.Info("JWT token verified with previous secret")
				return claims, nil
			}
			return nil, errors.New("invalid token claims")
		}

		lastError = err
	}

	return nil, lastError
}

// GetSecretInfo returns non-sensitive information about the JWT secrets
func (m *JWTManager) GetSecretInfo() map[string]interface{} {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return map[string]interface{}{
		"current_secret_created_at": m.secretCreatedAt,
		"next_rotation_at":          m.secretCreatedAt.Add(m.rotationInterval),
		"has_previous_secret":       m.previousSecret != nil,
		"rotation_interval_days":    m.rotationInterval / (time.Hour * 24),
	}
}

// Base64Secret returns a base64 encoded version of the current secret
// This is useful for sharing the secret with other services if needed
func (m *JWTManager) Base64Secret() string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return base64.StdEncoding.EncodeToString(m.currentSecret)
}
