package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/creatrid/creatrid/internal/middleware"
	"github.com/creatrid/creatrid/internal/store"
	"github.com/go-chi/chi/v5"
)

type ModerationHandler struct {
	store *store.Store
}

func NewModerationHandler(st *store.Store) *ModerationHandler {
	return &ModerationHandler{store: st}
}

// List handles GET /api/admin/moderation — returns paginated moderation flags.
// Query params: status (optional), limit (default 20), offset (default 0).
func (h *ModerationHandler) List(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	flags, total, err := h.store.ListModerationFlags(r.Context(), status, limit, offset)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch moderation flags"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"flags": flags,
		"total": total,
	})
}

// Resolve handles POST /api/admin/moderation/{id}/resolve — resolves a moderation flag.
func (h *ModerationHandler) Resolve(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}

	flagID := chi.URLParam(r, "id")

	var req struct {
		Status string `json:"status"` // e.g. "resolved", "dismissed"
		Notes  string `json:"notes"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	if req.Status == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "status is required"})
		return
	}

	if err := h.store.ResolveModerationFlag(r.Context(), flagID, user.ID, req.Status, req.Notes); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to resolve flag"})
		return
	}

	// Log to admin audit
	adminAudit(h.store, r, "resolve_moderation_flag", "moderation_flag", flagID, map[string]interface{}{
		"status": req.Status,
		"notes":  req.Notes,
	})

	writeJSON(w, http.StatusOK, map[string]bool{"success": true})
}
