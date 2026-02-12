package platform

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
)

type BehanceProvider struct {
	config *oauth2.Config
}

func NewBehanceProvider(clientID, clientSecret, redirectURL string) *BehanceProvider {
	return &BehanceProvider{
		config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes:       []string{"openid", "AdobeID"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://ims-na1.adobelogin.com/ims/authorize/v2",
				TokenURL: "https://ims-na1.adobelogin.com/ims/token/v3",
			},
		},
	}
}

func (b *BehanceProvider) Name() string { return "behance" }

func (b *BehanceProvider) AuthURL(state string) string {
	return b.config.AuthCodeURL(state)
}

func (b *BehanceProvider) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	return b.config.Exchange(ctx, code)
}

func (b *BehanceProvider) RefreshToken(ctx context.Context, refreshToken string) (*oauth2.Token, error) {
	return refreshWithConfig(ctx, b.config, refreshToken)
}

func (b *BehanceProvider) FetchProfile(ctx context.Context, token *oauth2.Token) (*Profile, error) {
	client := b.config.Client(ctx, token)

	req, err := http.NewRequestWithContext(ctx, "GET", "https://www.behance.net/v2/users/me", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create behance request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("behance user request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("behance user returned status %d", resp.StatusCode)
	}

	var result struct {
		User struct {
			ID         int    `json:"id"`
			Username   string `json:"username"`
			DisplayName string `json:"display_name"`
			Images     struct {
				Size138 string `json:"138"`
			} `json:"images"`
			URL       string `json:"url"`
			Stats     struct {
				Followers int `json:"followers"`
			} `json:"stats"`
			Fields   []string `json:"fields"`
		} `json:"user"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode behance response: %w", err)
	}

	return &Profile{
		PlatformUserID: fmt.Sprintf("%d", result.User.ID),
		Username:       result.User.Username,
		DisplayName:    result.User.DisplayName,
		AvatarURL:      result.User.Images.Size138,
		ProfileURL:     result.User.URL,
		FollowerCount:  result.User.Stats.Followers,
		Metadata: map[string]interface{}{
			"fields": result.User.Fields,
		},
	}, nil
}
