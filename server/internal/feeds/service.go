package feeds

import (
	"database/sql"
	"errors"

	"github.com/Harshitttttttt/Swayamsevak/server/internal/models"
	"github.com/google/uuid"
)

var (
	ErrFeedAlreadyExists = errors.New("feed already exists")
	ErrFeedNotFound      = errors.New("feed not found")
)

// FeedService provides feed-related functionality
type FeedService struct {
	feedRepo *models.FeedRepository
}

// NewFeedService creates a new feed service
func NewFeedService(feedRepo *models.FeedRepository) *FeedService {
	return &FeedService{
		feedRepo: feedRepo,
	}
}

// AddFeed adds a new feed
func (s *FeedService) AddFeed(feedURL, siteURL, title, description string) (*models.Feed, error) {
	// Check if the feed already exists could be added here
	_, err := s.feedRepo.GetFeedByURL(feedURL)
	if err == nil {
		return nil, ErrFeedAlreadyExists
	}

	// Only proceed if the error was "feed not found"
	if !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	// Create the feed
	feed, err := s.feedRepo.CreateFeed(feedURL, siteURL, title, description)
	if err != nil {
		return nil, err
	}

	return feed, nil
}

// GetFeedByID retrieves a feed by its ID
func (s *FeedService) GetFeedByID(id uuid.UUID) (*models.Feed, error) {
	feed, err := s.feedRepo.GetFeedByID(id)
	if err != nil {
		return nil, err
	}

	return feed, nil
}

// GetFeedByURL retrieves a feed by its URL
func (s *FeedService) GetFeedByURL(feedURL string) (*models.Feed, error) {
	feed, err := s.feedRepo.GetFeedByURL(feedURL)
	if err != nil {
		return nil, err

	}

	return feed, nil
}

// UpdateFeedLastFetchedAt updates the LastFetchedAt timestamp of a feed
func (s *FeedService) UpdateFeedLastFetchedAt(id uuid.UUID) error {
	err := s.feedRepo.UpdateLastFetchedAt(id)
	if err != nil {
		return err
	}

	return nil
}

// ListFeeds retrieves all feeds
func (s *FeedService) ListFeeds() ([]*models.Feed, error) {
	feeds, err := s.feedRepo.GetAllFeeds()
	if err != nil {
		return nil, err
	}

	return feeds, nil
}
