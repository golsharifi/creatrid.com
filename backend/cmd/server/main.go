package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"sort"
	"syscall"
	"time"

	"github.com/creatrid/creatrid/internal/auth"
	"github.com/creatrid/creatrid/internal/config"
	"github.com/creatrid/creatrid/internal/handler"
	"github.com/creatrid/creatrid/internal/middleware"
	"github.com/creatrid/creatrid/internal/platform"
	"github.com/creatrid/creatrid/internal/store"
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// Load config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to database
	pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(context.Background()); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Connected to database")

	// Run migrations
	if err := runMigrations(cfg.DatabaseURL); err != nil {
		log.Printf("Migration warning: %v", err)
	}

	// Init services
	st := store.New(pool)
	jwtSvc := auth.NewJWTService(cfg.JWTSecret)
	googleSvc := auth.NewGoogleService(cfg.GoogleClientID, cfg.GoogleSecret, cfg.GoogleRedirect)

	// Init platform providers
	var providers []platform.Provider
	providers = append(providers, platform.NewYouTubeProvider(cfg.GoogleClientID, cfg.GoogleSecret, cfg.YouTubeRedirect))
	if cfg.GitHubClientID != "" && cfg.GitHubSecret != "" {
		providers = append(providers, platform.NewGitHubProvider(cfg.GitHubClientID, cfg.GitHubSecret, cfg.GitHubRedirect))
	}
	if cfg.TwitterClientID != "" && cfg.TwitterSecret != "" {
		providers = append(providers, platform.NewTwitterProvider(cfg.TwitterClientID, cfg.TwitterSecret, cfg.TwitterRedirect))
	}
	if cfg.LinkedInClientID != "" && cfg.LinkedInSecret != "" {
		providers = append(providers, platform.NewLinkedInProvider(cfg.LinkedInClientID, cfg.LinkedInSecret, cfg.LinkedInRedirect))
	}
	if cfg.InstagramClientID != "" && cfg.InstagramSecret != "" {
		providers = append(providers, platform.NewInstagramProvider(cfg.InstagramClientID, cfg.InstagramSecret, cfg.InstagramRedirect))
	}
	if cfg.DribbbleClientID != "" && cfg.DribbbleSecret != "" {
		providers = append(providers, platform.NewDribbbleProvider(cfg.DribbbleClientID, cfg.DribbbleSecret, cfg.DribbbleRedirect))
	}
	if cfg.BehanceClientID != "" && cfg.BehanceSecret != "" {
		providers = append(providers, platform.NewBehanceProvider(cfg.BehanceClientID, cfg.BehanceSecret, cfg.BehanceRedirect))
	}

	// Init handlers
	authHandler := handler.NewAuthHandler(googleSvc, jwtSvc, st, cfg)
	userHandler := handler.NewUserHandler(st)
	connHandler := handler.NewConnectionHandler(st, cfg, providers...)
	ogHandler := handler.NewOGHandler(st, cfg)

	// Setup router
	r := chi.NewRouter()
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(chimw.RealIP)
	r.Use(middleware.CORS(cfg.FrontendURL))

	// Public routes
	r.Get("/api/auth/google", authHandler.GoogleLogin)
	r.Get("/api/auth/google/callback", authHandler.GoogleCallback)
	r.Get("/api/users/{username}", userHandler.PublicProfile)
	r.Get("/api/users/{username}/connections", connHandler.PublicList)
	r.Get("/p/{username}", ogHandler.ProfilePage)

	// Health check
	r.Get("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	// Connection OAuth routes (redirect-based auth)
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireAuthRedirect(jwtSvc, st, cfg.FrontendURL))
		r.Get("/api/connections/{platform}/connect", connHandler.Connect)
		r.Get("/api/connections/{platform}/callback", connHandler.Callback)
	})

	// Protected routes (JSON API)
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireAuth(jwtSvc, st))
		r.Get("/api/auth/me", authHandler.Me)
		r.Post("/api/auth/logout", authHandler.Logout)
		r.Post("/api/users/onboard", userHandler.Onboard)
		r.Patch("/api/users/profile", userHandler.UpdateProfile)
		r.Get("/api/connections", connHandler.List)
		r.Delete("/api/connections/{platform}", connHandler.Disconnect)
	})

	// Start server
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	go func() {
		log.Printf("Server starting on :%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
	log.Println("Server stopped")
}

func runMigrations(databaseURL string) error {
	pool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		return err
	}
	defer pool.Close()

	files, err := filepath.Glob("migrations/*.up.sql")
	if err != nil {
		return err
	}
	sort.Strings(files)

	for _, f := range files {
		data, err := os.ReadFile(f)
		if err != nil {
			return err
		}
		_, err = pool.Exec(context.Background(), string(data))
		if err != nil {
			log.Printf("Migration note (%s): %v (may already be applied)", f, err)
		}
	}

	log.Println("Migrations applied successfully")
	return nil
}
