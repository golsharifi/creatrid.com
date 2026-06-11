package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/creatrid/creatrid/internal/middleware"
	"github.com/creatrid/creatrid/internal/model"
	"github.com/creatrid/creatrid/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/nrednav/cuid2"
)

// CommentHandler — community comments on PUBLIC vault content. Receiving a
// comment from another creator credits the owner a couple of arena points
// ("community-based asset value": engagement feeds reputation, closed-loop).
type CommentHandler struct {
	store *store.Store
}

func NewCommentHandler(st *store.Store) *CommentHandler {
	return &CommentHandler{store: st}
}

// List handles GET /api/content/{id}/comments — public, for public items only.
func (h *CommentHandler) List(w http.ResponseWriter, r *http.Request) {
	contentID := chi.URLParam(r, "id")
	item, err := h.store.FindContentItemByID(r.Context(), contentID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Database error"})
		return
	}
	if item == nil || !item.IsPublic {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "Content not found"})
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if offset < 0 {
		offset = 0
	}

	comments, total, err := h.store.ListContentComments(r.Context(), contentID, limit, offset)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to load comments"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"comments": comments, "total": total})
}

// Create handles POST /api/content/{id}/comments — authenticated.
func (h *CommentHandler) Create(w http.ResponseWriter, r *http.Request) {
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
	if item == nil || !item.IsPublic {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "Content not found or not public"})
		return
	}

	var req struct {
		Body string `json:"body"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}
	body := strings.TrimSpace(req.Body)
	if body == "" || len(body) > 2000 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Comment must be 1-2000 characters"})
		return
	}

	id := cuid2.Generate()
	if err := h.store.CreateContentComment(r.Context(), id, contentID, user.ID, body); err != nil {
		log.Printf("Comment create failed: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to post comment"})
		return
	}

	// Engagement feeds the owner's arena reputation (not self-comments).
	if item.UserID != user.ID {
		if err := h.store.AddArenaPoints(r.Context(), cuid2.Generate(), item.UserID, 2, "comment_received"); err != nil {
			log.Printf("Comment points credit failed: %v", err)
		}
	}

	writeJSON(w, http.StatusCreated, map[string]string{"id": id})
}

// Delete handles DELETE /api/comments/{id} — author, content owner, or admin.
func (h *CommentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}

	commentID := chi.URLParam(r, "id")
	deleted, err := h.store.DeleteContentComment(r.Context(), commentID, user.ID, user.Role == model.RoleAdmin)
	if err != nil {
		log.Printf("Comment delete failed: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to delete comment"})
		return
	}
	if !deleted {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "Comment not found or not yours"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"success": true})
}
