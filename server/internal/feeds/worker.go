package feeds

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/Harshitttttttt/Swayamsevak/server/internal/models"
)

type Worker struct {
	feedService *FeedService
	interval    time.Duration
	concurrency int
}

func NewWorker(feedService *FeedService, interval time.Duration, concurrency int) *Worker {
	return &Worker{
		feedService: feedService,
		interval:    interval,
		concurrency: concurrency,
	}
}

func (w *Worker) Start(ctx context.Context) {
	log.Printf("Scraping on %v goroutines every %v\n duration", w.concurrency, w.interval)
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Feed Worker Shutting Down")
			return

		case <-ticker.C:
			w.runOnce(ctx)
		}
	}
}

func (w *Worker) runOnce(ctx context.Context) {
	feeds, err := w.feedService.GetNextFeedsToFetch(ctx, w.concurrency, w.interval)
	if err != nil {
		log.Println("Error fetching feeds: ", err)
		return
	}

	sem := make(chan struct{}, w.concurrency)
	var wg sync.WaitGroup

	for _, feed := range feeds {
		wg.Add(1)
		sem <- struct{}{}

		go func(f *models.Feed) {
			defer wg.Done()
			defer func() {
				<-sem
			}()

			if err := w.feedService.FetchAndStoreFeed(ctx, f); err != nil {
				log.Printf("Error processing feed: %v, with link: %s", err, f.FeedURL)
			} else {
				log.Printf("Successfully fetched and stored feed: %s", f.FeedURL)
			}
		}(feed)
	}

	wg.Wait()
}
