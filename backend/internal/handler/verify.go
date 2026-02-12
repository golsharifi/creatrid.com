package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/creatrid/creatrid/internal/store"
	"github.com/go-chi/chi/v5"
)

type VerifyHandler struct {
	store *store.Store
}

func NewVerifyHandler(st *store.Store) *VerifyHandler {
	return &VerifyHandler{store: st}
}

// Verify returns full verification data for a creator
func (h *VerifyHandler) Verify(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	if username == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Username is required"})
		return
	}

	user, err := h.store.FindUserByUsername(r.Context(), strings.ToLower(username))
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Database error"})
		return
	}
	if user == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "User not found"})
		return
	}

	connections, err := h.store.FindConnectionsByUserID(r.Context(), user.ID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch connections"})
		return
	}

	publicConns := make([]interface{}, 0, len(connections))
	for _, c := range connections {
		publicConns = append(publicConns, c.ToPublic())
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"user":            user.ToPublic(),
		"connections":     publicConns,
		"connectionCount": len(connections),
	})
}

// Score returns just the creator score for a username
func (h *VerifyHandler) Score(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	if username == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Username is required"})
		return
	}

	user, err := h.store.FindUserByUsername(r.Context(), strings.ToLower(username))
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Database error"})
		return
	}
	if user == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "User not found"})
		return
	}

	var score int
	if user.CreatorScore != nil {
		score = *user.CreatorScore
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"username": user.Username,
		"score":    score,
		"verified": user.IsVerified,
	})
}

// Search finds creators with filters (uses existing DiscoverCreators)
func (h *VerifyHandler) Search(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	minScore, _ := strconv.Atoi(r.URL.Query().Get("minScore"))
	platform := r.URL.Query().Get("platform")
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	if limit <= 0 || limit > 50 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	creators, total, err := h.store.DiscoverCreators(r.Context(), minScore, platform, q, limit, offset)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to search creators"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"creators": creators,
		"total":    total,
	})
}
