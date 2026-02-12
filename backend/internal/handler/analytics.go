package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/creatrid/creatrid/internal/middleware"
	"github.com/creatrid/creatrid/internal/store"
	"github.com/go-chi/chi/v5"
)

type AnalyticsHandler struct {
	store *store.Store
}

func NewAnalyticsHandler(st *store.Store) *AnalyticsHandler {
	return &AnalyticsHandler{store: st}
}

// TrackView records a public profile view (called from frontend)
func (h *AnalyticsHandler) TrackView(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	user, err := h.store.FindUserByUsername(r.Context(), strings.ToLower(username))
	if err != nil || user == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "User not found"})
		return
	}

	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.RemoteAddr
	}
	referrer := r.Header.Get("Referer")

	_ = h.store.RecordProfileView(r.Context(), user.ID, ip, referrer)
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

// TrackClick records a link click (called from frontend)
func (h *AnalyticsHandler) TrackClick(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	user, err := h.store.FindUserByUsername(r.Context(), strings.ToLower(username))
	if err != nil || user == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "User not found"})
		return
	}

	var req struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		return
	}

	validTypes := map[string]bool{"connection": true, "custom": true, "share": true}
	if !validTypes[req.Type] {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid click type"})
		return
	}

	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.RemoteAddr
	}

	_ = h.store.RecordLinkClick(r.Context(), user.ID, req.Type, req.Value, ip)
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

// Summary returns analytics for the authenticated user
func (h *AnalyticsHandler) Summary(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}

	summary, err := h.store.GetAnalyticsSummary(r.Context(), user.ID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch analytics"})
		return
	}

	writeJSON(w, http.StatusOK, summary)
}
