package dto

import (
	"time"

	"github.com/google/uuid"
)

// AddFeedRequest represents the payload to add a new feed
type AddFeedRequest struct {
	FeedURL     string `json:"feed_url" example:"https://example.com/feed"`
	SiteURL     string `json:"site_url" example:"https://example.com"`
	Title       string `json:"title" example:"Example Feed"`
	Description string `json:"description" example:"This is an example feed description."`
}

// AddFeedResponse represents the response after adding a new feed
type AddFeedResponse struct {
	ID uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
}

// FeedResponse represents the feed details in responses
type FeedResponse struct {
	ID            uuid.UUID  `json:"id"`
	FeedURL       string     `json:"feed_url"`
	SiteURL       string     `json:"site_url"`
	Title         string     `json:"title"`
	Description   string     `json:"description"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	LastFetchedAt *time.Time `json:"last_fetched_at,omitempty"`
}

// ListFeedsResponse represents the response for listing all feeds
type ListFeedsResponse struct {
	Feeds []FeedResponse `json:"feeds"`
}

// SubscribeFeedRequest represents the payload to subscribe to a feed
type SubscribeFeedRequest struct {
	FeedID      uuid.UUID `json:"feed_id" example:"17b3a6f1-1617-4104-b914-fffba0236bd9"`
	CustomTitle string    `json:"custom_title" example:"My RSS Feed"`
}

// SubscribeFeedResponse represents the response after subscribing to a feed
type SubscribeFeedResponse struct {
	Message string `json:"message" example:"Successfully subscribed to the feed"`
}
