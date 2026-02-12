package handler

import (
	"fmt"
	"html"
	"net/http"
	"strings"

	"github.com/creatrid/creatrid/internal/config"
	"github.com/creatrid/creatrid/internal/store"
	"github.com/go-chi/chi/v5"
)

type OGHandler struct {
	store  *store.Store
	config *config.Config
}

func NewOGHandler(store *store.Store, cfg *config.Config) *OGHandler {
	return &OGHandler{store: store, config: cfg}
}

func (h *OGHandler) ProfilePage(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	user, err := h.store.FindUserByUsername(r.Context(), strings.ToLower(username))
	if err != nil || user == nil {
		http.Redirect(w, r, h.config.FrontendURL, http.StatusTemporaryRedirect)
		return
	}

	profileURL := fmt.Sprintf("%s/profile?u=%s", h.config.FrontendURL, html.EscapeString(*user.Username))

	name := "Creator"
	if user.Name != nil {
		name = html.EscapeString(*user.Name)
	}

	description := fmt.Sprintf("@%s on Creatrid — Verified Creator Passport", html.EscapeString(*user.Username))
	if user.Bio != nil && *user.Bio != "" {
		description = html.EscapeString(*user.Bio)
	}

	image := ""
	if user.Image != nil {
		image = html.EscapeString(*user.Image)
	}

	scoreTag := ""
	if user.CreatorScore != nil {
		scoreTag = fmt.Sprintf(`<meta property="og:label1" content="Creator Score" /><meta property="og:data1" content="%d / 100" />`, *user.CreatorScore)
	}

	imageTag := ""
	if image != "" {
		imageTag = fmt.Sprintf(`<meta property="og:image" content="%s" /><meta name="twitter:image" content="%s" />`, image, image)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, `<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8" />
<title>%s — Creatrid</title>
<meta property="og:type" content="profile" />
<meta property="og:title" content="%s — Creator Passport" />
<meta property="og:description" content="%s" />
<meta property="og:url" content="%s" />
<meta property="og:site_name" content="Creatrid" />
%s
%s
<meta name="twitter:card" content="summary" />
<meta name="twitter:title" content="%s — Creator Passport" />
<meta name="twitter:description" content="%s" />
<meta http-equiv="refresh" content="0;url=%s" />
</head>
<body>
<p>Redirecting to <a href="%s">%s's profile</a>...</p>
</body>
</html>`, name, name, description, profileURL, imageTag, scoreTag,
		name, description, profileURL, profileURL, name)
}
