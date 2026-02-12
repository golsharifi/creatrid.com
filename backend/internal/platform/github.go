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

func (g *GitHubProvider) RefreshToken(ctx context.Context, refreshToken string) (*oauth2.Token, error) {
	return refreshWithConfig(ctx, g.config, refreshToken)
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

	metadata := map[string]interface{}{
		"public_repos": ghUser.PublicRepos,
		"bio":          ghUser.Bio,
	}

	// Fetch top 3 repos by stars for embedding
	repoReq, err := http.NewRequestWithContext(ctx, "GET", "https://api.github.com/users/"+ghUser.Login+"/repos?sort=stars&direction=desc&per_page=3", nil)
	if err == nil {
		repoReq.Header.Set("Accept", "application/vnd.github+json")
		repoResp, err := client.Do(repoReq)
		if err == nil {
			defer repoResp.Body.Close()
			if repoResp.StatusCode == http.StatusOK {
				var repos []struct {
					Name        string `json:"name"`
					Description string `json:"description"`
					HTMLURL     string `json:"html_url"`
					Stars       int    `json:"stargazers_count"`
					Language    string `json:"language"`
				}
				if json.NewDecoder(repoResp.Body).Decode(&repos) == nil && len(repos) > 0 {
					topRepos := make([]map[string]interface{}, 0, len(repos))
					for _, r := range repos {
						topRepos = append(topRepos, map[string]interface{}{
							"name":        r.Name,
							"description": r.Description,
							"url":         r.HTMLURL,
							"stars":       r.Stars,
							"language":    r.Language,
						})
					}
					metadata["top_repos"] = topRepos
				}
			}
		}
	}

	return &Profile{
		PlatformUserID: fmt.Sprintf("%d", ghUser.ID),
		Username:       ghUser.Login,
		DisplayName:    ghUser.Name,
		AvatarURL:      ghUser.AvatarURL,
		ProfileURL:     ghUser.HTMLURL,
		FollowerCount:  ghUser.Followers,
		Metadata:       metadata,
	}, nil
}
