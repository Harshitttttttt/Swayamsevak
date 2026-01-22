package app

import (
	"database/sql"
	"time"

	"github.com/Harshitttttttt/Swayamsevak/server/internal/auth"
	"github.com/Harshitttttttt/Swayamsevak/server/internal/feeds"
	"github.com/Harshitttttttt/Swayamsevak/server/internal/models"
)

type App struct {
	UserRepo         *models.UserRepository
	RefreshTokenRepo *models.RefreshTokenRepository

	AuthService *auth.AuthService

	FeedRepo    *models.FeedRepository
	FeedService *feeds.FeedService

	ArticleRepo *models.ArticleRepository
}

func NewApp(db *sql.DB, jwtSecret string, accessTokenTTL time.Duration) *App {
	userRepo := models.NewUserRepository(db)
	refreshRepo := models.NewRefreshTokenRepository(db)
	authService := auth.NewAuthService(userRepo, refreshRepo, jwtSecret, accessTokenTTL)

	feedRepo := models.NewFeedRepository(db)
	feedSubscriptionRepo := models.NewFeedSubscriptionRepository(db)
	articleRepo := models.NewArticleRepository(db)

	fetcher := feeds.NewFetcher()

	feedService := feeds.NewFeedService(feedRepo, feedSubscriptionRepo, articleRepo, fetcher)

	return &App{
		UserRepo:         userRepo,
		RefreshTokenRepo: refreshRepo,
		AuthService:      authService,
		FeedRepo:         feedRepo,
		FeedService:      feedService,
		ArticleRepo:      articleRepo,
	}
}
