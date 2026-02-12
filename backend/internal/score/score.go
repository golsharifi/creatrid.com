package score

import (
	"math"

	"github.com/creatrid/creatrid/internal/model"
)

// Calculate returns a creator score (0-100) based on profile, connections, and metrics.
func Calculate(user *model.User, connections []*model.Connection) int {
	score := 0

	// Profile completeness: 0-20 points
	if user.Name != nil && *user.Name != "" {
		score += 5
	}
	if user.Username != nil && *user.Username != "" {
		score += 5
	}
	if user.Bio != nil && *user.Bio != "" {
		score += 5
	}
	if user.Image != nil && *user.Image != "" {
		score += 5
	}

	// Email verified: 10 points
	if user.EmailVerified != nil {
		score += 10
	}

	// Connected platforms: 10 points each, max 50
	connectionPoints := len(connections) * 10
	if connectionPoints > 50 {
		connectionPoints = 50
	}
	score += connectionPoints

	// Follower bonus: 0-20 points (logarithmic scale)
	totalFollowers := 0
	for _, c := range connections {
		if c.FollowerCount != nil && *c.FollowerCount > 0 {
			totalFollowers += *c.FollowerCount
		}
	}
	score += followerBonus(totalFollowers)

	if score > 100 {
		score = 100
	}
	return score
}

// followerBonus returns 0-20 points based on total followers using log10 scale.
//
//	0         → 0
//	100       → 5
//	1,000     → 10
//	10,000    → 15
//	100,000+  → 20
func followerBonus(totalFollowers int) int {
	if totalFollowers <= 0 {
		return 0
	}
	// log10(100)=2, log10(100000)=5 → map [2,5] to [5,20]
	logVal := math.Log10(float64(totalFollowers))
	points := (logVal - 2) * 5
	if points < 0 {
		// Under 100 followers: partial credit
		points = logVal
	}
	bonus := int(math.Round(points))
	if bonus < 0 {
		return 0
	}
	if bonus > 20 {
		return 20
	}
	return bonus
}
