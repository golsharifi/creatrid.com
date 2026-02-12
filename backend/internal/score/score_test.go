package score

import (
	"testing"
	"time"

	"github.com/creatrid/creatrid/internal/model"
)

func strPtr(s string) *string { return &s }
func intPtr(i int) *int       { return &i }
func timePtr(t time.Time) *time.Time { return &t }

func TestCalculate(t *testing.T) {
	tests := []struct {
		name        string
		user        *model.User
		connections []*model.Connection
		wantMin     int
		wantMax     int
		wantExact   *int
	}{
		{
			name:        "empty user with no connections",
			user:        &model.User{},
			connections: nil,
			wantExact:   intPtr(0),
		},
		{
			name: "full profile gives 20 points",
			user: &model.User{
				Name:     strPtr("Test User"),
				Username: strPtr("testuser"),
				Bio:      strPtr("A short bio"),
				Image:    strPtr("https://example.com/img.png"),
			},
			connections: nil,
			wantExact:   intPtr(20),
		},
		{
			name: "email verified gives 10 points",
			user: &model.User{
				EmailVerified: timePtr(time.Now()),
			},
			connections: nil,
			wantExact:   intPtr(10),
		},
		{
			name: "full profile plus email verified gives 30 points",
			user: &model.User{
				Name:          strPtr("Test User"),
				Username:      strPtr("testuser"),
				Bio:           strPtr("A short bio"),
				Image:         strPtr("https://example.com/img.png"),
				EmailVerified: timePtr(time.Now()),
			},
			connections: nil,
			wantExact:   intPtr(30),
		},
		{
			name: "1 connection gives 10 points",
			user: &model.User{},
			connections: []*model.Connection{
				{Platform: "youtube"},
			},
			wantExact: intPtr(10),
		},
		{
			name: "5 connections gives 50 points (max connection points)",
			user: &model.User{},
			connections: []*model.Connection{
				{Platform: "youtube"},
				{Platform: "github"},
				{Platform: "twitter"},
				{Platform: "linkedin"},
				{Platform: "instagram"},
			},
			wantExact: intPtr(50),
		},
		{
			name: "6 connections still gives 50 points (capped)",
			user: &model.User{},
			connections: []*model.Connection{
				{Platform: "youtube"},
				{Platform: "github"},
				{Platform: "twitter"},
				{Platform: "linkedin"},
				{Platform: "instagram"},
				{Platform: "dribbble"},
			},
			wantExact: intPtr(50),
		},
		{
			name: "100 followers gives 0 follower bonus (boundary of log scale)",
			user: &model.User{},
			connections: []*model.Connection{
				{Platform: "youtube", FollowerCount: intPtr(100)},
			},
			// log10(100)=2, (2-2)*5=0, so no follower bonus at exactly 100
			wantExact: intPtr(10), // 10 (1 connection) + 0 (follower bonus)
		},
		{
			name: "100K followers gives 15 follower bonus points",
			user: &model.User{},
			connections: []*model.Connection{
				{Platform: "youtube", FollowerCount: intPtr(100000)},
			},
			// 10 (1 connection) + 15 (log10(100000)=5, (5-2)*5=15)
			wantExact: intPtr(25),
		},
		{
			name: "1M followers gives max 20 follower bonus",
			user: &model.User{},
			connections: []*model.Connection{
				{Platform: "youtube", FollowerCount: intPtr(1000000)},
			},
			wantExact: intPtr(30), // 10 (connection) + 20 (max follower bonus)
		},
		{
			name: "max score is capped at 100",
			user: &model.User{
				Name:          strPtr("Test User"),
				Username:      strPtr("testuser"),
				Bio:           strPtr("A short bio"),
				Image:         strPtr("https://example.com/img.png"),
				EmailVerified: timePtr(time.Now()),
			},
			connections: []*model.Connection{
				{Platform: "youtube", FollowerCount: intPtr(10000000)},
				{Platform: "github", FollowerCount: intPtr(10000000)},
				{Platform: "twitter", FollowerCount: intPtr(10000000)},
				{Platform: "linkedin", FollowerCount: intPtr(10000000)},
				{Platform: "instagram", FollowerCount: intPtr(10000000)},
			},
			wantExact: intPtr(100),
		},
		{
			name: "zero followers gives no bonus",
			user: &model.User{},
			connections: []*model.Connection{
				{Platform: "youtube", FollowerCount: intPtr(0)},
			},
			wantExact: intPtr(10), // just the 1 connection
		},
		{
			name: "nil follower count gives no bonus",
			user: &model.User{},
			connections: []*model.Connection{
				{Platform: "youtube", FollowerCount: nil},
			},
			wantExact: intPtr(10), // just the 1 connection
		},
		{
			name: "partial profile gives partial points",
			user: &model.User{
				Name:     strPtr("Test User"),
				Username: strPtr("testuser"),
			},
			connections: nil,
			wantExact:   intPtr(10),
		},
		{
			name: "empty string fields give no points",
			user: &model.User{
				Name:     strPtr(""),
				Username: strPtr(""),
				Bio:      strPtr(""),
				Image:    strPtr(""),
			},
			connections: nil,
			wantExact:   intPtr(0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Calculate(tt.user, tt.connections)

			if tt.wantExact != nil {
				if got != *tt.wantExact {
					t.Errorf("Calculate() = %d, want %d", got, *tt.wantExact)
				}
			} else {
				if got < tt.wantMin || got > tt.wantMax {
					t.Errorf("Calculate() = %d, want between %d and %d", got, tt.wantMin, tt.wantMax)
				}
			}
		})
	}
}

func TestFollowerBonus(t *testing.T) {
	tests := []struct {
		name      string
		followers int
		want      int
	}{
		{"zero followers", 0, 0},
		{"negative followers", -10, 0},
		{"1 follower", 1, 0},
		{"10 followers", 10, 1},
		{"100 followers", 100, 0},
		{"1000 followers", 1000, 5},
		{"10000 followers", 10000, 10},
		{"100000 followers", 100000, 15},
		{"1000000 followers hits 20", 1000000, 20},
		{"10000000 followers capped at 20", 10000000, 20},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := followerBonus(tt.followers)
			if got != tt.want {
				t.Errorf("followerBonus(%d) = %d, want %d", tt.followers, got, tt.want)
			}
		})
	}
}
