package platform

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
	ogithub "golang.org/x/oauth2/github"
)

type GitHubProvider struct {
	config *oauth2.Config
}

func NewGitHubProvider(clientID, clientSecret, redirectURL string) *GitHubProvider {
	return &GitHubProvider{
		config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes:       []string{"read:user"},
			Endpoint:     ogithub.Endpoint,
		},
	}
}

func (g *GitHubProvider) Name() string { return "github" }

func (g *GitHubProvider) AuthURL(state string) string {
	return g.config.AuthCodeURL(state)
}

func (g *GitHubProvider) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	return g.config.Exchange(ctx, code)
}

func (g *GitHubProvider) FetchProfile(ctx context.Context, token *oauth2.Token) (*Profile, error) {
	client := g.config.Client(ctx, token)

	req, err := http.NewRequestWithContext(ctx, "GET", "https://api.github.com/user", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create github request: %w", err)
	}
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("github user request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github user returned status %d", resp.StatusCode)
	}

	var ghUser struct {
		ID          int    `json:"id"`
		Login       string `json:"login"`
		Name        string `json:"name"`
		AvatarURL   string `json:"avatar_url"`
		HTMLURL     string `json:"html_url"`
		Followers   int    `json:"followers"`
		PublicRepos int    `json:"public_repos"`
		Bio         string `json:"bio"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&ghUser); err != nil {
		return nil, fmt.Errorf("failed to decode github response: %w", err)
	}

	return &Profile{
		PlatformUserID: fmt.Sprintf("%d", ghUser.ID),
		Username:       ghUser.Login,
		DisplayName:    ghUser.Name,
		AvatarURL:      ghUser.AvatarURL,
		ProfileURL:     ghUser.HTMLURL,
		FollowerCount:  ghUser.Followers,
		Metadata: map[string]interface{}{
			"public_repos": ghUser.PublicRepos,
			"bio":          ghUser.Bio,
		},
	}, nil
}
