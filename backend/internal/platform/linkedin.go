package platform

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
	olinkedin "golang.org/x/oauth2/linkedin"
)

type LinkedInProvider struct {
	config *oauth2.Config
}

func NewLinkedInProvider(clientID, clientSecret, redirectURL string) *LinkedInProvider {
	return &LinkedInProvider{
		config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes:       []string{"openid", "profile"},
			Endpoint:     olinkedin.Endpoint,
		},
	}
}

func (l *LinkedInProvider) Name() string { return "linkedin" }

func (l *LinkedInProvider) AuthURL(state string) string {
	return l.config.AuthCodeURL(state)
}

func (l *LinkedInProvider) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	return l.config.Exchange(ctx, code)
}

func (l *LinkedInProvider) RefreshToken(ctx context.Context, refreshToken string) (*oauth2.Token, error) {
	return refreshWithConfig(ctx, l.config, refreshToken)
}

func (l *LinkedInProvider) FetchProfile(ctx context.Context, token *oauth2.Token) (*Profile, error) {
	client := l.config.Client(ctx, token)

	req, err := http.NewRequestWithContext(ctx, "GET", "https://api.linkedin.com/v2/userinfo", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create linkedin request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("linkedin userinfo request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("linkedin userinfo returned status %d", resp.StatusCode)
	}

	var userInfo struct {
		Sub     string `json:"sub"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("failed to decode linkedin response: %w", err)
	}

	return &Profile{
		PlatformUserID: userInfo.Sub,
		Username:       "",
		DisplayName:    userInfo.Name,
		AvatarURL:      userInfo.Picture,
		ProfileURL:     "https://linkedin.com/in/" + userInfo.Sub,
		FollowerCount:  0,
		Metadata:       map[string]interface{}{},
	}, nil
}
