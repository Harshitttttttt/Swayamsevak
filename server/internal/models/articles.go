package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// Article represents an article in our system
type Article struct {
	ID          uuid.UUID
	FeedID      uuid.UUID
	GUID        string
	Title       string
	URL         string
	Author      string
	Content     string
	Summary     string
	PublishedAt sql.NullTime
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// ArticleRepository handles database operations for articles
type ArticleRepository struct {
	db *sql.DB
}

// NewArticleRepository creates a new article repository
func NewArticleRepository(db *sql.DB) *ArticleRepository {
	return &ArticleRepository{
		db: db,
	}
}

// CreateArticle adds a new article to the database
func (r *ArticleRepository) CreateArticle()
