package store

import (
	"context"
	"time"
)

// ModerationFlag represents a content moderation flag record.
type ModerationFlag struct {
	ID         string     `json:"id"`
	ContentID  string     `json:"contentId"`
	Reason     string     `json:"reason"`
	Details    *string    `json:"details"`
	Status     string     `json:"status"`
	ResolvedBy *string    `json:"resolvedBy"`
	ResolvedAt *time.Time `json:"resolvedAt"`
	Notes      *string    `json:"notes"`
	CreatedAt  time.Time  `json:"createdAt"`
}

// CreateModerationFlag inserts a new moderation flag into the database.
func (s *Store) CreateModerationFlag(ctx context.Context, id, contentID, reason, details string) error {
	var detailsPtr *string
	if details != "" {
		detailsPtr = &details
	}
	_, err := s.pool.Exec(ctx,
		`INSERT INTO moderation_flags (id, content_id, reason, details, status, created_at)
		 VALUES ($1, $2, $3, $4, 'pending', NOW())`,
		id, contentID, reason, detailsPtr,
	)
	return err
}

// ListModerationFlags returns a paginated list of moderation flags filtered by status.
// If status is empty, all flags are returned.
func (s *Store) ListModerationFlags(ctx context.Context, status string, limit, offset int) ([]ModerationFlag, int, error) {
	var total int

	if status != "" {
		err := s.pool.QueryRow(ctx,
			`SELECT COUNT(*) FROM moderation_flags WHERE status = $1`, status,
		).Scan(&total)
		if err != nil {
			return nil, 0, err
		}
	} else {
		err := s.pool.QueryRow(ctx,
			`SELECT COUNT(*) FROM moderation_flags`,
		).Scan(&total)
		if err != nil {
			return nil, 0, err
		}
	}

	var query string
	var args []interface{}
	if status != "" {
		query = `SELECT id, content_id, reason, details, status, resolved_by, resolved_at, notes, created_at
				 FROM moderation_flags
				 WHERE status = $1
				 ORDER BY created_at DESC
				 LIMIT $2 OFFSET $3`
		args = []interface{}{status, limit, offset}
	} else {
		query = `SELECT id, content_id, reason, details, status, resolved_by, resolved_at, notes, created_at
				 FROM moderation_flags
				 ORDER BY created_at DESC
				 LIMIT $1 OFFSET $2`
		args = []interface{}{limit, offset}
	}

	rows, err := s.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	flags := []ModerationFlag{}
	for rows.Next() {
		var f ModerationFlag
		if err := rows.Scan(&f.ID, &f.ContentID, &f.Reason, &f.Details, &f.Status,
			&f.ResolvedBy, &f.ResolvedAt, &f.Notes, &f.CreatedAt); err != nil {
			return nil, 0, err
		}
		flags = append(flags, f)
	}
	return flags, total, nil
}

// ResolveModerationFlag updates the status of a moderation flag and records who
// resolved it and any notes.
func (s *Store) ResolveModerationFlag(ctx context.Context, id, resolvedBy, status, notes string) error {
	var notesPtr *string
	if notes != "" {
		notesPtr = &notes
	}
	_, err := s.pool.Exec(ctx,
		`UPDATE moderation_flags
		 SET status = $1, resolved_by = $2, resolved_at = NOW(), notes = $3
		 WHERE id = $4`,
		status, resolvedBy, notesPtr, id,
	)
	return err
}
