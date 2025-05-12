package service

import (
	"context"

	"github.com/budhilaw/personal-website-backend/internal/model"
	"github.com/budhilaw/personal-website-backend/internal/repository"
)

// PortfolioService defines methods for portfolio service
type PortfolioService interface {
	Create(ctx context.Context, portfolio *model.PortfolioCreate, userID string) (string, error)
	Update(ctx context.Context, id string, portfolio *model.PortfolioUpdate) error
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*model.Portfolio, error)
	GetBySlug(ctx context.Context, slug string) (*model.Portfolio, error)
	List(ctx context.Context, page, perPage int, onlyPublished bool) ([]model.Portfolio, int, error)
	GetByAuthor(ctx context.Context, userID string, page, perPage int) ([]model.Portfolio, int, error)
	GetPortfolioWithAuthor(ctx context.Context, id string) (*model.PortfolioResponse, error)
	GetBySlugWithAuthor(ctx context.Context, slug string) (*model.PortfolioResponse, error)
}

// portfolioService is the implementation of PortfolioService
type portfolioService struct {
	portfolioRepo repository.PortfolioRepository
	userRepo      repository.UserRepository
}

// NewPortfolioService creates a new PortfolioService
func NewPortfolioService(portfolioRepo repository.PortfolioRepository, userRepo repository.UserRepository) PortfolioService {
	return &portfolioService{
		portfolioRepo: portfolioRepo,
		userRepo:      userRepo,
	}
}

// Create creates a new portfolio
func (s *portfolioService) Create(ctx context.Context, portfolio *model.PortfolioCreate, userID string) (string, error) {
	return s.portfolioRepo.Create(ctx, portfolio, userID)
}

// Update updates a portfolio
func (s *portfolioService) Update(ctx context.Context, id string, portfolio *model.PortfolioUpdate) error {
	return s.portfolioRepo.Update(ctx, id, portfolio)
}

// Delete deletes a portfolio
func (s *portfolioService) Delete(ctx context.Context, id string) error {
	return s.portfolioRepo.Delete(ctx, id)
}

// GetByID gets a portfolio by ID
func (s *portfolioService) GetByID(ctx context.Context, id string) (*model.Portfolio, error) {
	return s.portfolioRepo.GetByID(ctx, id)
}

// GetBySlug gets a portfolio by slug
func (s *portfolioService) GetBySlug(ctx context.Context, slug string) (*model.Portfolio, error) {
	return s.portfolioRepo.GetBySlug(ctx, slug)
}

// List lists portfolios with pagination
func (s *portfolioService) List(ctx context.Context, page, perPage int, onlyPublished bool) ([]model.Portfolio, int, error) {
	return s.portfolioRepo.List(ctx, page, perPage, onlyPublished)
}

// GetByAuthor gets portfolios by author ID with pagination
func (s *portfolioService) GetByAuthor(ctx context.Context, userID string, page, perPage int) ([]model.Portfolio, int, error) {
	return s.portfolioRepo.GetByAuthor(ctx, userID, page, perPage)
}

// GetPortfolioWithAuthor gets a portfolio with author information
func (s *portfolioService) GetPortfolioWithAuthor(ctx context.Context, id string) (*model.PortfolioResponse, error) {
	portfolio, err := s.portfolioRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	author, err := s.userRepo.GetByID(ctx, portfolio.UserID)
	if err != nil {
		return nil, err
	}

	response := &model.PortfolioResponse{
		ID:           portfolio.ID,
		Title:        portfolio.Title,
		Slug:         portfolio.Slug,
		Description:  portfolio.Description,
		Image:        portfolio.Image,
		ProjectURL:   portfolio.ProjectURL,
		GithubURL:    portfolio.GithubURL,
		Technologies: portfolio.Technologies,
		IsPublished:  portfolio.IsPublished,
		CreatedAt:    portfolio.CreatedAt,
		UpdatedAt:    portfolio.UpdatedAt,
	}

	response.Author.ID = author.ID
	response.Author.Username = author.Username
	response.Author.FirstName = author.FirstName
	response.Author.LastName = author.LastName
	response.Author.Avatar = author.Avatar

	return response, nil
}

// GetBySlugWithAuthor gets a portfolio by slug with author information
func (s *portfolioService) GetBySlugWithAuthor(ctx context.Context, slug string) (*model.PortfolioResponse, error) {
	portfolio, err := s.portfolioRepo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	author, err := s.userRepo.GetByID(ctx, portfolio.UserID)
	if err != nil {
		return nil, err
	}

	response := &model.PortfolioResponse{
		ID:           portfolio.ID,
		Title:        portfolio.Title,
		Slug:         portfolio.Slug,
		Description:  portfolio.Description,
		Image:        portfolio.Image,
		ProjectURL:   portfolio.ProjectURL,
		GithubURL:    portfolio.GithubURL,
		Technologies: portfolio.Technologies,
		IsPublished:  portfolio.IsPublished,
		CreatedAt:    portfolio.CreatedAt,
		UpdatedAt:    portfolio.UpdatedAt,
	}

	response.Author.ID = author.ID
	response.Author.Username = author.Username
	response.Author.FirstName = author.FirstName
	response.Author.LastName = author.LastName
	response.Author.Avatar = author.Avatar

	return response, nil
}
