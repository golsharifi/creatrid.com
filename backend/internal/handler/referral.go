package handler

import (
	"net/http"

	"github.com/creatrid/creatrid/internal/middleware"
	"github.com/creatrid/creatrid/internal/store"
	"github.com/nrednav/cuid2"
)

type ReferralHandler struct {
	store *store.Store
}

func NewReferralHandler(st *store.Store) *ReferralHandler {
	return &ReferralHandler{store: st}
}

// GetCode returns the user's referral code, generating one if it doesn't exist.
func (h *ReferralHandler) GetCode(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}

	code, err := h.store.GetUserReferralCode(r.Context(), user.ID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Database error"})
		return
	}

	if code == nil || *code == "" {
		// Generate a new 8-char referral code
		newCode := cuid2.Generate()[:8]
		if err := h.store.SetUserReferralCode(r.Context(), user.ID, newCode); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to generate code"})
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"code": newCode})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"code": *code})
}

// List returns the user's referrals and stats.
func (h *ReferralHandler) List(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}

	rewards, err := h.store.ListReferralsByUser(r.Context(), user.ID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch referrals"})
		return
	}

	stats, err := h.store.GetReferralStats(r.Context(), user.ID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch stats"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"referrals": rewards,
		"stats":     stats,
	})
}
