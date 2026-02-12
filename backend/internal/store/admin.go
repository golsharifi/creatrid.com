package store

import (
	"context"

	"github.com/creatrid/creatrid/internal/model"
)

type AdminStats struct {
	TotalUsers       int `json:"totalUsers"`
	OnboardedUsers   int `json:"onboardedUsers"`
	VerifiedUsers    int `json:"verifiedUsers"`
	TotalConnections int `json:"totalConnections"`
	TotalViews       int `json:"totalViews"`
	TotalClicks      int `json:"totalClicks"`
}

type AdminUser struct {
	ID           string       `json:"id"`
	Name         *string      `json:"name"`
	Email        string       `json:"email"`
	Username     *string      `json:"username"`
	Image        *string      `json:"image"`
	Role         model.UserRole `json:"role"`
	CreatorScore *int         `json:"creatorScore"`
	IsVerified   bool         `json:"isVerified"`
	Onboarded    bool         `json:"onboarded"`
	Connections  int          `json:"connections"`
}

func (s *Store) AdminGetStats(ctx context.Context) (*AdminStats, error) {
	stats := &AdminStats{}

	err := s.pool.QueryRow(ctx, `SELECT COUNT(*) FROM users`).Scan(&stats.TotalUsers)
	if err != nil {
		return nil, err
	}
	err = s.pool.QueryRow(ctx, `SELECT COUNT(*) FROM users WHERE onboarded = true`).Scan(&stats.OnboardedUsers)
	if err != nil {
		return nil, err
	}
	err = s.pool.QueryRow(ctx, `SELECT COUNT(*) FROM users WHERE is_verified = true`).Scan(&stats.VerifiedUsers)
	if err != nil {
		return nil, err
	}
	err = s.pool.QueryRow(ctx, `SELECT COUNT(*) FROM connections`).Scan(&stats.TotalConnections)
	if err != nil {
		return nil, err
	}
	_ = s.pool.QueryRow(ctx, `SELECT COUNT(*) FROM profile_views`).Scan(&stats.TotalViews)
	_ = s.pool.QueryRow(ctx, `SELECT COUNT(*) FROM link_clicks`).Scan(&stats.TotalClicks)

	return stats, nil
}

func (s *Store) AdminListUsers(ctx context.Context, limit, offset int) ([]AdminUser, int, error) {
	var total int
	err := s.pool.QueryRow(ctx, `SELECT COUNT(*) FROM users`).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := s.pool.Query(ctx,
		`SELECT u.id, u.name, u.email, u.username, u.image, u.role,
		        u.creator_score, u.is_verified, u.onboarded,
		        COALESCE((SELECT COUNT(*) FROM connections WHERE user_id = u.id), 0)
		 FROM users u
		 ORDER BY u.created_at DESC
		 LIMIT $1 OFFSET $2`,
		limit, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []AdminUser
	for rows.Next() {
		var u AdminUser
		if err := rows.Scan(
			&u.ID, &u.Name, &u.Email, &u.Username, &u.Image, &u.Role,
			&u.CreatorScore, &u.IsVerified, &u.Onboarded, &u.Connections,
		); err != nil {
			return nil, 0, err
		}
		users = append(users, u)
	}

	if users == nil {
		users = []AdminUser{}
	}

	return users, total, nil
}

func (s *Store) AdminSetVerified(ctx context.Context, userID string, verified bool) error {
	_, err := s.pool.Exec(ctx,
		`UPDATE users SET is_verified = $1 WHERE id = $2`,
		verified, userID,
	)
	return err
}
