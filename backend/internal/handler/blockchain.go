package handler

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/creatrid/creatrid/internal/blockchain"
	"github.com/creatrid/creatrid/internal/middleware"
	"github.com/creatrid/creatrid/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/nrednav/cuid2"
)

type BlockchainHandler struct {
	store     *store.Store
	anchorSvc *blockchain.AnchorService
}

func NewBlockchainHandler(st *store.Store, anchorSvc *blockchain.AnchorService) *BlockchainHandler {
	return &BlockchainHandler{
		store:     st,
		anchorSvc: anchorSvc,
	}
}

// Anchor handles POST /api/content/{id}/anchor — anchor content on blockchain.
func (h *BlockchainHandler) Anchor(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}

	contentID := chi.URLParam(r, "id")
	if contentID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Content ID is required"})
		return
	}

	// Fetch content item
	item, err := h.store.FindContentItemByID(r.Context(), contentID)
	if err != nil {
		log.Printf("Blockchain anchor DB error: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Database error"})
		return
	}
	if item == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "Content not found"})
		return
	}
	if item.UserID != user.ID {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "Not authorized"})
		return
	}

	// Check if already anchored
	existing, err := h.store.FindAnchorByContentID(r.Context(), contentID)
	if err != nil {
		log.Printf("Blockchain anchor lookup error: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Database error"})
		return
	}
	if existing != nil {
		writeJSON(w, http.StatusConflict, map[string]string{"error": "Content is already anchored"})
		return
	}

	// Perform anchoring
	txHash, blockNumber, contractAddress, err := h.anchorSvc.AnchorHash(item.HashSHA256)
	if err != nil {
		log.Printf("Blockchain anchor error: %v", err)
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "Blockchain anchoring failed: " + err.Error()})
		return
	}

	now := time.Now()
	anchor := &store.ContentAnchor{
		ID:              cuid2.Generate(),
		ContentID:       contentID,
		UserID:          user.ID,
		ContentHash:     item.HashSHA256,
		TxHash:          &txHash,
		Chain:           "polygon",
		BlockNumber:     &blockNumber,
		ContractAddress: &contractAddress,
		AnchorStatus:    "confirmed",
		CreatedAt:       now,
		ConfirmedAt:     &now,
	}

	if err := h.store.CreateContentAnchor(r.Context(), anchor); err != nil {
		log.Printf("Blockchain anchor save error: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to save anchor record"})
		return
	}

	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"anchor":    anchor,
		"simulated": h.anchorSvc.IsSimulated(),
	})
}

// GetAnchor handles GET /api/content/{id}/anchor — get anchor status for content.
func (h *BlockchainHandler) GetAnchor(w http.ResponseWriter, r *http.Request) {
	contentID := chi.URLParam(r, "id")
	if contentID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Content ID is required"})
		return
	}

	anchor, err := h.store.FindAnchorByContentID(r.Context(), contentID)
	if err != nil {
		log.Printf("Get anchor DB error: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Database error"})
		return
	}
	if anchor == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "No anchor found for this content"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"anchor": anchor,
	})
}

// VerifyByHash handles GET /api/verify/{hash} — public verification by content hash.
func (h *BlockchainHandler) VerifyByHash(w http.ResponseWriter, r *http.Request) {
	hash := chi.URLParam(r, "hash")
	if hash == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Hash is required"})
		return
	}

	anchor, err := h.store.FindAnchorByHash(r.Context(), hash)
	if err != nil {
		log.Printf("Verify by hash DB error: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Database error"})
		return
	}
	if anchor == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "No blockchain anchor found for this hash"})
		return
	}

	// Fetch the content item for additional context
	item, err := h.store.FindContentItemByID(r.Context(), anchor.ContentID)
	if err != nil {
		log.Printf("Verify content lookup error: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Database error"})
		return
	}

	var contentInfo map[string]interface{}
	if item != nil {
		contentInfo = map[string]interface{}{
			"id":          item.ID,
			"title":       item.Title,
			"contentType": item.ContentType,
			"createdAt":   item.CreatedAt,
		}
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"anchor":  anchor,
		"content": contentInfo,
	})
}

// ListAnchors handles GET /api/anchors — list user's anchored content.
func (h *BlockchainHandler) ListAnchors(w http.ResponseWriter, r *http.Request) {
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

	anchors, total, err := h.store.ListAnchorsByUser(r.Context(), user.ID, limit, offset)
	if err != nil {
		log.Printf("List anchors DB error: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch anchors"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"anchors": anchors,
		"total":   total,
	})
}
