package handler

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/creatrid/creatrid/internal/auth"
	"github.com/creatrid/creatrid/internal/config"
	"github.com/creatrid/creatrid/internal/middleware"
	"github.com/creatrid/creatrid/internal/model"
	"github.com/creatrid/creatrid/internal/store"
	"github.com/nrednav/cuid2"
)

type AuthHandler struct {
	google *auth.GoogleService
	jwt    *auth.JWTService
	store  *store.Store
	config *config.Config
}

func NewAuthHandler(google *auth.GoogleService, jwt *auth.JWTService, store *store.Store, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		google: google,
		jwt:    jwt,
		store:  store,
		config: cfg,
	}
}

func (h *AuthHandler) GoogleLogin(w http.ResponseWriter, r *http.Request) {
	state := randomState()

	sameSite := http.SameSiteLaxMode
	if h.config.CookieSecure {
		sameSite = http.SameSiteNoneMode
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		Path:     "/",
		MaxAge:   300,
		HttpOnly: true,
		Secure:   h.config.CookieSecure,
		SameSite: sameSite,
	})

	http.Redirect(w, r, h.google.AuthURL(state), http.StatusTemporaryRedirect)
}

func (h *AuthHandler) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	// Validate state
	stateCookie, err := r.Cookie("oauth_state")
	if err != nil || stateCookie.Value != r.URL.Query().Get("state") {
		http.Redirect(w, r, h.config.FrontendURL+"/sign-in?error=invalid_state", http.StatusTemporaryRedirect)
		return
	}

	// Clear state cookie
	http.SetCookie(w, &http.Cookie{
		Name:   "oauth_state",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	// Exchange code for token
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Redirect(w, r, h.config.FrontendURL+"/sign-in?error=no_code", http.StatusTemporaryRedirect)
		return
	}

	token, err := h.google.Exchange(r.Context(), code)
	if err != nil {
		log.Printf("OAuth exchange error: %v", err)
		http.Redirect(w, r, h.config.FrontendURL+"/sign-in?error=exchange_failed", http.StatusTemporaryRedirect)
		return
	}

	// Get Google user info
	googleUser, err := h.google.UserInfo(r.Context(), token)
	if err != nil {
		log.Printf("Google userinfo error: %v", err)
		http.Redirect(w, r, h.config.FrontendURL+"/sign-in?error=userinfo_failed", http.StatusTemporaryRedirect)
		return
	}

	// Find or create user
	account, err := h.store.FindAccountByProvider(r.Context(), "google", googleUser.ID)
	if err != nil {
		log.Printf("DB error finding account: %v", err)
		http.Redirect(w, r, h.config.FrontendURL+"/sign-in?error=db_error", http.StatusTemporaryRedirect)
		return
	}

	var user *model.User

	if account != nil {
		// Existing user
		user, err = h.store.FindUserByID(r.Context(), account.UserID)
		if err != nil || user == nil {
			log.Printf("DB error finding user: %v", err)
			http.Redirect(w, r, h.config.FrontendURL+"/sign-in?error=db_error", http.StatusTemporaryRedirect)
			return
		}
	} else {
		// Check if user exists by email
		user, err = h.store.FindUserByEmail(r.Context(), googleUser.Email)
		if err != nil {
			log.Printf("DB error finding user by email: %v", err)
			http.Redirect(w, r, h.config.FrontendURL+"/sign-in?error=db_error", http.StatusTemporaryRedirect)
			return
		}

		if user == nil {
			// Create new user
			now := time.Now()
			user = &model.User{
				ID:    cuid2.Generate(),
				Name:  &googleUser.Name,
				Email: googleUser.Email,
				Image: &googleUser.Picture,
				Role:  model.RoleCreator,
			}
			if googleUser.VerifiedEmail {
				user.EmailVerified = &now
			}

			if err := h.store.CreateUser(r.Context(), user); err != nil {
				log.Printf("DB error creating user: %v", err)
				http.Redirect(w, r, h.config.FrontendURL+"/sign-in?error=db_error", http.StatusTemporaryRedirect)
				return
			}
		}
	}

	// Upsert account
	expiresAt := int(token.Expiry.Unix())
	idToken, _ := token.Extra("id_token").(string)
	acct := &model.Account{
		UserID:            user.ID,
		Type:              "oauth",
		Provider:          "google",
		ProviderAccountID: googleUser.ID,
		AccessToken:       &token.AccessToken,
		RefreshToken:      &token.RefreshToken,
		ExpiresAt:         &expiresAt,
		TokenType:         strPtr(token.TokenType),
		IDToken:           strPtr(idToken),
	}
	if err := h.store.UpsertAccount(r.Context(), acct); err != nil {
		log.Printf("DB error upserting account: %v", err)
	}

	// Check if user has 2FA enabled
	if user.TOTPEnabled {
		// Generate temp JWT (5 min, 2fa:pending)
		tempToken, err := h.jwt.GenerateTempToken(user.ID)
		if err != nil {
			log.Printf("Temp JWT generation error: %v", err)
			http.Redirect(w, r, h.config.FrontendURL+"/sign-in?error=token_error", http.StatusTemporaryRedirect)
			return
		}
		h.jwt.SetTempTokenCookie(w, tempToken, h.config.CookieDomain, h.config.CookieSecure)
		http.Redirect(w, r, h.config.FrontendURL+"/sign-in/2fa", http.StatusTemporaryRedirect)
		return
	}

	// Generate JWT
	jwtToken, err := h.jwt.Generate(user.ID)
	if err != nil {
		log.Printf("JWT generation error: %v", err)
		http.Redirect(w, r, h.config.FrontendURL+"/sign-in?error=token_error", http.StatusTemporaryRedirect)
		return
	}

	h.jwt.SetTokenCookie(w, jwtToken, h.config.CookieDomain, h.config.CookieSecure)

	// Redirect based on onboarding status
	if user.Onboarded {
		http.Redirect(w, r, h.config.FrontendURL+"/dashboard", http.StatusTemporaryRedirect)
	} else {
		http.Redirect(w, r, h.config.FrontendURL+"/onboarding", http.StatusTemporaryRedirect)
	}
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"user": user})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	h.jwt.ClearTokenCookie(w, h.config.CookieDomain, h.config.CookieSecure)
	writeJSON(w, http.StatusOK, map[string]bool{"success": true})
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func randomState() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
