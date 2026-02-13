package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/creatrid/creatrid/internal/config"
	"github.com/creatrid/creatrid/internal/middleware"
	"github.com/creatrid/creatrid/internal/store"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/account"
	"github.com/stripe/stripe-go/v81/accountlink"
)

type PayoutHandler struct {
	store  *store.Store
	config *config.Config
}

func NewPayoutHandler(st *store.Store, cfg *config.Config) *PayoutHandler {
	stripe.Key = cfg.StripeSecretKey
	return &PayoutHandler{store: st, config: cfg}
}

// ConnectOnboard starts the Stripe Connect onboarding flow.
func (h *PayoutHandler) ConnectOnboard(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}

	existingID, _, err := h.store.GetUserStripeConnectID(r.Context(), user.ID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Database error"})
		return
	}

	var accountID string
	if existingID != nil && *existingID != "" {
		accountID = *existingID
	} else {
		params := &stripe.AccountParams{
			Type: stripe.String("standard"),
		}
		params.AddMetadata("user_id", user.ID)

		acct, err := account.New(params)
		if err != nil {
			log.Printf("Stripe account creation error: %v", err)
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to create Stripe account"})
			return
		}
		accountID = acct.ID

		if err := h.store.UpdateUserStripeConnect(r.Context(), user.ID, accountID, false); err != nil {
			log.Printf("Failed to store Stripe Connect account: %v", err)
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to save account"})
			return
		}
	}

	linkParams := &stripe.AccountLinkParams{
		Account:    stripe.String(accountID),
		Type:       stripe.String("account_onboarding"),
		RefreshURL: stripe.String(h.config.FrontendURL + "/earnings?refresh=true"),
		ReturnURL:  stripe.String(h.config.BackendURL + "/api/payouts/connect/callback"),
	}

	link, err := accountlink.New(linkParams)
	if err != nil {
		log.Printf("Stripe account link error: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to create onboarding link"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"url": link.URL})
}

// ConnectCallback handles the return from Stripe Connect onboarding.
func (h *PayoutHandler) ConnectCallback(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		http.Redirect(w, r, h.config.FrontendURL+"/sign-in", http.StatusTemporaryRedirect)
		return
	}

	accountID, _, err := h.store.GetUserStripeConnectID(r.Context(), user.ID)
	if err == nil && accountID != nil && *accountID != "" {
		_ = h.store.UpdateUserStripeConnect(r.Context(), user.ID, *accountID, true)
	}

	http.Redirect(w, r, h.config.FrontendURL+"/earnings", http.StatusTemporaryRedirect)
}

// ConnectStatus returns the Stripe Connect status for the authenticated user.
func (h *PayoutHandler) ConnectStatus(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}

	accountID, onboarded, err := h.store.GetUserStripeConnectID(r.Context(), user.ID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Database error"})
		return
	}

	connected := accountID != nil && *accountID != ""

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"connected": connected,
		"onboarded": onboarded,
	})
}

// Dashboard returns the payout dashboard for the authenticated user.
func (h *PayoutHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}

	dashboard, err := h.store.GetPayoutDashboard(r.Context(), user.ID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch payout dashboard"})
		return
	}

	writeJSON(w, http.StatusOK, dashboard)
}

// ListPayouts returns paginated payouts for the authenticated user.
func (h *PayoutHandler) ListPayouts(w http.ResponseWriter, r *http.Request) {
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

	payouts, total, err := h.store.ListPayoutsByUser(r.Context(), user.ID, limit, offset)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch payouts"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"payouts": payouts,
		"total":   total,
	})
}
