package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/creatrid/creatrid/internal/analytics"
	"github.com/creatrid/creatrid/internal/geoip"
	"github.com/creatrid/creatrid/internal/middleware"
	"github.com/creatrid/creatrid/internal/store"
	"github.com/go-chi/chi/v5"
)

type AnalyticsHandler struct {
	store  *store.Store
	geoSvc *geoip.Service
}

func NewAnalyticsHandler(st *store.Store, geoSvc *geoip.Service) *AnalyticsHandler {
	return &AnalyticsHandler{store: st, geoSvc: geoSvc}
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
	ua := analytics.ParseUserAgent(r.Header.Get("User-Agent"))

	var country, city string
	if h.geoSvc != nil {
		geo := h.geoSvc.Lookup(ip)
		country = geo.Country
		city = geo.City
	}

	_ = h.store.RecordProfileView(r.Context(), user.ID, ip, referrer, ua.Browser, ua.OS, ua.DeviceType, country, city)

	// Dispatch webhook event for profile view
	dispatchWebhook(user.ID, "profile.viewed", map[string]interface{}{
		"username": username,
		"referrer": referrer,
	})

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

// ExportCSV exports profile analytics as CSV.
func (h *AnalyticsHandler) ExportCSV(w http.ResponseWriter, r *http.Request) {
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

	username := "creator"
	if user.Username != nil {
		username = *user.Username
	}
	filename := fmt.Sprintf("analytics-%s-%s.csv", username, time.Now().Format("2006-01-02"))

	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))

	// Write CSV header and data
	fmt.Fprintln(w, "date,views,clicks,engagement_rate")
	for i, ds := range summary.ViewsByDay {
		clicks := 0
		if i < len(summary.ClicksByDay) {
			clicks = summary.ClicksByDay[i].Count
		}
		rate := 0.0
		if ds.Count > 0 {
			rate = float64(clicks) / float64(ds.Count) * 100
		}
		fmt.Fprintf(w, "%s,%d,%d,%.1f\n", ds.Date, ds.Count, clicks, rate)
	}
}
