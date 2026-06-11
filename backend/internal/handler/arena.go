package handler

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/creatrid/creatrid/internal/middleware"
	"github.com/creatrid/creatrid/internal/store"
	"github.com/nrednav/cuid2"
)

// ArenaHandler implements the VR Community Arena: a closed-loop game where
// creators collect passport stickers and earn points. Points/stickers are
// intentionally NOT cashable or externally tradeable (securities posture).
type ArenaHandler struct {
	store *store.Store
}

func NewArenaHandler(st *store.Store) *ArenaHandler {
	return &ArenaHandler{store: st}
}

const exploreCooldown = 4 * time.Hour

// Rarity weights and point bonuses for the explore drop.
var rarityWeights = []struct {
	rarity string
	weight int
	bonus  int
}{
	{"common", 60, 0},
	{"uncommon", 25, 5},
	{"rare", 12, 15},
	{"legendary", 3, 40},
}

const exploreBasePoints = 10

// Me handles GET /api/arena/me — points, stickers, cooldown, achievements.
func (h *ArenaHandler) Me(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}

	state, err := h.store.GetArenaState(r.Context(), user.ID)
	if err != nil || state == nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to load arena state"})
		return
	}
	stickers, err := h.store.ListUserStickers(r.Context(), user.ID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to load stickers"})
		return
	}
	catalog, err := h.store.ListStickers(r.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to load sticker catalog"})
		return
	}

	var nextExplore *time.Time
	if state.LastExplore != nil {
		t := state.LastExplore.Add(exploreCooldown)
		if t.After(time.Now()) {
			nextExplore = &t
		}
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"points":       state.Points,
		"stickers":     stickers,
		"catalog":      catalog,
		"nextExplore":  nextExplore,
		"achievements": achievements(state.Points, stickers),
	})
}

// Explore handles POST /api/arena/explore — the core game action: visit the
// farm, find a random sticker, earn points. Cooldown-limited.
func (h *ArenaHandler) Explore(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}

	state, err := h.store.GetArenaState(r.Context(), user.ID)
	if err != nil || state == nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to load arena state"})
		return
	}
	if state.LastExplore != nil {
		next := state.LastExplore.Add(exploreCooldown)
		if next.After(time.Now()) {
			writeJSON(w, http.StatusTooManyRequests, map[string]any{
				"error":       "The farm needs time to regrow",
				"nextExplore": next,
			})
			return
		}
	}

	catalog, err := h.store.ListStickers(r.Context())
	if err != nil || len(catalog) == 0 {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Sticker catalog unavailable"})
		return
	}

	rarity, bonus := rollRarity()
	var pool []*store.Sticker
	for _, s := range catalog {
		if s.Rarity == rarity {
			pool = append(pool, s)
		}
	}
	if len(pool) == 0 {
		pool = catalog
	}
	drop := pool[rand.Intn(len(pool))]
	points := exploreBasePoints + bonus

	if err := h.store.RecordExplore(r.Context(), cuid2.Generate(), user.ID, drop.ID, points); err != nil {
		log.Printf("Arena explore failed: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Explore failed"})
		return
	}

	next := time.Now().Add(exploreCooldown)
	writeJSON(w, http.StatusOK, map[string]any{
		"sticker":     drop,
		"points":      points,
		"totalPoints": state.Points + points,
		"nextExplore": next,
	})
}

func rollRarity() (string, int) {
	total := 0
	for _, rw := range rarityWeights {
		total += rw.weight
	}
	roll := rand.Intn(total)
	for _, rw := range rarityWeights {
		if roll < rw.weight {
			return rw.rarity, rw.bonus
		}
		roll -= rw.weight
	}
	return "common", 0
}

// Leaderboard handles GET /api/arena/leaderboard — top performers (public).
func (h *ArenaHandler) Leaderboard(w http.ResponseWriter, r *http.Request) {
	entries, err := h.store.ArenaLeaderboard(r.Context(), 25)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to load leaderboard"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"leaderboard": entries})
}

// Gift handles POST /api/arena/gift — give one of your stickers to another
// creator. Closed-loop "trading": stickers move between passports, never out.
func (h *ArenaHandler) Gift(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}

	var req struct {
		Username  string `json:"username"`
		StickerID string `json:"stickerId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Username == "" || req.StickerID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "username and stickerId are required"})
		return
	}

	target, err := h.store.FindUserByUsername(r.Context(), strings.ToLower(strings.TrimSpace(req.Username)))
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Database error"})
		return
	}
	if target == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "Creator not found"})
		return
	}
	if target.ID == user.ID {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "You can't gift a sticker to yourself"})
		return
	}

	if err := h.store.GiftSticker(r.Context(), user.ID, target.ID, req.StickerID); err != nil {
		if err == store.ErrInsufficientStickers {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "You don't have that sticker to give"})
			return
		}
		log.Printf("Arena gift failed: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Gift failed"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"success": true})
}

type achievement struct {
	Key   string `json:"key"`
	Name  string `json:"name"`
	Emoji string `json:"emoji"`
}

// achievements derives badge state from counts — no extra storage needed.
func achievements(points int, stickers []*store.UserSticker) []achievement {
	unique := len(stickers)
	totalCount := 0
	hasLegendary := false
	for _, s := range stickers {
		totalCount += s.Count
		if s.Rarity == "legendary" {
			hasLegendary = true
		}
	}

	var out []achievement
	if totalCount >= 1 {
		out = append(out, achievement{"first-find", "First Find", "🌱"})
	}
	if unique >= 5 {
		out = append(out, achievement{"collector", "Collector", "🧺"})
	}
	if unique >= 12 {
		out = append(out, achievement{"completionist", "Full Barn", "🏆"})
	}
	if points >= 100 {
		out = append(out, achievement{"century", "Century Club", "💯"})
	}
	if points >= 500 {
		out = append(out, achievement{"farmhand", "Master Farmhand", "🚜"})
	}
	if hasLegendary {
		out = append(out, achievement{"legend", "Legend Hunter", "🐉"})
	}
	if out == nil {
		out = []achievement{}
	}
	return out
}
