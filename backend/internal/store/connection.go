package store

import (
	"context"
	"encoding/json"

	"github.com/creatrid/creatrid/internal/model"
	"github.com/jackc/pgx/v5"
)

func (s *Store) UpsertConnection(ctx context.Context, c *model.Connection) error {
	_, err := s.pool.Exec(ctx,
		`INSERT INTO connections (id, user_id, platform, platform_user_id,
			username, display_name, avatar_url, profile_url,
			follower_count, access_token, refresh_token, token_expires_at, metadata)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		 ON CONFLICT (user_id, platform)
		 DO UPDATE SET
			platform_user_id = EXCLUDED.platform_user_id,
			username = EXCLUDED.username,
			display_name = EXCLUDED.display_name,
			avatar_url = EXCLUDED.avatar_url,
			profile_url = EXCLUDED.profile_url,
			follower_count = EXCLUDED.follower_count,
			access_token = EXCLUDED.access_token,
			refresh_token = COALESCE(EXCLUDED.refresh_token, connections.refresh_token),
			token_expires_at = EXCLUDED.token_expires_at,
			metadata = EXCLUDED.metadata`,
		c.ID, c.UserID, c.Platform, c.PlatformUserID,
		c.Username, c.DisplayName, c.AvatarURL, c.ProfileURL,
		c.FollowerCount, c.AccessToken, c.RefreshToken, c.TokenExpiresAt, c.Metadata,
	)
	return err
}

func (s *Store) FindConnectionsByUserID(ctx context.Context, userID string) ([]*model.Connection, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT id, user_id, platform, platform_user_id,
			username, display_name, avatar_url, profile_url,
			follower_count, access_token, refresh_token, token_expires_at,
			metadata, connected_at, updated_at
		 FROM connections WHERE user_id = $1
		 ORDER BY connected_at DESC`, userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var connections []*model.Connection
	for rows.Next() {
		var c model.Connection
		var metaBytes []byte
		err := rows.Scan(
			&c.ID, &c.UserID, &c.Platform, &c.PlatformUserID,
			&c.Username, &c.DisplayName, &c.AvatarURL, &c.ProfileURL,
			&c.FollowerCount, &c.AccessToken, &c.RefreshToken, &c.TokenExpiresAt,
			&metaBytes, &c.ConnectedAt, &c.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		c.Metadata = json.RawMessage(metaBytes)
		connections = append(connections, &c)
	}
	return connections, rows.Err()
}

func (s *Store) FindConnectionByUserAndPlatform(ctx context.Context, userID, platform string) (*model.Connection, error) {
	var c model.Connection
	var metaBytes []byte
	err := s.pool.QueryRow(ctx,
		`SELECT id, user_id, platform, platform_user_id,
			username, display_name, avatar_url, profile_url,
			follower_count, access_token, refresh_token, token_expires_at,
			metadata, connected_at, updated_at
		 FROM connections WHERE user_id = $1 AND platform = $2`,
		userID, platform,
	).Scan(
		&c.ID, &c.UserID, &c.Platform, &c.PlatformUserID,
		&c.Username, &c.DisplayName, &c.AvatarURL, &c.ProfileURL,
		&c.FollowerCount, &c.AccessToken, &c.RefreshToken, &c.TokenExpiresAt,
		&metaBytes, &c.ConnectedAt, &c.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	c.Metadata = json.RawMessage(metaBytes)
	return &c, err
}

func (s *Store) DeleteConnection(ctx context.Context, userID, platform string) error {
	_, err := s.pool.Exec(ctx,
		`DELETE FROM connections WHERE user_id = $1 AND platform = $2`,
		userID, platform,
	)
	return err
}

func (s *Store) CountConnectionsByUserID(ctx context.Context, userID string) (int, error) {
	var count int
	err := s.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM connections WHERE user_id = $1`, userID,
	).Scan(&count)
	return count, err
}
