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
	"fmt"
	"log"
	"net/http"

	_ "github.com/Harshitttttttt/Swayamsevak/server/docs"
	"github.com/Harshitttttttt/Swayamsevak/server/internal/app"
	"github.com/Harshitttttttt/Swayamsevak/server/internal/config"
	"github.com/Harshitttttttt/Swayamsevak/server/internal/database"
	"github.com/Harshitttttttt/Swayamsevak/server/internal/handlers"
	"github.com/Harshitttttttt/Swayamsevak/server/internal/middleware"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func main() {
	// Load config
	cfg := config.LoadEnv()

	// Connect to the database
	database, err := database.Connect(cfg.DBDSN)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	// Create the App
	app := app.NewApp(database, cfg.JWTSecret)

	// Create the handlers
	authHandler := handlers.NewAuthHandler(app.AuthService)
	userHandler := handlers.NewUserHandler(app.UserRepo)

	mux := http.NewServeMux()

	// Swagger documentation route
	mux.HandleFunc("GET /swagger/",
		httpSwagger.Handler(httpSwagger.URL("http://localhost:8080/swagger/doc.json")),
	)

	// Create the routes
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	// Public Routes
	mux.HandleFunc("POST /api/auth/register", authHandler.Register)
	mux.HandleFunc("POST /api/auth/login", authHandler.Login)
	mux.HandleFunc("POST /api/auth/refresh", authHandler.RefreshToken)

	// Protected Routes
	protectedProfile := middleware.AuthMiddleware(app.AuthService)(
		http.HandlerFunc(userHandler.Profile),
	)

	mux.Handle("GET /api/profile", protectedProfile)

	fmt.Printf("Server running on \x1b[91mhttp://localhost:%s\x1b[0m\n", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, mux))
}
