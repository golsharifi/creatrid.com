package scheduler

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/creatrid/creatrid/internal/platform"
	"github.com/creatrid/creatrid/internal/score"
	"github.com/creatrid/creatrid/internal/store"
	"golang.org/x/oauth2"
)

type Scheduler struct {
	store     *store.Store
	providers map[string]platform.Provider
	interval  time.Duration
}

func New(st *store.Store, providers map[string]platform.Provider, interval time.Duration) *Scheduler {
	return &Scheduler{
		store:     st,
		providers: providers,
		interval:  interval,
	}
}

func (s *Scheduler) Start(ctx context.Context) {
	log.Printf("Connection refresh scheduler started (interval: %s)", s.interval)
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Connection refresh scheduler stopped")
			return
		case <-ticker.C:
			s.refreshAll(ctx)
		}
	}
}

func (s *Scheduler) refreshAll(ctx context.Context) {
	connections, err := s.store.FindStaleConnections(ctx, 24*time.Hour, 50)
	if err != nil {
		log.Printf("Scheduler: failed to find stale connections: %v", err)
		return
	}

	if len(connections) == 0 {
		return
	}

	log.Printf("Scheduler: refreshing %d stale connections", len(connections))
	refreshed := 0

	for _, conn := range connections {
		provider, ok := s.providers[conn.Platform]
		if !ok {
			continue
		}

		// Try to refresh token
		var accessToken string
		if conn.RefreshToken != nil && *conn.RefreshToken != "" {
			newToken, err := provider.RefreshToken(ctx, *conn.RefreshToken)
			if err != nil {
				log.Printf("Scheduler: token refresh failed for %s/%s: %v", conn.UserID, conn.Platform, err)
				// Touch updated_at to avoid retrying too soon
				_ = s.store.UpdateConnectionProfile(ctx, conn.UserID, conn.Platform, conn.FollowerCount, conn.Metadata)
				continue
			}
			accessToken = newToken.AccessToken

			if !newToken.Expiry.IsZero() {
				exp := newToken.Expiry
				_ = s.store.UpdateConnectionTokens(ctx, conn.UserID, conn.Platform, newToken.AccessToken, newToken.RefreshToken, &exp)
			} else {
				_ = s.store.UpdateConnectionTokens(ctx, conn.UserID, conn.Platform, newToken.AccessToken, newToken.RefreshToken, nil)
			}
		} else if conn.AccessToken != nil {
			accessToken = *conn.AccessToken
		} else {
			continue
		}

		// Re-fetch profile
		profile, err := provider.FetchProfile(ctx, &oauth2.Token{AccessToken: accessToken})
		if err != nil {
			log.Printf("Scheduler: profile fetch failed for %s/%s: %v", conn.UserID, conn.Platform, err)
			continue
		}

		metadataJSON, _ := json.Marshal(profile.Metadata)
		_ = s.store.UpdateConnectionProfile(ctx, conn.UserID, conn.Platform, &profile.FollowerCount, metadataJSON)

		// Recalculate score
		user, err := s.store.FindUserByID(ctx, conn.UserID)
		if err != nil || user == nil {
			continue
		}
		allConns, err := s.store.FindConnectionsByUserID(ctx, conn.UserID)
		if err != nil {
			continue
		}
		newScore := score.Calculate(user, allConns)
		_ = s.store.UpdateUserScore(ctx, conn.UserID, newScore)

		refreshed++
	}

	if refreshed > 0 {
		log.Printf("Scheduler: refreshed %d connections", refreshed)
	}
}
