package handler

import (
	"net/http"
	"strconv"

	"github.com/creatrid/creatrid/internal/store"
	"github.com/go-chi/chi/v5"
)

type MarketplaceHandler struct {
	store *store.Store
}

func NewMarketplaceHandler(st *store.Store) *MarketplaceHandler {
	return &MarketplaceHandler{store: st}
}

// Browse returns a paginated list of public content items that have active license offerings.
// GET /api/marketplace
func (h *MarketplaceHandler) Browse(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	contentType := r.URL.Query().Get("type")
	query := r.URL.Query().Get("q")
	sort := r.URL.Query().Get("sort")

	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	if sort != "newest" && sort != "popular" {
		sort = "newest"
	}

	items, total, err := h.store.ListMarketplaceContent(r.Context(), contentType, query, sort, limit, offset)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch marketplace items"})
		return
	}

	// Fetch offerings for each item
	type marketplaceResponse struct {
		store.MarketplaceItem
		Offerings []*store.LicenseOffering `json:"offerings"`
	}

	results := make([]marketplaceResponse, 0, len(items))
	for _, item := range items {
		offerings, err := h.store.ListOfferingsByContent(r.Context(), item.ID)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch offerings"})
			return
		}
		results = append(results, marketplaceResponse{
			MarketplaceItem: item,
			Offerings:       offerings,
		})
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"items": results,
		"total": total,
	})
}

// Detail returns a single public content item with its license offerings and creator info.
// GET /api/marketplace/{id}
func (h *MarketplaceHandler) Detail(w http.ResponseWriter, r *http.Request) {
	contentID := chi.URLParam(r, "id")

	content, err := h.store.FindContentItemByID(r.Context(), contentID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Database error"})
		return
	}
	if content == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "Content not found"})
		return
	}
	if !content.IsPublic {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "Content not found"})
		return
	}

	offerings, err := h.store.ListOfferingsByContent(r.Context(), contentID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch offerings"})
		return
	}

	// Fetch creator info using a single-item marketplace query is overkill;
	// instead, use ListMarketplaceContent with a direct content lookup approach.
	// For simplicity, we fetch the user who owns this content via the store.
	// We already have the content with UserID, so we look up the user.
	creator, err := h.store.FindUserByID(r.Context(), content.UserID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Database error"})
		return
	}

	var creatorName, creatorUsername, creatorImage *string
	if creator != nil {
		creatorName = creator.Name
		creatorUsername = creator.Username
		creatorImage = creator.Image
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"content":         content,
		"offerings":       offerings,
		"creatorName":     creatorName,
		"creatorUsername":  creatorUsername,
		"creatorImage":    creatorImage,
	})
}
