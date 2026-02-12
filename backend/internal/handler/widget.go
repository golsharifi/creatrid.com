package handler

import (
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"strings"

	"github.com/creatrid/creatrid/internal/store"
	"github.com/go-chi/chi/v5"
)

type WidgetHandler struct {
	store *store.Store
}

func NewWidgetHandler(store *store.Store) *WidgetHandler {
	return &WidgetHandler{store: store}
}

func (h *WidgetHandler) JSON(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	user, err := h.store.FindUserByUsername(r.Context(), strings.ToLower(username))
	if err != nil || user == nil {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "User not found"})
		return
	}

	connections, _ := h.store.CountConnectionsByUserID(r.Context(), user.ID)

	score := 0
	if user.CreatorScore != nil {
		score = *user.CreatorScore
	}

	name := ""
	if user.Name != nil {
		name = *user.Name
	}

	image := ""
	if user.Image != nil {
		image = *user.Image
	}

	uname := ""
	if user.Username != nil {
		uname = *user.Username
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Cache-Control", "public, max-age=3600")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"name":        name,
		"username":    uname,
		"score":       score,
		"verified":    user.IsVerified,
		"image":       image,
		"connections": connections,
	})
}

func (h *WidgetHandler) SVGBadge(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	user, err := h.store.FindUserByUsername(r.Context(), strings.ToLower(username))
	if err != nil || user == nil {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "image/svg+xml")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`<svg xmlns="http://www.w3.org/2000/svg" width="200" height="60" viewBox="0 0 200 60">
  <rect width="200" height="60" rx="10" fill="#18181b"/>
  <text x="100" y="34" font-family="system-ui,sans-serif" font-size="12" fill="#a1a1aa" text-anchor="middle">User not found</text>
</svg>`))
		return
	}

	scoreStr := "\u2014"
	if user.CreatorScore != nil {
		scoreStr = fmt.Sprintf("%d", *user.CreatorScore)
	}

	verifiedMark := ""
	if user.IsVerified {
		verifiedMark = `<circle cx="188" cy="30" r="8" fill="#3b82f6"/><path d="M184 30l3 3 5-5" stroke="white" stroke-width="1.5" fill="none" stroke-linecap="round" stroke-linejoin="round"/>`
	}

	svg := fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="200" height="60" viewBox="0 0 200 60">
  <rect width="200" height="60" rx="10" fill="#18181b"/>
  <rect x="1" y="1" width="198" height="58" rx="9" fill="none" stroke="#3f3f46" stroke-width="1"/>
  <text x="14" y="24" font-family="system-ui,sans-serif" font-size="11" fill="#a1a1aa" font-weight="600">CREATRID</text>
  <text x="14" y="44" font-family="system-ui,sans-serif" font-size="14" fill="white" font-weight="700">Score: %s</text>
  %s
</svg>`, scoreStr, verifiedMark)

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Cache-Control", "public, max-age=3600")
	w.Header().Set("Content-Type", "image/svg+xml")
	w.Write([]byte(svg))
}

func (h *WidgetHandler) HTMLEmbed(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	user, err := h.store.FindUserByUsername(r.Context(), strings.ToLower(username))
	if err != nil || user == nil {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`<!DOCTYPE html><html><body><p>User not found</p></body></html>`))
		return
	}

	name := "Creator"
	if user.Name != nil {
		name = html.EscapeString(*user.Name)
	}

	uname := ""
	if user.Username != nil {
		uname = html.EscapeString(*user.Username)
	}

	scoreStr := "\u2014"
	if user.CreatorScore != nil {
		scoreStr = fmt.Sprintf("%d", *user.CreatorScore)
	}

	image := ""
	if user.Image != nil {
		image = html.EscapeString(*user.Image)
	}

	avatarHTML := ""
	if image != "" {
		avatarHTML = fmt.Sprintf(`<img src="%s" alt="%s" style="width:36px;height:36px;border-radius:50%%;object-fit:cover;flex-shrink:0;" />`, image, name)
	} else {
		avatarHTML = fmt.Sprintf(`<div style="width:36px;height:36px;border-radius:50%%;background:#3f3f46;display:flex;align-items:center;justify-content:center;flex-shrink:0;color:#a1a1aa;font-size:14px;font-weight:600;">%s</div>`, string([]rune(name)[0:1]))
	}

	verifiedHTML := ""
	if user.IsVerified {
		verifiedHTML = `<svg width="14" height="14" viewBox="0 0 24 24" fill="none" style="margin-left:4px;flex-shrink:0;"><circle cx="12" cy="12" r="12" fill="#3b82f6"/><path d="M8 12l3 3 5-5" stroke="white" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/></svg>`
	}

	profileURL := fmt.Sprintf("https://creatrid.com/profile?u=%s", uname)

	htmlContent := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8" />
<meta name="viewport" content="width=device-width, initial-scale=1" />
<style>
  * { margin: 0; padding: 0; box-sizing: border-box; }
  body { font-family: system-ui, -apple-system, sans-serif; background: transparent; }
  .widget {
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 10px 14px;
    background: #18181b;
    border: 1px solid #3f3f46;
    border-radius: 12px;
    text-decoration: none;
    color: white;
    max-width: 260px;
    transition: border-color 0.2s;
  }
  .widget:hover { border-color: #52525b; }
  .info { display: flex; flex-direction: column; gap: 2px; min-width: 0; }
  .name-row { display: flex; align-items: center; }
  .name {
    font-size: 13px;
    font-weight: 600;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }
  .meta {
    font-size: 11px;
    color: #a1a1aa;
    display: flex;
    align-items: center;
    gap: 6px;
  }
  .score { font-weight: 600; color: #e4e4e7; }
  .label { font-size: 9px; color: #71717a; letter-spacing: 0.5px; text-transform: uppercase; }
</style>
</head>
<body>
<a class="widget" href="%s" target="_blank" rel="noopener noreferrer">
  %s
  <div class="info">
    <div class="name-row">
      <span class="name">%s</span>
      %s
    </div>
    <div class="meta">
      <span class="score">Score: %s</span>
      <span class="label">Creatrid</span>
    </div>
  </div>
</a>
</body>
</html>`, profileURL, avatarHTML, name, verifiedHTML, scoreStr)

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Cache-Control", "public, max-age=3600")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(htmlContent))
}
