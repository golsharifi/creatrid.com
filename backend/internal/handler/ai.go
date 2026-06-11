package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/creatrid/creatrid/internal/ai"
	"github.com/creatrid/creatrid/internal/middleware"
	"github.com/creatrid/creatrid/internal/model"
	"github.com/creatrid/creatrid/internal/store"
	"github.com/nrednav/cuid2"
)

// AIHandler exposes the Brand Studio + Legal AI assistant endpoints.
// Only wired up when ANTHROPIC_API_KEY is configured.
type AIHandler struct {
	svc   *ai.Service
	store *store.Store
}

func NewAIHandler(svc *ai.Service, st *store.Store) *AIHandler {
	return &AIHandler{svc: svc, store: st}
}

const maxAIInputLen = 8000

// Monthly generation quotas per plan. AI calls cost real money; the free tier
// is a taste, paid tiers cover their own cost many times over.
var aiQuotas = map[string]int{
	"free":     15,
	"pro":      150,
	"business": 500,
}

func quotaForPlan(plan string) int {
	if q, ok := aiQuotas[plan]; ok {
		return q
	}
	return aiQuotas["free"]
}

// resolvePlan returns the user's effective plan for quota purposes.
func (h *AIHandler) resolvePlan(r *http.Request, user *model.User) string {
	if user.Role == model.RoleAdmin {
		return "business"
	}
	sub, err := h.store.FindSubscriptionByUserID(r.Context(), user.ID)
	if err != nil || sub == nil {
		return "free"
	}
	if sub.Status == "active" || sub.Status == "trialing" {
		return sub.Plan
	}
	return "free"
}

func monthStart() time.Time {
	now := time.Now().UTC()
	return time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
}

// checkQuota enforces the monthly cap. Returns (remaining, plan, ok); on !ok
// it has already written the 429 response.
func (h *AIHandler) checkQuota(w http.ResponseWriter, r *http.Request) (int, string, bool) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return 0, "", false
	}
	plan := h.resolvePlan(r, user)
	limit := quotaForPlan(plan)
	used, err := h.store.CountAIGenerationsSince(r.Context(), user.ID, monthStart())
	if err != nil {
		log.Printf("AI quota check failed: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Quota check failed"})
		return 0, "", false
	}
	if used >= limit {
		writeJSON(w, http.StatusTooManyRequests, map[string]any{
			"error":   "Monthly AI generation limit reached",
			"used":    used,
			"limit":   limit,
			"plan":    plan,
			"upgrade": plan == "free",
		})
		return 0, "", false
	}
	return limit - used, plan, true
}

func (h *AIHandler) recordUsage(r *http.Request, task string) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		return
	}
	if err := h.store.RecordAIGeneration(r.Context(), cuid2.Generate(), user.ID, task); err != nil {
		log.Printf("AI usage record failed: %v", err)
	}
}

// Quota handles GET /api/ai/quota — current usage for the quota meter UI.
func (h *AIHandler) Quota(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}
	plan := h.resolvePlan(r, user)
	limit := quotaForPlan(plan)
	used, err := h.store.CountAIGenerationsSince(r.Context(), user.ID, monthStart())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Quota check failed"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"used": used, "limit": limit, "plan": plan})
}

// GenerateLogos handles POST /api/ai/brand/logos.
func (h *AIHandler) GenerateLogos(w http.ResponseWriter, r *http.Request) {
	if _, _, ok := h.checkQuota(w, r); !ok {
		return
	}
	var req struct {
		Brief string `json:"brief"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Brief == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "brief is required"})
		return
	}
	if len(req.Brief) > maxAIInputLen {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "brief is too long"})
		return
	}

	concepts, err := h.svc.GenerateLogos(r.Context(), req.Brief)
	if err != nil {
		log.Printf("AI logo generation failed: %v", err)
		writeJSON(w, http.StatusBadGateway, map[string]string{"error": "Generation failed, please try again"})
		return
	}
	h.recordUsage(r, "logos")
	writeJSON(w, http.StatusOK, map[string]any{"concepts": concepts})
}

// GenerateCopy handles POST /api/ai/brand/copy.
func (h *AIHandler) GenerateCopy(w http.ResponseWriter, r *http.Request) {
	if _, _, ok := h.checkQuota(w, r); !ok {
		return
	}
	var req struct {
		Brief string `json:"brief"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Brief == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "brief is required"})
		return
	}
	if len(req.Brief) > maxAIInputLen {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "brief is too long"})
		return
	}

	text, err := h.svc.GenerateCopy(r.Context(), req.Brief)
	if err != nil {
		log.Printf("AI copy generation failed: %v", err)
		writeJSON(w, http.StatusBadGateway, map[string]string{"error": "Generation failed, please try again"})
		return
	}
	h.recordUsage(r, "copy")
	writeJSON(w, http.StatusOK, map[string]string{"text": text})
}

// RefineText handles POST /api/ai/brand/refine.
func (h *AIHandler) RefineText(w http.ResponseWriter, r *http.Request) {
	if _, _, ok := h.checkQuota(w, r); !ok {
		return
	}
	var req struct {
		Text        string `json:"text"`
		Instruction string `json:"instruction"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Text == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "text is required"})
		return
	}
	if req.Instruction == "" {
		req.Instruction = "Improve clarity and impact."
	}
	if len(req.Text) > maxAIInputLen || len(req.Instruction) > 1000 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "input is too long"})
		return
	}

	text, err := h.svc.RefineText(r.Context(), req.Text, req.Instruction)
	if err != nil {
		log.Printf("AI refine failed: %v", err)
		writeJSON(w, http.StatusBadGateway, map[string]string{"error": "Refinement failed, please try again"})
		return
	}
	h.recordUsage(r, "refine")
	writeJSON(w, http.StatusOK, map[string]string{"text": text})
}

// LegalAssist handles POST /api/ai/legal.
func (h *AIHandler) LegalAssist(w http.ResponseWriter, r *http.Request) {
	if _, _, ok := h.checkQuota(w, r); !ok {
		return
	}
	var req struct {
		Question string `json:"question"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Question == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "question is required"})
		return
	}
	if len(req.Question) > maxAIInputLen {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "question is too long"})
		return
	}

	text, err := h.svc.LegalAssist(r.Context(), req.Question)
	if err != nil {
		log.Printf("AI legal assist failed: %v", err)
		writeJSON(w, http.StatusBadGateway, map[string]string{"error": "Request failed, please try again"})
		return
	}
	h.recordUsage(r, "legal")
	writeJSON(w, http.StatusOK, map[string]string{"text": text})
}
