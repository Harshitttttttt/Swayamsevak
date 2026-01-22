package feeds

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Harshitttttttt/Swayamsevak/server/internal/models"
	"github.com/google/uuid"
)

var (
	ErrFeedAlreadyExists = errors.New("feed already exists")
	ErrFeedNotFound      = errors.New("feed not found")
)

// FeedService provides feed-related functionality
type FeedService struct {
	feedRepo             *models.FeedRepository
	feedSubscriptionRepo *models.FeedSubscriptionRepository
	articleRepo          *models.ArticleRepository

	fetcher *Fetcher
}

// NewFeedService creates a new feed service
func NewFeedService(feedRepo *models.FeedRepository, feedSubscriptionRepo *models.FeedSubscriptionRepository, articleRepo *models.ArticleRepository, fetcher *Fetcher) *FeedService {
	return &FeedService{
		feedRepo:             feedRepo,
		feedSubscriptionRepo: feedSubscriptionRepo,
		articleRepo:          articleRepo,
		fetcher:              fetcher,
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

// SubscribeToFeed allows a user to subscribe to a feed
func (s *FeedService) SubscribeToFeed(userID, feedID uuid.UUID, customTitle string) error {
	// Get the feed from the db
	feed, err := s.feedRepo.GetFeedByID(feedID)
	if err != nil {
		return err
	}

	// Handle customTitle being optional
	if customTitle == "" {
		customTitle = feed.Title
	}

	_, err = s.feedSubscriptionRepo.SubscribeUserToFeed(userID, feedID, customTitle)
	if err != nil {
		return err
	}

	return nil
}

// FetchAndStoreFeed fetches a feed and stores its articles
func (s *FeedService) FetchAndStoreFeed(ctx context.Context, feed *models.Feed) error {
	// Claim First
	if err := s.feedRepo.UpdateLastFetchedAt(feed.ID); err != nil {
		return err
	}

	rss, err := s.fetcher.Fetch(ctx, feed.FeedURL)
	if err != nil {
		return err
	}

	articles := NormalizeItems(rss, feed.ID)

	return s.articleRepo.InsertManyArticlesIgnoreDuplicates(ctx, articles)
}

// GetNextFeedToFetch retrieves the next feed that needs to be fetched
func (s *FeedService) GetNextFeedsToFetch(ctx context.Context, limit int, olderThan time.Duration) ([]*models.Feed, error) {
	feeds, err := s.feedRepo.GetNextFeedsToFetch(ctx, limit, olderThan)
	if err != nil {
		return nil, err
	}

	return feeds, nil
}

// FetchUsersSubscribedFeeds is a business logic function that retrieves a users articles for their subscribed feeds
func (s *FeedService) FetchUserSubscribedFeeds(ctx context.Context, userID uuid.UUID, offset, limit int) ([]*models.Article, error) {
	articles, err := s.articleRepo.GetUserSubscribedArticles(ctx, userID, offset, limit)
	if err != nil {
		return nil, err
	}

	return articles, err
}
