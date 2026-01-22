package feeds

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"time"

	"github.com/Harshitttttttt/Swayamsevak/server/internal/models"
	"github.com/google/uuid"
	"github.com/mmcdole/gofeed"
)

func NormalizeItems(items []*gofeed.Item, feedID uuid.UUID) []*models.Article {
	articles := make([]*models.Article, 0, len(items))
	for _, item := range items {
		if item == nil {
			continue
		}

		var authorName string
		if item.Author != nil {
			authorName = item.Author.Name
		}

		article := &models.Article{
			FeedID:      feedID,
			GUID:        resolveGUID(item),
			Title:       sanitize(item.Title),
			URL:         item.Link,
			Author:      authorName,
			Content:     item.Content,
			Summary:     sanitize(item.Description),
			PublishedAt: resolvePublishedAt(item, time.Now()),
		}

		// Skip articles without a valid GUID or URL
		if article.GUID == "" || article.URL == "" {
			continue
		}

		articles = append(articles, article)
	}

	return articles
}

func resolvePublishedAt(item *gofeed.Item, fallback time.Time) time.Time {
	if item.PublishedParsed != nil {
		return *item.PublishedParsed
	}
	if item.UpdatedParsed != nil {
		return *item.UpdatedParsed
	}
	return fallback
}

func resolveGUID(item *gofeed.Item) string {
	if item.GUID != "" {
		return item.GUID
	}

	if item.Link != "" {
		return hash(item.Link)
	}

	return ""
}

func sanitize(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}
	return s
}

func hash(value string) string {
	sum := sha256.Sum256([]byte(value))
	return hex.EncodeToString(sum[:])
}
