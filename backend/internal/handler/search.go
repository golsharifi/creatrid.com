package handler

import (
	"net/http"
	"strconv"

	"github.com/creatrid/creatrid/internal/store"
)

type SearchHandler struct {
	store *store.Store
}

func NewSearchHandler(st *store.Store) *SearchHandler {
	return &SearchHandler{store: st}
}

// Search handles GET /api/search?q=&type=&limit=&offset=
func (h *SearchHandler) Search(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Query parameter 'q' is required"})
		return
	}

	searchType := r.URL.Query().Get("type")
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	result := map[string]interface{}{}

	if searchType == "" || searchType == "content" {
		content, contentTotal, err := h.store.SearchContent(r.Context(), query, limit, offset)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Search failed"})
			return
		}
		result["content"] = content
		result["contentTotal"] = contentTotal
	}

	if searchType == "" || searchType == "creators" {
		users, usersTotal, err := h.store.SearchUsers(r.Context(), query, limit, offset)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Search failed"})
			return
		}
		result["creators"] = users
		result["creatorsTotal"] = usersTotal
	}

	writeJSON(w, http.StatusOK, result)
}

// Suggestions handles GET /api/search/suggestions?q=
func (h *SearchHandler) Suggestions(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		writeJSON(w, http.StatusOK, map[string]interface{}{"suggestions": []string{}})
		return
	}

	suggestions, err := h.store.SearchSuggestions(r.Context(), query, 10)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch suggestions"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{"suggestions": suggestions})
}
