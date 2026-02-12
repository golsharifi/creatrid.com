package platform

import (
	"context"
	"fmt"
	"strconv"

	"golang.org/x/oauth2"
)

func parseIntOrZero(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}

func refreshWithConfig(ctx context.Context, cfg *oauth2.Config, refreshToken string) (*oauth2.Token, error) {
	if refreshToken == "" {
		return nil, fmt.Errorf("no refresh token available")
	}
	src := cfg.TokenSource(ctx, &oauth2.Token{RefreshToken: refreshToken})
	newToken, err := src.Token()
	if err != nil {
		return nil, fmt.Errorf("token refresh failed: %w", err)
	}
	return newToken, nil
}
