package store

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

type APIKey struct {
	ID         string     `json:"id"`
	UserID     string     `json:"userId"`
	Name       string     `json:"name"`
	KeyPrefix  string     `json:"keyPrefix"`
	KeyHash    string     `json:"-"`
	LastUsedAt *time.Time `json:"lastUsedAt"`
	CreatedAt  time.Time  `json:"createdAt"`
}

func (s *Store) CreateAPIKey(ctx context.Context, key *APIKey) error {
	_, err := s.pool.Exec(ctx,
		`INSERT INTO api_keys (id, user_id, name, key_prefix, key_hash, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		key.ID, key.UserID, key.Name, key.KeyPrefix, key.KeyHash, key.CreatedAt,
	)
	return err
}

func (s *Store) FindAPIKeysByUserID(ctx context.Context, userID string) ([]*APIKey, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT id, user_id, name, key_prefix, key_hash, last_used_at, created_at
		 FROM api_keys WHERE user_id = $1
		 ORDER BY created_at DESC`, userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var keys []*APIKey
	for rows.Next() {
		var k APIKey
		if err := rows.Scan(&k.ID, &k.UserID, &k.Name, &k.KeyPrefix, &k.KeyHash, &k.LastUsedAt, &k.CreatedAt); err != nil {
			return nil, err
		}
		keys = append(keys, &k)
	}
	if keys == nil {
		keys = []*APIKey{}
	}
	return keys, rows.Err()
}

func (s *Store) FindAPIKeyByHash(ctx context.Context, keyHash string) (*APIKey, error) {
	var k APIKey
	err := s.pool.QueryRow(ctx,
		`SELECT id, user_id, name, key_prefix, key_hash, last_used_at, created_at
		 FROM api_keys WHERE key_hash = $1`, keyHash,
	).Scan(&k.ID, &k.UserID, &k.Name, &k.KeyPrefix, &k.KeyHash, &k.LastUsedAt, &k.CreatedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return &k, err
}

func (s *Store) DeleteAPIKey(ctx context.Context, id, userID string) error {
	_, err := s.pool.Exec(ctx,
		`DELETE FROM api_keys WHERE id = $1 AND user_id = $2`,
		id, userID,
	)
	return err
}

func (s *Store) UpdateAPIKeyLastUsed(ctx context.Context, id string, usedAt time.Time) error {
	_, err := s.pool.Exec(ctx,
		`UPDATE api_keys SET last_used_at = $1 WHERE id = $2`,
		usedAt, id,
	)
	return err
}
