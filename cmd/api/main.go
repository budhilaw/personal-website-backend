package main

import (
	"os"

	"github.com/budhilaw/personal-website-backend/config"
	"github.com/budhilaw/personal-website-backend/db"
	"github.com/budhilaw/personal-website-backend/internal/controller"
	"github.com/budhilaw/personal-website-backend/internal/middleware"
	"github.com/budhilaw/personal-website-backend/internal/repository"
	"github.com/budhilaw/personal-website-backend/internal/router"
	"github.com/budhilaw/personal-website-backend/internal/service"
	"github.com/budhilaw/personal-website-backend/pkg/logger"
	"github.com/gofiber/fiber/v2"
	fiberRecover "github.com/gofiber/fiber/v2/middleware/recover"
	"go.uber.org/zap"
)

func main() {
	// Load config
	cfg := config.InitConfig()

	// Initialize logger
	log := logger.InitLogger(cfg.IsProduction())
	defer log.Sync()

	// Initialize JWT Manager for secret rotation
	middleware.InitJWTManager(cfg)

	// Check if this is a database command
	if len(os.Args) > 1 {
		handleDBCommand()
	}

	// Initialize database
	database, err := db.InitDB(cfg)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer database.Close()

	// Run migrations
	if err := db.RunMigrations(database); err != nil {
		logger.Fatal("Failed to run migrations", zap.Error(err))
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(database)
	articleRepo := repository.NewArticleRepository(database)
	portfolioRepo := repository.NewPortfolioRepository(database)
	telegramRepo := repository.NewTelegramRepository(cfg, log)

	// Initialize services
	telegramService := service.NewTelegramService(telegramRepo, cfg, log)
	authService := service.NewAuthService(userRepo, telegramService, cfg)
	articleService := service.NewArticleService(articleRepo, userRepo)
	portfolioService := service.NewPortfolioService(portfolioRepo, userRepo)

	// Initialize controllers
	authController := controller.NewAuthController(authService, cfg)
	articleController := controller.NewArticleController(articleService)
	portfolioController := controller.NewPortfolioController(portfolioService)

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		// Enable case sensitive routing
		CaseSensitive: true,
		// Enable strict routing
		StrictRouting: true,
		// Set server name
		ServerHeader: "Personal Website API",
		// Set application name
		AppName: cfg.AppName,
	})

	// Use global middlewares
	app.Use(fiberRecover.New())
	app.Use(middleware.ZapLogger())

	// Security middleware
	app.Use(middleware.Security(cfg.FrontendURL))
	app.Use(middleware.Helmet())
	app.Use(middleware.RateLimiter())

	// Brute force protection
	app.Use(middleware.BruteForceProtection())
	app.Use(middleware.TrackLoginAttempt())

	// Setup routes
	router.SetupRoutes(app, authController, articleController, portfolioController, cfg)

	// Start server
	logger.Info("Starting server", zap.String("port", cfg.Port))
	if err := app.Listen(":" + cfg.Port); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}
