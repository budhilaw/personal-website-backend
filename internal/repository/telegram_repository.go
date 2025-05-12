package repository

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/budhilaw/personal-website-backend/config"
	"go.uber.org/zap"
)

// TelegramRepository handles API calls to Telegram
type TelegramRepository struct {
	botToken   string
	chatID     string
	topicID    int
	httpClient *http.Client
	logger     *zap.Logger
}

// NewTelegramRepository creates a new Telegram repository
func NewTelegramRepository(cfg config.Config, logger *zap.Logger) *TelegramRepository {
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	return &TelegramRepository{
		botToken:   cfg.TelegramBotToken,
		chatID:     cfg.TelegramChatID,
		topicID:    cfg.TelegramTopicID,
		httpClient: httpClient,
		logger:     logger,
	}
}

// TelegramMessage is a Telegram message to be sent
type TelegramMessage struct {
	ChatID              string `json:"chat_id"`
	Text                string `json:"text"`
	ParseMode           string `json:"parse_mode"`
	DisableNotification bool   `json:"disable_notification"`
	MessageThreadID     int    `json:"message_thread_id,omitempty"`
}

// SendMessage sends a message to Telegram
func (r *TelegramRepository) SendMessage(message string, disableNotification bool) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", r.botToken)

	msg := TelegramMessage{
		ChatID:              r.chatID,
		Text:                message,
		ParseMode:           "Markdown",
		DisableNotification: disableNotification,
	}

	// Add topic ID if it's set
	if r.topicID > 0 {
		msg.MessageThreadID = r.topicID
	}

	jsonData, err := json.Marshal(msg)
	if err != nil {
		r.logger.Error("Failed to marshal Telegram message", zap.Error(err))
		return err
	}

	resp, err := r.httpClient.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		r.logger.Error("Failed to send Telegram message", zap.Error(err))
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		r.logger.Error("Telegram API returned non-OK status",
			zap.Int("status_code", resp.StatusCode),
			zap.String("status", resp.Status))
		return fmt.Errorf("telegram API returned status %d", resp.StatusCode)
	}

	r.logger.Debug("Telegram message sent successfully")
	return nil
}
