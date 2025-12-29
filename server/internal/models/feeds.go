package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// Feed represents a feed in our system
type Feed struct {
	ID            uuid.UUID
	FeedURL       string
	SiteURL       string
	Title         string
	Description   string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	LastFetchedAt sql.NullTime
}

// FeedRepository handles database operations for feeds
type FeedRepository struct {
	db *sql.DB
}

// NewFeedRepository creates a new feed repository
func NewFeedRepository(db *sql.DB) *FeedRepository {
	return &FeedRepository{
		db: db,
	}
}

// CreateFeed adds a new feed to the database
func (r *FeedRepository) CreateFeed(feedURL, siteURL, title, description string) (*Feed, error) {
	feed := &Feed{
		FeedURL:     feedURL,
		SiteURL:     siteURL,
		Title:       title,
		Description: description,
	}

	query :=
		`
		INSERT INTO feeds (feed_url, site_url, title, description)
		VALUES ($1, $2, $3, $4) 
		RETURNING id, created_at, updated_at, last_fetched_at;
	`

	if err := r.db.QueryRow(query, feed.FeedURL, feed.SiteURL, feed.Title, feed.Description).Scan(&feed.ID, &feed.CreatedAt, &feed.UpdatedAt, &feed.LastFetchedAt); err != nil {
		return nil, err
	}

	return feed, nil
}

// GetAllFeeds retrieves all feeds from the database
func (r *FeedRepository) GetAllFeeds() ([]*Feed, error) {
	query :=
		`
		SELECT id, feed_url, site_url, title, description, created_at, updated_at, last_fetched_at
		FROM feeds;
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var feeds []*Feed

	for rows.Next() {
		var feed Feed
		if err := rows.Scan(&feed.ID, &feed.FeedURL, &feed.SiteURL, &feed.Title, &feed.Description, &feed.CreatedAt, &feed.UpdatedAt, &feed.LastFetchedAt); err != nil {
			return nil, err
		}
		feeds = append(feeds, &feed)
	}

	return feeds, nil
}

// GetFeedByID retrieves a feed by its ID
func (r *FeedRepository) GetFeedByID(id uuid.UUID) (*Feed, error) {
	query :=
		`
		SELECT id, feed_url, site_url, title, description, created_at, updated_at, last_fetched_at
		FROM feeds
		WHERE id = $1;
	`

	var feed Feed
	if err := r.db.QueryRow(query, id).Scan(&feed.ID, &feed.FeedURL, &feed.SiteURL, &feed.Title, &feed.Description, &feed.CreatedAt, &feed.UpdatedAt, &feed.LastFetchedAt); err != nil {
		return nil, err
	}

	return &feed, nil
}

// UpdateLastFetchedAt updates the LastFetchedAt timestamp of a feed
func (r *FeedRepository) UpdateLastFetchedAt(id uuid.UUID) error {
	query :=
		`
		UPDATE feeds
		SET last_fetched_at = now(), updated_at = now()
		WHERE id = $1;
	`

	res, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// GetFeedByURL gets a feed by its unique URL
func (r *FeedRepository) GetFeedByURL(feedURL string) (*Feed, error) {
	query :=
		`
	SELECT id, feed_url, site_url, title, description, created_at, updated_at, last_fetched_at
	FROM feeds
	WHERE feed_url = $1;
	`

	var feed Feed
	if err := r.db.QueryRow(query, feedURL).Scan(&feed.ID, &feed.FeedURL, &feed.SiteURL, feed.Title, feed.Description, feed.CreatedAt, feed.UpdatedAt, feed.LastFetchedAt); err != nil {
		return nil, err
	}

	return &feed, nil
}
