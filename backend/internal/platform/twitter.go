package platform

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
)

type TwitterProvider struct {
	config *oauth2.Config
}

func NewTwitterProvider(clientID, clientSecret, redirectURL string) *TwitterProvider {
	return &TwitterProvider{
		config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes:       []string{"tweet.read", "users.read"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://twitter.com/i/oauth2/authorize",
				TokenURL: "https://api.twitter.com/2/oauth2/token",
			},
		},
	}
}

func (t *TwitterProvider) Name() string { return "twitter" }

func (t *TwitterProvider) AuthURL(state string) string {
	return t.config.AuthCodeURL(state,
		oauth2.SetAuthURLParam("code_challenge", "challenge"),
		oauth2.SetAuthURLParam("code_challenge_method", "plain"),
	)
}

func (t *TwitterProvider) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	return t.config.Exchange(ctx, code,
		oauth2.SetAuthURLParam("code_verifier", "challenge"),
	)
}

func (t *TwitterProvider) FetchProfile(ctx context.Context, token *oauth2.Token) (*Profile, error) {
	client := t.config.Client(ctx, token)

	req, err := http.NewRequestWithContext(ctx, "GET", "https://api.twitter.com/2/users/me?user.fields=profile_image_url,public_metrics,description", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create twitter request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("twitter user request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("twitter user returned status %d", resp.StatusCode)
	}

	var result struct {
		Data struct {
			ID              string `json:"id"`
			Username        string `json:"username"`
			Name            string `json:"name"`
			ProfileImageURL string `json:"profile_image_url"`
			Description     string `json:"description"`
			PublicMetrics   struct {
				FollowersCount int `json:"followers_count"`
				TweetCount     int `json:"tweet_count"`
			} `json:"public_metrics"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode twitter response: %w", err)
	}

	return &Profile{
		PlatformUserID: result.Data.ID,
		Username:       result.Data.Username,
		DisplayName:    result.Data.Name,
		AvatarURL:      result.Data.ProfileImageURL,
		ProfileURL:     "https://x.com/" + result.Data.Username,
		FollowerCount:  result.Data.PublicMetrics.FollowersCount,
		Metadata: map[string]interface{}{
			"tweet_count": result.Data.PublicMetrics.TweetCount,
			"bio":         result.Data.Description,
		},
	}, nil
}
