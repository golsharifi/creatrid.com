package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/creatrid/creatrid/internal/store"
)

// statusCapture wraps http.ResponseWriter to capture the status code
type statusCapture struct {
	http.ResponseWriter
	statusCode int
}

func (w *statusCapture) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

// RecordAPIUsage records API usage for requests authenticated via API key
func RecordAPIUsage(st *store.Store) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Only record if this request was authenticated via API key
			// The apiKeyUser context is set by RequireAPIKey middleware
			apiKeyID := apiKeyIDFromContext(r.Context())
			if apiKeyID == "" {
				next.ServeHTTP(w, r)
				return
			}

			start := time.Now()
			sc := &statusCapture{ResponseWriter: w, statusCode: http.StatusOK}
			next.ServeHTTP(sc, r)

			// Record usage asynchronously
			elapsed := time.Since(start).Milliseconds()
			go func() {
				_ = st.RecordAPIUsage(
					context.Background(),
					apiKeyID,
					r.URL.Path,
					r.Method,
					sc.statusCode,
					int(elapsed),
				)
			}()
		})
	}
}

// We need to store the API key ID (not just user ID) in context for usage tracking.
// The existing RequireAPIKey only stores the user ID.
// We'll add an extra context key for the API key ID.

const apiKeyIDContextKey contextKey = "apiKeyID"

// SetAPIKeyIDContext adds the API key ID to the context. Called by RequireAPIKey.
func SetAPIKeyIDContext(ctx context.Context, keyID string) context.Context {
	return context.WithValue(ctx, apiKeyIDContextKey, keyID)
}

func apiKeyIDFromContext(ctx context.Context) string {
	id, _ := ctx.Value(apiKeyIDContextKey).(string)
	return id
}
