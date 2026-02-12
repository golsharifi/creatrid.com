package config

import (
	"fmt"
	"log"
	"os"
)

type Config struct {
	Port           string
	DatabaseURL    string
	JWTSecret      string
	GoogleClientID string
	GoogleSecret   string
	GoogleRedirect string
	BackendURL     string
	FrontendURL    string
	CookieDomain    string
	CookieSecure    bool
	GitHubClientID  string
	GitHubSecret    string
	YouTubeRedirect string
	GitHubRedirect  string
}

func Load() (*Config, error) {
	cfg := &Config{
		Port:           getEnv("PORT", "8080"),
		DatabaseURL:    os.Getenv("DATABASE_URL"),
		JWTSecret:      os.Getenv("JWT_SECRET"),
		GoogleClientID: os.Getenv("GOOGLE_CLIENT_ID"),
		GoogleSecret:   os.Getenv("GOOGLE_CLIENT_SECRET"),
		BackendURL:     getEnv("BACKEND_URL", "http://localhost:8080"),
		FrontendURL:    getEnv("FRONTEND_URL", "http://localhost:3000"),
		CookieDomain:   os.Getenv("COOKIE_DOMAIN"),
		CookieSecure:   os.Getenv("COOKIE_SECURE") == "true",
		GitHubClientID: os.Getenv("GITHUB_CLIENT_ID"),
		GitHubSecret:   os.Getenv("GITHUB_CLIENT_SECRET"),
	}

	cfg.GoogleRedirect = cfg.BackendURL + "/api/auth/google/callback"
	cfg.YouTubeRedirect = cfg.BackendURL + "/api/connections/youtube/callback"
	cfg.GitHubRedirect = cfg.BackendURL + "/api/connections/github/callback"

	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}
	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}
	if cfg.GoogleClientID == "" {
		return nil, fmt.Errorf("GOOGLE_CLIENT_ID is required")
	}
	if cfg.GoogleSecret == "" {
		return nil, fmt.Errorf("GOOGLE_CLIENT_SECRET is required")
	}

	if cfg.GitHubClientID == "" || cfg.GitHubSecret == "" {
		log.Println("Warning: GITHUB_CLIENT_ID or GITHUB_CLIENT_SECRET not set. GitHub connections will be unavailable.")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
