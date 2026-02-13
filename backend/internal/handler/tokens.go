package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/creatrid/creatrid/internal/config"
	"github.com/creatrid/creatrid/internal/middleware"
	"github.com/creatrid/creatrid/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/nrednav/cuid2"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/paymentintent"
)

// TokenHandler manages creator token endpoints.
type TokenHandler struct {
	store  *store.Store
	config *config.Config
}

// NewTokenHandler creates a new TokenHandler.
func NewTokenHandler(st *store.Store, cfg *config.Config) *TokenHandler {
	return &TokenHandler{store: st, config: cfg}
}

var symbolRegex = regexp.MustCompile(`^[A-Z0-9]{2,10}$`)

type createTokenRequest struct {
	Name        string `json:"name"`
	Symbol      string `json:"symbol"`
	Description string `json:"description"`
	PriceCents  int    `json:"priceCents"`
}

// Create creates a new creator token.
func (h *TokenHandler) Create(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}

	var req createTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Symbol = strings.ToUpper(strings.TrimSpace(req.Symbol))

	if req.Name == "" || req.Symbol == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Name and symbol are required"})
		return
	}

	if !symbolRegex.MatchString(req.Symbol) {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Symbol must be 2-10 uppercase letters/numbers"})
		return
	}

	if req.PriceCents < 1 {
		req.PriceCents = 100
	}

	// Check if user already has a token
	existing, err := h.store.FindTokenByUserID(r.Context(), user.ID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Database error"})
		return
	}
	if existing != nil {
		writeJSON(w, http.StatusConflict, map[string]string{"error": "You already have a creator token"})
		return
	}

	// Check if symbol is taken
	existingSymbol, err := h.store.FindTokenBySymbol(r.Context(), req.Symbol)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Database error"})
		return
	}
	if existingSymbol != nil {
		writeJSON(w, http.StatusConflict, map[string]string{"error": "Symbol already taken"})
		return
	}

	var desc *string
	if req.Description != "" {
		d := strings.TrimSpace(req.Description)
		desc = &d
	}

	now := time.Now()
	token := &store.CreatorToken{
		ID:          cuid2.Generate(),
		UserID:      user.ID,
		Name:        req.Name,
		Symbol:      req.Symbol,
		Description: desc,
		TotalSupply: 0,
		PriceCents:  req.PriceCents,
		IsActive:    true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := h.store.CreateCreatorToken(r.Context(), token); err != nil {
		log.Printf("Failed to create token: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to create token"})
		return
	}

	writeJSON(w, http.StatusCreated, map[string]interface{}{"token": token})
}

// Get returns the current user's token.
func (h *TokenHandler) Get(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}

	token, err := h.store.FindTokenByUserID(r.Context(), user.ID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Database error"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{"token": token})
}

type updateTokenRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	PriceCents  *int    `json:"priceCents"`
	IsActive    *bool   `json:"isActive"`
}

// Update modifies the current user's token.
func (h *TokenHandler) Update(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}

	token, err := h.store.FindTokenByUserID(r.Context(), user.ID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Database error"})
		return
	}
	if token == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "No token found"})
		return
	}

	var req updateTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	name := token.Name
	if req.Name != nil {
		name = strings.TrimSpace(*req.Name)
	}
	desc := token.Description
	if req.Description != nil {
		d := strings.TrimSpace(*req.Description)
		desc = &d
	}
	priceCents := token.PriceCents
	if req.PriceCents != nil && *req.PriceCents > 0 {
		priceCents = *req.PriceCents
	}
	isActive := token.IsActive
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	if err := h.store.UpdateCreatorToken(r.Context(), token.ID, name, desc, priceCents, isActive); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to update token"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]bool{"success": true})
}

// Holders lists the holders of a specific token.
func (h *TokenHandler) Holders(w http.ResponseWriter, r *http.Request) {
	tokenID := chi.URLParam(r, "id")
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	holders, total, err := h.store.ListTokenHolders(r.Context(), tokenID, limit, offset)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Database error"})
		return
	}
	if holders == nil {
		holders = []store.TokenBalance{}
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"holders": holders,
		"total":   total,
	})
}

// Transactions lists the transactions for a specific token.
func (h *TokenHandler) Transactions(w http.ResponseWriter, r *http.Request) {
	tokenID := chi.URLParam(r, "id")
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	txs, total, err := h.store.ListTokenTransactions(r.Context(), tokenID, limit, offset)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Database error"})
		return
	}
	if txs == nil {
		txs = []store.TokenTransaction{}
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"transactions": txs,
		"total":        total,
	})
}

type purchaseTokenRequest struct {
	Amount int `json:"amount"`
}

// Purchase buys tokens — creates a Stripe payment intent.
func (h *TokenHandler) Purchase(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}

	tokenID := chi.URLParam(r, "id")

	var req purchaseTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	if req.Amount <= 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Amount must be positive"})
		return
	}

	token, err := h.store.FindTokenByID(r.Context(), tokenID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Database error"})
		return
	}
	if token == nil || !token.IsActive {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "Token not found or inactive"})
		return
	}

	// Cannot buy your own token
	if token.UserID == user.ID {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Cannot purchase your own token"})
		return
	}

	totalCents := int64(token.PriceCents) * int64(req.Amount)

	if h.config.StripeSecretKey == "" {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "Payment processing is not configured"})
		return
	}

	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(totalCents),
		Currency: stripe.String(string(stripe.CurrencyUSD)),
	}
	params.AddMetadata("type", "token_purchase")
	params.AddMetadata("token_id", tokenID)
	params.AddMetadata("buyer_user_id", user.ID)
	params.AddMetadata("amount", strconv.Itoa(req.Amount))

	pi, err := paymentintent.New(params)
	if err != nil {
		log.Printf("Stripe payment intent error: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Payment failed"})
		return
	}

	purchaseID := cuid2.Generate()
	if err := h.store.MintTokens(r.Context(), tokenID, user.ID, req.Amount, "purchase", purchaseID); err != nil {
		log.Printf("Failed to mint tokens: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to mint tokens"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success":      true,
		"clientSecret": pi.ClientSecret,
		"amountCents":  totalCents,
		"tokensMinted": req.Amount,
	})
}

// PublicToken returns a creator's token info by username.
func (h *TokenHandler) PublicToken(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")

	user, err := h.store.FindUserByUsername(r.Context(), username)
	if err != nil || user == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "User not found"})
		return
	}

	token, err := h.store.FindTokenByUserID(r.Context(), user.ID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Database error"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{"token": token})
}
