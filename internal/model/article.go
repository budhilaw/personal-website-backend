package model

import (
	"time"
)

type Article struct {
	ID            string    `json:"id"`
	Title         string    `json:"title"`
	Slug          string    `json:"slug"`
	Content       string    `json:"content"`
	Excerpt       string    `json:"excerpt,omitempty"`
	FeaturedImage string    `json:"featured_image,omitempty"`
	IsPublished   bool      `json:"is_published"`
	UserID        string    `json:"user_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	PublishedAt   time.Time `json:"published_at,omitempty"`
}

// ArticleCreate represents article creation request body
type ArticleCreate struct {
	Title         string `json:"title" validate:"required"`
	Content       string `json:"content" validate:"required"`
	Excerpt       string `json:"excerpt"`
	FeaturedImage string `json:"featured_image"`
	IsPublished   bool   `json:"is_published"`
}

// ArticleUpdate represents article update request body
type ArticleUpdate struct {
	Title         string `json:"title" validate:"required"`
	Content       string `json:"content" validate:"required"`
	Excerpt       string `json:"excerpt"`
	FeaturedImage string `json:"featured_image"`
	IsPublished   bool   `json:"is_published"`
}

// ArticleResponse represents article response with author information
type ArticleResponse struct {
	ID            string `json:"id"`
	Title         string `json:"title"`
	Slug          string `json:"slug"`
	Content       string `json:"content"`
	Excerpt       string `json:"excerpt,omitempty"`
	FeaturedImage string `json:"featured_image,omitempty"`
	IsPublished   bool   `json:"is_published"`
	Author        struct {
		ID        string `json:"id"`
		Username  string `json:"username"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name,omitempty"`
		Avatar    string `json:"avatar,omitempty"`
	} `json:"author"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	PublishedAt time.Time `json:"published_at,omitempty"`
}

// ArticleList represents a list of articles with pagination
type ArticleList struct {
	Articles []ArticleResponse `json:"articles"`
	Total    int               `json:"total"`
	Page     int               `json:"page"`
	PerPage  int               `json:"per_page"`
}
