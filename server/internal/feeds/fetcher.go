package feeds

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/mmcdole/gofeed"
)

var (
	ErrFeedUnavailable = errors.New("feed unavailable")
)

type Fetcher struct {
	parser *gofeed.Parser
	client *http.Client
}

func NewFetcher() *Fetcher {
	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	parser := gofeed.NewParser()
	parser.Client = client
	parser.UserAgent = "Swayamsevak/1.0 (+https://github.com/Harshitttttttt/Swayamsevak)"

	return &Fetcher{
		parser: parser,
		client: client,
	}
}

// Fetch receives and parses a feed from the given URL.
func (f *Fetcher) Fetch(ctx context.Context, feedURL string) ([]*gofeed.Item, error) {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	feed, err := f.parser.ParseURLWithContext(feedURL, ctx)
	if err != nil {
		return nil, ErrFeedUnavailable
	}

	if feed == nil || len(feed.Items) == 0 {
		return nil, nil
	}

	return feed.Items, nil
}
