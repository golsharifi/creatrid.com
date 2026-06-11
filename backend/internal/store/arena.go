package store

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
)

// VR Community Arena storage. Points and stickers are closed-loop: they carry
// in-platform status only and are never cashable or externally transferable.

type Sticker struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Emoji       string `json:"emoji"`
	Rarity      string `json:"rarity"`
	Description string `json:"description"`
}

type UserSticker struct {
	Sticker
	Count         int       `json:"count"`
	FirstEarnedAt time.Time `json:"firstEarnedAt"`
}

type ArenaState struct {
	Points      int        `json:"points"`
	LastExplore *time.Time `json:"lastExplore"`
}

type LeaderboardEntry struct {
	UserID     string  `json:"userId"`
	Username   *string `json:"username"`
	Name       *string `json:"name"`
	Image      *string `json:"image"`
	Points     int     `json:"points"`
	Stickers   int     `json:"stickers"`
	IsVerified bool    `json:"isVerified"`
}

var ErrInsufficientStickers = errors.New("sticker not owned")

func (s *Store) GetArenaState(ctx context.Context, userID string) (*ArenaState, error) {
	var st ArenaState
	err := s.pool.QueryRow(ctx,
		`SELECT arena_points, arena_last_explore FROM users WHERE id = $1`, userID,
	).Scan(&st.Points, &st.LastExplore)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return &st, err
}

func (s *Store) ListStickers(ctx context.Context) ([]*Sticker, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT id, name, emoji, rarity, description FROM stickers ORDER BY
		 CASE rarity WHEN 'common' THEN 0 WHEN 'uncommon' THEN 1 WHEN 'rare' THEN 2 ELSE 3 END, name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*Sticker
	for rows.Next() {
		var st Sticker
		if err := rows.Scan(&st.ID, &st.Name, &st.Emoji, &st.Rarity, &st.Description); err != nil {
			return nil, err
		}
		out = append(out, &st)
	}
	if out == nil {
		out = []*Sticker{}
	}
	return out, nil
}

func (s *Store) ListUserStickers(ctx context.Context, userID string) ([]*UserSticker, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT st.id, st.name, st.emoji, st.rarity, st.description, us.count, us.first_earned_at
		 FROM user_stickers us JOIN stickers st ON st.id = us.sticker_id
		 WHERE us.user_id = $1
		 ORDER BY CASE st.rarity WHEN 'common' THEN 0 WHEN 'uncommon' THEN 1 WHEN 'rare' THEN 2 ELSE 3 END, st.name`,
		userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*UserSticker
	for rows.Next() {
		var us UserSticker
		if err := rows.Scan(&us.ID, &us.Name, &us.Emoji, &us.Rarity, &us.Description, &us.Count, &us.FirstEarnedAt); err != nil {
			return nil, err
		}
		out = append(out, &us)
	}
	if out == nil {
		out = []*UserSticker{}
	}
	return out, nil
}

// RecordExplore awards a sticker + points and stamps the cooldown, atomically.
func (s *Store) RecordExplore(ctx context.Context, eventID, userID, stickerID string, points int) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx,
		`UPDATE users SET arena_points = arena_points + $2, arena_last_explore = NOW(), updated_at = NOW() WHERE id = $1`,
		userID, points); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx,
		`INSERT INTO user_stickers (user_id, sticker_id, count) VALUES ($1, $2, 1)
		 ON CONFLICT (user_id, sticker_id) DO UPDATE SET count = user_stickers.count + 1`,
		userID, stickerID); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx,
		`INSERT INTO point_events (id, user_id, delta, reason) VALUES ($1, $2, $3, 'explore')`,
		eventID, userID, points); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// GiftSticker transfers one copy of a sticker from one user to another.
func (s *Store) GiftSticker(ctx context.Context, fromUserID, toUserID, stickerID string) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	res, err := tx.Exec(ctx,
		`UPDATE user_stickers SET count = count - 1
		 WHERE user_id = $1 AND sticker_id = $2 AND count > 0`,
		fromUserID, stickerID)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return ErrInsufficientStickers
	}
	if _, err := tx.Exec(ctx,
		`DELETE FROM user_stickers WHERE user_id = $1 AND sticker_id = $2 AND count <= 0`,
		fromUserID, stickerID); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx,
		`INSERT INTO user_stickers (user_id, sticker_id, count) VALUES ($1, $2, 1)
		 ON CONFLICT (user_id, sticker_id) DO UPDATE SET count = user_stickers.count + 1`,
		toUserID, stickerID); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (s *Store) ArenaLeaderboard(ctx context.Context, limit int) ([]*LeaderboardEntry, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT u.id, u.username, u.name, u.image, u.arena_points, u.is_verified,
		        COALESCE((SELECT SUM(count) FROM user_stickers us WHERE us.user_id = u.id), 0)
		 FROM users u
		 WHERE u.arena_points > 0 AND u.username IS NOT NULL
		 ORDER BY u.arena_points DESC
		 LIMIT $1`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*LeaderboardEntry
	for rows.Next() {
		var e LeaderboardEntry
		if err := rows.Scan(&e.UserID, &e.Username, &e.Name, &e.Image, &e.Points, &e.IsVerified, &e.Stickers); err != nil {
			return nil, err
		}
		out = append(out, &e)
	}
	if out == nil {
		out = []*LeaderboardEntry{}
	}
	return out, nil
}

// AddArenaPoints credits points with a ledger entry (e.g. comment_received).
func (s *Store) AddArenaPoints(ctx context.Context, eventID, userID string, delta int, reason string) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if _, err := tx.Exec(ctx,
		`UPDATE users SET arena_points = GREATEST(0, arena_points + $2), updated_at = NOW() WHERE id = $1`,
		userID, delta); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx,
		`INSERT INTO point_events (id, user_id, delta, reason) VALUES ($1, $2, $3, $4)`,
		eventID, userID, delta, reason); err != nil {
		return err
	}
	return tx.Commit(ctx)
}
