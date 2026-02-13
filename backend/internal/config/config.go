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

	TwitterClientID  string
	TwitterSecret    string
	TwitterRedirect  string
	LinkedInClientID string
	LinkedInSecret   string
	LinkedInRedirect string
	InstagramClientID string
	InstagramSecret   string
	InstagramRedirect string
	DribbbleClientID string
	DribbbleSecret   string
	DribbbleRedirect string
	BehanceClientID  string
	BehanceSecret    string
	BehanceRedirect  string

	AzureStorageAccount          string
	AzureStorageKey              string
	AzureStorageContainer        string
	AzureStorageContentContainer string

	SMTPHost     string
	SMTPPort     string
	SMTPUsername string
	SMTPPassword string
	SMTPFrom     string

	RefreshInterval  string
	MaxMindDBPath    string

	StripeSecretKey     string
	StripeWebhookSecret string
	StripePricePro      string
	StripePriceBusiness string

	BlockchainRPCURL     string
	BlockchainPrivateKey string
	BlockchainChainID    string
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

		TwitterClientID:  os.Getenv("TWITTER_CLIENT_ID"),
		TwitterSecret:    os.Getenv("TWITTER_CLIENT_SECRET"),
		LinkedInClientID: os.Getenv("LINKEDIN_CLIENT_ID"),
		LinkedInSecret:   os.Getenv("LINKEDIN_CLIENT_SECRET"),
		InstagramClientID: os.Getenv("INSTAGRAM_CLIENT_ID"),
		InstagramSecret:   os.Getenv("INSTAGRAM_CLIENT_SECRET"),
		DribbbleClientID: os.Getenv("DRIBBBLE_CLIENT_ID"),
		DribbbleSecret:   os.Getenv("DRIBBBLE_CLIENT_SECRET"),
		BehanceClientID:  os.Getenv("BEHANCE_CLIENT_ID"),
		BehanceSecret:    os.Getenv("BEHANCE_CLIENT_SECRET"),

		AzureStorageAccount:   os.Getenv("AZURE_STORAGE_ACCOUNT"),
		AzureStorageKey:       os.Getenv("AZURE_STORAGE_KEY"),
		AzureStorageContainer:        getEnv("AZURE_STORAGE_CONTAINER", "avatars"),
		AzureStorageContentContainer: getEnv("AZURE_STORAGE_CONTENT_CONTAINER", "vault"),

		SMTPHost:     os.Getenv("SMTP_HOST"),
		SMTPPort:     getEnv("SMTP_PORT", "587"),
		SMTPUsername: os.Getenv("SMTP_USERNAME"),
		SMTPPassword: os.Getenv("SMTP_PASSWORD"),
		SMTPFrom:     getEnv("SMTP_FROM", "info@creatrid.com"),

		RefreshInterval: getEnv("REFRESH_INTERVAL", "6h"),
		MaxMindDBPath:   os.Getenv("MAXMIND_DB_PATH"),

		StripeSecretKey:     os.Getenv("STRIPE_SECRET_KEY"),
		StripeWebhookSecret: os.Getenv("STRIPE_WEBHOOK_SECRET"),
		StripePricePro:      os.Getenv("STRIPE_PRICE_PRO"),
		StripePriceBusiness: os.Getenv("STRIPE_PRICE_BUSINESS"),

		BlockchainRPCURL:     os.Getenv("BLOCKCHAIN_RPC_URL"),
		BlockchainPrivateKey: os.Getenv("BLOCKCHAIN_PRIVATE_KEY"),
		BlockchainChainID:    getEnv("BLOCKCHAIN_CHAIN_ID", "137"),
	}

	cfg.GoogleRedirect = cfg.BackendURL + "/api/auth/google/callback"
	cfg.YouTubeRedirect = cfg.BackendURL + "/api/connections/youtube/callback"
	cfg.GitHubRedirect = cfg.BackendURL + "/api/connections/github/callback"
	cfg.TwitterRedirect = cfg.BackendURL + "/api/connections/twitter/callback"
	cfg.LinkedInRedirect = cfg.BackendURL + "/api/connections/linkedin/callback"
	cfg.InstagramRedirect = cfg.BackendURL + "/api/connections/instagram/callback"
	cfg.DribbbleRedirect = cfg.BackendURL + "/api/connections/dribbble/callback"
	cfg.BehanceRedirect = cfg.BackendURL + "/api/connections/behance/callback"

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
