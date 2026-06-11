package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/creatrid/creatrid/internal/middleware"
	"github.com/creatrid/creatrid/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/nrednav/cuid2"
)

// TradeHandler — the closed-loop sticker exchange. Stickers trade for arena
// points, entirely in-platform: points are not cashable and stickers cannot
// leave Creatrid. External/on-chain trading (transferable crypto IDs, NFT
// markets, Polymarket compatibility) is intentionally NOT implemented here —
// it requires securities/CFTC legal sign-off and music/IP-style licensing the
// client is obtaining; see the Trade page notes and NEXT_PUBLIC_ENABLE_TRADE.
type TradeHandler struct {
	store *store.Store
}

func NewTradeHandler(st *store.Store) *TradeHandler {
	return &TradeHandler{store: st}
}

// Wallet handles GET /api/arena/wallet — balance + transaction history.
func (h *TradeHandler) Wallet(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}
	state, err := h.store.GetArenaState(r.Context(), user.ID)
	if err != nil || state == nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to load wallet"})
		return
	}
	events, err := h.store.ListPointEvents(r.Context(), user.ID, 50)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to load history"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"points": state.Points, "events": events})
}

// Listings handles GET /api/trade/listings — the open exchange floor (public).
func (h *TradeHandler) Listings(w http.ResponseWriter, r *http.Request) {
	listings, err := h.store.ListOpenListings(r.Context(), 50)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to load listings"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"listings": listings})
}

// MyListings handles GET /api/trade/listings/mine.
func (h *TradeHandler) MyListings(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}
	listings, err := h.store.ListUserListings(r.Context(), user.ID, 50)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to load listings"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"listings": listings})
}

// CreateListing handles POST /api/trade/listings {stickerId, pricePoints}.
func (h *TradeHandler) CreateListing(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}
	var req struct {
		StickerID   string `json:"stickerId"`
		PricePoints int    `json:"pricePoints"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.StickerID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "stickerId and pricePoints are required"})
		return
	}
	if req.PricePoints < 1 || req.PricePoints > 100000 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Price must be between 1 and 100,000 points"})
		return
	}

	id := cuid2.Generate()
	if err := h.store.CreateListing(r.Context(), id, user.ID, req.StickerID, req.PricePoints); err != nil {
		if err == store.ErrInsufficientStickers {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "You don't have that sticker to sell"})
			return
		}
		log.Printf("Create listing failed: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to create listing"})
		return
	}
	writeJSON(w, http.StatusCreated, map[string]string{"id": id})
}

// CancelListing handles DELETE /api/trade/listings/{id} — returns the sticker.
func (h *TradeHandler) CancelListing(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}
	id := chi.URLParam(r, "id")
	if err := h.store.CancelListing(r.Context(), id, user.ID); err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "Listing not found, not yours, or already closed"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"success": true})
}

// Buy handles POST /api/trade/listings/{id}/buy.
func (h *TradeHandler) Buy(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}
	id := chi.URLParam(r, "id")
	listing, err := h.store.BuyListing(r.Context(), cuid2.Generate(), cuid2.Generate(), id, user.ID)
	if err != nil {
		switch err {
		case store.ErrListingUnavailable:
			writeJSON(w, http.StatusConflict, map[string]string{"error": "Listing is no longer available (or it's your own)"})
		case store.ErrInsufficientPoints:
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Not enough points — explore the farm or share work to earn more"})
		default:
			log.Printf("Buy listing failed: %v", err)
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Purchase failed"})
		}
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"success": true, "listing": listing})
}
