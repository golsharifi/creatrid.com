package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/creatrid/creatrid/internal/middleware"
	"github.com/creatrid/creatrid/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/nrednav/cuid2"
)

type CollectionHandler struct {
	store *store.Store
}

func NewCollectionHandler(st *store.Store) *CollectionHandler {
	return &CollectionHandler{store: st}
}

type createCollectionRequest struct {
	Title       string  `json:"title"`
	Description *string `json:"description"`
	IsPublic    bool    `json:"isPublic"`
}

// Create handles POST /api/collections
func (h *CollectionHandler) Create(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}

	var req createCollectionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	title := strings.TrimSpace(req.Title)
	if title == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Title is required"})
		return
	}

	now := time.Now()
	coll := &store.Collection{
		ID:          cuid2.Generate(),
		UserID:      user.ID,
		Title:       title,
		Description: req.Description,
		IsPublic:    req.IsPublic,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := h.store.CreateCollection(r.Context(), coll); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to create collection"})
		return
	}

	writeJSON(w, http.StatusCreated, coll)
}

// List handles GET /api/collections
func (h *CollectionHandler) List(w http.ResponseWriter, r *http.Request) {
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

	collections, total, err := h.store.ListCollectionsByUser(r.Context(), user.ID, limit, offset)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch collections"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"collections": collections,
		"total":       total,
	})
}

// Get handles GET /api/collections/{id}
func (h *CollectionHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	coll, err := h.store.FindCollectionByID(r.Context(), id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Database error"})
		return
	}
	if coll == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "Collection not found"})
		return
	}

	// If not public, must be owner
	if !coll.IsPublic {
		user := middleware.UserFromContext(r.Context())
		if user == nil || user.ID != coll.UserID {
			writeJSON(w, http.StatusForbidden, map[string]string{"error": "Not authorized"})
			return
		}
	}

	items, err := h.store.ListCollectionItems(r.Context(), id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch items"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"collection": coll,
		"items":      items,
	})
}

type updateCollectionRequest struct {
	Title       string  `json:"title"`
	Description *string `json:"description"`
	IsPublic    bool    `json:"isPublic"`
}

// Update handles PATCH /api/collections/{id}
func (h *CollectionHandler) Update(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}

	id := chi.URLParam(r, "id")
	coll, err := h.store.FindCollectionByID(r.Context(), id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Database error"})
		return
	}
	if coll == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "Collection not found"})
		return
	}
	if coll.UserID != user.ID {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "Not authorized"})
		return
	}

	var req updateCollectionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	title := strings.TrimSpace(req.Title)
	if title == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Title is required"})
		return
	}

	if err := h.store.UpdateCollection(r.Context(), id, title, req.Description, req.IsPublic); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to update collection"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]bool{"success": true})
}

// Delete handles DELETE /api/collections/{id}
func (h *CollectionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}

	id := chi.URLParam(r, "id")
	coll, err := h.store.FindCollectionByID(r.Context(), id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Database error"})
		return
	}
	if coll == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "Collection not found"})
		return
	}
	if coll.UserID != user.ID {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "Not authorized"})
		return
	}

	if err := h.store.DeleteCollection(r.Context(), id); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to delete collection"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]bool{"deleted": true})
}

type addItemRequest struct {
	ContentID string `json:"contentId"`
	Position  int    `json:"position"`
}

// AddItem handles POST /api/collections/{id}/items
func (h *CollectionHandler) AddItem(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}

	id := chi.URLParam(r, "id")
	coll, err := h.store.FindCollectionByID(r.Context(), id)
	if err != nil || coll == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "Collection not found"})
		return
	}
	if coll.UserID != user.ID {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "Not authorized"})
		return
	}

	var req addItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	if err := h.store.AddItemToCollection(r.Context(), id, req.ContentID, req.Position); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to add item"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]bool{"success": true})
}

// RemoveItem handles DELETE /api/collections/{id}/items/{contentId}
func (h *CollectionHandler) RemoveItem(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}

	id := chi.URLParam(r, "id")
	contentID := chi.URLParam(r, "contentId")

	coll, err := h.store.FindCollectionByID(r.Context(), id)
	if err != nil || coll == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "Collection not found"})
		return
	}
	if coll.UserID != user.ID {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "Not authorized"})
		return
	}

	if err := h.store.RemoveItemFromCollection(r.Context(), id, contentID); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to remove item"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]bool{"success": true})
}

// ListItems handles GET /api/collections/{id}/items
func (h *CollectionHandler) ListItems(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	coll, err := h.store.FindCollectionByID(r.Context(), id)
	if err != nil || coll == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "Collection not found"})
		return
	}

	if !coll.IsPublic {
		user := middleware.UserFromContext(r.Context())
		if user == nil || user.ID != coll.UserID {
			writeJSON(w, http.StatusForbidden, map[string]string{"error": "Not authorized"})
			return
		}
	}

	items, err := h.store.ListCollectionItems(r.Context(), id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch items"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{"items": items})
}

// PublicList handles GET /api/users/{username}/collections
func (h *CollectionHandler) PublicList(w http.ResponseWriter, r *http.Request) {
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

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	collections, total, err := h.store.ListPublicCollectionsByUser(r.Context(), user.ID, limit, offset)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch collections"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"collections": collections,
		"total":       total,
	})
}
