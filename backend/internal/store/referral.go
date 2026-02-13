package store

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

type ReferralReward struct {
	ID             string    `json:"id"`
	ReferrerUserID string    `json:"referrerUserId"`
	ReferredUserID string    `json:"referredUserId"`
	RewardType     string    `json:"rewardType"`
	RewardValue    int       `json:"rewardValue"`
	CreatedAt      time.Time `json:"createdAt"`
	ReferredName   *string   `json:"referredName,omitempty"`
}

type ReferralStats struct {
	TotalReferred int `json:"totalReferred"`
	TotalBonus    int `json:"totalBonus"`
}

func (s *Store) GetUserReferralCode(ctx context.Context, userID string) (*string, error) {
	var code *string
	err := s.pool.QueryRow(ctx, `SELECT referral_code FROM users WHERE id = $1`, userID).Scan(&code)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return code, err
}

func (s *Store) SetUserReferralCode(ctx context.Context, userID, code string) error {
	_, err := s.pool.Exec(ctx, `UPDATE users SET referral_code = $1 WHERE id = $2`, code, userID)
	return err
}

func (s *Store) FindUserByReferralCode(ctx context.Context, code string) (string, error) {
	var userID string
	err := s.pool.QueryRow(ctx, `SELECT id FROM users WHERE referral_code = $1`, code).Scan(&userID)
	if err == pgx.ErrNoRows {
		return "", nil
	}
	return userID, err
}

func (s *Store) SetUserReferredBy(ctx context.Context, userID, referrerID string) error {
	_, err := s.pool.Exec(ctx, `UPDATE users SET referred_by = $1 WHERE id = $2`, referrerID, userID)
	return err
}

func (s *Store) CreateReferralReward(ctx context.Context, id, referrerID, referredID, rewardType string, rewardValue int) error {
	_, err := s.pool.Exec(ctx,
		`INSERT INTO referral_rewards (id, referrer_user_id, referred_user_id, reward_type, reward_value, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		id, referrerID, referredID, rewardType, rewardValue, time.Now(),
	)
	return err
}

func (s *Store) ListReferralsByUser(ctx context.Context, userID string) ([]*ReferralReward, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT r.id, r.referrer_user_id, r.referred_user_id, r.reward_type, r.reward_value, r.created_at, u.name
		 FROM referral_rewards r
		 LEFT JOIN users u ON u.id = r.referred_user_id
		 WHERE r.referrer_user_id = $1
		 ORDER BY r.created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rewards []*ReferralReward
	for rows.Next() {
		var r ReferralReward
		if err := rows.Scan(&r.ID, &r.ReferrerUserID, &r.ReferredUserID, &r.RewardType, &r.RewardValue, &r.CreatedAt, &r.ReferredName); err != nil {
			return nil, err
		}
		rewards = append(rewards, &r)
	}
	if rewards == nil {
		rewards = []*ReferralReward{}
	}
	return rewards, nil
}

func (s *Store) GetReferralStats(ctx context.Context, userID string) (*ReferralStats, error) {
	stats := &ReferralStats{}
	err := s.pool.QueryRow(ctx,
		`SELECT COUNT(*), COALESCE(SUM(reward_value), 0) FROM referral_rewards WHERE referrer_user_id = $1`,
		userID,
	).Scan(&stats.TotalReferred, &stats.TotalBonus)
	return stats, err
}
