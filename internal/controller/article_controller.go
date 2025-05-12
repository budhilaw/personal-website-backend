package controller

import (
	"strconv"

	"github.com/budhilaw/personal-website-backend/internal/model"
	"github.com/budhilaw/personal-website-backend/internal/service"
	"github.com/gofiber/fiber/v2"
)

// ArticleController handles article-related requests
type ArticleController struct {
	articleService service.ArticleService
}

// NewArticleController creates a new ArticleController
func NewArticleController(articleService service.ArticleService) *ArticleController {
	return &ArticleController{
		articleService: articleService,
	}
}

// CreateArticle handles create article requests
func (c *ArticleController) CreateArticle(ctx *fiber.Ctx) error {
	userID := ctx.Locals("user_id").(string)

	var articleReq model.ArticleCreate
	if err := ctx.BodyParser(&articleReq); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate request
	if articleReq.Title == "" || articleReq.Content == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Title and content are required",
		})
	}

	id, err := c.articleService.Create(ctx.Context(), &articleReq, userID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create article",
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id":      id,
		"message": "Article created successfully",
	})
}

// UpdateArticle handles update article requests
func (c *ArticleController) UpdateArticle(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	var articleReq model.ArticleUpdate
	if err := ctx.BodyParser(&articleReq); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate request
	if articleReq.Title == "" || articleReq.Content == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Title and content are required",
		})
	}

	if err := c.articleService.Update(ctx.Context(), id, &articleReq); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update article",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Article updated successfully",
	})
}

// DeleteArticle handles delete article requests
func (c *ArticleController) DeleteArticle(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	if err := c.articleService.Delete(ctx.Context(), id); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete article",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Article deleted successfully",
	})
}

// GetArticle handles get article by ID requests
func (c *ArticleController) GetArticle(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	article, err := c.articleService.GetArticleWithAuthor(ctx.Context(), id)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Article not found",
		})
	}

	return ctx.JSON(article)
}

// GetArticleBySlug handles get article by slug requests
func (c *ArticleController) GetArticleBySlug(ctx *fiber.Ctx) error {
	slug := ctx.Params("slug")
	if slug == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Slug is required",
		})
	}

	article, err := c.articleService.GetBySlugWithAuthor(ctx.Context(), slug)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Article not found",
		})
	}

	return ctx.JSON(article)
}

// ListArticles handles list articles requests
func (c *ArticleController) ListArticles(ctx *fiber.Ctx) error {
	// Parse query parameters
	page, err := strconv.Atoi(ctx.Query("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	perPage, err := strconv.Atoi(ctx.Query("per_page", "10"))
	if err != nil || perPage < 1 {
		perPage = 10
	}

	// Only list published articles for public
	articles, total, err := c.articleService.List(ctx.Context(), page, perPage, true)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to list articles",
		})
	}

	// Convert to response
	var responseArticles []model.ArticleResponse
	for _, article := range articles {
		articleResp, err := c.articleService.GetArticleWithAuthor(ctx.Context(), article.ID)
		if err != nil {
			continue
		}
		responseArticles = append(responseArticles, *articleResp)
	}

	return ctx.JSON(model.ArticleList{
		Articles: responseArticles,
		Total:    total,
		Page:     page,
		PerPage:  perPage,
	})
}

// ListAdminArticles handles list articles for admin
func (c *ArticleController) ListAdminArticles(ctx *fiber.Ctx) error {
	userID := ctx.Locals("user_id").(string)

	// Parse query parameters
	page, err := strconv.Atoi(ctx.Query("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	perPage, err := strconv.Atoi(ctx.Query("per_page", "10"))
	if err != nil || perPage < 1 {
		perPage = 10
	}

	// Get parameter
	onlyMine := ctx.Query("only_mine", "false") == "true"
	if onlyMine {
		// List only user's articles
		articles, total, err := c.articleService.GetByAuthor(ctx.Context(), userID, page, perPage)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to list articles",
			})
		}

		// Convert to response
		var responseArticles []model.ArticleResponse
		for _, article := range articles {
			articleResp, err := c.articleService.GetArticleWithAuthor(ctx.Context(), article.ID)
			if err != nil {
				continue
			}
			responseArticles = append(responseArticles, *articleResp)
		}

		return ctx.JSON(model.ArticleList{
			Articles: responseArticles,
			Total:    total,
			Page:     page,
			PerPage:  perPage,
		})
	}

	// List all articles for admin
	articles, total, err := c.articleService.List(ctx.Context(), page, perPage, false)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to list articles",
		})
	}

	// Convert to response
	var responseArticles []model.ArticleResponse
	for _, article := range articles {
		articleResp, err := c.articleService.GetArticleWithAuthor(ctx.Context(), article.ID)
		if err != nil {
			continue
		}
		responseArticles = append(responseArticles, *articleResp)
	}

	return ctx.JSON(model.ArticleList{
		Articles: responseArticles,
		Total:    total,
		Page:     page,
		PerPage:  perPage,
	})
}
