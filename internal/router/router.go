package router

import (
	"github.com/budhilaw/personal-website-backend/config"
	"github.com/budhilaw/personal-website-backend/internal/controller"
	"github.com/budhilaw/personal-website-backend/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

// SetupRoutes sets up the API routes
func SetupRoutes(
	app *fiber.App,
	authController *controller.AuthController,
	articleController *controller.ArticleController,
	portfolioController *controller.PortfolioController,
	cfg config.Config,
) {
	// API v1 group
	v1 := app.Group("/api/v1")

	// Public routes
	public := v1.Group("/public")
	setupPublicRoutes(public, articleController, portfolioController)

	// Admin routes (protected)
	admin := v1.Group("/admin")
	admin.Use(middleware.Protected(cfg))
	admin.Use(middleware.AdminOnly())
	setupAdminRoutes(admin, authController, articleController, portfolioController)

	// Auth routes
	auth := v1.Group("/auth")
	setupAuthRoutes(auth, authController, cfg)
}

// setupPublicRoutes sets up public routes
func setupPublicRoutes(
	router fiber.Router,
	articleController *controller.ArticleController,
	portfolioController *controller.PortfolioController,
) {
	// Articles
	articles := router.Group("/articles")
	articles.Get("/", articleController.ListArticles)
	articles.Get("/:id", articleController.GetArticle)
	articles.Get("/slug/:slug", articleController.GetArticleBySlug)

	// Portfolios
	portfolios := router.Group("/portfolios")
	portfolios.Get("/", portfolioController.ListPortfolios)
	portfolios.Get("/:id", portfolioController.GetPortfolio)
	portfolios.Get("/slug/:slug", portfolioController.GetPortfolioBySlug)
}

// setupAdminRoutes sets up admin routes
func setupAdminRoutes(
	router fiber.Router,
	authController *controller.AuthController,
	articleController *controller.ArticleController,
	portfolioController *controller.PortfolioController,
) {
	// Profile
	profile := router.Group("/profile")
	profile.Get("/", authController.GetProfile)
	profile.Put("/", authController.UpdateProfile)
	profile.Put("/avatar", authController.UpdateAvatar)
	profile.Put("/password", authController.UpdatePassword)

	// Articles
	articles := router.Group("/articles")
	articles.Get("/", articleController.ListAdminArticles)
	articles.Post("/", articleController.CreateArticle)
	articles.Put("/:id", articleController.UpdateArticle)
	articles.Delete("/:id", articleController.DeleteArticle)
	articles.Get("/:id", articleController.GetArticle)

	// Portfolios
	portfolios := router.Group("/portfolios")
	portfolios.Get("/", portfolioController.ListAdminPortfolios)
	portfolios.Post("/", portfolioController.CreatePortfolio)
	portfolios.Put("/:id", portfolioController.UpdatePortfolio)
	portfolios.Delete("/:id", portfolioController.DeletePortfolio)
	portfolios.Get("/:id", portfolioController.GetPortfolio)
}

// setupAuthRoutes sets up authentication routes
func setupAuthRoutes(
	router fiber.Router,
	authController *controller.AuthController,
	cfg config.Config,
) {
	router.Post("/login", authController.Login)
} 