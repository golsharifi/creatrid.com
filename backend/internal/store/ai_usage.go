package store

import (
	"context"
	"time"
)

// RecordAIGeneration logs one successful AI generation for quota accounting.
func (s *Store) RecordAIGeneration(ctx context.Context, id, userID, task string) error {
	_, err := s.pool.Exec(ctx,
		`INSERT INTO ai_generations (id, user_id, task, created_at) VALUES ($1, $2, $3, NOW())`,
		id, userID, task,
	)
	return err
}

// CountAIGenerationsSince returns how many generations a user has made since t.
func (s *Store) CountAIGenerationsSince(ctx context.Context, userID string, t time.Time) (int, error) {
	var count int
	err := s.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM ai_generations WHERE user_id = $1 AND created_at >= $2`,
		userID, t,
	).Scan(&count)
	return count, err
}
