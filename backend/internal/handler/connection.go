package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/creatrid/creatrid/internal/config"
	"github.com/creatrid/creatrid/internal/middleware"
	"github.com/creatrid/creatrid/internal/model"
	"github.com/creatrid/creatrid/internal/platform"
	"github.com/creatrid/creatrid/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/nrednav/cuid2"
)

type ConnectionHandler struct {
	providers map[string]platform.Provider
	store     *store.Store
	config    *config.Config
}

func NewConnectionHandler(st *store.Store, cfg *config.Config, providers ...platform.Provider) *ConnectionHandler {
	pm := make(map[string]platform.Provider)
	for _, p := range providers {
		pm[p.Name()] = p
	}
	return &ConnectionHandler{
		providers: pm,
		store:     st,
		config:    cfg,
	}
}

func (h *ConnectionHandler) Connect(w http.ResponseWriter, r *http.Request) {
	platformName := chi.URLParam(r, "platform")

	provider, ok := h.providers[platformName]
	if !ok {
		http.Redirect(w, r, h.config.FrontendURL+"/connections?error=unsupported_platform", http.StatusTemporaryRedirect)
		return
	}

	state := randomState()

	sameSite := http.SameSiteLaxMode
	if h.config.CookieSecure {
		sameSite = http.SameSiteNoneMode
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "conn_state",
		Value:    state,
		Path:     "/",
		MaxAge:   300,
		HttpOnly: true,
		Secure:   h.config.CookieSecure,
		SameSite: sameSite,
	})

	http.Redirect(w, r, provider.AuthURL(state), http.StatusTemporaryRedirect)
}

func (h *ConnectionHandler) Callback(w http.ResponseWriter, r *http.Request) {
	platformName := chi.URLParam(r, "platform")

	provider, ok := h.providers[platformName]
	if !ok {
		http.Redirect(w, r, h.config.FrontendURL+"/connections?error=unsupported_platform", http.StatusTemporaryRedirect)
		return
	}

	user := middleware.UserFromContext(r.Context())
	if user == nil {
		http.Redirect(w, r, h.config.FrontendURL+"/sign-in?error=not_authenticated", http.StatusTemporaryRedirect)
		return
	}

	// Validate state
	stateCookie, err := r.Cookie("conn_state")
	if err != nil || stateCookie.Value != r.URL.Query().Get("state") {
		http.Redirect(w, r, h.config.FrontendURL+"/connections?error=invalid_state", http.StatusTemporaryRedirect)
		return
	}

	// Clear state cookie
	http.SetCookie(w, &http.Cookie{
		Name:   "conn_state",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	// Check for OAuth error (user denied)
	if errParam := r.URL.Query().Get("error"); errParam != "" {
		http.Redirect(w, r, h.config.FrontendURL+"/connections?error="+errParam, http.StatusTemporaryRedirect)
		return
	}

	// Exchange code
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Redirect(w, r, h.config.FrontendURL+"/connections?error=no_code", http.StatusTemporaryRedirect)
		return
	}

	token, err := provider.Exchange(r.Context(), code)
	if err != nil {
		log.Printf("Connection exchange error (%s): %v", platformName, err)
		http.Redirect(w, r, h.config.FrontendURL+"/connections?error=exchange_failed", http.StatusTemporaryRedirect)
		return
	}

	// Fetch profile
	profile, err := provider.FetchProfile(r.Context(), token)
	if err != nil {
		log.Printf("Connection fetch profile error (%s): %v", platformName, err)
		http.Redirect(w, r, h.config.FrontendURL+"/connections?error=fetch_failed", http.StatusTemporaryRedirect)
		return
	}

	// Marshal metadata
	metadataJSON, _ := json.Marshal(profile.Metadata)

	conn := &model.Connection{
		ID:             cuid2.Generate(),
		UserID:         user.ID,
		Platform:       platformName,
		PlatformUserID: profile.PlatformUserID,
		Username:       strPtr(profile.Username),
		DisplayName:    strPtr(profile.DisplayName),
		AvatarURL:      strPtr(profile.AvatarURL),
		ProfileURL:     strPtr(profile.ProfileURL),
		FollowerCount:  &profile.FollowerCount,
		AccessToken:    strPtr(token.AccessToken),
		RefreshToken:   strPtr(token.RefreshToken),
		Metadata:       metadataJSON,
	}

	if !token.Expiry.IsZero() {
		conn.TokenExpiresAt = &token.Expiry
	}

	if err := h.store.UpsertConnection(r.Context(), conn); err != nil {
		log.Printf("Connection upsert error (%s): %v", platformName, err)
		http.Redirect(w, r, h.config.FrontendURL+"/connections?error=save_failed", http.StatusTemporaryRedirect)
		return
	}

	recalcScore(r.Context(), h.store, user.ID)
	http.Redirect(w, r, h.config.FrontendURL+"/connections?connected="+platformName, http.StatusTemporaryRedirect)
}

func (h *ConnectionHandler) List(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}

	connections, err := h.store.FindConnectionsByUserID(r.Context(), user.ID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch connections"})
		return
	}

	public := make([]*model.PublicConnection, 0, len(connections))
	for _, c := range connections {
		public = append(public, c.ToPublic())
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{"connections": public})
}

func (h *ConnectionHandler) Disconnect(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}

	platformName := chi.URLParam(r, "platform")

	if err := h.store.DeleteConnection(r.Context(), user.ID, platformName); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to disconnect"})
		return
	}

	recalcScore(r.Context(), h.store, user.ID)
	writeJSON(w, http.StatusOK, map[string]bool{"success": true})
}

func (h *ConnectionHandler) PublicList(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	user, err := h.store.FindUserByUsername(r.Context(), strings.ToLower(username))
	if err != nil || user == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "User not found"})
		return
	}

	connections, err := h.store.FindConnectionsByUserID(r.Context(), user.ID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch connections"})
		return
	}

	public := make([]*model.PublicConnection, 0, len(connections))
	for _, c := range connections {
		public = append(public, c.ToPublic())
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{"connections": public})
}
