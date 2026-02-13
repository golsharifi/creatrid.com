package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/creatrid/creatrid/internal/middleware"
	"github.com/creatrid/creatrid/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/nrednav/cuid2"
)

type DMCAHandler struct {
	store *store.Store
}

func NewDMCAHandler(st *store.Store) *DMCAHandler {
	return &DMCAHandler{store: st}
}

// Report handles POST /api/content/{id}/report — submit a DMCA takedown request.
func (h *DMCAHandler) Report(w http.ResponseWriter, r *http.Request) {
	contentID := chi.URLParam(r, "id")

	// Verify content exists
	content, err := h.store.FindContentItemByID(r.Context(), contentID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Database error"})
		return
	}
	if content == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "Content not found"})
		return
	}

	var req struct {
		ReporterEmail string  `json:"reporterEmail"`
		ReporterName  string  `json:"reporterName"`
		Reason        string  `json:"reason"`
		EvidenceURL   *string `json:"evidenceUrl"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	if req.ReporterEmail == "" || req.ReporterName == "" || req.Reason == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "reporterEmail, reporterName, and reason are required"})
		return
	}

	takedown := &store.TakedownRequest{
		ID:            cuid2.Generate(),
		ReporterEmail: req.ReporterEmail,
		ReporterName:  req.ReporterName,
		ContentID:     contentID,
		Reason:        req.Reason,
		EvidenceURL:   req.EvidenceURL,
		Status:        "pending",
		CreatedAt:     time.Now(),
	}

	if err := h.store.CreateTakedownRequest(r.Context(), takedown); err != nil {
		log.Printf("Failed to create takedown request: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to submit report"})
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{"id": takedown.ID, "status": "pending"})
}

// ListTakedowns handles GET /api/admin/takedowns — list takedown requests (admin only).
func (h *DMCAHandler) ListTakedowns(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	requests, total, err := h.store.ListTakedownRequests(r.Context(), status, limit, offset)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch takedown requests"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"requests": requests,
		"total":    total,
	})
}

// ResolveTakedown handles POST /api/admin/takedowns/{id}/resolve — resolve a takedown request (admin only).
func (h *DMCAHandler) ResolveTakedown(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}

	takedownID := chi.URLParam(r, "id")

	var req struct {
		Status string `json:"status"`
		Notes  string `json:"notes"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	if req.Status != "approved" && req.Status != "denied" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Status must be 'approved' or 'denied'"})
		return
	}

	if err := h.store.ResolveTakedown(r.Context(), takedownID, user.ID, req.Status, req.Notes); err != nil {
		log.Printf("Failed to resolve takedown: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to resolve takedown"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": req.Status})
}
