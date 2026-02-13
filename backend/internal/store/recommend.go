package store

import (
	"context"
	"fmt"
	"strings"
)

type RecommendedCreator struct {
	ID              string  `json:"id"`
	Name            *string `json:"name"`
	Username        *string `json:"username"`
	Image           *string `json:"image"`
	Bio             *string `json:"bio"`
	CreatorScore    *int    `json:"creatorScore"`
	CreatorTier     string  `json:"creatorTier"`
	IsVerified      bool    `json:"isVerified"`
	SharedPlatforms int     `json:"sharedPlatforms"`
}

func (s *Store) RecommendCreators(ctx context.Context, userID string, userPlatforms []string, userScore int, limit int) ([]*RecommendedCreator, error) {
	if len(userPlatforms) == 0 {
		return []*RecommendedCreator{}, nil
	}

	// Build platform matching condition
	platformPlaceholders := make([]string, len(userPlatforms))
	args := []interface{}{userID}
	for i, p := range userPlatforms {
		args = append(args, p)
		platformPlaceholders[i] = fmt.Sprintf("$%d", i+2)
	}

	scoreMin := userScore - 20
	if scoreMin < 0 {
		scoreMin = 0
	}
	scoreMax := userScore + 20
	args = append(args, scoreMin, scoreMax, limit)
	scoreMinIdx := len(userPlatforms) + 2
	scoreMaxIdx := scoreMinIdx + 1
	limitIdx := scoreMaxIdx + 1

	query := fmt.Sprintf(`
		SELECT u.id, u.name, u.username, u.image, u.bio, u.creator_score, u.creator_tier, u.is_verified,
		       COUNT(c.platform) AS shared_platforms
		FROM users u
		JOIN connections c ON c.user_id = u.id AND c.platform IN (%s)
		WHERE u.id != $1
		  AND u.onboarded = true
		  AND u.creator_score IS NOT NULL
		  AND u.creator_score BETWEEN $%d AND $%d
		GROUP BY u.id
		ORDER BY shared_platforms DESC, u.creator_score DESC
		LIMIT $%d`,
		strings.Join(platformPlaceholders, ","),
		scoreMinIdx, scoreMaxIdx, limitIdx,
	)

	rows, err := s.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var creators []*RecommendedCreator
	for rows.Next() {
		var c RecommendedCreator
		if err := rows.Scan(&c.ID, &c.Name, &c.Username, &c.Image, &c.Bio, &c.CreatorScore, &c.CreatorTier, &c.IsVerified, &c.SharedPlatforms); err != nil {
			return nil, err
		}
		creators = append(creators, &c)
	}
	if creators == nil {
		creators = []*RecommendedCreator{}
	}
	return creators, nil
}
