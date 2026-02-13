package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/creatrid/creatrid/internal/middleware"
	"github.com/creatrid/creatrid/internal/store"
)

type ErrorLogHandler struct {
	store *store.Store
}

func NewErrorLogHandler(store *store.Store) *ErrorLogHandler {
	return &ErrorLogHandler{store: store}
}

type ReportErrorRequest struct {
	Source    string          `json:"source"`
	Level     string          `json:"level"`
	Message   string          `json:"message"`
	Stack     *string         `json:"stack,omitempty"`
	URL       *string         `json:"url,omitempty"`
	UserAgent *string         `json:"user_agent,omitempty"`
	Metadata  json.RawMessage `json:"metadata,omitempty"`
}

func (h *ErrorLogHandler) Report(w http.ResponseWriter, r *http.Request) {
	var req ReportErrorRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Message == "" {
		http.Error(w, "message is required", http.StatusBadRequest)
		return
	}

	if req.Source == "" {
		req.Source = "frontend"
	}

	if req.Level == "" {
		req.Level = "error"
	}

	// Use User-Agent from header if not provided in body
	if req.UserAgent == nil {
		ua := r.Header.Get("User-Agent")
		if ua != "" {
			req.UserAgent = &ua
		}
	}

	// Try to get user ID from context (if authenticated)
	var userID *string
	if user := middleware.UserFromContext(r.Context()); user != nil {
		userID = &user.ID
	}

	err := h.store.CreateErrorLog(r.Context(), req.Source, req.Level, req.Message, req.Stack, req.URL, req.UserAgent, userID, req.Metadata)
	if err != nil {
		http.Error(w, "Failed to log error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

type ListErrorsResponse struct {
	Entries []store.ErrorLogEntry `json:"entries"`
	Total   int                   `json:"total"`
}

func (h *ErrorLogHandler) List(w http.ResponseWriter, r *http.Request) {
	source := r.URL.Query().Get("source")

	limit := 50
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	offset := 0
	if o := r.URL.Query().Get("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	entries, total, err := h.store.ListErrorLog(r.Context(), source, limit, offset)
	if err != nil {
		http.Error(w, "Failed to fetch errors", http.StatusInternalServerError)
		return
	}

	if entries == nil {
		entries = []store.ErrorLogEntry{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ListErrorsResponse{
		Entries: entries,
		Total:   total,
	})
}
