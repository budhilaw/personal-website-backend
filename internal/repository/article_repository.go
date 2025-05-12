package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/budhilaw/personal-website-backend/internal/model"
	"github.com/budhilaw/personal-website-backend/internal/util"
	"github.com/jmoiron/sqlx"
)

// ArticleRepository defines methods for article repository
type ArticleRepository interface {
	Create(ctx context.Context, article *model.ArticleCreate, userID string) (string, error)
	Update(ctx context.Context, id string, article *model.ArticleUpdate) error
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*model.Article, error)
	GetBySlug(ctx context.Context, slug string) (*model.Article, error)
	List(ctx context.Context, page, perPage int, onlyPublished bool) ([]model.Article, int, error)
	GetByAuthor(ctx context.Context, userID string, page, perPage int) ([]model.Article, int, error)
}

// articleRepository is the implementation of ArticleRepository
type articleRepository struct {
	db *sqlx.DB
}

// NewArticleRepository creates a new ArticleRepository
func NewArticleRepository(db *sqlx.DB) ArticleRepository {
	return &articleRepository{db: db}
}

// Create creates a new article
func (r *articleRepository) Create(ctx context.Context, articleCreate *model.ArticleCreate, userID string) (string, error) {
	query := `INSERT INTO articles (title, slug, content, excerpt, featured_image, is_published, user_id, published_at) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8) 
			  RETURNING id`

	slug := util.GenerateSlug(articleCreate.Title)
	var publishedAt sql.NullTime
	if articleCreate.IsPublished {
		publishedAt = sql.NullTime{Time: time.Now(), Valid: true}
	}

	var id string
	err := r.db.QueryRowContext(
		ctx, query,
		articleCreate.Title,
		slug,
		articleCreate.Content,
		articleCreate.Excerpt,
		articleCreate.FeaturedImage,
		articleCreate.IsPublished,
		userID,
		publishedAt,
	).Scan(&id)
	if err != nil {
		return "", err
	}

	return id, nil
}

// Update updates an article
func (r *articleRepository) Update(ctx context.Context, id string, articleUpdate *model.ArticleUpdate) error {
	// Get current state to check if published state changed
	var currentState bool
	err := r.db.QueryRowContext(ctx, "SELECT is_published FROM articles WHERE id = $1", id).Scan(&currentState)
	if err != nil {
		return err
	}

	query := `UPDATE articles 
			  SET title = $2, slug = $3, content = $4, excerpt = $5, featured_image = $6, is_published = $7, updated_at = $8`

	params := []interface{}{
		id,
		articleUpdate.Title,
		util.GenerateSlug(articleUpdate.Title),
		articleUpdate.Content,
		articleUpdate.Excerpt,
		articleUpdate.FeaturedImage,
		articleUpdate.IsPublished,
		time.Now(),
	}

	// If article is being published now
	if !currentState && articleUpdate.IsPublished {
		query += ", published_at = $9 WHERE id = $1"
		params = append(params, time.Now())
	} else {
		query += " WHERE id = $1"
	}

	_, err = r.db.ExecContext(ctx, query, params...)
	return err
}

// Delete deletes an article
func (r *articleRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM articles WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// GetByID gets an article by ID
func (r *articleRepository) GetByID(ctx context.Context, id string) (*model.Article, error) {
	query := `SELECT id, title, slug, content, excerpt, featured_image, is_published, user_id, created_at, updated_at, published_at 
			  FROM articles 
			  WHERE id = $1`

	var article model.Article
	var publishedAt sql.NullTime
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&article.ID,
		&article.Title,
		&article.Slug,
		&article.Content,
		&article.Excerpt,
		&article.FeaturedImage,
		&article.IsPublished,
		&article.UserID,
		&article.CreatedAt,
		&article.UpdatedAt,
		&publishedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("article not found")
		}
		return nil, err
	}

	if publishedAt.Valid {
		article.PublishedAt = publishedAt.Time
	}

	return &article, nil
}

// GetBySlug gets an article by slug
func (r *articleRepository) GetBySlug(ctx context.Context, slug string) (*model.Article, error) {
	query := `SELECT id, title, slug, content, excerpt, featured_image, is_published, user_id, created_at, updated_at, published_at 
			  FROM articles 
			  WHERE slug = $1`

	var article model.Article
	var publishedAt sql.NullTime
	err := r.db.QueryRowContext(ctx, query, slug).Scan(
		&article.ID,
		&article.Title,
		&article.Slug,
		&article.Content,
		&article.Excerpt,
		&article.FeaturedImage,
		&article.IsPublished,
		&article.UserID,
		&article.CreatedAt,
		&article.UpdatedAt,
		&publishedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("article not found")
		}
		return nil, err
	}

	if publishedAt.Valid {
		article.PublishedAt = publishedAt.Time
	}

	return &article, nil
}

// List lists articles with pagination
func (r *articleRepository) List(ctx context.Context, page, perPage int, onlyPublished bool) ([]model.Article, int, error) {
	offset := (page - 1) * perPage

	// Count total
	countQuery := `SELECT COUNT(*) FROM articles`
	if onlyPublished {
		countQuery += ` WHERE is_published = true`
	}

	var total int
	err := r.db.QueryRowContext(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get articles
	query := `SELECT id, title, slug, content, excerpt, featured_image, is_published, user_id, created_at, updated_at, published_at 
			  FROM articles`
	if onlyPublished {
		query += ` WHERE is_published = true`
	}
	query += ` ORDER BY created_at DESC LIMIT $1 OFFSET $2`

	rows, err := r.db.QueryContext(ctx, query, perPage, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var articles []model.Article
	for rows.Next() {
		var article model.Article
		var publishedAt sql.NullTime
		err := rows.Scan(
			&article.ID,
			&article.Title,
			&article.Slug,
			&article.Content,
			&article.Excerpt,
			&article.FeaturedImage,
			&article.IsPublished,
			&article.UserID,
			&article.CreatedAt,
			&article.UpdatedAt,
			&publishedAt,
		)
		if err != nil {
			return nil, 0, err
		}

		if publishedAt.Valid {
			article.PublishedAt = publishedAt.Time
		}

		articles = append(articles, article)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return articles, total, nil
}

// GetByAuthor gets articles by author ID with pagination
func (r *articleRepository) GetByAuthor(ctx context.Context, userID string, page, perPage int) ([]model.Article, int, error) {
	offset := (page - 1) * perPage

	// Count total
	countQuery := `SELECT COUNT(*) FROM articles WHERE user_id = $1`
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, userID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get articles
	query := `SELECT id, title, slug, content, excerpt, featured_image, is_published, user_id, created_at, updated_at, published_at 
			  FROM articles 
			  WHERE user_id = $1 
			  ORDER BY created_at DESC 
			  LIMIT $2 OFFSET $3`

	rows, err := r.db.QueryContext(ctx, query, userID, perPage, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var articles []model.Article
	for rows.Next() {
		var article model.Article
		var publishedAt sql.NullTime
		err := rows.Scan(
			&article.ID,
			&article.Title,
			&article.Slug,
			&article.Content,
			&article.Excerpt,
			&article.FeaturedImage,
			&article.IsPublished,
			&article.UserID,
			&article.CreatedAt,
			&article.UpdatedAt,
			&publishedAt,
		)
		if err != nil {
			return nil, 0, err
		}

		if publishedAt.Valid {
			article.PublishedAt = publishedAt.Time
		}

		articles = append(articles, article)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return articles, total, nil
}
