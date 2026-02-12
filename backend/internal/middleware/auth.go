package middleware

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/creatrid/creatrid/internal/auth"
	"github.com/creatrid/creatrid/internal/model"
	"github.com/creatrid/creatrid/internal/store"
)

type contextKey string

const userContextKey contextKey = "user"

func RequireAuth(jwtSvc *auth.JWTService, st *store.Store) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("token")
			if err != nil {
				writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
				return
			}

			userID, err := jwtSvc.Validate(cookie.Value)
			if err != nil {
				writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
				return
			}

			user, err := st.FindUserByID(r.Context(), userID)
			if err != nil || user == nil {
				writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
				return
			}

			ctx := context.WithValue(r.Context(), userContextKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequireAuthRedirect(jwtSvc *auth.JWTService, st *store.Store, frontendURL string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("token")
			if err != nil {
				http.Redirect(w, r, frontendURL+"/sign-in?error=not_authenticated", http.StatusTemporaryRedirect)
				return
			}

			userID, err := jwtSvc.Validate(cookie.Value)
			if err != nil {
				http.Redirect(w, r, frontendURL+"/sign-in?error=session_expired", http.StatusTemporaryRedirect)
				return
			}

			user, err := st.FindUserByID(r.Context(), userID)
			if err != nil || user == nil {
				http.Redirect(w, r, frontendURL+"/sign-in?error=not_authenticated", http.StatusTemporaryRedirect)
				return
			}

			ctx := context.WithValue(r.Context(), userContextKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func UserFromContext(ctx context.Context) *model.User {
	user, _ := ctx.Value(userContextKey).(*model.User)
	return user
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
