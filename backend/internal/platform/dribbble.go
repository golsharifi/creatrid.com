package platform

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
)

type DribbbleProvider struct {
	config *oauth2.Config
}

func NewDribbbleProvider(clientID, clientSecret, redirectURL string) *DribbbleProvider {
	return &DribbbleProvider{
		config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes:       []string{"public"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://dribbble.com/oauth/authorize",
				TokenURL: "https://dribbble.com/oauth/token",
			},
		},
	}
}

func (d *DribbbleProvider) Name() string { return "dribbble" }

func (d *DribbbleProvider) AuthURL(state string) string {
	return d.config.AuthCodeURL(state)
}

func (d *DribbbleProvider) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	return d.config.Exchange(ctx, code)
}

func (d *DribbbleProvider) FetchProfile(ctx context.Context, token *oauth2.Token) (*Profile, error) {
	client := d.config.Client(ctx, token)

	req, err := http.NewRequestWithContext(ctx, "GET", "https://api.dribbble.com/v2/user", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create dribbble request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("dribbble user request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("dribbble user returned status %d", resp.StatusCode)
	}

	var user struct {
		ID            int    `json:"id"`
		Login         string `json:"login"`
		Name          string `json:"name"`
		AvatarURL     string `json:"avatar_url"`
		HTMLURL       string `json:"html_url"`
		FollowersCount int   `json:"followers_count"`
		ShotsCount    int    `json:"shots_count"`
		Bio           string `json:"bio"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("failed to decode dribbble response: %w", err)
	}

	return &Profile{
		PlatformUserID: fmt.Sprintf("%d", user.ID),
		Username:       user.Login,
		DisplayName:    user.Name,
		AvatarURL:      user.AvatarURL,
		ProfileURL:     user.HTMLURL,
		FollowerCount:  user.FollowersCount,
		Metadata: map[string]interface{}{
			"shots_count": user.ShotsCount,
			"bio":         user.Bio,
		},
	}, nil
}
