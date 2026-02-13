package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/creatrid/creatrid/internal/middleware"
	"github.com/creatrid/creatrid/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/nrednav/cuid2"
)

type CollaborationHandler struct {
	store *store.Store
}

func NewCollaborationHandler(st *store.Store) *CollaborationHandler {
	return &CollaborationHandler{store: st}
}

// Discover returns a paginated list of creators with optional filters
func (h *CollaborationHandler) Discover(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	minScore, _ := strconv.Atoi(r.URL.Query().Get("minScore"))
	platform := r.URL.Query().Get("platform")
	search := r.URL.Query().Get("q")

	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	users, total, err := h.store.DiscoverCreators(r.Context(), minScore, platform, search, limit, offset)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch creators"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"creators": users,
		"total":    total,
	})
}

// SendRequest creates a new collaboration request
func (h *CollaborationHandler) SendRequest(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}

	var req struct {
		ToUserID string `json:"toUserId"`
		Message  string `json:"message"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		return
	}

	if req.ToUserID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "toUserId is required"})
		return
	}
	if req.ToUserID == user.ID {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Cannot send request to yourself"})
		return
	}
	if len(req.Message) > 1000 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Message too long (max 1000 characters)"})
		return
	}

	id := cuid2.Generate()
	if err := h.store.CreateCollaborationRequest(r.Context(), id, user.ID, req.ToUserID, req.Message); err != nil {
		writeJSON(w, http.StatusConflict, map[string]string{"error": "Request already sent to this user"})
		return
	}

	// Notify the recipient
	senderName := "Someone"
	if user.Name != nil {
		senderName = *user.Name
	}
	notif := &store.Notification{
		ID:        cuid2.Generate(),
		UserID:    req.ToUserID,
		Type:      "collab_request",
		Title:     "New collaboration request",
		Message:   fmt.Sprintf("%s wants to collaborate with you", senderName),
		Data:      []byte(fmt.Sprintf(`{"requestId":"%s","fromUserId":"%s"}`, id, user.ID)),
		CreatedAt: time.Now(),
	}
	if err := h.store.CreateNotification(r.Context(), notif); err != nil {
		log.Printf("Failed to create collab notification: %v", err)
	}

	writeJSON(w, http.StatusCreated, map[string]string{"id": id})
}

// Inbox returns incoming collaboration requests
func (h *CollaborationHandler) Inbox(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}

	requests, err := h.store.IncomingCollaborations(r.Context(), user.ID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch requests"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{"requests": requests})
}

// Outbox returns outgoing collaboration requests
func (h *CollaborationHandler) Outbox(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}

	requests, err := h.store.OutgoingCollaborations(r.Context(), user.ID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch requests"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{"requests": requests})
}

// Respond accepts or declines a collaboration request
func (h *CollaborationHandler) Respond(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}

	requestID := chi.URLParam(r, "id")
	var req struct {
		Action string `json:"action"` // "accept" or "decline"
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		return
	}

	if req.Action != "accept" && req.Action != "decline" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Action must be 'accept' or 'decline'"})
		return
	}

	cr, err := h.store.FindCollaborationRequest(r.Context(), requestID)
	if err != nil || cr == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "Request not found"})
		return
	}

	if cr.ToUserID != user.ID {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "Not authorized"})
		return
	}

	if cr.Status != "pending" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Request already " + cr.Status})
		return
	}

	status := "accepted"
	if req.Action == "decline" {
		status = "declined"
	}

	if err := h.store.UpdateCollaborationStatus(r.Context(), requestID, status); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to update"})
		return
	}

	// Notify the original sender about the response
	responderName := "Someone"
	if user.Name != nil {
		responderName = *user.Name
	}
	notif := &store.Notification{
		ID:        cuid2.Generate(),
		UserID:    cr.FromUserID,
		Type:      "collab_response",
		Title:     fmt.Sprintf("Collaboration %s", status),
		Message:   fmt.Sprintf("%s %s your collaboration request", responderName, status),
		Data:      []byte(fmt.Sprintf(`{"requestId":"%s","status":"%s"}`, requestID, status)),
		CreatedAt: time.Now(),
	}
	if err := h.store.CreateNotification(r.Context(), notif); err != nil {
		log.Printf("Failed to create collab response notification: %v", err)
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": status})
}
