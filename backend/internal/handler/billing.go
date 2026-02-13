package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/creatrid/creatrid/internal/config"
	"github.com/creatrid/creatrid/internal/middleware"
	"github.com/creatrid/creatrid/internal/store"
	"github.com/nrednav/cuid2"
	"github.com/stripe/stripe-go/v81"
	portalsession "github.com/stripe/stripe-go/v81/billingportal/session"
	checkoutsession "github.com/stripe/stripe-go/v81/checkout/session"
	"github.com/stripe/stripe-go/v81/customer"
	"github.com/stripe/stripe-go/v81/webhook"
)

type BillingHandler struct {
	store  *store.Store
	config *config.Config
}

func NewBillingHandler(store *store.Store, cfg *config.Config) *BillingHandler {
	stripe.Key = cfg.StripeSecretKey
	return &BillingHandler{store: store, config: cfg}
}

type checkoutRequest struct {
	Plan string `json:"plan"`
}

func (h *BillingHandler) CreateCheckout(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}

	var req checkoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	var priceID string
	switch req.Plan {
	case "pro":
		priceID = h.config.StripePricePro
	case "business":
		priceID = h.config.StripePriceBusiness
	default:
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid plan. Must be 'pro' or 'business'"})
		return
	}

	if priceID == "" {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "Pricing not configured for this plan"})
		return
	}

	// Find or create Stripe customer
	customerID, err := h.findOrCreateCustomer(r, user.Email, user.ID)
	if err != nil {
		log.Printf("Stripe customer error: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to create checkout session"})
		return
	}

	// Create Checkout Session
	params := &stripe.CheckoutSessionParams{
		Mode:     stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		Customer: stripe.String(customerID),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(priceID),
				Quantity: stripe.Int64(1),
			},
		},
		SuccessURL: stripe.String(h.config.FrontendURL + "/billing?success=true"),
		CancelURL:  stripe.String(h.config.FrontendURL + "/billing?canceled=true"),
		SubscriptionData: &stripe.CheckoutSessionSubscriptionDataParams{
			Metadata: map[string]string{
				"user_id": user.ID,
				"plan":    req.Plan,
			},
		},
	}
	params.AddMetadata("user_id", user.ID)
	params.AddMetadata("plan", req.Plan)

	session, err := checkoutsession.New(params)
	if err != nil {
		log.Printf("Stripe checkout error: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to create checkout session"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"url": session.URL})
}

func (h *BillingHandler) GetSubscription(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}

	sub, err := h.store.FindSubscriptionByUserID(r.Context(), user.ID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Database error"})
		return
	}

	if sub == nil {
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"plan":             "free",
			"status":           "active",
			"currentPeriodEnd": nil,
		})
		return
	}

	writeJSON(w, http.StatusOK, sub)
}

func (h *BillingHandler) CreatePortal(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}

	sub, err := h.store.FindSubscriptionByUserID(r.Context(), user.ID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Database error"})
		return
	}

	if sub == nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "No active subscription found"})
		return
	}

	params := &stripe.BillingPortalSessionParams{
		Customer:  stripe.String(sub.StripeCustomerID),
		ReturnURL: stripe.String(h.config.FrontendURL + "/billing"),
	}

	session, err := portalsession.New(params)
	if err != nil {
		log.Printf("Stripe portal error: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to create portal session"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"url": session.URL})
}

func (h *BillingHandler) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Failed to read body"})
		return
	}

	event, err := webhook.ConstructEvent(body, r.Header.Get("Stripe-Signature"), h.config.StripeWebhookSecret)
	if err != nil {
		log.Printf("Stripe webhook signature verification failed: %v", err)
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid signature"})
		return
	}

	// Idempotency check
	exists, err := h.store.FindPaymentEventByStripeID(r.Context(), event.ID)
	if err != nil {
		log.Printf("Stripe webhook idempotency check error: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Database error"})
		return
	}
	if exists {
		writeJSON(w, http.StatusOK, map[string]string{"status": "already_processed"})
		return
	}

	switch event.Type {
	case "checkout.session.completed":
		h.handleCheckoutCompleted(r, event)
	case "customer.subscription.updated":
		h.handleSubscriptionUpdated(r, event)
	case "customer.subscription.deleted":
		h.handleSubscriptionDeleted(r, event)
	case "invoice.payment_failed":
		h.handlePaymentFailed(r, event)
	}

	// Save event for audit trail
	userID := h.extractUserIDFromEvent(event)
	if saveErr := h.store.CreatePaymentEvent(r.Context(), cuid2.Generate(), userID, event.ID, string(event.Type), body); saveErr != nil {
		log.Printf("Failed to save payment event: %v", saveErr)
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *BillingHandler) handleCheckoutCompleted(r *http.Request, event stripe.Event) {
	var session stripe.CheckoutSession
	if err := json.Unmarshal(event.Data.Raw, &session); err != nil {
		log.Printf("Stripe webhook: failed to parse checkout session: %v", err)
		return
	}

	// Check if this is a license purchase
	if session.Metadata != nil && session.Metadata["type"] == "license" {
		h.handleLicensePurchaseCompleted(r, session)
		return
	}

	customerID := session.Customer.ID
	subscriptionID := ""
	if session.Subscription != nil {
		subscriptionID = session.Subscription.ID
	}

	plan := "pro"
	userID := ""
	if session.Metadata != nil {
		if p, ok := session.Metadata["plan"]; ok {
			plan = p
		}
		if u, ok := session.Metadata["user_id"]; ok {
			userID = u
		}
	}
	// Fallback: look up by customer ID
	if userID == "" {
		sub, _ := h.store.FindSubscriptionByStripeCustomerID(r.Context(), customerID)
		if sub != nil {
			userID = sub.UserID
		}
	}

	if userID == "" {
		log.Printf("Stripe webhook: checkout.session.completed but no user_id found")
		return
	}

	now := time.Now()
	sub := &store.Subscription{
		ID:                   cuid2.Generate(),
		UserID:               userID,
		StripeCustomerID:     customerID,
		StripeSubscriptionID: &subscriptionID,
		Plan:                 plan,
		Status:               "active",
		CreatedAt:            now,
		UpdatedAt:            now,
	}

	if err := h.store.UpsertSubscription(r.Context(), sub); err != nil {
		log.Printf("Stripe webhook: failed to upsert subscription: %v", err)
	}
}

func (h *BillingHandler) handleLicensePurchaseCompleted(r *http.Request, session stripe.CheckoutSession) {
	offeringID := session.Metadata["offering_id"]
	contentID := session.Metadata["content_id"]
	buyerUserID := session.Metadata["buyer_user_id"]

	if offeringID == "" || contentID == "" || buyerUserID == "" {
		log.Printf("Stripe webhook: license purchase missing metadata (offering=%s, content=%s, buyer=%s)", offeringID, contentID, buyerUserID)
		return
	}

	// Get offering for price details
	offering, err := h.store.FindOfferingByID(r.Context(), offeringID)
	if err != nil || offering == nil {
		log.Printf("Stripe webhook: failed to find offering %s: %v", offeringID, err)
		return
	}

	// Get buyer email
	buyerEmail := ""
	buyer, err := h.store.FindUserByID(r.Context(), buyerUserID)
	if err == nil && buyer != nil {
		buyerEmail = buyer.Email
	}

	amountCents := offering.PriceCents
	platformFeeCents := amountCents * 15 / 100
	creatorPayoutCents := amountCents - platformFeeCents

	stripeSessionID := session.ID
	purchase := &store.LicensePurchase{
		ID:                 cuid2.Generate(),
		OfferingID:         offeringID,
		ContentID:          contentID,
		BuyerUserID:        &buyerUserID,
		BuyerEmail:         buyerEmail,
		StripeSessionID:    &stripeSessionID,
		AmountCents:        amountCents,
		PlatformFeeCents:   platformFeeCents,
		CreatorPayoutCents: creatorPayoutCents,
		Status:             "completed",
		CreatedAt:          time.Now(),
	}

	if err := h.store.CreateLicensePurchase(r.Context(), purchase); err != nil {
		log.Printf("Stripe webhook: failed to create license purchase: %v", err)
		return
	}

	log.Printf("License purchase created: buyer=%s content=%s offering=%s amount=%d", buyerUserID, contentID, offeringID, amountCents)

	// Notify the content creator about the sale
	content, err := h.store.FindContentItemByID(r.Context(), contentID)
	if err == nil && content != nil {
		notif := &store.Notification{
			ID:        cuid2.Generate(),
			UserID:    content.UserID,
			Type:      "license_sale",
			Title:     "New license sale!",
			Message:   fmt.Sprintf("Your content \"%s\" was licensed for $%.2f", content.Title, float64(amountCents)/100),
			Data:      []byte(fmt.Sprintf(`{"contentId":"%s","purchaseId":"%s","amountCents":%d}`, contentID, purchase.ID, amountCents)),
			CreatedAt: time.Now(),
		}
		if err := h.store.CreateNotification(r.Context(), notif); err != nil {
			log.Printf("Failed to create sale notification: %v", err)
		}
	}
}

func (h *BillingHandler) handleSubscriptionUpdated(r *http.Request, event stripe.Event) {
	var sub stripe.Subscription
	if err := json.Unmarshal(event.Data.Raw, &sub); err != nil {
		log.Printf("Stripe webhook: failed to parse subscription: %v", err)
		return
	}

	customerID := sub.Customer.ID
	existing, _ := h.store.FindSubscriptionByStripeCustomerID(r.Context(), customerID)
	if existing == nil {
		log.Printf("Stripe webhook: subscription.updated but no local subscription found for customer %s", customerID)
		return
	}

	plan := existing.Plan
	if len(sub.Items.Data) > 0 {
		// Determine plan from price ID
		priceID := sub.Items.Data[0].Price.ID
		switch priceID {
		case h.config.StripePricePro:
			plan = "pro"
		case h.config.StripePriceBusiness:
			plan = "business"
		}
	}

	status := string(sub.Status)
	var periodEnd *time.Time
	if sub.CurrentPeriodEnd > 0 {
		t := time.Unix(sub.CurrentPeriodEnd, 0)
		periodEnd = &t
	}

	if err := h.store.UpdateSubscriptionPlan(r.Context(), existing.UserID, plan, status, periodEnd); err != nil {
		log.Printf("Stripe webhook: failed to update subscription plan: %v", err)
	}
}

func (h *BillingHandler) handleSubscriptionDeleted(r *http.Request, event stripe.Event) {
	var sub stripe.Subscription
	if err := json.Unmarshal(event.Data.Raw, &sub); err != nil {
		log.Printf("Stripe webhook: failed to parse subscription: %v", err)
		return
	}

	customerID := sub.Customer.ID
	existing, _ := h.store.FindSubscriptionByStripeCustomerID(r.Context(), customerID)
	if existing == nil {
		return
	}

	if err := h.store.UpdateSubscriptionPlan(r.Context(), existing.UserID, "free", "canceled", nil); err != nil {
		log.Printf("Stripe webhook: failed to cancel subscription: %v", err)
	}
}

func (h *BillingHandler) handlePaymentFailed(r *http.Request, event stripe.Event) {
	var invoice stripe.Invoice
	if err := json.Unmarshal(event.Data.Raw, &invoice); err != nil {
		log.Printf("Stripe webhook: failed to parse invoice: %v", err)
		return
	}

	customerID := invoice.Customer.ID
	existing, _ := h.store.FindSubscriptionByStripeCustomerID(r.Context(), customerID)
	if existing == nil {
		return
	}

	if err := h.store.UpdateSubscriptionPlan(r.Context(), existing.UserID, existing.Plan, "past_due", existing.CurrentPeriodEnd); err != nil {
		log.Printf("Stripe webhook: failed to set past_due status: %v", err)
	}
}

func (h *BillingHandler) findOrCreateCustomer(r *http.Request, email, userID string) (string, error) {
	// Check if we already have a subscription record with a customer ID
	sub, err := h.store.FindSubscriptionByUserID(r.Context(), userID)
	if err != nil {
		return "", err
	}
	if sub != nil {
		return sub.StripeCustomerID, nil
	}

	// Create a new Stripe customer
	params := &stripe.CustomerParams{
		Email: stripe.String(email),
	}
	params.AddMetadata("user_id", userID)

	c, err := customer.New(params)
	if err != nil {
		return "", err
	}

	// Store the customer ID for future lookups
	now := time.Now()
	newSub := &store.Subscription{
		ID:               cuid2.Generate(),
		UserID:           userID,
		StripeCustomerID: c.ID,
		Plan:             "free",
		Status:           "active",
		CreatedAt:        now,
		UpdatedAt:        now,
	}
	if err := h.store.UpsertSubscription(r.Context(), newSub); err != nil {
		log.Printf("Warning: failed to store customer ID: %v", err)
	}

	return c.ID, nil
}

func (h *BillingHandler) extractUserIDFromEvent(event stripe.Event) string {
	// Try to extract user_id from various event types
	var raw map[string]interface{}
	if err := json.Unmarshal(event.Data.Raw, &raw); err != nil {
		return ""
	}

	// Try metadata first
	if meta, ok := raw["metadata"].(map[string]interface{}); ok {
		if uid, ok := meta["user_id"].(string); ok {
			return uid
		}
	}

	// Try customer ID lookup
	if customerID, ok := raw["customer"].(string); ok && customerID != "" {
		sub, _ := h.store.FindSubscriptionByStripeCustomerID(context.Background(), customerID)
		if sub != nil {
			return sub.UserID
		}
	}

	return ""
}
