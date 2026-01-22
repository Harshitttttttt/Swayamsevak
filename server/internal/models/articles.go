package models

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
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
	PublishedAt time.Time
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
func (r *ArticleRepository) InsertManyArticlesIgnoreDuplicates(ctx context.Context, articles []*Article) error {
	if len(articles) == 0 {
		return nil
	}

	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Build VALUES clause dynamically
	valueStrings := make([]string, 0, len(articles))
	valueArgs := make([]any, 0, len(articles)*8)

	for i, a := range articles {
		// ($1, $2, $3, ...)
		offset := i * 8
		valueStrings = append(valueStrings,
			fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d)",
				offset+1,
				offset+2,
				offset+3,
				offset+4,
				offset+5,
				offset+6,
				offset+7,
				offset+8,
			),
		)

		valueArgs = append(valueArgs,
			a.FeedID,
			a.GUID,
			a.Title,
			a.URL,
			a.Author,
			a.Content,
			a.Summary,
			a.PublishedAt,
		)
	}

	query := fmt.Sprintf(`
		INSERT INTO articles (
			feed_id,
			guid,
			title,
			url,
			author,
			content,
			summary,
			published_at
		)
		VALUES %s
		ON CONFLICT (feed_id, guid) DO NOTHING
	`, strings.Join(valueStrings, ","))

	if _, err := tx.ExecContext(ctx, query, valueArgs...); err != nil {
		return err
	}

	return tx.Commit()
}

func (r *ArticleRepository) GetUserSubscribedArticles(ctx context.Context, userID uuid.UUID, offset, limit int) ([]*Article, error) {
	query :=
		`
		SELECT a.id, a.feed_id, a.guid, a.title, a.url, a.author, a.content, a.summary, a.published_at, a.created_at, a.updated_at
		FROM articles a
		JOIN feed_subscriptions fs ON a.feed_id = fs.feed_id
		WHERE fs.user_id = $1
		ORDER BY a.published_at DESC
		LIMIT $2
		OFFSET $3;
	`

	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	articles := make([]*Article, 0)
	for rows.Next() {
		var article Article

		if err := rows.Scan(
			&article.ID,
			&article.FeedID,
			&article.GUID,
			&article.Title,
			&article.URL,
			&article.Author,
			&article.Content,
			&article.Summary,
			&article.PublishedAt,
			&article.CreatedAt,
			&article.UpdatedAt,
		); err != nil {
			return nil, err
		}

		articles = append(articles, &article)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return articles, nil
}
