package model

import (
	"encoding/json"
	"time"
)

type UserRole string

const (
	RoleCreator UserRole = "CREATOR"
	RoleBrand   UserRole = "BRAND"
	RoleAdmin   UserRole = "ADMIN"
)

type User struct {
	ID            string     `json:"id"`
	Name          *string    `json:"name"`
	Email         string     `json:"email"`
	EmailVerified *time.Time `json:"emailVerified,omitempty"`
	Image         *string    `json:"image"`
	Username      *string    `json:"username"`
	Bio           *string    `json:"bio"`
	Role          UserRole   `json:"role"`
	CreatorScore  *int             `json:"creatorScore"`
	IsVerified    bool             `json:"isVerified"`
	Onboarded     bool             `json:"onboarded"`
	Theme         string           `json:"theme"`
	CustomLinks   json.RawMessage  `json:"customLinks"`
	CreatedAt     time.Time        `json:"createdAt"`
	UpdatedAt     time.Time        `json:"updatedAt"`
}

type PublicUser struct {
	ID           string          `json:"id"`
	Name         *string         `json:"name"`
	Username     *string         `json:"username"`
	Bio          *string         `json:"bio"`
	Image        *string         `json:"image"`
	Role         UserRole        `json:"role"`
	CreatorScore *int            `json:"creatorScore"`
	IsVerified   bool            `json:"isVerified"`
	Theme        string          `json:"theme"`
	CustomLinks  json.RawMessage `json:"customLinks"`
	CreatedAt    time.Time       `json:"createdAt"`
}

func (u *User) ToPublic() *PublicUser {
	return &PublicUser{
		ID:           u.ID,
		Name:         u.Name,
		Username:     u.Username,
		Bio:          u.Bio,
		Image:        u.Image,
		Role:         u.Role,
		CreatorScore: u.CreatorScore,
		IsVerified:   u.IsVerified,
		Theme:        u.Theme,
		CustomLinks:  u.CustomLinks,
		CreatedAt:    u.CreatedAt,
	}
}

type Account struct {
	UserID            string  `json:"userId"`
	Type              string  `json:"type"`
	Provider          string  `json:"provider"`
	ProviderAccountID string  `json:"providerAccountId"`
	RefreshToken      *string `json:"-"`
	AccessToken       *string `json:"-"`
	ExpiresAt         *int    `json:"-"`
	TokenType         *string `json:"-"`
	Scope             *string `json:"-"`
	IDToken           *string `json:"-"`
}
