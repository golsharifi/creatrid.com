package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/creatrid/creatrid/internal/auth"
	"github.com/creatrid/creatrid/internal/config"
	"github.com/creatrid/creatrid/internal/middleware"
	"github.com/creatrid/creatrid/internal/store"
)

type TwoFAHandler struct {
	store  *store.Store
	totp   *auth.TOTPService
	jwt    *auth.JWTService
	config *config.Config
}

func NewTwoFAHandler(st *store.Store, totpSvc *auth.TOTPService, jwtSvc *auth.JWTService, cfg *config.Config) *TwoFAHandler {
	return &TwoFAHandler{
		store:  st,
		totp:   totpSvc,
		jwt:    jwtSvc,
		config: cfg,
	}
}

// Setup generates a TOTP secret, stores it, and returns the QR code URL + secret.
// POST /api/auth/2fa/setup (requires auth)
func (h *TwoFAHandler) Setup(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}

	// Check if already enabled
	if user.TOTPEnabled {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "2FA is already enabled"})
		return
	}

	// Generate secret
	key, err := h.totp.GenerateSecret(user.Email)
	if err != nil {
		log.Printf("TOTP generate error: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to generate 2FA secret"})
		return
	}

	// Store secret (not yet enabled)
	if err := h.store.SetTOTPSecret(r.Context(), user.ID, key.Secret()); err != nil {
		log.Printf("TOTP store error: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to save 2FA secret"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"secret": key.Secret(),
		"qrUrl":  key.URL(),
	})
}

// Verify validates a TOTP code and enables 2FA, returning backup codes.
// POST /api/auth/2fa/verify (requires auth)
func (h *TwoFAHandler) Verify(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}

	var body struct {
		Code string `json:"code"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Code == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Code is required"})
		return
	}

	// Get stored secret
	secret, enabled, err := h.store.GetTOTPSecret(r.Context(), user.ID)
	if err != nil || secret == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "2FA not set up. Please call setup first."})
		return
	}
	if enabled {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "2FA is already enabled"})
		return
	}

	// Validate code
	if !h.totp.ValidateCode(secret, body.Code) {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid code. Please try again."})
		return
	}

	// Generate backup codes
	backupCodes, err := h.totp.GenerateBackupCodes()
	if err != nil {
		log.Printf("Backup codes error: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to generate backup codes"})
		return
	}

	// Enable 2FA
	if err := h.store.EnableTOTP(r.Context(), user.ID, backupCodes); err != nil {
		log.Printf("Enable TOTP error: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to enable 2FA"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success":     true,
		"backupCodes": backupCodes,
	})
}

// Validate checks a TOTP code (or backup code) from the temp token flow and issues a full JWT.
// POST /api/auth/2fa/validate (no auth — uses temp token from cookie)
func (h *TwoFAHandler) Validate(w http.ResponseWriter, r *http.Request) {
	// Get temp token from cookie
	cookie, err := r.Cookie("token")
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}

	userID, err := h.jwt.ValidateTempToken(cookie.Value)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Session expired. Please sign in again."})
		return
	}

	var body struct {
		Code string `json:"code"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Code == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Code is required"})
		return
	}

	// Get TOTP secret
	secret, enabled, err := h.store.GetTOTPSecret(r.Context(), userID)
	if err != nil || !enabled || secret == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "2FA is not enabled for this account"})
		return
	}

	// Try TOTP code first
	valid := h.totp.ValidateCode(secret, body.Code)

	// If not valid, try backup code
	if !valid {
		consumed, err := h.store.ConsumeTOTPBackupCode(r.Context(), userID, body.Code)
		if err != nil {
			log.Printf("Backup code check error: %v", err)
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Verification failed"})
			return
		}
		valid = consumed
	}

	if !valid {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid code. Please try again."})
		return
	}

	// Issue full JWT
	fullToken, err := h.jwt.Generate(userID)
	if err != nil {
		log.Printf("JWT generation error: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to create session"})
		return
	}

	h.jwt.SetTokenCookie(w, fullToken, h.config.CookieDomain, h.config.CookieSecure)
	writeJSON(w, http.StatusOK, map[string]bool{"success": true})
}

// Disable disables 2FA after validating the current TOTP code.
// POST /api/auth/2fa/disable (requires auth)
func (h *TwoFAHandler) Disable(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}

	var body struct {
		Code string `json:"code"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Code == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Code is required"})
		return
	}

	// Get stored secret
	secret, enabled, err := h.store.GetTOTPSecret(r.Context(), user.ID)
	if err != nil || !enabled || secret == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "2FA is not enabled"})
		return
	}

	// Validate code
	if !h.totp.ValidateCode(secret, body.Code) {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid code. Please try again."})
		return
	}

	// Disable 2FA
	if err := h.store.DisableTOTP(r.Context(), user.ID); err != nil {
		log.Printf("Disable TOTP error: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to disable 2FA"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]bool{"success": true})
}
