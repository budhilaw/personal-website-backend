package service

import (
	"context"

	"github.com/budhilaw/personal-website-backend/internal/model"
	"github.com/budhilaw/personal-website-backend/internal/repository"
)

// ArticleService defines methods for article service
type ArticleService interface {
	Create(ctx context.Context, article *model.ArticleCreate, userID string) (string, error)
	Update(ctx context.Context, id string, article *model.ArticleUpdate) error
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*model.Article, error)
	GetBySlug(ctx context.Context, slug string) (*model.Article, error)
	List(ctx context.Context, page, perPage int, onlyPublished bool) ([]model.Article, int, error)
	GetByAuthor(ctx context.Context, userID string, page, perPage int) ([]model.Article, int, error)
	GetArticleWithAuthor(ctx context.Context, id string) (*model.ArticleResponse, error)
	GetBySlugWithAuthor(ctx context.Context, slug string) (*model.ArticleResponse, error)
}

// articleService is the implementation of ArticleService
type articleService struct {
	articleRepo repository.ArticleRepository
	userRepo    repository.UserRepository
}

// NewArticleService creates a new ArticleService
func NewArticleService(articleRepo repository.ArticleRepository, userRepo repository.UserRepository) ArticleService {
	return &articleService{
		articleRepo: articleRepo,
		userRepo:    userRepo,
	}
}

// Create creates a new article
func (s *articleService) Create(ctx context.Context, article *model.ArticleCreate, userID string) (string, error) {
	return s.articleRepo.Create(ctx, article, userID)
}

// Update updates an article
func (s *articleService) Update(ctx context.Context, id string, article *model.ArticleUpdate) error {
	return s.articleRepo.Update(ctx, id, article)
}

// Delete deletes an article
func (s *articleService) Delete(ctx context.Context, id string) error {
	return s.articleRepo.Delete(ctx, id)
}

// GetByID gets an article by ID
func (s *articleService) GetByID(ctx context.Context, id string) (*model.Article, error) {
	return s.articleRepo.GetByID(ctx, id)
}

// GetBySlug gets an article by slug
func (s *articleService) GetBySlug(ctx context.Context, slug string) (*model.Article, error) {
	return s.articleRepo.GetBySlug(ctx, slug)
}

// List lists articles with pagination
func (s *articleService) List(ctx context.Context, page, perPage int, onlyPublished bool) ([]model.Article, int, error) {
	return s.articleRepo.List(ctx, page, perPage, onlyPublished)
}

// GetByAuthor gets articles by author ID with pagination
func (s *articleService) GetByAuthor(ctx context.Context, userID string, page, perPage int) ([]model.Article, int, error) {
	return s.articleRepo.GetByAuthor(ctx, userID, page, perPage)
}

// GetArticleWithAuthor gets an article with author information
func (s *articleService) GetArticleWithAuthor(ctx context.Context, id string) (*model.ArticleResponse, error) {
	article, err := s.articleRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	author, err := s.userRepo.GetByID(ctx, article.UserID)
	if err != nil {
		return nil, err
	}

	response := &model.ArticleResponse{
		ID:            article.ID,
		Title:         article.Title,
		Slug:          article.Slug,
		Content:       article.Content,
		Excerpt:       article.Excerpt,
		FeaturedImage: article.FeaturedImage,
		IsPublished:   article.IsPublished,
		CreatedAt:     article.CreatedAt,
		UpdatedAt:     article.UpdatedAt,
		PublishedAt:   article.PublishedAt,
	}

	response.Author.ID = author.ID
	response.Author.Username = author.Username
	response.Author.FirstName = author.FirstName
	response.Author.LastName = author.LastName
	response.Author.Avatar = author.Avatar

	return response, nil
}

// GetBySlugWithAuthor gets an article by slug with author information
func (s *articleService) GetBySlugWithAuthor(ctx context.Context, slug string) (*model.ArticleResponse, error) {
	article, err := s.articleRepo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	author, err := s.userRepo.GetByID(ctx, article.UserID)
	if err != nil {
		return nil, err
	}

	response := &model.ArticleResponse{
		ID:            article.ID,
		Title:         article.Title,
		Slug:          article.Slug,
		Content:       article.Content,
		Excerpt:       article.Excerpt,
		FeaturedImage: article.FeaturedImage,
		IsPublished:   article.IsPublished,
		CreatedAt:     article.CreatedAt,
		UpdatedAt:     article.UpdatedAt,
		PublishedAt:   article.PublishedAt,
	}

	response.Author.ID = author.ID
	response.Author.Username = author.Username
	response.Author.FirstName = author.FirstName
	response.Author.LastName = author.LastName
	response.Author.Avatar = author.Avatar

	return response, nil
}
