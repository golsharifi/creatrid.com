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
	"github.com/creatrid/creatrid/internal/email"
	"github.com/creatrid/creatrid/internal/handler"
	"github.com/creatrid/creatrid/internal/middleware"
	"github.com/creatrid/creatrid/internal/platform"
	"github.com/creatrid/creatrid/internal/scheduler"
	"github.com/creatrid/creatrid/internal/storage"
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

	// Init blob storage (optional)
	var blobStore *storage.BlobStorage
	if cfg.AzureStorageAccount != "" && cfg.AzureStorageKey != "" {
		blobStore, err = storage.NewBlobStorage(cfg.AzureStorageAccount, cfg.AzureStorageKey, cfg.AzureStorageContainer)
		if err != nil {
			log.Printf("Warning: Failed to init blob storage: %v", err)
		} else {
			log.Println("Azure Blob Storage connected")
		}
	} else {
		log.Println("Warning: AZURE_STORAGE_ACCOUNT or AZURE_STORAGE_KEY not set. Image upload will be unavailable.")
	}

	// Init email service (optional)
	var emailSvc *email.Service
	if cfg.SMTPHost != "" && cfg.SMTPUsername != "" {
		emailSvc = email.NewService(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUsername, cfg.SMTPPassword, cfg.SMTPFrom)
		log.Println("Email service configured")
	} else {
		log.Println("Warning: SMTP not configured. Email notifications will be unavailable.")
	}

	// Init SSE hub for real-time notifications
	sseHub := handler.NewSSEHub()

	// Init handlers
	authHandler := handler.NewAuthHandler(googleSvc, jwtSvc, st, cfg)
	userHandler := handler.NewUserHandler(st, blobStore, emailSvc, jwtSvc, cfg)
	connHandler := handler.NewConnectionHandler(st, cfg, emailSvc, providers...)
	ogHandler := handler.NewOGHandler(st, cfg)
	analyticsHandler := handler.NewAnalyticsHandler(st)
	adminHandler := handler.NewAdminHandler(st)
	digestHandler := handler.NewDigestHandler(st, emailSvc)
	collabHandler := handler.NewCollaborationHandler(st, sseHub)
	widgetHandler := handler.NewWidgetHandler(st)
	apiKeyHandler := handler.NewAPIKeyHandler(st)
	verifyHandler := handler.NewVerifyHandler(st)
	billingHandler := handler.NewBillingHandler(st, cfg, sseHub)
	contentHandler := handler.NewContentHandler(st, blobStore, cfg)
	licenseHandler := handler.NewLicenseHandler(st, cfg)
	marketplaceHandler := handler.NewMarketplaceHandler(st)
	dmcaHandler := handler.NewDMCAHandler(st)
	notificationHandler := handler.NewNotificationHandler(st, sseHub)
	contentAnalyticsHandler := handler.NewContentAnalyticsHandler(st)
	payoutHandler := handler.NewPayoutHandler(st, cfg)
	collectionHandler := handler.NewCollectionHandler(st)
	searchHandler := handler.NewSearchHandler(st)
	webhookHandler := handler.NewWebhookHandler(st)
	referralHandler := handler.NewReferralHandler(st)
	recommendHandler := handler.NewRecommendHandler(st)
	moderationHandler := handler.NewModerationHandler(st)
	errorLogHandler := handler.NewErrorLogHandler(st)

	// Start connection refresh scheduler
	providerMap := make(map[string]platform.Provider)
	for _, p := range providers {
		providerMap[p.Name()] = p
	}
	refreshInterval, _ := time.ParseDuration(cfg.RefreshInterval)
	if refreshInterval == 0 {
		refreshInterval = 6 * time.Hour
	}
	sched := scheduler.New(st, providerMap, refreshInterval)
	go sched.Start(context.Background())

	// Start weekly digest cron
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			now := time.Now().UTC()
			if now.Weekday() == time.Monday && now.Hour() == 9 {
				lastSent, _ := st.GetSetting(context.Background(), "last_digest_sent")
				today := now.Format("2006-01-02")
				if lastSent != today {
					sent := digestHandler.RunDigestCron(context.Background())
					_ = st.SetSetting(context.Background(), "last_digest_sent", today)
					log.Printf("Digest cron: sent %d emails", sent)
				}
			}
		}
	}()

	// Start error log cleanup (delete entries older than 30 days)
	go func() {
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			deleted, err := st.DeleteOldErrors(context.Background(), time.Now().Add(-30*24*time.Hour))
			if err != nil {
				log.Printf("Error log cleanup failed: %v", err)
			} else if deleted > 0 {
				log.Printf("Error log cleanup: deleted %d old entries", deleted)
			}
		}
	}()

	// Setup router
	r := chi.NewRouter()
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(chimw.RealIP)
	r.Use(middleware.CORS(cfg.FrontendURL))

	// Error reporting (stricter rate limit)
	r.Group(func(r chi.Router) {
		r.Use(middleware.RateLimit(2, 5))
		r.Post("/api/errors", errorLogHandler.Report)
	})

	// Auth routes (stricter rate limit)
	r.Group(func(r chi.Router) {
		r.Use(middleware.RateLimit(5, 10)) // 5 req/s per IP, burst 10
		r.Get("/api/auth/google", authHandler.GoogleLogin)
		r.Get("/api/auth/google/callback", authHandler.GoogleCallback)
	})

	// Public routes (standard rate limit)
	r.Group(func(r chi.Router) {
		r.Use(middleware.RateLimit(20, 40)) // 20 req/s per IP, burst 40
		r.Get("/api/auth/verify-email/{token}", userHandler.VerifyEmail)
		r.Get("/api/users/{username}", userHandler.PublicProfile)
		r.Get("/api/users/{username}/connections", connHandler.PublicList)
		r.Post("/api/users/{username}/view", analyticsHandler.TrackView)
		r.Post("/api/users/{username}/click", analyticsHandler.TrackClick)
		r.Get("/p/{username}", ogHandler.ProfilePage)
		r.Get("/api/discover", collabHandler.Discover)
		r.Get("/api/widget/{username}", widgetHandler.JSON)
		r.Get("/api/widget/{username}/svg", widgetHandler.SVGBadge)
		r.Get("/api/widget/{username}/html", widgetHandler.HTMLEmbed)
		r.Post("/api/billing/webhook", billingHandler.HandleWebhook)
		r.Get("/api/users/{username}/content", contentHandler.PublicList)
		r.Get("/api/content/{id}/proof", contentHandler.Proof)
		r.Get("/api/content/{id}/licenses", licenseHandler.ListOfferings)
		r.Get("/api/marketplace", marketplaceHandler.Browse)
		r.Get("/api/marketplace/{id}", marketplaceHandler.Detail)
		r.Post("/api/content/{id}/report", dmcaHandler.Report)
		r.Post("/api/content/{id}/view", contentAnalyticsHandler.TrackView)
		r.Get("/api/search", searchHandler.Search)
		r.Get("/api/search/suggestions", searchHandler.Suggestions)
		r.Get("/api/users/{username}/collections", collectionHandler.PublicList)

		// Health check
		r.Get("/api/health", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"status":"ok"}`))
		})
	})

	// Connection OAuth routes (redirect-based auth, stricter rate limit)
	r.Group(func(r chi.Router) {
		r.Use(middleware.RateLimit(5, 10)) // 5 req/s per IP, burst 10
		r.Use(middleware.RequireAuthRedirect(jwtSvc, st, cfg.FrontendURL))
		r.Get("/api/connections/{platform}/connect", connHandler.Connect)
		r.Get("/api/connections/{platform}/callback", connHandler.Callback)
		r.Get("/api/payouts/connect/callback", payoutHandler.ConnectCallback)
	})

	// Protected routes (JSON API)
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireAuth(jwtSvc, st))
		r.Use(middleware.RateLimitUser(10, 20))
		r.Get("/api/auth/me", authHandler.Me)
		r.Post("/api/auth/logout", authHandler.Logout)
		r.Post("/api/auth/verify-email/send", userHandler.SendEmailVerification)
		r.Post("/api/users/onboard", userHandler.Onboard)
		r.Patch("/api/users/profile", userHandler.UpdateProfile)
		r.Post("/api/users/profile/image", userHandler.UploadImage)
		r.Delete("/api/users/account", userHandler.DeleteAccount)
		r.Get("/api/users/export", userHandler.ExportProfile)
		r.Get("/api/connections", connHandler.List)
		r.Delete("/api/connections/{platform}", connHandler.Disconnect)
		r.Post("/api/connections/{platform}/refresh", connHandler.Refresh)
		r.Get("/api/analytics", analyticsHandler.Summary)
		r.Post("/api/collaborations", collabHandler.SendRequest)
		r.Get("/api/collaborations/inbox", collabHandler.Inbox)
		r.Get("/api/collaborations/outbox", collabHandler.Outbox)
		r.Post("/api/collaborations/{id}/respond", collabHandler.Respond)
		r.Post("/api/keys", apiKeyHandler.Create)
		r.Get("/api/keys", apiKeyHandler.List)
		r.Delete("/api/keys/{id}", apiKeyHandler.Delete)
		r.Post("/api/billing/checkout", billingHandler.CreateCheckout)
		r.Get("/api/billing/subscription", billingHandler.GetSubscription)
		r.Post("/api/billing/portal", billingHandler.CreatePortal)

		// Content Vault
		r.Post("/api/content", contentHandler.Upload)
		r.Get("/api/content", contentHandler.List)
		r.Get("/api/content/{id}", contentHandler.Get)
		r.Patch("/api/content/{id}", contentHandler.Update)
		r.Delete("/api/content/{id}", contentHandler.Delete)
		r.Get("/api/content/{id}/download", contentHandler.Download)

		// Licensing
		r.Post("/api/content/{id}/licenses", licenseHandler.CreateOffering)
		r.Patch("/api/licenses/{id}", licenseHandler.UpdateOffering)
		r.Delete("/api/licenses/{id}", licenseHandler.DeleteOffering)
		r.Post("/api/licenses/{id}/checkout", licenseHandler.Checkout)
		r.Get("/api/licenses/purchases", licenseHandler.Purchases)
		r.Get("/api/licenses/sales", licenseHandler.Sales)

		// Notifications
		r.Get("/api/notifications", notificationHandler.List)
		r.Get("/api/notifications/unread-count", notificationHandler.UnreadCount)
		r.Post("/api/notifications/{id}/read", notificationHandler.MarkRead)
		r.Post("/api/notifications/read-all", notificationHandler.MarkAllRead)
		r.Get("/api/notifications/stream", notificationHandler.Stream)

		// Content Analytics
		r.Get("/api/content/{id}/analytics", contentAnalyticsHandler.ItemAnalytics)
		r.Get("/api/content-analytics", contentAnalyticsHandler.CreatorSummary)

		// Payouts / Stripe Connect
		r.Post("/api/payouts/connect", payoutHandler.ConnectOnboard)
		r.Get("/api/payouts/connect/status", payoutHandler.ConnectStatus)
		r.Get("/api/payouts/dashboard", payoutHandler.Dashboard)
		r.Get("/api/payouts", payoutHandler.ListPayouts)

		// Collections
		r.Post("/api/collections", collectionHandler.Create)
		r.Get("/api/collections", collectionHandler.List)
		r.Get("/api/collections/{id}", collectionHandler.Get)
		r.Patch("/api/collections/{id}", collectionHandler.Update)
		r.Delete("/api/collections/{id}", collectionHandler.Delete)
		r.Post("/api/collections/{id}/items", collectionHandler.AddItem)
		r.Delete("/api/collections/{id}/items/{contentId}", collectionHandler.RemoveItem)
		r.Get("/api/collections/{id}/items", collectionHandler.ListItems)

		// Webhooks
		r.Post("/api/webhooks", webhookHandler.Create)
		r.Get("/api/webhooks", webhookHandler.List)
		r.Patch("/api/webhooks/{id}", webhookHandler.Update)
		r.Delete("/api/webhooks/{id}", webhookHandler.Delete)
		r.Get("/api/webhooks/{id}/deliveries", webhookHandler.Deliveries)

		// Referrals
		r.Get("/api/referrals/code", referralHandler.GetCode)
		r.Get("/api/referrals", referralHandler.List)

		// Recommendations
		r.Get("/api/recommendations", recommendHandler.List)

		// Analytics Export
		r.Get("/api/analytics/export", analyticsHandler.ExportCSV)
		r.Get("/api/content-analytics/export", contentAnalyticsHandler.ExportCSV)
	})

	// Admin routes
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireAuth(jwtSvc, st))
		r.Use(handler.RequireAdmin)
		r.Use(middleware.RateLimit(30, 60))
		r.Get("/api/admin/stats", adminHandler.Stats)
		r.Get("/api/admin/users", adminHandler.ListUsers)
		r.Post("/api/admin/users/verify", adminHandler.SetVerified)
		r.Post("/api/admin/digest", digestHandler.SendWeeklyDigest)
		r.Get("/api/admin/audit", adminHandler.AuditLog)
		r.Get("/api/admin/takedowns", dmcaHandler.ListTakedowns)
		r.Post("/api/admin/takedowns/{id}/resolve", dmcaHandler.ResolveTakedown)
		r.Get("/api/admin/errors", errorLogHandler.List)
		r.Get("/api/admin/moderation", moderationHandler.List)
		r.Post("/api/admin/moderation/{id}/resolve", moderationHandler.Resolve)
	})

	// Third-party verification API (API key auth)
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireAPIKey(st))
		r.Get("/api/v1/verify/{username}", verifyHandler.Verify)
		r.Get("/api/v1/verify/{username}/score", verifyHandler.Score)
		r.Get("/api/v1/search", verifyHandler.Search)
	})

	// Start server — WriteTimeout increased to 5 minutes to support SSE
	// (Server-Sent Event) connections which are long-lived.
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 5 * time.Minute,
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
