package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DBDSN           string
	JWTSecret       string
	Port            string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
	CookieSecure    bool
}

// LoadEnv() loads environment variables from the .env file
func LoadEnv() *Config {
	// Load from .env if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Check required variables
	requiredVars := []string{"DB_DSN", "JWT_SECRET", "PORT"}
	for _, v := range requiredVars {
		if os.Getenv(v) == "" {
			log.Fatalf("Required Environment variable %s is not set", v)
		}
	}

	// Defaults
	accessTTL := 15 * time.Minute
	refreshTTL := 7 * 24 * time.Hour
	cookieSecure := true

	cfg := &Config{
		DBDSN:           os.Getenv("DB_DSN"),
		JWTSecret:       os.Getenv("JWT_SECRET"),
		Port:            os.Getenv("PORT"),
		AccessTokenTTL:  accessTTL,
		RefreshTokenTTL: refreshTTL,
		CookieSecure:    cookieSecure,
	}

	// Default port
	if cfg.Port == "" {
		cfg.Port = "8080"
	}

	return cfg
}
