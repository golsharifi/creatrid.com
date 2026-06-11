package store

import (
	"context"

	"github.com/jackc/pgx/v5"
)

// Signature track + watermark live in dedicated accessors instead of widening
// every users scan in user.go.

func (s *Store) SetSignatureTrack(ctx context.Context, userID string, contentID *string) error {
	_, err := s.pool.Exec(ctx,
		`UPDATE users SET signature_track_id = $2, updated_at = NOW() WHERE id = $1`,
		userID, contentID,
	)
	return err
}

func (s *Store) GetSignatureTrackID(ctx context.Context, userID string) (*string, error) {
	var id *string
	err := s.pool.QueryRow(ctx,
		`SELECT signature_track_id FROM users WHERE id = $1`, userID,
	).Scan(&id)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return id, err
}

func (s *Store) SetWatermarkURL(ctx context.Context, userID string, url *string) error {
	_, err := s.pool.Exec(ctx,
		`UPDATE users SET watermark_url = $2, updated_at = NOW() WHERE id = $1`,
		userID, url,
	)
	return err
}

func (s *Store) GetWatermarkURL(ctx context.Context, userID string) (*string, error) {
	var url *string
	err := s.pool.QueryRow(ctx,
		`SELECT watermark_url FROM users WHERE id = $1`, userID,
	).Scan(&url)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return url, err
}
