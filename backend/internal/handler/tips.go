package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/creatrid/creatrid/internal/config"
	"github.com/creatrid/creatrid/internal/middleware"
	"github.com/creatrid/creatrid/internal/store"
	"github.com/nrednav/cuid2"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/paymentintent"
)

// TipHandler manages tip endpoints.
type TipHandler struct {
	store  *store.Store
	config *config.Config
}

// NewTipHandler creates a new TipHandler.
func NewTipHandler(st *store.Store, cfg *config.Config) *TipHandler {
	return &TipHandler{store: st, config: cfg}
}

type sendTipRequest struct {
	ToUserID    string  `json:"toUserId"`
	AmountCents int     `json:"amountCents"`
	Message     *string `json:"message"`
}

// Send creates a tip payment.
func (h *TipHandler) Send(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}

	var req sendTipRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	if req.ToUserID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "toUserId is required"})
		return
	}
	if req.AmountCents < 100 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Minimum tip is $1.00 (100 cents)"})
		return
	}
	if req.ToUserID == user.ID {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Cannot tip yourself"})
		return
	}

	// Verify recipient exists
	recipient, err := h.store.FindUserByID(r.Context(), req.ToUserID)
	if err != nil || recipient == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "Recipient not found"})
		return
	}

	var msg *string
	if req.Message != nil {
		m := strings.TrimSpace(*req.Message)
		if m != "" {
			msg = &m
		}
	}

	tipID := cuid2.Generate()
	tip := &store.Tip{
		ID:          tipID,
		FromUserID:  user.ID,
		ToUserID:    req.ToUserID,
		AmountCents: req.AmountCents,
		Message:     msg,
		Status:      "pending",
		CreatedAt:   time.Now(),
	}

	if err := h.store.CreateTip(r.Context(), tip); err != nil {
		log.Printf("Failed to create tip: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to create tip"})
		return
	}

	if h.config.StripeSecretKey != "" {
		// Real Stripe flow
		params := &stripe.PaymentIntentParams{
			Amount:   stripe.Int64(int64(req.AmountCents)),
			Currency: stripe.String(string(stripe.CurrencyUSD)),
		}
		params.AddMetadata("type", "tip")
		params.AddMetadata("tip_id", tipID)
		params.AddMetadata("from_user_id", user.ID)
		params.AddMetadata("to_user_id", req.ToUserID)

		pi, err := paymentintent.New(params)
		if err != nil {
			log.Printf("Stripe payment intent error: %v", err)
			_ = h.store.UpdateTipStatus(r.Context(), tipID, "failed", nil)
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Payment failed"})
			return
		}

		// Auto-complete for simulated mode — in production, confirm via webhook
		stripeID := pi.ID
		_ = h.store.UpdateTipStatus(r.Context(), tipID, "completed", &stripeID)

		writeJSON(w, http.StatusOK, map[string]interface{}{
			"id":           tipID,
			"clientSecret": pi.ClientSecret,
		})
		return
	}

	// Simulated mode — immediately mark completed
	_ = h.store.UpdateTipStatus(r.Context(), tipID, "completed", nil)

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"id":        tipID,
		"simulated": true,
	})
}

// Received lists tips the current user has received.
func (h *TipHandler) Received(w http.ResponseWriter, r *http.Request) {
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

	tips, total, err := h.store.ListTipsReceived(r.Context(), user.ID, limit, offset)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Database error"})
		return
	}
	if tips == nil {
		tips = []store.Tip{}
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"tips":  tips,
		"total": total,
	})
}

// Sent lists tips the current user has sent.
func (h *TipHandler) Sent(w http.ResponseWriter, r *http.Request) {
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

	tips, total, err := h.store.ListTipsSent(r.Context(), user.ID, limit, offset)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Database error"})
		return
	}
	if tips == nil {
		tips = []store.Tip{}
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"tips":  tips,
		"total": total,
	})
}

// Stats returns aggregate tip statistics for the current user.
func (h *TipHandler) Stats(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}

	stats, err := h.store.GetTipStats(r.Context(), user.ID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Database error"})
		return
	}

	writeJSON(w, http.StatusOK, stats)
}
