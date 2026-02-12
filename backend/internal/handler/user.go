package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/creatrid/creatrid/internal/auth"
	"github.com/creatrid/creatrid/internal/config"
	"github.com/creatrid/creatrid/internal/email"
	"github.com/creatrid/creatrid/internal/middleware"
	"github.com/creatrid/creatrid/internal/storage"
	"github.com/creatrid/creatrid/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/nrednav/cuid2"
)

var (
	usernameRegex    = regexp.MustCompile(`^[a-zA-Z0-9_-]{3,30}$`)
	reservedUsernames = map[string]bool{
		"admin": true, "api": true, "auth": true, "dashboard": true,
		"settings": true, "connections": true, "onboarding": true,
		"sign-in": true, "sign-out": true, "about": true, "help": true,
		"support": true, "terms": true, "privacy": true, "blog": true,
		"profile": true,
	}
)

type UserHandler struct {
	store  *store.Store
	blob   *storage.BlobStorage
	email  *email.Service
	jwt    *auth.JWTService
	config *config.Config
}

func NewUserHandler(store *store.Store, blob *storage.BlobStorage, emailSvc *email.Service, jwtSvc *auth.JWTService, cfg *config.Config) *UserHandler {
	return &UserHandler{store: store, blob: blob, email: emailSvc, jwt: jwtSvc, config: cfg}
}

type onboardRequest struct {
	Username string `json:"username"`
	Name     string `json:"name"`
}

func (h *UserHandler) Onboard(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}

	var req onboardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	username := strings.ToLower(strings.TrimSpace(req.Username))
	name := strings.TrimSpace(req.Name)

	if name == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Name is required"})
		return
	}

	if err := validateUsername(username); err != "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err})
		return
	}

	// Check uniqueness
	existing, dbErr := h.store.FindUserByUsername(r.Context(), username)
	if dbErr != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Database error"})
		return
	}
	if existing != nil && existing.ID != user.ID {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "This username is already taken"})
		return
	}

	if err := h.store.UpdateUserOnboarding(r.Context(), user.ID, username, name); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to update profile"})
		return
	}

	recalcScore(r.Context(), h.store, user.ID)

	// Send welcome email (async, don't block response, respecting preferences)
	if h.email != nil && user.GetEmailPrefs().Welcome {
		go func() {
			profileURL := "https://creatrid.com/profile?u=" + username
			subj, body := email.WelcomeEmail(name, username, profileURL)
			if err := h.email.Send(user.Email, subj, body); err != nil {
				log.Printf("Failed to send welcome email to %s: %v", user.Email, err)
			}
		}()
	}

	writeJSON(w, http.StatusOK, map[string]bool{"success": true})
}

type customLink struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

type updateProfileRequest struct {
	Name        *string          `json:"name"`
	Bio         *string          `json:"bio"`
	Username    *string          `json:"username"`
	Theme       *string          `json:"theme"`
	CustomLinks []customLink     `json:"customLinks,omitempty"`
	EmailPrefs  *json.RawMessage `json:"emailPrefs,omitempty"`
}

func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}

	var req updateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	// Trim values
	if req.Name != nil {
		trimmed := strings.TrimSpace(*req.Name)
		req.Name = &trimmed
	}
	if req.Bio != nil {
		trimmed := strings.TrimSpace(*req.Bio)
		req.Bio = &trimmed
	}

	// Validate theme if provided
	if req.Theme != nil {
		validThemes := map[string]bool{
			"default": true, "ocean": true, "sunset": true,
			"forest": true, "midnight": true, "rose": true,
		}
		if !validThemes[*req.Theme] {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid theme"})
			return
		}
	}

	// Validate username if provided
	if req.Username != nil {
		username := strings.ToLower(strings.TrimSpace(*req.Username))
		req.Username = &username

		if err := validateUsername(username); err != "" {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err})
			return
		}

		existing, dbErr := h.store.FindUserByUsername(r.Context(), username)
		if dbErr != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Database error"})
			return
		}
		if existing != nil && existing.ID != user.ID {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "This username is already taken"})
			return
		}
	}

	if err := h.store.UpdateUserProfile(r.Context(), user.ID, req.Name, req.Bio, req.Username, req.Theme); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to update profile"})
		return
	}

	// Update custom links if provided
	if req.CustomLinks != nil {
		if len(req.CustomLinks) > 10 {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Maximum 10 custom links"})
			return
		}
		linksJSON, _ := json.Marshal(req.CustomLinks)
		if err := h.store.UpdateUserCustomLinks(r.Context(), user.ID, linksJSON); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to update links"})
			return
		}
	}

	// Update email preferences if provided
	if req.EmailPrefs != nil {
		if err := h.store.UpdateUserEmailPrefs(r.Context(), user.ID, *req.EmailPrefs); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to update email preferences"})
			return
		}
	}

	recalcScore(r.Context(), h.store, user.ID)
	writeJSON(w, http.StatusOK, map[string]bool{"success": true})
}

func (h *UserHandler) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}

	// Log deletion for audit trail
	_ = h.store.LogDeletion(r.Context(), cuid2.Generate(), user.Email)

	// Delete avatar from blob storage if exists
	if h.blob != nil && user.Image != nil && *user.Image != "" {
		_ = h.blob.Delete(r.Context(), *user.Image)
	}

	// Delete user (cascading delete handles connections, analytics, collaborations)
	if err := h.store.DeleteUser(r.Context(), user.ID); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to delete account"})
		return
	}

	// Clear auth cookie
	h.jwt.ClearTokenCookie(w, h.config.CookieDomain, h.config.CookieSecure)

	// Send confirmation email (async)
	if h.email != nil {
		go func() {
			name := "Creator"
			if user.Name != nil {
				name = *user.Name
			}
			subj, body := email.AccountDeletedEmail(name)
			_ = h.email.Send(user.Email, subj, body)
		}()
	}

	writeJSON(w, http.StatusOK, map[string]bool{"deleted": true})
}

func (h *UserHandler) ExportProfile(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}

	connections, _ := h.store.FindConnectionsByUserID(r.Context(), user.ID)
	publicConns := make([]interface{}, 0, len(connections))
	for _, c := range connections {
		publicConns = append(publicConns, c.ToPublic())
	}

	w.Header().Set("Content-Disposition", `attachment; filename="creatrid-profile.json"`)
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"user":        user.ToPublic(),
		"connections": publicConns,
		"exportedAt":  time.Now(),
	})
}

func (h *UserHandler) PublicProfile(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	if username == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Username is required"})
		return
	}

	user, err := h.store.FindUserByUsername(r.Context(), strings.ToLower(username))
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Database error"})
		return
	}
	if user == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "User not found"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{"user": user.ToPublic()})
}

var allowedImageTypes = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/webp": ".webp",
}

const maxImageSize = 5 << 20 // 5 MB

func (h *UserHandler) UploadImage(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}

	if h.blob == nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "Image upload is not configured"})
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxImageSize)
	if err := r.ParseMultipartForm(maxImageSize); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "File too large (max 5 MB)"})
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "No image file provided"})
		return
	}
	defer file.Close()

	contentType := header.Header.Get("Content-Type")
	ext, ok := allowedImageTypes[contentType]
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Only JPEG, PNG, and WebP images are allowed"})
		return
	}

	// Delete old image if it exists
	if user.Image != nil && *user.Image != "" {
		_ = h.blob.Delete(r.Context(), *user.Image)
	}

	blobName := fmt.Sprintf("avatars/%s%s", cuid2.Generate(), ext)
	imageURL, err := h.blob.Upload(r.Context(), blobName, file, contentType)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to upload image"})
		return
	}

	if err := h.store.UpdateUserImage(r.Context(), user.ID, imageURL); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to save image"})
		return
	}

	recalcScore(r.Context(), h.store, user.ID)
	writeJSON(w, http.StatusOK, map[string]string{"image": imageURL})
}

func validateUsername(username string) string {
	if !usernameRegex.MatchString(username) {
		return "Username must be 3-30 characters (letters, numbers, hyphens, underscores)"
	}
	if reservedUsernames[username] {
		return "This username is reserved"
	}
	return ""
}
