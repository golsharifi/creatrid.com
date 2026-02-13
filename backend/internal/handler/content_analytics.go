package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/creatrid/creatrid/internal/middleware"
	"github.com/creatrid/creatrid/internal/store"
	"github.com/go-chi/chi/v5"
)

type ContentAnalyticsHandler struct {
	store *store.Store
}

func NewContentAnalyticsHandler(st *store.Store) *ContentAnalyticsHandler {
	return &ContentAnalyticsHandler{store: st}
}

// TrackView records a public content view.
func (h *ContentAnalyticsHandler) TrackView(w http.ResponseWriter, r *http.Request) {
	contentID := chi.URLParam(r, "id")
	if contentID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Content ID is required"})
		return
	}

	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.RemoteAddr
	}
	referrer := r.Referer()

	_ = h.store.RecordContentView(r.Context(), contentID, ip, referrer)
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

// ItemAnalytics returns analytics for a single content item owned by the user.
func (h *ContentAnalyticsHandler) ItemAnalytics(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}

	contentID := chi.URLParam(r, "id")

	item, err := h.store.FindContentItemByID(r.Context(), contentID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Database error"})
		return
	}
	if item == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "Content not found"})
		return
	}
	if item.UserID != user.ID {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "Not authorized"})
		return
	}

	summary, err := h.store.GetContentAnalytics(r.Context(), contentID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch analytics"})
		return
	}

	writeJSON(w, http.StatusOK, summary)
}

// CreatorSummary returns analytics for all content owned by the authenticated user.
func (h *ContentAnalyticsHandler) CreatorSummary(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}

	items, err := h.store.GetCreatorContentAnalytics(r.Context(), user.ID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch content analytics"})
		return
	}

	totalRevenue := 0
	for _, item := range items {
		totalRevenue += item.RevenueCents
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"items":        items,
		"totalRevenue": totalRevenue,
	})
}

// ExportCSV exports content analytics as CSV.
func (h *ContentAnalyticsHandler) ExportCSV(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}

	items, err := h.store.GetCreatorContentAnalytics(r.Context(), user.ID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch content analytics"})
		return
	}

	username := "creator"
	if user.Username != nil {
		username = *user.Username
	}
	filename := fmt.Sprintf("content-analytics-%s-%s.csv", username, time.Now().Format("2006-01-02"))

	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))

	fmt.Fprintln(w, "content_id,title,views,downloads,revenue_cents")
	for _, item := range items {
		fmt.Fprintf(w, "%s,%q,%d,%d,%d\n", item.ContentID, item.Title, item.TotalViews, item.TotalDownloads, item.RevenueCents)
	}
}
