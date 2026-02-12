package model

import (
	"encoding/json"
	"time"
)

type Connection struct {
	ID             string          `json:"id"`
	UserID         string          `json:"userId"`
	Platform       string          `json:"platform"`
	PlatformUserID string          `json:"platformUserId"`
	Username       *string         `json:"username"`
	DisplayName    *string         `json:"displayName"`
	AvatarURL      *string         `json:"avatarUrl"`
	ProfileURL     *string         `json:"profileUrl"`
	FollowerCount  *int            `json:"followerCount"`
	AccessToken    *string         `json:"-"`
	RefreshToken   *string         `json:"-"`
	TokenExpiresAt *time.Time      `json:"-"`
	Metadata       json.RawMessage `json:"metadata"`
	ConnectedAt    time.Time       `json:"connectedAt"`
	UpdatedAt      time.Time       `json:"updatedAt"`
}

type PublicConnection struct {
	Platform      string          `json:"platform"`
	Username      *string         `json:"username"`
	DisplayName   *string         `json:"displayName"`
	AvatarURL     *string         `json:"avatarUrl"`
	ProfileURL    *string         `json:"profileUrl"`
	FollowerCount *int            `json:"followerCount"`
	Metadata      json.RawMessage `json:"metadata"`
	ConnectedAt   time.Time       `json:"connectedAt"`
}

func (c *Connection) ToPublic() *PublicConnection {
	return &PublicConnection{
		Platform:      c.Platform,
		Username:      c.Username,
		DisplayName:   c.DisplayName,
		AvatarURL:     c.AvatarURL,
		ProfileURL:    c.ProfileURL,
		FollowerCount: c.FollowerCount,
		Metadata:      c.Metadata,
		ConnectedAt:   c.ConnectedAt,
	}
}
