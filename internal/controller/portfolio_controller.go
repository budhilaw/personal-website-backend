package controller

import (
	"strconv"

	"github.com/budhilaw/personal-website-backend/internal/model"
	"github.com/budhilaw/personal-website-backend/internal/service"
	"github.com/gofiber/fiber/v2"
)

// PortfolioController handles portfolio-related requests
type PortfolioController struct {
	portfolioService service.PortfolioService
}

// NewPortfolioController creates a new PortfolioController
func NewPortfolioController(portfolioService service.PortfolioService) *PortfolioController {
	return &PortfolioController{
		portfolioService: portfolioService,
	}
}

// CreatePortfolio handles create portfolio requests
func (c *PortfolioController) CreatePortfolio(ctx *fiber.Ctx) error {
	userID := ctx.Locals("user_id").(string)

	var portfolioReq model.PortfolioCreate
	if err := ctx.BodyParser(&portfolioReq); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate request
	if portfolioReq.Title == "" || portfolioReq.Description == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Title and description are required",
		})
	}

	id, err := c.portfolioService.Create(ctx.Context(), &portfolioReq, userID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create portfolio",
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id":      id,
		"message": "Portfolio created successfully",
	})
}

// UpdatePortfolio handles update portfolio requests
func (c *PortfolioController) UpdatePortfolio(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	var portfolioReq model.PortfolioUpdate
	if err := ctx.BodyParser(&portfolioReq); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate request
	if portfolioReq.Title == "" || portfolioReq.Description == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Title and description are required",
		})
	}

	if err := c.portfolioService.Update(ctx.Context(), id, &portfolioReq); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update portfolio",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Portfolio updated successfully",
	})
}

// DeletePortfolio handles delete portfolio requests
func (c *PortfolioController) DeletePortfolio(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	if err := c.portfolioService.Delete(ctx.Context(), id); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete portfolio",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Portfolio deleted successfully",
	})
}

// GetPortfolio handles get portfolio by ID requests
func (c *PortfolioController) GetPortfolio(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	portfolio, err := c.portfolioService.GetPortfolioWithAuthor(ctx.Context(), id)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Portfolio not found",
		})
	}

	return ctx.JSON(portfolio)
}

// GetPortfolioBySlug handles get portfolio by slug requests
func (c *PortfolioController) GetPortfolioBySlug(ctx *fiber.Ctx) error {
	slug := ctx.Params("slug")
	if slug == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Slug is required",
		})
	}

	portfolio, err := c.portfolioService.GetBySlugWithAuthor(ctx.Context(), slug)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Portfolio not found",
		})
	}

	return ctx.JSON(portfolio)
}

// ListPortfolios handles list portfolios requests
func (c *PortfolioController) ListPortfolios(ctx *fiber.Ctx) error {
	// Parse query parameters
	page, err := strconv.Atoi(ctx.Query("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	perPage, err := strconv.Atoi(ctx.Query("per_page", "10"))
	if err != nil || perPage < 1 {
		perPage = 10
	}

	// Only list published portfolios for public
	portfolios, total, err := c.portfolioService.List(ctx.Context(), page, perPage, true)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to list portfolios",
		})
	}

	// Convert to response
	var responsePortfolios []model.PortfolioResponse
	for _, portfolio := range portfolios {
		portfolioResp, err := c.portfolioService.GetPortfolioWithAuthor(ctx.Context(), portfolio.ID)
		if err != nil {
			continue
		}
		responsePortfolios = append(responsePortfolios, *portfolioResp)
	}

	return ctx.JSON(model.PortfolioList{
		Portfolios: responsePortfolios,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
	})
}

// ListAdminPortfolios handles list portfolios for admin
func (c *PortfolioController) ListAdminPortfolios(ctx *fiber.Ctx) error {
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
		// List only user's portfolios
		portfolios, total, err := c.portfolioService.GetByAuthor(ctx.Context(), userID, page, perPage)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to list portfolios",
			})
		}

		// Convert to response
		var responsePortfolios []model.PortfolioResponse
		for _, portfolio := range portfolios {
			portfolioResp, err := c.portfolioService.GetPortfolioWithAuthor(ctx.Context(), portfolio.ID)
			if err != nil {
				continue
			}
			responsePortfolios = append(responsePortfolios, *portfolioResp)
		}

		return ctx.JSON(model.PortfolioList{
			Portfolios: responsePortfolios,
			Total:      total,
			Page:       page,
			PerPage:    perPage,
		})
	}

	// List all portfolios for admin (both published and unpublished)
	portfolios, total, err := c.portfolioService.List(ctx.Context(), page, perPage, false)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to list portfolios",
		})
	}

	// Convert to response
	var responsePortfolios []model.PortfolioResponse
	for _, portfolio := range portfolios {
		portfolioResp, err := c.portfolioService.GetPortfolioWithAuthor(ctx.Context(), portfolio.ID)
		if err != nil {
			continue
		}
		responsePortfolios = append(responsePortfolios, *portfolioResp)
	}

	return ctx.JSON(model.PortfolioList{
		Portfolios: responsePortfolios,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
	})
}
