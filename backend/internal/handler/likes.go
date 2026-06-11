package handler

import (
	"log"
	"net/http"

	"github.com/creatrid/creatrid/internal/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/nrednav/cuid2"
)

// Likes on public works — the "community ratings" half of reputation: each
// like from another creator credits the owner one arena point (closed-loop).

// Like handles POST /api/content/{id}/like.
func (h *CommentHandler) Like(w http.ResponseWriter, r *http.Request) {
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

	created, err := h.store.LikeContent(r.Context(), contentID, user.ID)
	if err != nil {
		log.Printf("Like failed: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to like"})
		return
	}
	if created && item.UserID != user.ID {
		if err := h.store.AddArenaPoints(r.Context(), cuid2.Generate(), item.UserID, 1, "like_received"); err != nil {
			log.Printf("Like points credit failed: %v", err)
		}
	}

	count, _, _ := h.store.ContentLikeStats(r.Context(), contentID, user.ID)
	writeJSON(w, http.StatusOK, map[string]any{"liked": true, "likes": count})
}

// Unlike handles DELETE /api/content/{id}/like.
func (h *CommentHandler) Unlike(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}

	contentID := chi.URLParam(r, "id")
	item, err := h.store.FindContentItemByID(r.Context(), contentID)
	if err != nil || item == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "Content not found"})
		return
	}

	removed, err := h.store.UnlikeContent(r.Context(), contentID, user.ID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to unlike"})
		return
	}
	if removed && item.UserID != user.ID {
		if err := h.store.AddArenaPoints(r.Context(), cuid2.Generate(), item.UserID, -1, "like_removed"); err != nil {
			log.Printf("Unlike points debit failed: %v", err)
		}
	}

	count, _, _ := h.store.ContentLikeStats(r.Context(), contentID, user.ID)
	writeJSON(w, http.StatusOK, map[string]any{"liked": false, "likes": count})
}

// LikeStats handles GET /api/content/{id}/likes — public count; "liked" is
// always false here because this route carries no auth context.
func (h *CommentHandler) LikeStats(w http.ResponseWriter, r *http.Request) {
	contentID := chi.URLParam(r, "id")
	item, err := h.store.FindContentItemByID(r.Context(), contentID)
	if err != nil || item == nil || !item.IsPublic {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "Content not found"})
		return
	}
	count, _, err := h.store.ContentLikeStats(r.Context(), contentID, "")
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to load likes"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"likes": count})
}

// MyLikeStats handles GET /api/content/{id}/likes/me — count + whether the
// authenticated user liked it.
func (h *CommentHandler) MyLikeStats(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}
	contentID := chi.URLParam(r, "id")
	count, liked, err := h.store.ContentLikeStats(r.Context(), contentID, user.ID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to load likes"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"likes": count, "liked": liked})
}
