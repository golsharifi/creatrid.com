package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/creatrid/creatrid/internal/middleware"
	"github.com/creatrid/creatrid/internal/model"
	"github.com/creatrid/creatrid/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/nrednav/cuid2"
)

type AgencyHandler struct {
	store *store.Store
}

func NewAgencyHandler(st *store.Store) *AgencyHandler {
	return &AgencyHandler{store: st}
}

// Create creates an agency profile and converts the user to BRAND role
func (h *AgencyHandler) Create(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}

	// Check if user already has an agency
	existing, err := h.store.FindAgencyByUserID(r.Context(), user.ID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to check agency"})
		return
	}
	if existing != nil {
		writeJSON(w, http.StatusConflict, map[string]string{"error": "Agency already exists"})
		return
	}

	var req struct {
		Name        string  `json:"name"`
		Website     *string `json:"website"`
		Description *string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		return
	}
	if req.Name == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Agency name is required"})
		return
	}

	now := time.Now()
	agency := &store.Agency{
		ID:          cuid2.Generate(),
		UserID:      user.ID,
		Name:        req.Name,
		Website:     req.Website,
		Description: req.Description,
		IsVerified:  false,
		MaxCreators: 10,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := h.store.CreateAgency(r.Context(), agency); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to create agency"})
		return
	}

	// Update user role to BRAND
	if user.Role != model.RoleAdmin {
		if err := h.store.UpdateUserRole(r.Context(), user.ID, string(model.RoleBrand)); err != nil {
			// Non-fatal — agency was created
		}
	}

	writeJSON(w, http.StatusCreated, map[string]interface{}{"agency": agency})
}

// Get returns the current user's agency profile with stats
func (h *AgencyHandler) Get(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}

	agency, err := h.store.FindAgencyByUserID(r.Context(), user.ID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch agency"})
		return
	}
	if agency == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "No agency found"})
		return
	}

	stats, err := h.store.GetAgencyStats(r.Context(), agency.ID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch stats"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"agency": agency,
		"stats":  stats,
	})
}

// Update updates the agency profile
func (h *AgencyHandler) Update(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}

	agency, err := h.store.FindAgencyByUserID(r.Context(), user.ID)
	if err != nil || agency == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "No agency found"})
		return
	}

	var req struct {
		Name        *string `json:"name"`
		Website     *string `json:"website"`
		Description *string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		return
	}

	name := agency.Name
	if req.Name != nil && *req.Name != "" {
		name = *req.Name
	}
	website := agency.Website
	if req.Website != nil {
		website = req.Website
	}
	description := agency.Description
	if req.Description != nil {
		description = req.Description
	}

	if err := h.store.UpdateAgency(r.Context(), agency.ID, name, website, description); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to update agency"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]bool{"success": true})
}

// Invite invites a creator by username to the agency
func (h *AgencyHandler) Invite(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}

	agency, err := h.store.FindAgencyByUserID(r.Context(), user.ID)
	if err != nil || agency == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "No agency found"})
		return
	}

	var req struct {
		Username string `json:"username"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		return
	}
	if req.Username == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Username is required"})
		return
	}

	// Find creator by username
	creator, err := h.store.FindUserByUsername(r.Context(), req.Username)
	if err != nil || creator == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "Creator not found"})
		return
	}

	// Check max creators limit
	count, err := h.store.CountActiveAgencyCreators(r.Context(), agency.ID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to check limits"})
		return
	}
	if count >= agency.MaxCreators {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Maximum creators limit reached"})
		return
	}

	id := cuid2.Generate()
	if err := h.store.InviteCreator(r.Context(), id, agency.ID, creator.ID); err != nil {
		writeJSON(w, http.StatusConflict, map[string]string{"error": "Creator already invited or managed"})
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{"id": id})
}

// ListCreators lists managed creators with scores
func (h *AgencyHandler) ListCreators(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}

	agency, err := h.store.FindAgencyByUserID(r.Context(), user.ID)
	if err != nil || agency == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "No agency found"})
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

	creators, total, err := h.store.ListAgencyCreators(r.Context(), agency.ID, limit, offset)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch creators"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"creators": creators,
		"total":    total,
	})
}

// RemoveCreator removes a creator from the agency
func (h *AgencyHandler) RemoveCreator(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}

	agency, err := h.store.FindAgencyByUserID(r.Context(), user.ID)
	if err != nil || agency == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "No agency found"})
		return
	}

	creatorID := chi.URLParam(r, "creatorId")
	if creatorID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Creator ID is required"})
		return
	}

	if err := h.store.RemoveAgencyCreator(r.Context(), agency.ID, creatorID); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to remove creator"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]bool{"success": true})
}

// Analytics returns aggregate analytics across managed creators
func (h *AgencyHandler) Analytics(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}

	agency, err := h.store.FindAgencyByUserID(r.Context(), user.ID)
	if err != nil || agency == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "No agency found"})
		return
	}

	analytics, err := h.store.GetAgencyAnalytics(r.Context(), agency.ID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch analytics"})
		return
	}

	writeJSON(w, http.StatusOK, analytics)
}

// APIUsage returns API usage stats for the current user
func (h *AgencyHandler) APIUsage(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}

	days, _ := strconv.Atoi(r.URL.Query().Get("days"))
	if days <= 0 || days > 365 {
		days = 30
	}

	summary, err := h.store.GetAPIUsageSummary(r.Context(), user.ID, days)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch usage"})
		return
	}

	writeJSON(w, http.StatusOK, summary)
}

// ListInvites lists pending agency invites for the current user (creator side)
func (h *AgencyHandler) ListInvites(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}

	invites, err := h.store.ListCreatorInvites(r.Context(), user.ID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch invites"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{"invites": invites})
}

// RespondToInvite accepts or declines an agency invite
func (h *AgencyHandler) RespondToInvite(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}

	inviteID := chi.URLParam(r, "id")
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

	// Verify invite belongs to this user
	invite, err := h.store.FindAgencyInvite(r.Context(), inviteID)
	if err != nil || invite == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "Invite not found"})
		return
	}
	if invite.CreatorID != user.ID {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "Not authorized"})
		return
	}
	if invite.Status != "pending" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invite already " + invite.Status})
		return
	}

	accept := req.Action == "accept"
	if err := h.store.RespondToAgencyInvite(r.Context(), inviteID, accept); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to respond"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]bool{"success": true})
}

// BulkVerify verifies multiple creators at once (admin only)
func (h *AgencyHandler) BulkVerify(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil || user.Role != model.RoleAdmin {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "Admin access required"})
		return
	}

	var req struct {
		CreatorIDs []string `json:"creatorIds"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		return
	}
	if len(req.CreatorIDs) == 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "creatorIds is required"})
		return
	}

	if err := h.store.BulkSetVerified(r.Context(), req.CreatorIDs, true); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to verify creators"})
		return
	}

	adminAudit(h.store, r, "bulk_verify_creators", "user", "", map[string]interface{}{
		"creatorIds": req.CreatorIDs,
		"count":      len(req.CreatorIDs),
	})

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success":  true,
		"verified": len(req.CreatorIDs),
	})
}
