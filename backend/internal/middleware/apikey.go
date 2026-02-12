package middleware

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/creatrid/creatrid/internal/store"
)

const apiKeyUserContextKey contextKey = "apiKeyUser"

type apiKeyRateLimiter struct {
	mu      sync.Mutex
	counts  map[string]*apiKeyCounter
}

type apiKeyCounter struct {
	count    int
	windowAt time.Time
}

func newAPIKeyRateLimiter() *apiKeyRateLimiter {
	rl := &apiKeyRateLimiter{
		counts: make(map[string]*apiKeyCounter),
	}
	go rl.cleanup()
	return rl
}

func (rl *apiKeyRateLimiter) cleanup() {
	for {
		time.Sleep(time.Minute)
		rl.mu.Lock()
		now := time.Now()
		for id, c := range rl.counts {
			if now.Sub(c.windowAt) > 2*time.Minute {
				delete(rl.counts, id)
			}
		}
		rl.mu.Unlock()
	}
}

func (rl *apiKeyRateLimiter) allow(keyID string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	c, exists := rl.counts[keyID]

	if !exists || now.Sub(c.windowAt) >= time.Minute {
		rl.counts[keyID] = &apiKeyCounter{count: 1, windowAt: now}
		return true
	}

	if c.count >= 100 {
		return false
	}

	c.count++
	return true
}

func RequireAPIKey(st *store.Store) func(http.Handler) http.Handler {
	limiter := newAPIKeyRateLimiter()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Missing API key"})
				return
			}

			rawKey := strings.TrimPrefix(authHeader, "Bearer ")
			if rawKey == "" {
				writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Invalid API key"})
				return
			}

			// Hash the key with SHA-256
			hash := sha256.Sum256([]byte(rawKey))
			keyHash := hex.EncodeToString(hash[:])

			apiKey, err := st.FindAPIKeyByHash(r.Context(), keyHash)
			if err != nil {
				writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
				return
			}
			if apiKey == nil {
				writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Invalid API key"})
				return
			}

			// Per-key rate limiting
			if !limiter.allow(apiKey.ID) {
				w.Header().Set("Retry-After", "60")
				writeJSON(w, http.StatusTooManyRequests, map[string]string{"error": "Rate limit exceeded (100 requests per minute)"})
				return
			}

			// Update last_used_at asynchronously
			go func() {
				_ = st.UpdateAPIKeyLastUsed(context.Background(), apiKey.ID, time.Now())
			}()

			ctx := context.WithValue(r.Context(), apiKeyUserContextKey, apiKey.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func APIKeyUserFromContext(ctx context.Context) string {
	userID, _ := ctx.Value(apiKeyUserContextKey).(string)
	return userID
}
