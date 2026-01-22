package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/Harshitttttttt/Swayamsevak/server/internal/feeds"
	"github.com/Harshitttttttt/Swayamsevak/server/internal/handlers/dto"
	"github.com/Harshitttttttt/Swayamsevak/server/internal/middleware"
	"github.com/google/uuid"
)

// FeedHandler contains HTTP handlers for feed-related operations
type FeedHandler struct {
	feedService *feeds.FeedService
}

// NewFeedHandler creates a new Feed handler
func NewFeedHandler(feedService *feeds.FeedService) *FeedHandler {
	return &FeedHandler{
		feedService: feedService,
	}
}

// AddFeedHandler godoc
// @Summary      Add a new RSS feed
// @Description  Register a new RSS/Atom feed URL in the system for aggregation
// @Tags         Feeds
// @Accept       json
// @Produce      json
// @Param        request body dto.AddFeedRequest true "Feed registration details"
// @Security     BearerAuth
// @Success      201 {object} dto.AddFeedResponse "Feed successfully registered"
// @Failure      400 {string} string "Invalid Request Body"
// @Failure      409 {string} string "Feed already exists, aborting"
// @Failure      500 {string} string "Internal Server Error"
// @Router       /feed [post]
func (h *FeedHandler) AddFeedHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// Parse the request body
	var req dto.AddFeedRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	// Validate Input
	if req.FeedURL == "" || req.SiteURL == "" || req.Title == "" || req.Description == "" {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}

	// Call the FeedService to add the feed
	feed, err := h.feedService.AddFeed(req.FeedURL, req.SiteURL, req.Title, req.Description)
	if err != nil {
		if errors.Is(err, feeds.ErrFeedAlreadyExists) {
			http.Error(w, "Feed already exists, aborting", http.StatusConflict)
			return
		}

		http.Error(w, "Error Creating Feed", http.StatusInternalServerError)
		return
	}

	// Respond with the created feed
	response := &dto.AddFeedResponse{
		ID: feed.ID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// GetAllFeedsHandler godoc
// @Summary      Get all RSS feeds
// @Description  Retrieve all registered RSS/Atom feeds in the system
// @Tags         Feeds
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} dto.ListFeedsResponse "List of all registered feeds"
// @Failure      500 {string} string "Internal Server Error"
// @Router       /feeds [get]
func (h *FeedHandler) GetAllFeedsHandler(w http.ResponseWriter, r *http.Request) {
	feeds, err := h.feedService.ListFeeds()
	if err != nil {
		http.Error(w, "Failed to fetch feeds", http.StatusInternalServerError)
		return
	}

	response := make([]dto.FeedResponse, 0, len(feeds))
	for _, feed := range feeds {
		var lastFetched *time.Time
		if feed.LastFetchedAt.Valid {
			t := feed.LastFetchedAt.Time
			lastFetched = &t
		}

		response = append(response, dto.FeedResponse{
			ID:            feed.ID,
			FeedURL:       feed.FeedURL,
			SiteURL:       feed.SiteURL,
			Title:         feed.Title,
			Description:   feed.Description,
			CreatedAt:     feed.CreatedAt,
			UpdatedAt:     feed.UpdatedAt,
			LastFetchedAt: lastFetched,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(dto.ListFeedsResponse{
		Feeds: response,
	})
}

// SubscribeToFeedHandler godoc
// @Summary      Subscribe to an RSS feed
// @Description  Subscribe the authenticated user to a specific RSS/Atom feed
// @Tags         Feeds
// @Accept			 json
// @Produce      json
// @Param        request body dto.SubscribeFeedRequest true "Subscription details"
// @Security     BearerAuth
// @Success      200 {object} dto.SubscribeFeedResponse "Successfully subscribed to the feed"
// @Failure      400 {string} string "Invalid Request Body"
// @Failure      500 {string} string "Internal Server Error"
// @Router       /feed/subscribe [post]
func (h *FeedHandler) SubscribeToFeedHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// Get the user ID from the context (set by authentication middleware)
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse the request body
	var req dto.SubscribeFeedRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	// Validate Input
	if req.FeedID == uuid.Nil {
		http.Error(w, "UserID and FeedID are required", http.StatusBadRequest)
		return
	}

	// Call the FeedService to subscribe to the feed
	err := h.feedService.SubscribeToFeed(userID, req.FeedID, req.CustomTitle)
	if err != nil {
		http.Error(w, "Error Subscribing to Feed", http.StatusInternalServerError)
		return
	}

	// Respond with success message
	response := &dto.SubscribeFeedResponse{
		Message: "Successfully subscribed to the feed",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// FetchUserArticlesHandler godoc
// @Summary      Fetch articles for user subscribed feeds
// @Description  Fetch the articles for feeds subscribed to by a user
// @Tags         Feeds
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} dto.GetUserArticlesResponse "Successfully fetched all articles"
// @Failure      400 {string} string "Invalid Request Body"
// @Failure      500 {string} string "Internal Server Error"
// @Router       /feed/articles [get]
func (h *FeedHandler) GetUserArticlesHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// Get the userID from the context
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Manage offset and limit here
	articles, err := h.feedService.FetchUserSubscribedFeeds(r.Context(), userID, 0, 100)
	if err != nil {
		http.Error(w, "Error fetching articles: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := make([]dto.ArticlesResponse, 0, len(articles))
	for _, article := range articles {
		response = append(response, dto.ArticlesResponse{
			ID:          article.ID,
			FeedID:      article.FeedID,
			GUID:        article.GUID,
			Title:       article.Title,
			URL:         article.URL,
			Author:      article.Author,
			Content:     article.Content,
			Summary:     article.Summary,
			PublishedAt: article.PublishedAt,
			CreatedAt:   article.CreatedAt,
			UpdatedAt:   article.UpdatedAt,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(dto.GetUserArticlesResponse{
		Articles: response,
	})
}
