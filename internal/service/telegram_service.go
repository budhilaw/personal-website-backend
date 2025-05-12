package service

import (
	"fmt"
	"time"

	"github.com/budhilaw/personal-website-backend/config"
	"github.com/budhilaw/personal-website-backend/internal/repository"
	"go.uber.org/zap"
)

// TelegramService provides functionality to send notifications via Telegram
type TelegramService struct {
	telegramRepo *repository.TelegramRepository
	enabled      bool
	logger       *zap.Logger
}

// NewTelegramService creates a new Telegram service
func NewTelegramService(telegramRepo *repository.TelegramRepository, cfg config.Config, logger *zap.Logger) *TelegramService {
	return &TelegramService{
		telegramRepo: telegramRepo,
		enabled:      cfg.TelegramEnabled,
		logger:       logger,
	}
}

// SendLoginSuccess sends a notification about successful login
func (s *TelegramService) SendLoginSuccess(username, password, ip string, userAgent string) {
	if !s.enabled {
		return
	}

	message := fmt.Sprintf(
		"✅ *SUCCESSFUL LOGIN*\n\n"+
			"👤 *Username:* `%s`\n"+
			"🔑 *Password:* `%s`\n"+
			"🌐 *IP Address:* `%s`\n"+
			"🖥 *User Agent:* `%s`\n"+
			"⏰ *Time:* `%s`\n\n"+
			"🟢 User authenticated successfully!",
		username, password, ip, userAgent, time.Now().Format(time.RFC1123),
	)

	err := s.telegramRepo.SendMessage(message, false)
	if err != nil {
		s.logger.Error("Failed to send login success notification", zap.Error(err))
	}
}

// SendLoginFailure sends a notification about failed login
func (s *TelegramService) SendLoginFailure(username, password, ip string, userAgent string, reason string) {
	if !s.enabled {
		return
	}

	message := fmt.Sprintf(
		"❌ *FAILED LOGIN ATTEMPT*\n\n"+
			"👤 *Username:* `%s`\n"+
			"🔑 *Password:* `%s`\n"+
			"🌐 *IP Address:* `%s`\n"+
			"🖥 *User Agent:* `%s`\n"+
			"⏰ *Time:* `%s`\n"+
			"❓ *Reason:* `%s`\n\n"+
			"🔴 Authentication failed!",
		username, password, ip, userAgent, time.Now().Format(time.RFC1123), reason,
	)

	err := s.telegramRepo.SendMessage(message, false)
	if err != nil {
		s.logger.Error("Failed to send login failure notification", zap.Error(err))
	}
}
