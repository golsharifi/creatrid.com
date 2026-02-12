package handler

import (
	"log"
	"net/http"

	"github.com/creatrid/creatrid/internal/email"
	"github.com/creatrid/creatrid/internal/store"
)

type DigestHandler struct {
	store *store.Store
	email *email.Service
}

func NewDigestHandler(st *store.Store, emailSvc *email.Service) *DigestHandler {
	return &DigestHandler{store: st, email: emailSvc}
}

// SendWeeklyDigest sends digest emails to all onboarded users. Intended to be called by admin or cron.
func (h *DigestHandler) SendWeeklyDigest(w http.ResponseWriter, r *http.Request) {
	if h.email == nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "Email service not configured"})
		return
	}

	users, _, err := h.store.AdminListUsers(r.Context(), 10000, 0)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch users"})
		return
	}

	sent := 0
	for _, u := range users {
		if !u.Onboarded || u.Email == "" || !u.WeeklyDigestOptIn {
			continue
		}

		summary, err := h.store.GetAnalyticsSummary(r.Context(), u.ID)
		if err != nil {
			continue
		}

		connCount := u.Connections
		name := "Creator"
		if u.Name != nil {
			name = *u.Name
		}

		subj, body := email.WeeklyDigestEmail(name, summary.TotalViews, summary.ViewsThisWeek, summary.TotalClicks, connCount, u.CreatorScore)
		if err := h.email.Send(u.Email, subj, body); err != nil {
			log.Printf("Failed to send digest to %s: %v", u.Email, err)
			continue
		}
		sent++
	}

	writeJSON(w, http.StatusOK, map[string]int{"sent": sent})
}
