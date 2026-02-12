package platform

import (
	"context"

	"golang.org/x/oauth2"
)

type Profile struct {
	PlatformUserID string
	Username       string
	DisplayName    string
	AvatarURL      string
	ProfileURL     string
	FollowerCount  int
	Metadata       map[string]interface{}
}

type Provider interface {
	Name() string
	AuthURL(state string) string
	Exchange(ctx context.Context, code string) (*oauth2.Token, error)
	FetchProfile(ctx context.Context, token *oauth2.Token) (*Profile, error)
	RefreshToken(ctx context.Context, refreshToken string) (*oauth2.Token, error)
}
