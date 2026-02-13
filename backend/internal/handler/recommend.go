package handler

import (
	"net/http"

	"github.com/creatrid/creatrid/internal/middleware"
	"github.com/creatrid/creatrid/internal/store"
)

type RecommendHandler struct {
	store *store.Store
}

func NewRecommendHandler(st *store.Store) *RecommendHandler {
	return &RecommendHandler{store: st}
}

func (h *RecommendHandler) List(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}

	// Get user's connected platforms
	connections, err := h.store.FindConnectionsByUserID(r.Context(), user.ID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch connections"})
		return
	}

	if len(connections) == 0 {
		writeJSON(w, http.StatusOK, map[string]interface{}{"creators": []interface{}{}})
		return
	}

	var platforms []string
	for _, c := range connections {
		platforms = append(platforms, c.Platform)
	}

	userScore := 0
	if user.CreatorScore != nil {
		userScore = *user.CreatorScore
	}

	creators, err := h.store.RecommendCreators(r.Context(), user.ID, platforms, userScore, 10)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch recommendations"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{"creators": creators})
}
