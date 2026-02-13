package store

import (
	"context"
	"encoding/json"
	"time"

)

type ErrorLogEntry struct {
	ID        int64           `json:"id"`
	Source    string          `json:"source"`
	Level     string          `json:"level"`
	Message   string          `json:"message"`
	Stack     *string         `json:"stack,omitempty"`
	URL       *string         `json:"url,omitempty"`
	UserAgent *string         `json:"user_agent,omitempty"`
	UserID    *string         `json:"user_id,omitempty"`
	Metadata  json.RawMessage `json:"metadata,omitempty"`
	CreatedAt time.Time       `json:"created_at"`
}

func (s *Store) CreateErrorLog(ctx context.Context, source, level, message string, stack, url, userAgent, userID *string, metadata json.RawMessage) error {
	if metadata == nil {
		metadata = json.RawMessage("{}")
	}

	query := `
		INSERT INTO error_log (source, level, message, stack, url, user_agent, user_id, metadata)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := s.pool.Exec(ctx, query, source, level, message, stack, url, userAgent, userID, metadata)
	return err
}

func (s *Store) ListErrorLog(ctx context.Context, source string, limit, offset int) ([]ErrorLogEntry, int, error) {
	var total int
	var countQuery string
	var listQuery string
	var countArgs []interface{}
	var listArgs []interface{}

	if source != "" {
		countQuery = `SELECT COUNT(*) FROM error_log WHERE source = $1`
		countArgs = []interface{}{source}
		listQuery = `
			SELECT id, source, level, message, stack, url, user_agent, user_id, metadata, created_at
			FROM error_log
			WHERE source = $1
			ORDER BY created_at DESC
			LIMIT $2 OFFSET $3
		`
		listArgs = []interface{}{source, limit, offset}
	} else {
		countQuery = `SELECT COUNT(*) FROM error_log`
		countArgs = []interface{}{}
		listQuery = `
			SELECT id, source, level, message, stack, url, user_agent, user_id, metadata, created_at
			FROM error_log
			ORDER BY created_at DESC
			LIMIT $1 OFFSET $2
		`
		listArgs = []interface{}{limit, offset}
	}

	err := s.pool.QueryRow(ctx, countQuery, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := s.pool.Query(ctx, listQuery, listArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var entries []ErrorLogEntry
	for rows.Next() {
		var e ErrorLogEntry
		err := rows.Scan(&e.ID, &e.Source, &e.Level, &e.Message, &e.Stack, &e.URL, &e.UserAgent, &e.UserID, &e.Metadata, &e.CreatedAt)
		if err != nil {
			return nil, 0, err
		}
		entries = append(entries, e)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return entries, total, nil
}

func (s *Store) DeleteOldErrors(ctx context.Context, olderThan time.Time) (int64, error) {
	query := `DELETE FROM error_log WHERE created_at < $1`
	result, err := s.pool.Exec(ctx, query, olderThan)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}
