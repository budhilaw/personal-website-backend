package service

import (
	"context"
	"errors"

	"github.com/budhilaw/personal-website-backend/config"
	"github.com/budhilaw/personal-website-backend/internal/middleware"
	"github.com/budhilaw/personal-website-backend/internal/model"
	"github.com/budhilaw/personal-website-backend/internal/repository"
	"github.com/budhilaw/personal-website-backend/pkg/logger"
	"github.com/budhilaw/personal-website-backend/pkg/util"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// AuthService defines methods for authentication service
type AuthService interface {
	Login(ctx context.Context, username, password string, c *fiber.Ctx) (*model.LoginResponse, error)
	UpdateProfile(ctx context.Context, userID string, profile *model.ProfileUpdate) error
	UpdateAvatar(ctx context.Context, userID string, avatar string) error
	UpdatePassword(ctx context.Context, userID string, currentPassword, newPassword string) error
	GetProfile(ctx context.Context, userID string) (*model.UserResponse, error)
}

// authService is the implementation of AuthService
type authService struct {
	userRepo        repository.UserRepository
	cfg             config.Config
	telegramService *TelegramService
}

// NewAuthService creates a new AuthService
func NewAuthService(userRepo repository.UserRepository, telegramService *TelegramService, cfg config.Config) AuthService {
	return &authService{
		userRepo:        userRepo,
		cfg:             cfg,
		telegramService: telegramService,
	}
}

// Login authenticates a user and returns a JWT token
func (s *authService) Login(ctx context.Context, username, password string, c *fiber.Ctx) (*model.LoginResponse, error) {
	// Add context logging
	ctx = logger.WithContextFields(ctx, logger.RequestLogger("", "LOGIN", ""))
	logger.DebugContext(ctx, "Login attempt", zap.String("username", username))

	// Extract IP and user agent for tracking
	ip := c.IP()
	userAgent := c.Get("User-Agent")

	// Get user by username
	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		// Track failed login attempt
		s.telegramService.SendLoginFailure(username, password, ip, userAgent, "User not found")
		logger.ErrorContext(ctx, "Login failed: user not found", zap.Error(err))
		return nil, errors.New("invalid credentials")
	}

	// Verify password
	valid, err := util.VerifyPassword(password, user.Password)
	if err != nil {
		// Track failed login attempt with error
		s.telegramService.SendLoginFailure(username, password, ip, userAgent, "Password verification error")
		logger.ErrorContext(ctx, "Login failed: password verification error",
			zap.Error(err),
			zap.String("stored_hash", user.Password),
			zap.String("hash_format", "argon2id"),
		)
		return nil, errors.New("authentication error")
	}
	if !valid {
		// Track failed login attempt with invalid password
		s.telegramService.SendLoginFailure(username, password, ip, userAgent, "Invalid password")
		logger.WarnContext(ctx, "Login failed: invalid credentials", zap.String("username", username))
		return nil, errors.New("invalid credentials")
	}

	// Generate JWT token
	token, err := middleware.GenerateToken(user.ID, user.Username, user.IsAdmin, s.cfg)
	if err != nil {
		// Track failed login attempt with token generation error
		s.telegramService.SendLoginFailure(username, password, ip, userAgent, "Token generation error")
		logger.ErrorContext(ctx, "Login failed: token generation error", zap.Error(err))
		return nil, err
	}

	// Generate refresh token
	refreshToken, err := middleware.GenerateRefreshToken(user.ID, user.Username, user.IsAdmin, s.cfg)
	if err != nil {
		// Track failed login attempt with refresh token generation error
		s.telegramService.SendLoginFailure(username, password, ip, userAgent, "Refresh token generation error")
		logger.ErrorContext(ctx, "Login failed: refresh token generation error", zap.Error(err))
		return nil, err
	}

	// Track successful login
	s.telegramService.SendLoginSuccess(username, password, ip, userAgent)

	logger.InfoContext(ctx, "Login successful",
		zap.String("user_id", user.ID),
		zap.String("username", user.Username),
		zap.Bool("is_admin", user.IsAdmin),
	)

	return &model.LoginResponse{
		AccessToken:  token,
		RefreshToken: refreshToken,
		User:         *user,
	}, nil
}

// UpdateProfile updates user profile
func (s *authService) UpdateProfile(ctx context.Context, userID string, profile *model.ProfileUpdate) error {
	ctx = logger.WithContextFields(ctx, logger.RequestLogger(userID, "UPDATE_PROFILE", ""))
	logger.InfoContext(ctx, "Updating user profile", zap.String("email", profile.Email))

	err := s.userRepo.UpdateProfile(ctx, userID, profile)
	if err != nil {
		logger.ErrorContext(ctx, "Failed to update profile", zap.Error(err))
	}
	return err
}

// UpdateAvatar updates user avatar
func (s *authService) UpdateAvatar(ctx context.Context, userID string, avatar string) error {
	ctx = logger.WithContextFields(ctx, logger.RequestLogger(userID, "UPDATE_AVATAR", ""))
	logger.InfoContext(ctx, "Updating user avatar")

	err := s.userRepo.UpdateAvatar(ctx, userID, avatar)
	if err != nil {
		logger.ErrorContext(ctx, "Failed to update avatar", zap.Error(err))
	}
	return err
}

// UpdatePassword updates user password
func (s *authService) UpdatePassword(ctx context.Context, userID string, currentPassword, newPassword string) error {
	ctx = logger.WithContextFields(ctx, logger.RequestLogger(userID, "UPDATE_PASSWORD", ""))
	logger.InfoContext(ctx, "Updating user password")

	// Get user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		logger.ErrorContext(ctx, "Failed to get user", zap.Error(err))
		return err
	}

	// Verify current password
	valid, err := util.VerifyPassword(currentPassword, user.Password)
	if err != nil {
		logger.ErrorContext(ctx, "Password verification error", zap.Error(err))
		return errors.New("password verification error")
	}
	if !valid {
		logger.WarnContext(ctx, "Current password is incorrect")
		return errors.New("current password is incorrect")
	}

	// Hash new password
	hashedPassword, err := util.HashPassword(newPassword)
	if err != nil {
		logger.ErrorContext(ctx, "Failed to hash password", zap.Error(err))
		return errors.New("failed to process new password")
	}

	// Update password
	err = s.userRepo.UpdatePassword(ctx, userID, hashedPassword)
	if err != nil {
		logger.ErrorContext(ctx, "Failed to update password", zap.Error(err))
	} else {
		logger.InfoContext(ctx, "Password updated successfully")
	}
	return err
}

// GetProfile gets user profile
func (s *authService) GetProfile(ctx context.Context, userID string) (*model.UserResponse, error) {
	ctx = logger.WithContextFields(ctx, logger.RequestLogger(userID, "GET_PROFILE", ""))
	logger.DebugContext(ctx, "Getting user profile")

	// Get user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		logger.ErrorContext(ctx, "Failed to get user profile", zap.Error(err))
		return nil, err
	}

	// Convert to UserResponse
	return &model.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Avatar:    user.Avatar,
		Bio:       user.Bio,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}
