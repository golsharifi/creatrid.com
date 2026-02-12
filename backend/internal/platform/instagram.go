package platform

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
	ofacebook "golang.org/x/oauth2/facebook"
)

type InstagramProvider struct {
	config *oauth2.Config
}

func NewInstagramProvider(clientID, clientSecret, redirectURL string) *InstagramProvider {
	return &InstagramProvider{
		config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes:       []string{"instagram_basic"},
			Endpoint:     ofacebook.Endpoint,
		},
	}
}

func (i *InstagramProvider) Name() string { return "instagram" }

func (i *InstagramProvider) AuthURL(state string) string {
	return i.config.AuthCodeURL(state)
}

func (i *InstagramProvider) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	return i.config.Exchange(ctx, code)
}

func (i *InstagramProvider) FetchProfile(ctx context.Context, token *oauth2.Token) (*Profile, error) {
	client := i.config.Client(ctx, token)

	req, err := http.NewRequestWithContext(ctx, "GET",
		"https://graph.instagram.com/me?fields=id,username,account_type,media_count",
		nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create instagram request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("instagram user request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("instagram user returned status %d", resp.StatusCode)
	}

	var igUser struct {
		ID          string `json:"id"`
		Username    string `json:"username"`
		AccountType string `json:"account_type"`
		MediaCount  int    `json:"media_count"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&igUser); err != nil {
		return nil, fmt.Errorf("failed to decode instagram response: %w", err)
	}

	return &Profile{
		PlatformUserID: igUser.ID,
		Username:       igUser.Username,
		DisplayName:    igUser.Username,
		AvatarURL:      "",
		ProfileURL:     "https://instagram.com/" + igUser.Username,
		FollowerCount:  0,
		Metadata: map[string]interface{}{
			"account_type": igUser.AccountType,
			"media_count":  igUser.MediaCount,
		},
	}, nil
}
