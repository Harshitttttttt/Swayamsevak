package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// FeedSubsciption represents an a Feed subsctiption in our system
type FeedSubscription struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	FeedID      uuid.UUID
	CustomTitle string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// FeedSubscriptionRepository handles database operations for feed_subscriptions
type FeedSubscriptionRepository struct {
	db *sql.DB
}

// NewFeedSubscriptionRepository creates a new FeedSubscription Repository
func NewFeedSubscriptionRepository(db *sql.DB) *FeedSubscriptionRepository {
	return &FeedSubscriptionRepository{db: db}
}

// SubscribeUserToFeed creates a new Feed Subscription
func (r *FeedSubscriptionRepository) SubscribeUserToFeed(userID, feedID uuid.UUID, customTitle string) (*FeedSubscription, error) {
	feedSubscription := &FeedSubscription{
		UserID:      userID,
		FeedID:      feedID,
		CustomTitle: customTitle,
	}

	query :=
		`
		INSERT INTO feed_subscriptions (user_id, feed_id, custom_title)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at;
	`

	if err := r.db.QueryRow(query, userID, feedID, customTitle).Scan(&feedSubscription.ID, &feedSubscription.CreatedAt, &feedSubscription.UpdatedAt); err != nil {
		return nil, err
	}

	return feedSubscription, nil
}

func (r *FeedSubscriptionRepository) DeleteSubscription(userID, feedID uuid.UUID) error {
	query :=
		`
		DELETE FROM feed_subscriptions
		WHERE user_id = $1 AND feed_id = $2;
	`

	res, err := r.db.Exec(query, userID, feedID)
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

func (r *FeedSubscriptionRepository) GetSubscriptionsByUser(userID uuid.UUID) ([]*FeedSubscription, error) {
	query :=
		`
		SELECT id, user_id, feed_id, custom_title, created_at, updated_at
		FROM feed_subscriptions
		WHERE user_id = $1
		ORDER BY created_at DESC;
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// return an empty slice (not nil) to avoid callers needing to nil-check
	feedSubscriptions := make([]*FeedSubscription, 0)

	for rows.Next() {
		var feedSubscription FeedSubscription
		if err := rows.Scan(&feedSubscription.ID, &feedSubscription.UserID, &feedSubscription.FeedID, &feedSubscription.CustomTitle, &feedSubscription.CreatedAt, &feedSubscription.UpdatedAt); err != nil {
			return nil, err
		}

		feedSubscriptions = append(feedSubscriptions, &feedSubscription)
	}

	// check for errors encountered during iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return feedSubscriptions, nil
}

func (r *FeedSubscriptionRepository) Exists(userID, feedID uuid.UUID) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM feed_subscriptions
			WHERE user_id = $1 AND feed_id = $2
		);
	`

	var exists bool
	err := r.db.QueryRow(query, userID, feedID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
