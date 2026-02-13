package handler

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/creatrid/creatrid/internal/middleware"
	"github.com/creatrid/creatrid/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/nrednav/cuid2"
)

var validWebhookEvents = map[string]bool{
	"license.sold":          true,
	"content.uploaded":      true,
	"profile.viewed":        true,
	"collaboration.received": true,
	"payout.completed":      true,
}

type WebhookHandler struct {
	store *store.Store
}

func NewWebhookHandler(st *store.Store) *WebhookHandler {
	return &WebhookHandler{store: st}
}

func (h *WebhookHandler) Create(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}

	var req struct {
		URL    string   `json:"url"`
		Events []string `json:"events"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		return
	}

	if req.URL == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "URL is required"})
		return
	}

	for _, e := range req.Events {
		if !validWebhookEvents[e] {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid event: " + e})
			return
		}
	}

	// Generate a signing secret
	secretBytes := make([]byte, 32)
	if _, err := rand.Read(secretBytes); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to generate secret"})
		return
	}
	secret := "whsec_" + hex.EncodeToString(secretBytes)

	ep := &store.WebhookEndpoint{
		ID:       cuid2.Generate(),
		UserID:   user.ID,
		URL:      req.URL,
		Secret:   secret,
		Events:   req.Events,
		IsActive: true,
	}

	if err := h.store.CreateWebhookEndpoint(r.Context(), ep); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to create webhook"})
		return
	}

	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"id":     ep.ID,
		"url":    ep.URL,
		"events": ep.Events,
		"secret": secret,
	})
}

func (h *WebhookHandler) List(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}

	endpoints, err := h.store.ListWebhookEndpointsByUser(r.Context(), user.ID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch webhooks"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{"endpoints": endpoints})
}

func (h *WebhookHandler) Update(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}

	id := chi.URLParam(r, "id")
	ep, err := h.store.FindWebhookEndpointByID(r.Context(), id)
	if err != nil || ep == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "Webhook not found"})
		return
	}
	if ep.UserID != user.ID {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "Not authorized"})
		return
	}

	var req struct {
		URL      string   `json:"url"`
		Events   []string `json:"events"`
		IsActive bool     `json:"isActive"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		return
	}

	if err := h.store.UpdateWebhookEndpoint(r.Context(), id, req.URL, req.Events, req.IsActive); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to update"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]bool{"success": true})
}

func (h *WebhookHandler) Delete(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}

	id := chi.URLParam(r, "id")
	ep, err := h.store.FindWebhookEndpointByID(r.Context(), id)
	if err != nil || ep == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "Webhook not found"})
		return
	}
	if ep.UserID != user.ID {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "Not authorized"})
		return
	}

	if err := h.store.DeleteWebhookEndpoint(r.Context(), id); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to delete"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]bool{"success": true})
}

func (h *WebhookHandler) Deliveries(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}

	id := chi.URLParam(r, "id")
	ep, err := h.store.FindWebhookEndpointByID(r.Context(), id)
	if err != nil || ep == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "Webhook not found"})
		return
	}
	if ep.UserID != user.ID {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "Not authorized"})
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	deliveries, total, err := h.store.ListWebhookDeliveries(r.Context(), id, limit, offset)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch deliveries"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"deliveries": deliveries,
		"total":      total,
	})
}
