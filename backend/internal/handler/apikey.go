package handler

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"time"

	"github.com/creatrid/creatrid/internal/middleware"
	"github.com/creatrid/creatrid/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/nrednav/cuid2"
)

type APIKeyHandler struct {
	store *store.Store
}

func NewAPIKeyHandler(st *store.Store) *APIKeyHandler {
	return &APIKeyHandler{store: st}
}

type createAPIKeyRequest struct {
	Name string `json:"name"`
}

// Create generates a new API key for the authenticated user
func (h *APIKeyHandler) Create(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}

	var req createAPIKeyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	if req.Name == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Name is required"})
		return
	}
	if len(req.Name) > 100 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Name must be under 100 characters"})
		return
	}

	// Generate random 32-byte key, hex encode with crk_ prefix
	rawBytes := make([]byte, 32)
	if _, err := rand.Read(rawBytes); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to generate key"})
		return
	}
	rawKey := "crk_" + hex.EncodeToString(rawBytes)

	// Hash with SHA-256 for storage
	hash := sha256.Sum256([]byte(rawKey))
	keyHash := hex.EncodeToString(hash[:])

	// First 8 chars as prefix (after crk_)
	keyPrefix := rawKey[:12] // "crk_" + 8 hex chars

	apiKey := &store.APIKey{
		ID:        cuid2.Generate(),
		UserID:    user.ID,
		Name:      req.Name,
		KeyPrefix: keyPrefix,
		KeyHash:   keyHash,
		CreatedAt: time.Now(),
	}

	if err := h.store.CreateAPIKey(r.Context(), apiKey); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to create API key"})
		return
	}

	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"id":   apiKey.ID,
		"name": apiKey.Name,
		"key":  rawKey,
	})
}

// List returns all API keys for the authenticated user
func (h *APIKeyHandler) List(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}

	keys, err := h.store.FindAPIKeysByUserID(r.Context(), user.ID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch API keys"})
		return
	}

	// Return safe fields only (no hash)
	type safeKey struct {
		ID         string     `json:"id"`
		Name       string     `json:"name"`
		KeyPrefix  string     `json:"keyPrefix"`
		LastUsedAt *time.Time `json:"lastUsedAt"`
		CreatedAt  time.Time  `json:"createdAt"`
	}

	safeKeys := make([]safeKey, len(keys))
	for i, k := range keys {
		safeKeys[i] = safeKey{
			ID:         k.ID,
			Name:       k.Name,
			KeyPrefix:  k.KeyPrefix,
			LastUsedAt: k.LastUsedAt,
			CreatedAt:  k.CreatedAt,
		}
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{"keys": safeKeys})
}

// Delete revokes an API key
func (h *APIKeyHandler) Delete(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}

	keyID := chi.URLParam(r, "id")
	if keyID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Key ID is required"})
		return
	}

	if err := h.store.DeleteAPIKey(r.Context(), keyID, user.ID); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to delete API key"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]bool{"success": true})
}
