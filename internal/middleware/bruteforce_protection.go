package middleware

import (
	"sync"
	"time"

	"github.com/budhilaw/personal-website-backend/internal/logger"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// Brute force protection configuration
const (
	maxFailedAttempts     = 5     // Maximum consecutive failed attempts before blocking
	initialBlockDuration  = 30    // Initial block duration in seconds
	blockMultiplier       = 2     // Multiplier for each subsequent block
	maxBlockDuration      = 86400 // Maximum block duration in seconds (24 hours)
	cleanupInterval       = 3600  // Cleanup interval in seconds
	failedAttemptsTimeout = 1800  // Clear failed attempts after this many seconds
)

// LoginAttempt tracks information about login attempts
type LoginAttempt struct {
	IP             string
	Username       string
	FailedAttempts int
	LastFailedAt   time.Time
	BlockedUntil   time.Time
}

// BruteForceProtector manages brute force protection
type BruteForceProtector struct {
	attempts   map[string]*LoginAttempt // Key is IP + username
	ipAttempts map[string]*LoginAttempt // Key is IP only (for IP-based blocking)
	mutex      sync.RWMutex
}

var (
	bruteForceProtector *BruteForceProtector
	once                sync.Once
)

// GetBruteForceProtector returns the singleton brute force protector
func GetBruteForceProtector() *BruteForceProtector {
	once.Do(func() {
		bruteForceProtector = &BruteForceProtector{
			attempts:   make(map[string]*LoginAttempt),
			ipAttempts: make(map[string]*LoginAttempt),
		}
		go bruteForceProtector.startCleanupTask()
	})
	return bruteForceProtector
}

// startCleanupTask periodically cleans up old login attempts
func (b *BruteForceProtector) startCleanupTask() {
	ticker := time.NewTicker(time.Second * cleanupInterval)
	defer ticker.Stop()

	for range ticker.C {
		b.cleanup()
	}
}

// cleanup removes expired login attempts
func (b *BruteForceProtector) cleanup() {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	now := time.Now()

	// Clean up account-specific attempts
	for k, attempt := range b.attempts {
		// If block has expired and no recent failed attempts, remove the entry
		if attempt.BlockedUntil.Before(now) &&
			attempt.LastFailedAt.Add(time.Second*failedAttemptsTimeout).Before(now) {
			delete(b.attempts, k)
		}
	}

	// Clean up IP-based attempts
	for k, attempt := range b.ipAttempts {
		// If block has expired and no recent failed attempts, remove the entry
		if attempt.BlockedUntil.Before(now) &&
			attempt.LastFailedAt.Add(time.Second*failedAttemptsTimeout).Before(now) {
			delete(b.ipAttempts, k)
		}
	}

	logger.Debug("Cleaned up brute force protection cache",
		zap.Int("remaining_attempts", len(b.attempts)),
		zap.Int("remaining_ip_attempts", len(b.ipAttempts)))
}

// IsBlocked checks if a login attempt is blocked
func (b *BruteForceProtector) IsBlocked(ip, username string) (bool, time.Time) {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	// Check account-specific block
	key := ip + ":" + username
	if attempt, exists := b.attempts[key]; exists && attempt.BlockedUntil.After(time.Now()) {
		return true, attempt.BlockedUntil
	}

	// Check IP-based block (regardless of username)
	if attempt, exists := b.ipAttempts[ip]; exists && attempt.BlockedUntil.After(time.Now()) {
		return true, attempt.BlockedUntil
	}

	return false, time.Time{}
}

// RecordFailedAttempt records a failed login attempt
func (b *BruteForceProtector) RecordFailedAttempt(ip, username string) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	now := time.Now()
	key := ip + ":" + username

	// Update account-specific attempts
	attempt, exists := b.attempts[key]
	if !exists {
		attempt = &LoginAttempt{
			IP:             ip,
			Username:       username,
			FailedAttempts: 0,
		}
		b.attempts[key] = attempt
	}

	attempt.FailedAttempts++
	attempt.LastFailedAt = now

	// Update IP-based attempts
	ipAttempt, exists := b.ipAttempts[ip]
	if !exists {
		ipAttempt = &LoginAttempt{
			IP:             ip,
			FailedAttempts: 0,
		}
		b.ipAttempts[ip] = ipAttempt
	}

	ipAttempt.FailedAttempts++
	ipAttempt.LastFailedAt = now

	// Check if account should be blocked
	if attempt.FailedAttempts >= maxFailedAttempts {
		blockDuration := time.Duration(initialBlockDuration) * time.Second

		// If already blocked, increase duration exponentially
		if attempt.BlockedUntil.After(now) {
			prevDuration := attempt.BlockedUntil.Sub(attempt.LastFailedAt)
			blockDuration = prevDuration * blockMultiplier

			// Cap at maximum duration
			if blockDuration > time.Duration(maxBlockDuration)*time.Second {
				blockDuration = time.Duration(maxBlockDuration) * time.Second
			}
		}

		attempt.BlockedUntil = now.Add(blockDuration)

		logger.Warn("Account temporarily blocked due to too many failed attempts",
			zap.String("username", username),
			zap.String("ip", ip),
			zap.Time("blocked_until", attempt.BlockedUntil),
			zap.Duration("block_duration", blockDuration))
	}

	// Check if IP should be blocked (more severe threshold)
	if ipAttempt.FailedAttempts >= maxFailedAttempts*2 {
		blockDuration := time.Duration(initialBlockDuration*2) * time.Second

		// If already blocked, increase duration exponentially
		if ipAttempt.BlockedUntil.After(now) {
			prevDuration := ipAttempt.BlockedUntil.Sub(ipAttempt.LastFailedAt)
			blockDuration = prevDuration * blockMultiplier

			// Cap at maximum duration
			if blockDuration > time.Duration(maxBlockDuration)*time.Second {
				blockDuration = time.Duration(maxBlockDuration) * time.Second
			}
		}

		ipAttempt.BlockedUntil = now.Add(blockDuration)

		logger.Warn("IP temporarily blocked due to too many failed attempts",
			zap.String("ip", ip),
			zap.Time("blocked_until", ipAttempt.BlockedUntil),
			zap.Duration("block_duration", blockDuration))
	}
}

// RecordSuccessfulAttempt resets failed login attempts counter
func (b *BruteForceProtector) RecordSuccessfulAttempt(ip, username string) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	// Reset account-specific attempts
	key := ip + ":" + username
	delete(b.attempts, key)

	// We don't reset IP-based attempts on success as one account success
	// shouldn't clear attempts on other accounts from the same IP
}

// BruteForceProtection middleware checks for brute force attacks
func BruteForceProtection() fiber.Handler {
	protector := GetBruteForceProtector()

	return func(c *fiber.Ctx) error {
		// Only apply to login endpoints
		if c.Path() == "/api/v1/auth/login" && c.Method() == "POST" {
			ip := c.IP()

			// Get username from body (we need to check before login attempt)
			body := make(map[string]interface{})
			if err := c.BodyParser(&body); err == nil {
				username, ok := body["username"].(string)
				if ok {
					// Check if this login attempt is blocked
					blocked, blockedUntil := protector.IsBlocked(ip, username)
					if blocked {
						// Calculate remaining block time
						remaining := blockedUntil.Sub(time.Now()).Seconds()

						logger.Warn("Blocked login attempt",
							zap.String("ip", ip),
							zap.String("username", username),
							zap.Float64("seconds_remaining", remaining))

						return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
							"error":             "Too many failed login attempts, please try again later",
							"seconds_remaining": int(remaining),
						})
					}
				}
			}
		}

		return c.Next()
	}
}

// TrackLoginAttempt middleware to track login success/failure
func TrackLoginAttempt() fiber.Handler {
	protector := GetBruteForceProtector()

	return func(c *fiber.Ctx) error {
		// Only apply to login endpoints
		if c.Path() == "/api/v1/auth/login" && c.Method() == "POST" {
			// Store original path, method and username for later
			ip := c.IP()
			path := c.Path()
			method := c.Method()

			var username string
			body := make(map[string]interface{})
			if parseErr := c.BodyParser(&body); parseErr == nil {
				if un, ok := body["username"].(string); ok {
					username = un
				}
			}

			// Process the request
			err := c.Next()

			// After the request is processed, check the status code
			if username != "" {
				statusCode := c.Response().StatusCode()

				if statusCode == fiber.StatusOK {
					// Successful login
					protector.RecordSuccessfulAttempt(ip, username)
					logger.Debug("Successful login attempt",
						zap.String("username", username),
						zap.String("ip", ip),
						zap.String("path", path),
						zap.String("method", method))
				} else if statusCode == fiber.StatusUnauthorized {
					// Failed login
					protector.RecordFailedAttempt(ip, username)
					logger.Warn("Failed login attempt",
						zap.String("username", username),
						zap.String("ip", ip),
						zap.String("path", path),
						zap.String("method", method))
				}
			}

			return err
		}

		return c.Next()
	}
}
