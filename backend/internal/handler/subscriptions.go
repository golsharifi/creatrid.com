package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/creatrid/creatrid/internal/config"
	"github.com/creatrid/creatrid/internal/middleware"
	"github.com/creatrid/creatrid/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/nrednav/cuid2"
	"github.com/stripe/stripe-go/v81"
	stripesub "github.com/stripe/stripe-go/v81/subscription"
)

// FanSubscriptionHandler manages fan subscription and content gating endpoints.
type FanSubscriptionHandler struct {
	store  *store.Store
	config *config.Config
}

// NewFanSubscriptionHandler creates a new FanSubscriptionHandler.
func NewFanSubscriptionHandler(st *store.Store, cfg *config.Config) *FanSubscriptionHandler {
	return &FanSubscriptionHandler{store: st, config: cfg}
}

// Tier pricing in cents
var tierPricing = map[string]int{
	"supporter": 300,
	"superfan":  1000,
	"patron":    2500,
}

type subscribeRequest struct {
	CreatorUserID string `json:"creatorUserId"`
	Tier          string `json:"tier"`
}

// Subscribe creates a fan subscription to a creator.
func (h *FanSubscriptionHandler) Subscribe(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}

	var req subscribeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	if req.CreatorUserID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "creatorUserId is required"})
		return
	}
	if req.CreatorUserID == user.ID {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Cannot subscribe to yourself"})
		return
	}

	priceCents, ok := tierPricing[req.Tier]
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid tier. Must be supporter, superfan, or patron"})
		return
	}

	// Check if already subscribed
	existing, err := h.store.FindFanSubscription(r.Context(), user.ID, req.CreatorUserID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Database error"})
		return
	}
	if existing != nil && existing.Status == "active" {
		writeJSON(w, http.StatusConflict, map[string]string{"error": "Already subscribed to this creator"})
		return
	}

	// Verify creator exists
	creator, err := h.store.FindUserByID(r.Context(), req.CreatorUserID)
	if err != nil || creator == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "Creator not found"})
		return
	}

	sub := &store.FanSubscription{
		ID:            cuid2.Generate(),
		FanUserID:     user.ID,
		CreatorUserID: req.CreatorUserID,
		Tier:          req.Tier,
		PriceCents:    priceCents,
		Status:        "active",
		StartedAt:     time.Now(),
	}

	if h.config.StripeSecretKey != "" {
		// Real Stripe flow — create a subscription
		// For simplicity, we simulate it via metadata. In production,
		// you'd create a Stripe subscription with a proper price ID.
		params := &stripe.SubscriptionParams{
			// In production, use proper customer/price IDs
		}
		_ = params // unused in test mode

		// Since we don't have fan subscription price IDs configured,
		// we still mark as active immediately and record for tracking
		stripeSubID := "sim_" + cuid2.Generate()
		sub.StripeSubscriptionID = &stripeSubID
	}

	if err := h.store.CreateFanSubscription(r.Context(), sub); err != nil {
		log.Printf("Failed to create fan subscription: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to create subscription"})
		return
	}

	writeJSON(w, http.StatusCreated, map[string]interface{}{"subscription": sub})
}

// MySubscriptions lists the current user's subscriptions as a fan.
func (h *FanSubscriptionHandler) MySubscriptions(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}

	subs, err := h.store.ListSubscriptionsByFan(r.Context(), user.ID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Database error"})
		return
	}
	if subs == nil {
		subs = []store.FanSubscription{}
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{"subscriptions": subs})
}

// MyFans lists the current user's fans (as creator).
func (h *FanSubscriptionHandler) MyFans(w http.ResponseWriter, r *http.Request) {
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

	fans, total, err := h.store.ListFansByCreator(r.Context(), user.ID, limit, offset)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Database error"})
		return
	}
	if fans == nil {
		fans = []store.FanSubscription{}
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"fans":  fans,
		"total": total,
	})
}

// Cancel cancels a fan subscription.
func (h *FanSubscriptionHandler) Cancel(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}

	subID := chi.URLParam(r, "id")

	// Verify ownership — only the fan can cancel
	subs, err := h.store.ListSubscriptionsByFan(r.Context(), user.ID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Database error"})
		return
	}

	var found *store.FanSubscription
	for i := range subs {
		if subs[i].ID == subID {
			found = &subs[i]
			break
		}
	}
	if found == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "Subscription not found"})
		return
	}
	if found.Status != "active" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Subscription is not active"})
		return
	}

	// Cancel on Stripe if configured
	if h.config.StripeSecretKey != "" && found.StripeSubscriptionID != nil {
		_, err := stripesub.Cancel(*found.StripeSubscriptionID, nil)
		if err != nil {
			log.Printf("Stripe subscription cancel error: %v", err)
			// Continue anyway to update local status
		}
	}

	if err := h.store.UpdateFanSubscriptionStatus(r.Context(), subID, "canceled"); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to cancel subscription"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]bool{"success": true})
}

type gateContentRequest struct {
	TokenID          *string `json:"tokenId"`
	MinTokens        *int    `json:"minTokens"`
	SubscriptionTier *string `json:"subscriptionTier"`
}

// GateContent adds token/subscription gating to a content item.
func (h *FanSubscriptionHandler) GateContent(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}

	contentID := chi.URLParam(r, "id")

	// Verify ownership of content
	content, err := h.store.FindContentItemByID(r.Context(), contentID)
	if err != nil || content == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "Content not found"})
		return
	}
	if content.UserID != user.ID {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "Not your content"})
		return
	}

	var req gateContentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	if req.TokenID == nil && req.SubscriptionTier == nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Must specify tokenId or subscriptionTier"})
		return
	}

	minTokens := 1
	if req.MinTokens != nil && *req.MinTokens > 0 {
		minTokens = *req.MinTokens
	}

	gated := &store.GatedContent{
		ID:               cuid2.Generate(),
		ContentID:        contentID,
		TokenID:          req.TokenID,
		MinTokens:        minTokens,
		SubscriptionTier: req.SubscriptionTier,
		CreatedAt:        time.Now(),
	}

	if err := h.store.SetGatedContent(r.Context(), gated); err != nil {
		log.Printf("Failed to gate content: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to gate content"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]bool{"success": true})
}

// RemoveGate removes token/subscription gating from a content item.
func (h *FanSubscriptionHandler) RemoveGate(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}

	contentID := chi.URLParam(r, "id")

	// Verify ownership
	content, err := h.store.FindContentItemByID(r.Context(), contentID)
	if err != nil || content == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "Content not found"})
		return
	}
	if content.UserID != user.ID {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "Not your content"})
		return
	}

	if err := h.store.RemoveGatedContent(r.Context(), contentID); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to remove gate"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]bool{"success": true})
}

// CheckAccess checks if the current user has access to gated content.
func (h *FanSubscriptionHandler) CheckAccess(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}

	contentID := chi.URLParam(r, "id")

	hasAccess, err := h.store.CheckContentAccess(r.Context(), contentID, user.ID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Database error"})
		return
	}

	reason := ""
	if !hasAccess {
		gated, _ := h.store.FindGatedContent(r.Context(), contentID)
		if gated != nil {
			if gated.TokenID != nil {
				reason = "Requires " + strconv.Itoa(gated.MinTokens) + " token(s)"
			}
			if gated.SubscriptionTier != nil {
				if reason != "" {
					reason += " or "
				}
				reason += *gated.SubscriptionTier + " subscription"
			}
		}
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"hasAccess": hasAccess,
		"reason":    reason,
	})
}
