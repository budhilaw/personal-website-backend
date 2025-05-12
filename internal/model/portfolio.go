package model

import (
	"encoding/json"
	"time"
)

type Portfolio struct {
	ID           string          `json:"id"`
	Title        string          `json:"title"`
	Slug         string          `json:"slug"`
	Description  string          `json:"description"`
	Image        string          `json:"image,omitempty"`
	ProjectURL   string          `json:"project_url,omitempty"`
	GithubURL    string          `json:"github_url,omitempty"`
	Technologies json.RawMessage `json:"technologies,omitempty"`
	IsPublished  bool            `json:"is_published"`
	UserID       string          `json:"user_id"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}

// PortfolioCreate represents portfolio creation request body
type PortfolioCreate struct {
	Title        string   `json:"title" validate:"required"`
	Description  string   `json:"description" validate:"required"`
	Image        string   `json:"image"`
	ProjectURL   string   `json:"project_url"`
	GithubURL    string   `json:"github_url"`
	Technologies []string `json:"technologies"`
	IsPublished  bool     `json:"is_published"`
}

// PortfolioUpdate represents portfolio update request body
type PortfolioUpdate struct {
	Title        string   `json:"title" validate:"required"`
	Description  string   `json:"description" validate:"required"`
	Image        string   `json:"image"`
	ProjectURL   string   `json:"project_url"`
	GithubURL    string   `json:"github_url"`
	Technologies []string `json:"technologies"`
	IsPublished  bool     `json:"is_published"`
}

// PortfolioResponse represents portfolio response with author information
type PortfolioResponse struct {
	ID           string          `json:"id"`
	Title        string          `json:"title"`
	Slug         string          `json:"slug"`
	Description  string          `json:"description"`
	Image        string          `json:"image,omitempty"`
	ProjectURL   string          `json:"project_url,omitempty"`
	GithubURL    string          `json:"github_url,omitempty"`
	Technologies json.RawMessage `json:"technologies,omitempty"`
	IsPublished  bool            `json:"is_published"`
	Author       struct {
		ID        string `json:"id"`
		Username  string `json:"username"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name,omitempty"`
		Avatar    string `json:"avatar,omitempty"`
	} `json:"author"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// PortfolioList represents a list of portfolios with pagination
type PortfolioList struct {
	Portfolios []PortfolioResponse `json:"portfolios"`
	Total      int                 `json:"total"`
	Page       int                 `json:"page"`
	PerPage    int                 `json:"per_page"`
}
