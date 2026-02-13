package handler

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/creatrid/creatrid/internal/config"
	"github.com/creatrid/creatrid/internal/middleware"
	"github.com/creatrid/creatrid/internal/storage"
	"github.com/creatrid/creatrid/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/nrednav/cuid2"
)

const maxContentSize = 500 << 20 // 500 MB

type ContentHandler struct {
	store  *store.Store
	blob   *storage.BlobStorage
	config *config.Config
}

func NewContentHandler(st *store.Store, blob *storage.BlobStorage, cfg *config.Config) *ContentHandler {
	return &ContentHandler{
		store:  st,
		blob:   blob,
		config: cfg,
	}
}

// mimeToExt maps MIME types to file extensions.
var mimeToExt = map[string]string{
	"image/jpeg":      ".jpg",
	"image/png":       ".png",
	"image/webp":      ".webp",
	"image/gif":       ".gif",
	"video/mp4":       ".mp4",
	"audio/mpeg":      ".mp3",
	"application/pdf": ".pdf",
}

// mimeToContentType maps MIME type prefixes to content type categories.
func mimeToContentType(mime string) string {
	switch {
	case strings.HasPrefix(mime, "image/"):
		return "image"
	case strings.HasPrefix(mime, "video/"):
		return "video"
	case strings.HasPrefix(mime, "audio/"):
		return "audio"
	case strings.HasPrefix(mime, "text/"):
		return "text"
	default:
		return "other"
	}
}

// extFromMime returns a file extension for the given MIME type.
func extFromMime(mime string) string {
	if ext, ok := mimeToExt[mime]; ok {
		return ext
	}
	// Fall back to mime subtype
	parts := strings.SplitN(mime, "/", 2)
	if len(parts) == 2 {
		return "." + parts[1]
	}
	return ".bin"
}

// Upload handles POST /api/content — multipart file upload.
func (h *ContentHandler) Upload(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}

	if h.blob == nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "File upload is not configured"})
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxContentSize)
	if err := r.ParseMultipartForm(maxContentSize); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "File too large (max 500 MB)"})
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "No file provided"})
		return
	}
	defer file.Close()

	title := strings.TrimSpace(r.FormValue("title"))
	if title == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Title is required"})
		return
	}

	var description *string
	if desc := strings.TrimSpace(r.FormValue("description")); desc != "" {
		description = &desc
	}

	var tags []string
	if tagsStr := strings.TrimSpace(r.FormValue("tags")); tagsStr != "" {
		for _, t := range strings.Split(tagsStr, ",") {
			trimmed := strings.TrimSpace(t)
			if trimmed != "" {
				tags = append(tags, trimmed)
			}
		}
	}
	if tags == nil {
		tags = []string{}
	}

	isPublic := true
	if ipStr := r.FormValue("is_public"); ipStr != "" {
		isPublic = ipStr == "true"
	}

	mimeType := header.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}
	contentType := mimeToContentType(mimeType)
	ext := extFromMime(mimeType)

	// Read file into buffer so we can compute hash and then upload
	buf, err := io.ReadAll(file)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Failed to read file"})
		return
	}

	// Compute SHA-256 hash
	hashBytes := sha256.Sum256(buf)
	hashHex := hex.EncodeToString(hashBytes[:])

	contentID := cuid2.Generate()
	blobName := "vault/" + user.ID + "/" + contentID + ext

	fileURL, err := h.blob.Upload(r.Context(), blobName, bytes.NewReader(buf), mimeType)
	if err != nil {
		log.Printf("Content upload blob error: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to upload file"})
		return
	}

	now := time.Now()
	item := &store.ContentItem{
		ID:          contentID,
		UserID:      user.ID,
		Title:       title,
		Description: description,
		ContentType: contentType,
		MimeType:    mimeType,
		FileSize:    int64(len(buf)),
		FileURL:     fileURL,
		HashSHA256:  hashHex,
		IsPublic:    isPublic,
		Tags:        tags,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := h.store.CreateContentItem(r.Context(), item); err != nil {
		log.Printf("Content upload DB error: %v", err)
		// Attempt to clean up the uploaded blob
		_ = h.blob.Delete(r.Context(), fileURL)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to save content record"})
		return
	}

	writeJSON(w, http.StatusCreated, item)
}

// List handles GET /api/content — list authenticated user's content.
func (h *ContentHandler) List(w http.ResponseWriter, r *http.Request) {
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

	items, total, err := h.store.ListContentItemsByUser(r.Context(), user.ID, limit, offset)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch content"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"items": items,
		"total": total,
	})
}

// Get handles GET /api/content/{id} — get a single content item (owner only).
func (h *ContentHandler) Get(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}

	id := chi.URLParam(r, "id")
	item, err := h.store.FindContentItemByID(r.Context(), id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Database error"})
		return
	}
	if item == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "Content not found"})
		return
	}
	if item.UserID != user.ID {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "Not authorized"})
		return
	}

	writeJSON(w, http.StatusOK, item)
}

type updateContentRequest struct {
	Title       string   `json:"title"`
	Description *string  `json:"description"`
	Tags        []string `json:"tags"`
	IsPublic    bool     `json:"isPublic"`
}

// Update handles PATCH /api/content/{id} — update content metadata.
func (h *ContentHandler) Update(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}

	id := chi.URLParam(r, "id")
	item, err := h.store.FindContentItemByID(r.Context(), id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Database error"})
		return
	}
	if item == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "Content not found"})
		return
	}
	if item.UserID != user.ID {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "Not authorized"})
		return
	}

	var req updateContentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	title := strings.TrimSpace(req.Title)
	if title == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Title is required"})
		return
	}

	tags := req.Tags
	if tags == nil {
		tags = []string{}
	}

	if err := h.store.UpdateContentItem(r.Context(), id, title, req.Description, tags, req.IsPublic); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to update content"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]bool{"success": true})
}

// Delete handles DELETE /api/content/{id} — delete content and its blob.
func (h *ContentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}

	id := chi.URLParam(r, "id")
	item, err := h.store.FindContentItemByID(r.Context(), id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Database error"})
		return
	}
	if item == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "Content not found"})
		return
	}
	if item.UserID != user.ID {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "Not authorized"})
		return
	}

	// Delete blob file
	if h.blob != nil {
		_ = h.blob.Delete(r.Context(), item.FileURL)
	}

	if err := h.store.DeleteContentItem(r.Context(), id); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to delete content"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]bool{"deleted": true})
}

// Download handles GET /api/content/{id}/download — redirect to file URL if authorized.
func (h *ContentHandler) Download(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}

	id := chi.URLParam(r, "id")
	item, err := h.store.FindContentItemByID(r.Context(), id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Database error"})
		return
	}
	if item == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "Content not found"})
		return
	}

	// Check authorization: owner or has a license
	if item.UserID != user.ID {
		hasLicense, err := h.store.HasLicense(r.Context(), user.ID, item.ID)
		if err != nil || !hasLicense {
			writeJSON(w, http.StatusForbidden, map[string]string{"error": "Not authorized to download this content"})
			return
		}
	}

	// Set Content-Disposition header with the original filename
	filename := item.Title + filepath.Ext(item.FileURL)
	w.Header().Set("Content-Disposition", `attachment; filename="`+filename+`"`)
	http.Redirect(w, r, item.FileURL, http.StatusFound)
}

// PublicList handles GET /api/users/{username}/content — list a user's public content.
func (h *ContentHandler) PublicList(w http.ResponseWriter, r *http.Request) {
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

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	contentType := r.URL.Query().Get("type")
	query := r.URL.Query().Get("q")

	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	// Use ListPublicContent filtered by this user's public content
	// We fetch public content and filter by user in the query
	items, total, err := h.store.ListContentItemsByUser(r.Context(), user.ID, limit, offset)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch content"})
		return
	}

	// Filter to only public items and apply optional type/query filters
	var publicItems []map[string]interface{}
	for _, item := range items {
		if !item.IsPublic {
			continue
		}
		if contentType != "" && item.ContentType != contentType {
			continue
		}
		if query != "" {
			q := strings.ToLower(query)
			titleMatch := strings.Contains(strings.ToLower(item.Title), q)
			descMatch := item.Description != nil && strings.Contains(strings.ToLower(*item.Description), q)
			if !titleMatch && !descMatch {
				continue
			}
		}
		publicItems = append(publicItems, map[string]interface{}{
			"id":          item.ID,
			"title":       item.Title,
			"description": item.Description,
			"contentType": item.ContentType,
			"mimeType":    item.MimeType,
			"fileSize":    item.FileSize,
			"thumbnailUrl": item.ThumbnailURL,
			"hashSha256":  item.HashSHA256,
			"isPublic":    item.IsPublic,
			"tags":        item.Tags,
			"createdAt":   item.CreatedAt,
			"updatedAt":   item.UpdatedAt,
		})
	}
	if publicItems == nil {
		publicItems = []map[string]interface{}{}
	}

	_ = total // total from ListContentItemsByUser includes private; use len for public count
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"items": publicItems,
		"total": len(publicItems),
	})
}

// Proof handles GET /api/content/{id}/proof — return ownership proof for a content item.
func (h *ContentHandler) Proof(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())

	id := chi.URLParam(r, "id")
	item, err := h.store.FindContentItemByID(r.Context(), id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Database error"})
		return
	}
	if item == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "Content not found"})
		return
	}

	// Content must be public or the user must be the owner
	if !item.IsPublic {
		if user == nil || item.UserID != user.ID {
			writeJSON(w, http.StatusForbidden, map[string]string{"error": "Not authorized"})
			return
		}
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"id":         item.ID,
		"hashSha256": item.HashSHA256,
		"createdAt":  item.CreatedAt,
		"title":      item.Title,
	})
}
