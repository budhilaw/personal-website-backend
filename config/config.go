package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	AppName string `mapstructure:"APP_NAME"`
	AppEnv  string `mapstructure:"APP_ENV"`
	Port    string `mapstructure:"PORT"`

	PostgresHost     string `mapstructure:"POSTGRES_HOST"`
	PostgresPort     string `mapstructure:"POSTGRES_PORT"`
	PostgresUser     string `mapstructure:"POSTGRES_USER"`
	PostgresPassword string `mapstructure:"POSTGRES_PASSWORD"`
	PostgresDB       string `mapstructure:"POSTGRES_DB"`
	PostgresSSLMode  string `mapstructure:"POSTGRES_SSL_MODE"`

	JWTSecret            string        `mapstructure:"JWT_SECRET"`
	JWTExpiration        time.Duration `mapstructure:"JWT_EXPIRATION"`
	JWTRefreshSecret     string        `mapstructure:"JWT_REFRESH_SECRET"`
	JWTRefreshExpiration time.Duration `mapstructure:"JWT_REFRESH_EXPIRATION"`

	FrontendURL string `mapstructure:"FRONTEND_URL"`

	// Telegram configuration for login activity tracking
	TelegramEnabled  bool   `mapstructure:"TELEGRAM_ENABLED"`
	TelegramBotToken string `mapstructure:"TELEGRAM_BOT_TOKEN"`
	TelegramChatID   string `mapstructure:"TELEGRAM_CHAT_ID"`
	TelegramTopicID  int    `mapstructure:"TELEGRAM_TOPIC_ID"`
}

// IsProduction returns true if the application is running in production mode
func (c *Config) IsProduction() bool {
	return c.AppEnv == "production"
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (config Config, err error) {
	// Load .env file if it exists
	_ = godotenv.Load()

	viper.AutomaticEnv()

	// Set default values
	viper.SetDefault("APP_NAME", "Personal Website API")
	viper.SetDefault("APP_ENV", "development")
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("POSTGRES_HOST", "localhost")
	viper.SetDefault("POSTGRES_PORT", "5432")
	viper.SetDefault("POSTGRES_USER", "postgres")
	viper.SetDefault("POSTGRES_PASSWORD", "postgres")
	viper.SetDefault("POSTGRES_DB", "personal_website")
	viper.SetDefault("POSTGRES_SSL_MODE", "disable")
	viper.SetDefault("JWT_SECRET", "your-secret-key")
	viper.SetDefault("JWT_EXPIRATION", time.Hour*24)
	viper.SetDefault("JWT_REFRESH_SECRET", "your-refresh-secret-key")
	viper.SetDefault("JWT_REFRESH_EXPIRATION", time.Hour*24*7)
	viper.SetDefault("FRONTEND_URL", "http://localhost:3000")

	// Default Telegram settings
	viper.SetDefault("TELEGRAM_ENABLED", false)
	viper.SetDefault("TELEGRAM_BOT_TOKEN", "")
	viper.SetDefault("TELEGRAM_CHAT_ID", "")
	viper.SetDefault("TELEGRAM_TOPIC_ID", 0)

	err = viper.Unmarshal(&config)
	if err != nil {
		return
	}

	// Parse TELEGRAM_TOPIC_ID manually in case it's not set correctly
	if telegramTopicID := os.Getenv("TELEGRAM_TOPIC_ID"); telegramTopicID != "" {
		if topicID, err := strconv.Atoi(telegramTopicID); err == nil {
			config.TelegramTopicID = topicID
		}
	}

	// Print configuration for debugging in development mode
	if config.AppEnv == "development" {
		fmt.Printf("Configuration loaded: %+v\n", config)
	}

	return
}

// GetPostgresConnString returns a PostgreSQL connection string
func (c *Config) GetPostgresConnString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.PostgresHost, c.PostgresPort, c.PostgresUser, c.PostgresPassword, c.PostgresDB, c.PostgresSSLMode)
}

// InitConfig initializes the configuration
func InitConfig() Config {
	cfg, err := LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
		os.Exit(1)
	}
	return cfg
}
