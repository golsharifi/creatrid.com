package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/creatrid/creatrid/internal/config"
	"github.com/creatrid/creatrid/internal/middleware"
	"github.com/creatrid/creatrid/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/nrednav/cuid2"
	"github.com/stripe/stripe-go/v81"
	checkoutsession "github.com/stripe/stripe-go/v81/checkout/session"
)

type LicenseHandler struct {
	store  *store.Store
	config *config.Config
}

func NewLicenseHandler(st *store.Store, cfg *config.Config) *LicenseHandler {
	stripe.Key = cfg.StripeSecretKey
	return &LicenseHandler{store: st, config: cfg}
}

// CreateOffering creates a new license offering for a content item.
// POST /api/content/{id}/licenses
func (h *LicenseHandler) CreateOffering(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}

	contentID := chi.URLParam(r, "id")

	// Verify user owns the content
	content, err := h.store.FindContentItemByID(r.Context(), contentID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Database error"})
		return
	}
	if content == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "Content not found"})
		return
	}
	if content.UserID != user.ID {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "Not authorized"})
		return
	}

	var req struct {
		LicenseType string  `json:"licenseType"`
		PriceCents  int     `json:"priceCents"`
		TermsText   *string `json:"termsText"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	validTypes := map[string]bool{"personal": true, "commercial": true, "editorial": true, "ai_training": true}
	if !validTypes[req.LicenseType] {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid license type. Must be personal, commercial, editorial, or ai_training"})
		return
	}

	if req.PriceCents <= 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Price must be greater than zero"})
		return
	}

	offering := &store.LicenseOffering{
		ID:          cuid2.Generate(),
		ContentID:   contentID,
		LicenseType: req.LicenseType,
		PriceCents:  req.PriceCents,
		Currency:    "usd",
		IsActive:    true,
		TermsText:   req.TermsText,
		CreatedAt:   time.Now(),
	}

	if err := h.store.CreateLicenseOffering(r.Context(), offering); err != nil {
		log.Printf("Failed to create license offering: %v", err)
		writeJSON(w, http.StatusConflict, map[string]string{"error": "License offering for this type already exists on this content"})
		return
	}

	writeJSON(w, http.StatusCreated, offering)
}

// ListOfferings returns all active license offerings for a content item.
// GET /api/content/{id}/licenses
func (h *LicenseHandler) ListOfferings(w http.ResponseWriter, r *http.Request) {
	contentID := chi.URLParam(r, "id")

	offerings, err := h.store.ListOfferingsByContent(r.Context(), contentID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch offerings"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{"offerings": offerings})
}

// UpdateOffering updates a license offering.
// PATCH /api/licenses/{id}
func (h *LicenseHandler) UpdateOffering(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}

	offeringID := chi.URLParam(r, "id")

	offering, err := h.store.FindOfferingByID(r.Context(), offeringID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Database error"})
		return
	}
	if offering == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "Offering not found"})
		return
	}

	// Verify user owns the content behind this offering
	content, err := h.store.FindContentItemByID(r.Context(), offering.ContentID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Database error"})
		return
	}
	if content == nil || content.UserID != user.ID {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "Not authorized"})
		return
	}

	var req struct {
		PriceCents *int    `json:"priceCents"`
		IsActive   *bool   `json:"isActive"`
		TermsText  *string `json:"termsText"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	// Apply updates, keeping existing values for fields not provided
	priceCents := offering.PriceCents
	if req.PriceCents != nil {
		priceCents = *req.PriceCents
	}
	isActive := offering.IsActive
	if req.IsActive != nil {
		isActive = *req.IsActive
	}
	terms := offering.TermsText
	if req.TermsText != nil {
		terms = req.TermsText
	}

	if priceCents <= 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Price must be greater than zero"})
		return
	}

	if err := h.store.UpdateOffering(r.Context(), offeringID, priceCents, isActive, terms); err != nil {
		log.Printf("Failed to update offering: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to update offering"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "updated"})
}

// DeleteOffering deletes a license offering.
// DELETE /api/licenses/{id}
func (h *LicenseHandler) DeleteOffering(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}

	offeringID := chi.URLParam(r, "id")

	offering, err := h.store.FindOfferingByID(r.Context(), offeringID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Database error"})
		return
	}
	if offering == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "Offering not found"})
		return
	}

	// Verify user owns the content behind this offering
	content, err := h.store.FindContentItemByID(r.Context(), offering.ContentID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Database error"})
		return
	}
	if content == nil || content.UserID != user.ID {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "Not authorized"})
		return
	}

	if err := h.store.DeleteOffering(r.Context(), offeringID); err != nil {
		log.Printf("Failed to delete offering: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to delete offering"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

// Checkout creates a Stripe checkout session for a one-time license purchase.
// POST /api/licenses/{id}/checkout
func (h *LicenseHandler) Checkout(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}

	offeringID := chi.URLParam(r, "id")

	offering, err := h.store.FindOfferingByID(r.Context(), offeringID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Database error"})
		return
	}
	if offering == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "Offering not found"})
		return
	}
	if !offering.IsActive {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "This license offering is no longer available"})
		return
	}

	// Find content for the title
	content, err := h.store.FindContentItemByID(r.Context(), offering.ContentID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Database error"})
		return
	}
	if content == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "Content not found"})
		return
	}

	// Don't allow buying your own content
	if content.UserID == user.ID {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Cannot purchase a license for your own content"})
		return
	}

	// Check if user already has a license for this content
	hasLicense, err := h.store.HasLicense(r.Context(), user.ID, content.ID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Database error"})
		return
	}
	if hasLicense {
		writeJSON(w, http.StatusConflict, map[string]string{"error": "You already have a license for this content"})
		return
	}

	successURL := h.config.FrontendURL + "/purchases?success=true"
	cancelURL := h.config.FrontendURL + "/marketplace/item?id=" + content.ID + "&canceled=true"

	params := &stripe.CheckoutSessionParams{
		Mode: stripe.String(string(stripe.CheckoutSessionModePayment)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency:   stripe.String(offering.Currency),
					UnitAmount: stripe.Int64(int64(offering.PriceCents)),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String(content.Title + " - " + offering.LicenseType + " License"),
					},
				},
				Quantity: stripe.Int64(1),
			},
		},
		SuccessURL: stripe.String(successURL),
		CancelURL:  stripe.String(cancelURL),
	}
	params.AddMetadata("type", "license")
	params.AddMetadata("offering_id", offering.ID)
	params.AddMetadata("content_id", offering.ContentID)
	params.AddMetadata("buyer_user_id", user.ID)

	session, err := checkoutsession.New(params)
	if err != nil {
		log.Printf("Stripe checkout error: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to create checkout session"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"url": session.URL})
}

// Purchases lists license purchases made by the current user.
// GET /api/licenses/purchases
func (h *LicenseHandler) Purchases(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}

	purchases, err := h.store.ListPurchasesByBuyer(r.Context(), user.ID)
	if err != nil {
		log.Printf("Failed to fetch purchases: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch purchases"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{"purchases": purchases})
}

// Sales lists license sales for the current user's content.
// GET /api/licenses/sales
func (h *LicenseHandler) Sales(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}

	sales, err := h.store.ListSalesByCreator(r.Context(), user.ID)
	if err != nil {
		log.Printf("Failed to fetch sales: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch sales"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{"sales": sales})
}
