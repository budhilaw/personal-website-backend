package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/budhilaw/personal-website-backend/internal/model"
	"github.com/budhilaw/personal-website-backend/internal/util"
)

// PortfolioRepository defines methods for portfolio repository
type PortfolioRepository interface {
	Create(ctx context.Context, portfolio *model.PortfolioCreate, userID string) (string, error)
	Update(ctx context.Context, id string, portfolio *model.PortfolioUpdate) error
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*model.Portfolio, error)
	GetBySlug(ctx context.Context, slug string) (*model.Portfolio, error)
	List(ctx context.Context, page, perPage int, onlyPublished bool) ([]model.Portfolio, int, error)
	GetByAuthor(ctx context.Context, userID string, page, perPage int) ([]model.Portfolio, int, error)
}

// portfolioRepository is the implementation of PortfolioRepository
type portfolioRepository struct {
	db *sql.DB
}

// NewPortfolioRepository creates a new PortfolioRepository
func NewPortfolioRepository(db *sql.DB) PortfolioRepository {
	return &portfolioRepository{db: db}
}

// Create creates a new portfolio
func (r *portfolioRepository) Create(ctx context.Context, portfolioCreate *model.PortfolioCreate, userID string) (string, error) {
	query := `INSERT INTO portfolios (title, slug, description, image, project_url, github_url, technologies, is_published, user_id) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) 
			  RETURNING id`

	slug := util.GenerateSlug(portfolioCreate.Title)

	// Convert technologies slice to JSON
	var technologiesJSON []byte
	var err error
	if len(portfolioCreate.Technologies) > 0 {
		technologiesJSON, err = json.Marshal(portfolioCreate.Technologies)
		if err != nil {
			return "", err
		}
	}

	var id string
	err = r.db.QueryRowContext(
		ctx, query,
		portfolioCreate.Title,
		slug,
		portfolioCreate.Description,
		portfolioCreate.Image,
		portfolioCreate.ProjectURL,
		portfolioCreate.GithubURL,
		technologiesJSON,
		portfolioCreate.IsPublished,
		userID,
	).Scan(&id)
	if err != nil {
		return "", err
	}

	return id, nil
}

// Update updates a portfolio
func (r *portfolioRepository) Update(ctx context.Context, id string, portfolioUpdate *model.PortfolioUpdate) error {
	query := `UPDATE portfolios 
			  SET title = $2, slug = $3, description = $4, image = $5, project_url = $6, github_url = $7, technologies = $8, is_published = $9, updated_at = $10
			  WHERE id = $1`

	// Convert technologies slice to JSON
	var technologiesJSON []byte
	var err error
	if len(portfolioUpdate.Technologies) > 0 {
		technologiesJSON, err = json.Marshal(portfolioUpdate.Technologies)
		if err != nil {
			return err
		}
	}

	_, err = r.db.ExecContext(
		ctx, query,
		id,
		portfolioUpdate.Title,
		util.GenerateSlug(portfolioUpdate.Title),
		portfolioUpdate.Description,
		portfolioUpdate.Image,
		portfolioUpdate.ProjectURL,
		portfolioUpdate.GithubURL,
		technologiesJSON,
		portfolioUpdate.IsPublished,
		time.Now(),
	)
	return err
}

// Delete deletes a portfolio
func (r *portfolioRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM portfolios WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// GetByID gets a portfolio by ID
func (r *portfolioRepository) GetByID(ctx context.Context, id string) (*model.Portfolio, error) {
	query := `SELECT id, title, slug, description, image, project_url, github_url, technologies, is_published, user_id, created_at, updated_at 
			  FROM portfolios 
			  WHERE id = $1`

	var portfolio model.Portfolio
	var technologiesJSON sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&portfolio.ID,
		&portfolio.Title,
		&portfolio.Slug,
		&portfolio.Description,
		&portfolio.Image,
		&portfolio.ProjectURL,
		&portfolio.GithubURL,
		&technologiesJSON,
		&portfolio.IsPublished,
		&portfolio.UserID,
		&portfolio.CreatedAt,
		&portfolio.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("portfolio not found")
		}
		return nil, err
	}

	if technologiesJSON.Valid {
		portfolio.Technologies = json.RawMessage(technologiesJSON.String)
	}

	return &portfolio, nil
}

// GetBySlug gets a portfolio by slug
func (r *portfolioRepository) GetBySlug(ctx context.Context, slug string) (*model.Portfolio, error) {
	query := `SELECT id, title, slug, description, image, project_url, github_url, technologies, is_published, user_id, created_at, updated_at 
			  FROM portfolios 
			  WHERE slug = $1`

	var portfolio model.Portfolio
	var technologiesJSON sql.NullString

	err := r.db.QueryRowContext(ctx, query, slug).Scan(
		&portfolio.ID,
		&portfolio.Title,
		&portfolio.Slug,
		&portfolio.Description,
		&portfolio.Image,
		&portfolio.ProjectURL,
		&portfolio.GithubURL,
		&technologiesJSON,
		&portfolio.IsPublished,
		&portfolio.UserID,
		&portfolio.CreatedAt,
		&portfolio.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("portfolio not found")
		}
		return nil, err
	}

	if technologiesJSON.Valid {
		portfolio.Technologies = json.RawMessage(technologiesJSON.String)
	}

	return &portfolio, nil
}

// List lists portfolios with pagination
func (r *portfolioRepository) List(ctx context.Context, page, perPage int, onlyPublished bool) ([]model.Portfolio, int, error) {
	offset := (page - 1) * perPage

	// Count total
	countQuery := `SELECT COUNT(*) FROM portfolios`
	if onlyPublished {
		countQuery += ` WHERE is_published = true`
	}

	var total int
	err := r.db.QueryRowContext(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get portfolios
	query := `SELECT id, title, slug, description, image, project_url, github_url, technologies, is_published, user_id, created_at, updated_at 
			  FROM portfolios`
	if onlyPublished {
		query += ` WHERE is_published = true`
	}
	query += ` ORDER BY created_at DESC LIMIT $1 OFFSET $2`

	rows, err := r.db.QueryContext(ctx, query, perPage, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var portfolios []model.Portfolio
	for rows.Next() {
		var portfolio model.Portfolio
		var technologiesJSON sql.NullString

		err := rows.Scan(
			&portfolio.ID,
			&portfolio.Title,
			&portfolio.Slug,
			&portfolio.Description,
			&portfolio.Image,
			&portfolio.ProjectURL,
			&portfolio.GithubURL,
			&technologiesJSON,
			&portfolio.IsPublished,
			&portfolio.UserID,
			&portfolio.CreatedAt,
			&portfolio.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}

		if technologiesJSON.Valid {
			portfolio.Technologies = json.RawMessage(technologiesJSON.String)
		}

		portfolios = append(portfolios, portfolio)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return portfolios, total, nil
}

// GetByAuthor gets portfolios by author ID with pagination
func (r *portfolioRepository) GetByAuthor(ctx context.Context, userID string, page, perPage int) ([]model.Portfolio, int, error) {
	offset := (page - 1) * perPage

	// Count total
	countQuery := `SELECT COUNT(*) FROM portfolios WHERE user_id = $1`
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, userID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get portfolios
	query := `SELECT id, title, slug, description, image, project_url, github_url, technologies, is_published, user_id, created_at, updated_at 
			  FROM portfolios 
			  WHERE user_id = $1 
			  ORDER BY created_at DESC 
			  LIMIT $2 OFFSET $3`

	rows, err := r.db.QueryContext(ctx, query, userID, perPage, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var portfolios []model.Portfolio
	for rows.Next() {
		var portfolio model.Portfolio
		var technologiesJSON sql.NullString

		err := rows.Scan(
			&portfolio.ID,
			&portfolio.Title,
			&portfolio.Slug,
			&portfolio.Description,
			&portfolio.Image,
			&portfolio.ProjectURL,
			&portfolio.GithubURL,
			&technologiesJSON,
			&portfolio.IsPublished,
			&portfolio.UserID,
			&portfolio.CreatedAt,
			&portfolio.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}

		if technologiesJSON.Valid {
			portfolio.Technologies = json.RawMessage(technologiesJSON.String)
		}

		portfolios = append(portfolios, portfolio)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return portfolios, total, nil
}
