package platform

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type YouTubeProvider struct {
	config *oauth2.Config
}

func NewYouTubeProvider(clientID, clientSecret, redirectURL string) *YouTubeProvider {
	return &YouTubeProvider{
		config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes:       []string{"https://www.googleapis.com/auth/youtube.readonly"},
			Endpoint:     google.Endpoint,
		},
	}
}

func (y *YouTubeProvider) Name() string { return "youtube" }

func (y *YouTubeProvider) AuthURL(state string) string {
	return y.config.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.SetAuthURLParam("prompt", "consent"))
}

func (y *YouTubeProvider) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	return y.config.Exchange(ctx, code)
}

func (y *YouTubeProvider) FetchProfile(ctx context.Context, token *oauth2.Token) (*Profile, error) {
	client := y.config.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/youtube/v3/channels?part=snippet,statistics&mine=true")
	if err != nil {
		return nil, fmt.Errorf("youtube channels request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("youtube channels returned status %d", resp.StatusCode)
	}

	var result struct {
		Items []struct {
			ID      string `json:"id"`
			Snippet struct {
				Title      string `json:"title"`
				CustomURL  string `json:"customUrl"`
				Thumbnails struct {
					Default struct {
						URL string `json:"url"`
					} `json:"default"`
				} `json:"thumbnails"`
			} `json:"snippet"`
			Statistics struct {
				ViewCount       string `json:"viewCount"`
				SubscriberCount string `json:"subscriberCount"`
				VideoCount      string `json:"videoCount"`
			} `json:"statistics"`
		} `json:"items"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode youtube response: %w", err)
	}

	if len(result.Items) == 0 {
		return nil, fmt.Errorf("no youtube channel found for this account")
	}

	ch := result.Items[0]
	subscriberCount := parseIntOrZero(ch.Statistics.SubscriberCount)
	videoCount := parseIntOrZero(ch.Statistics.VideoCount)
	viewCount := parseIntOrZero(ch.Statistics.ViewCount)

	username := ch.Snippet.CustomURL
	profileURL := "https://youtube.com/channel/" + ch.ID
	if username != "" {
		profileURL = "https://youtube.com/" + username
	}

	return &Profile{
		PlatformUserID: ch.ID,
		Username:       username,
		DisplayName:    ch.Snippet.Title,
		AvatarURL:      ch.Snippet.Thumbnails.Default.URL,
		ProfileURL:     profileURL,
		FollowerCount:  subscriberCount,
		Metadata: map[string]interface{}{
			"subscriber_count": subscriberCount,
			"video_count":      videoCount,
			"view_count":       viewCount,
		},
	}, nil
}
