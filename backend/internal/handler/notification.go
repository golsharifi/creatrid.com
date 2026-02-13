package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/creatrid/creatrid/internal/middleware"
	"github.com/creatrid/creatrid/internal/store"
	"github.com/go-chi/chi/v5"
)

type NotificationHandler struct {
	store *store.Store
	hub   *SSEHub
}

func NewNotificationHandler(st *store.Store, hub *SSEHub) *NotificationHandler {
	return &NotificationHandler{store: st, hub: hub}
}

// Stream handles GET /api/notifications/stream — Server-Sent Events endpoint
// that pushes real-time notifications to the authenticated user.
func (h *NotificationHandler) Stream(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Streaming not supported"})
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no") // Disable nginx buffering

	ch := h.hub.Subscribe(user.ID)
	defer h.hub.Unsubscribe(user.ID, ch)

	// Send initial connected comment
	fmt.Fprint(w, ": connected\n\n")
	flusher.Flush()

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-r.Context().Done():
			return
		case msg, ok := <-ch:
			if !ok {
				return
			}
			fmt.Fprintf(w, "data: %s\n\n", msg)
			flusher.Flush()
		case <-ticker.C:
			fmt.Fprint(w, ": keepalive\n\n")
			flusher.Flush()
		}
	}
}

// List returns paginated notifications for the authenticated user
func (h *NotificationHandler) List(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	notifications, total, err := h.store.ListNotificationsByUser(r.Context(), user.ID, limit, offset)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch notifications"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"notifications": notifications,
		"total":         total,
	})
}

// UnreadCount returns the number of unread notifications for the authenticated user
func (h *NotificationHandler) UnreadCount(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}

	count, err := h.store.CountUnreadNotifications(r.Context(), user.ID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch unread count"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"count": count,
	})
}

// MarkRead marks a single notification as read
func (h *NotificationHandler) MarkRead(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}

	id := chi.URLParam(r, "id")
	if err := h.store.MarkNotificationRead(r.Context(), id, user.ID); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to mark notification as read"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]bool{"success": true})
}

// MarkAllRead marks all unread notifications as read for the authenticated user
func (h *NotificationHandler) MarkAllRead(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}

	if err := h.store.MarkAllNotificationsRead(r.Context(), user.ID); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to mark notifications as read"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]bool{"success": true})
}
