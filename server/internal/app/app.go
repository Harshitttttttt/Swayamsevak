package app

import (
	"database/sql"
	"time"

	"github.com/Harshitttttttt/Swayamsevak/server/internal/auth"
	"github.com/Harshitttttttt/Swayamsevak/server/internal/models"
)

type App struct {
	UserRepo         *models.UserRepository
	RefreshTokenRepo *models.RefreshTokenRepository
	AuthService      *auth.AuthService
}

func NewApp(db *sql.DB, jwtSecret string) *App {
	userRepo := models.NewUserRepository(db)
	refreshRepo := models.NewRefreshTokenRepository(db)
	authService := auth.NewAuthService(userRepo, refreshRepo, jwtSecret, 15*time.Minute)

	return &App{
		UserRepo:         userRepo,
		RefreshTokenRepo: refreshRepo,
		AuthService:      authService,
	}
}
