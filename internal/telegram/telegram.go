package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/budhilaw/personal-website-backend/config"
	"github.com/budhilaw/personal-website-backend/internal/logger"
	"go.uber.org/zap"
)

// TelegramService handles sending messages to Telegram
type TelegramService struct {
	botToken   string
	chatID     string
	topicID    int
	httpClient *http.Client
	enabled    bool
}

// NewTelegramService creates a new Telegram service
func NewTelegramService(cfg config.Config) *TelegramService {
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	return &TelegramService{
		botToken:   cfg.TelegramBotToken,
		chatID:     cfg.TelegramChatID,
		topicID:    cfg.TelegramTopicID,
		httpClient: httpClient,
		enabled:    cfg.TelegramEnabled,
	}
}

// Message is a Telegram message to be sent
type Message struct {
	ChatID              string `json:"chat_id"`
	Text                string `json:"text"`
	ParseMode           string `json:"parse_mode"`
	DisableNotification bool   `json:"disable_notification"`
	MessageThreadID     int    `json:"message_thread_id,omitempty"`
}

// SendLoginSuccess sends a notification about successful login
func (t *TelegramService) SendLoginSuccess(username, password, ip string, userAgent string) {
	if !t.enabled {
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

	t.sendMessage(message)
}

// SendLoginFailure sends a notification about failed login
func (t *TelegramService) SendLoginFailure(username, password, ip string, userAgent string, reason string) {
	if !t.enabled {
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

	t.sendMessage(message)
}

// sendMessage sends a message to Telegram
func (t *TelegramService) sendMessage(text string) {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", t.botToken)

	msg := Message{
		ChatID:              t.chatID,
		Text:                text,
		ParseMode:           "Markdown",
		DisableNotification: false,
	}

	// Add topic ID if it's set
	if t.topicID > 0 {
		msg.MessageThreadID = t.topicID
	}

	jsonData, err := json.Marshal(msg)
	if err != nil {
		logger.Error("Failed to marshal Telegram message", zap.Error(err))
		return
	}

	resp, err := t.httpClient.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		logger.Error("Failed to send Telegram message", zap.Error(err))
		return
	}
	defer resp.Body.Close()

	// request url
	fmt.Println(url)

	// request body
	fmt.Println(string(jsonData))

	if resp.StatusCode != http.StatusOK {
		logger.Error("Telegram API returned non-OK status",
			zap.Int("status_code", resp.StatusCode),
			zap.String("status", resp.Status))
		return
	}

	logger.Debug("Telegram message sent successfully")
}
