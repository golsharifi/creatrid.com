package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/creatrid/creatrid/internal/imaging"
	"github.com/creatrid/creatrid/internal/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/nrednav/cuid2"
)

// Passport extras on UserHandler: the watermark mark (used to stamp Vault
// image uploads) and the signature track (a public Vault audio item that
// plays on the creator's public profile).

// GetWatermark handles GET /api/users/watermark.
func (h *UserHandler) GetWatermark(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}
	url, err := h.store.GetWatermarkURL(r.Context(), user.ID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Database error"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"watermarkUrl": url})
}

// UploadWatermark handles POST /api/users/watermark (multipart "image").
func (h *UserHandler) UploadWatermark(w http.ResponseWriter, r *http.Request) {
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

	buf, err := io.ReadAll(file)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Failed to read image"})
		return
	}
	if resized, _, resizeErr := imaging.ResizeImage(buf, 512, 512); resizeErr == nil {
		buf = resized
	}

	// Replace any previous watermark blob.
	if old, _ := h.store.GetWatermarkURL(r.Context(), user.ID); old != nil && *old != "" {
		_ = h.blob.Delete(r.Context(), *old)
	}

	blobName := fmt.Sprintf("watermarks/%s%s", cuid2.Generate(), ext)
	url, err := h.blob.Upload(r.Context(), blobName, bytes.NewReader(buf), contentType)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to upload watermark"})
		return
	}
	if err := h.store.SetWatermarkURL(r.Context(), user.ID, &url); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to save watermark"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"watermarkUrl": url})
}

// DeleteWatermark handles DELETE /api/users/watermark.
func (h *UserHandler) DeleteWatermark(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}
	if old, _ := h.store.GetWatermarkURL(r.Context(), user.ID); old != nil && *old != "" && h.blob != nil {
		_ = h.blob.Delete(r.Context(), *old)
	}
	if err := h.store.SetWatermarkURL(r.Context(), user.ID, nil); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to remove watermark"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"success": true})
}

// SetSignatureTrack handles POST /api/users/signature-track {contentId|null}.
func (h *UserHandler) SetSignatureTrack(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}

	var req struct {
		ContentID *string `json:"contentId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	if req.ContentID != nil {
		item, err := h.store.FindContentItemByID(r.Context(), *req.ContentID)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Database error"})
			return
		}
		switch {
		case item == nil || item.UserID != user.ID:
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "Track not found in your Vault"})
			return
		case item.ContentType != "audio":
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Signature track must be an audio file"})
			return
		case item.IsEncrypted:
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Encrypted files can't be a signature track"})
			return
		case !item.IsPublic:
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Make the track public first — visitors need to hear it"})
			return
		}
	}

	if err := h.store.SetSignatureTrack(r.Context(), user.ID, req.ContentID); err != nil {
		log.Printf("Set signature track failed: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to save signature track"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"success": true})
}

// StreamSignatureTrack handles GET /api/users/{username}/signature-track —
// public: redirects to the audio file if the creator set a public track.
func (h *UserHandler) StreamSignatureTrack(w http.ResponseWriter, r *http.Request) {
	username := strings.ToLower(chi.URLParam(r, "username"))
	user, err := h.store.FindUserByUsername(r.Context(), username)
	if err != nil || user == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "User not found"})
		return
	}
	trackID, err := h.store.GetSignatureTrackID(r.Context(), user.ID)
	if err != nil || trackID == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "No signature track"})
		return
	}
	item, err := h.store.FindContentItemByID(r.Context(), *trackID)
	if err != nil || item == nil || !item.IsPublic || item.ContentType != "audio" {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "No signature track"})
		return
	}
	http.Redirect(w, r, item.FileURL, http.StatusFound)
}
