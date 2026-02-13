package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/creatrid/creatrid/internal/middleware"
	"github.com/creatrid/creatrid/internal/model"
	"github.com/creatrid/creatrid/internal/store"
)

func adminAudit(st *store.Store, r *http.Request, action, targetType, targetID string, details map[string]interface{}) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		return
	}
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.RemoteAddr
	}
	var detailsJSON []byte
	if details != nil {
		detailsJSON, _ = json.Marshal(details)
	}
	var targetPtr *string
	if targetID != "" {
		targetPtr = &targetID
	}
	_ = st.CreateAuditEntry(r.Context(), user.ID, action, targetType, targetPtr, detailsJSON, &ip)
}

type AdminHandler struct {
	store *store.Store
}

func NewAdminHandler(st *store.Store) *AdminHandler {
	return &AdminHandler{store: st}
}

func (h *AdminHandler) Stats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.store.AdminGetStats(r.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch stats"})
		return
	}
	writeJSON(w, http.StatusOK, stats)
}

func (h *AdminHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	users, total, err := h.store.AdminListUsers(r.Context(), limit, offset)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch users"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"users": users,
		"total": total,
	})
}

func (h *AdminHandler) SetVerified(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID   string `json:"userId"`
		Verified bool   `json:"verified"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		return
	}

	if req.UserID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "userId is required"})
		return
	}

	if err := h.store.AdminSetVerified(r.Context(), req.UserID, req.Verified); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to update"})
		return
	}

	action := "verify_user"
	if !req.Verified {
		action = "unverify_user"
	}
	adminAudit(h.store, r, action, "user", req.UserID, map[string]interface{}{"verified": req.Verified})

	writeJSON(w, http.StatusOK, map[string]bool{"success": true})
}

func (h *AdminHandler) AuditLog(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	entries, total, err := h.store.ListAuditLog(r.Context(), limit, offset)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch audit log"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{"entries": entries, "total": total})
}

// RequireAdmin middleware checks that the authenticated user has admin role
func RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := middleware.UserFromContext(r.Context())
		if user == nil || user.Role != model.RoleAdmin {
			writeJSON(w, http.StatusForbidden, map[string]string{"error": "Admin access required"})
			return
		}
		next.ServeHTTP(w, r)
	})
}
