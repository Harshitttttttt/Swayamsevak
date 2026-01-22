// @title           Swayamsevak API
// @version         1.0
// @description     RSS Aggregator API with user authentication
// @termsOfService  http://swagger.io/terms/

// @contact.name   Harshit Mestry
// @contact.url    https://github.com/Harshitttttttt/Swayamsevak
// @contact.email  harshitsmestry@gmail.com

// @license.name  MIT
// @license.url   https://github.com/Harshitttttttt/Swayamsevak/blob/main/LICENSE

// @host      localhost:8080
// @BasePath  /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/Harshitttttttt/Swayamsevak/server/docs"
	"github.com/Harshitttttttt/Swayamsevak/server/internal/app"
	"github.com/Harshitttttttt/Swayamsevak/server/internal/config"
	"github.com/Harshitttttttt/Swayamsevak/server/internal/database"
	"github.com/Harshitttttttt/Swayamsevak/server/internal/feeds"
	"github.com/Harshitttttttt/Swayamsevak/server/internal/handlers"
	"github.com/Harshitttttttt/Swayamsevak/server/internal/middleware"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func main() {
	// Load config
	cfg := config.LoadEnv()

	// Connect to the database
	db, err := database.Connect(cfg.DBDSN)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	// Create the App
	app := app.NewApp(db, cfg.JWTSecret, cfg.AccessTokenTTL)

	// Start the worker to scrape feeds
	feedWorker := feeds.NewWorker(app.FeedService, 10*time.Second, 10)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Run worker in a separate goroutine
	go feedWorker.Start(ctx)

	// Create the handlers
	authHandler := handlers.NewAuthHandler(app.AuthService, cfg.RefreshTokenTTL, cfg.CookieSecure)
	userHandler := handlers.NewUserHandler(app.UserRepo)
	feedHandler := handlers.NewFeedHandler(app.FeedService)

	mux := http.NewServeMux()

	// Swagger documentation route
	mux.HandleFunc("GET /swagger/",
		httpSwagger.Handler(httpSwagger.URL("/swagger/doc.json")),
	)

	// Create the routes
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	// Public Routes

	// Auth Routes
	mux.HandleFunc("POST /api/auth/register", authHandler.Register)
	mux.HandleFunc("POST /api/auth/login", authHandler.Login)
	mux.HandleFunc("POST /api/auth/refresh", authHandler.RefreshToken)
	mux.HandleFunc("POST /api/auth/logout", authHandler.Logout)

	// User Routes
	protectedProfile := middleware.AuthMiddleware(app.AuthService)(http.HandlerFunc(userHandler.Profile))
	mux.Handle("GET /api/profile", protectedProfile)

	// Feed Routes
	protectedFeedCreation := middleware.AuthMiddleware(app.AuthService)(http.HandlerFunc(feedHandler.AddFeedHandler))
	mux.Handle("POST /api/feed", protectedFeedCreation)

	protectedGetAllFeeds := middleware.AuthMiddleware(app.AuthService)(http.HandlerFunc(feedHandler.GetAllFeedsHandler))
	mux.Handle("GET /api/feeds", protectedGetAllFeeds)

	protectedSubscribeFeed := middleware.AuthMiddleware(app.AuthService)(http.HandlerFunc(feedHandler.SubscribeToFeedHandler))
	mux.Handle("POST /api/feed/subscribe", protectedSubscribeFeed)

	protectedGetArticlesForSubscribedFeeds := middleware.AuthMiddleware(app.AuthService)(http.HandlerFunc(feedHandler.GetUserArticlesHandler))
	mux.Handle("GET /api/feed/articles", protectedGetArticlesForSubscribedFeeds)

	// Apply CORS
	handler := enableCORS(mux)

	fmt.Printf("Server running on \x1b[91mhttp://localhost:%s\x1b[0m\n", cfg.Port)
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())

}

// enableCORS sets the necessary headers to allow React frontend communication
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set headers
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173") // Adjust port for your React app
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
