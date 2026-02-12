package store

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

type CollaborationRequest struct {
	ID         string    `json:"id"`
	FromUserID string    `json:"fromUserId"`
	ToUserID   string    `json:"toUserId"`
	Message    string    `json:"message"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`

	// Joined fields for listing
	FromName     *string `json:"fromName,omitempty"`
	FromUsername *string `json:"fromUsername,omitempty"`
	FromImage    *string `json:"fromImage,omitempty"`
	ToName       *string `json:"toName,omitempty"`
	ToUsername   *string `json:"toUsername,omitempty"`
	ToImage      *string `json:"toImage,omitempty"`
}

func (s *Store) CreateCollaborationRequest(ctx context.Context, id, fromUserID, toUserID, message string) error {
	_, err := s.pool.Exec(ctx,
		`INSERT INTO collaboration_requests (id, from_user_id, to_user_id, message)
		 VALUES ($1, $2, $3, $4)`,
		id, fromUserID, toUserID, message,
	)
	return err
}

func (s *Store) FindCollaborationRequest(ctx context.Context, id string) (*CollaborationRequest, error) {
	var cr CollaborationRequest
	err := s.pool.QueryRow(ctx,
		`SELECT id, from_user_id, to_user_id, message, status, created_at, updated_at
		 FROM collaboration_requests WHERE id = $1`, id,
	).Scan(&cr.ID, &cr.FromUserID, &cr.ToUserID, &cr.Message, &cr.Status, &cr.CreatedAt, &cr.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return &cr, err
}

func (s *Store) UpdateCollaborationStatus(ctx context.Context, id, status string) error {
	_, err := s.pool.Exec(ctx,
		`UPDATE collaboration_requests SET status = $1, updated_at = NOW() WHERE id = $2`,
		status, id,
	)
	return err
}

// IncomingCollaborations returns requests sent TO this user
func (s *Store) IncomingCollaborations(ctx context.Context, userID string) ([]CollaborationRequest, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT cr.id, cr.from_user_id, cr.to_user_id, cr.message, cr.status, cr.created_at, cr.updated_at,
		        f.name, f.username, f.image
		 FROM collaboration_requests cr
		 JOIN users f ON f.id = cr.from_user_id
		 WHERE cr.to_user_id = $1
		 ORDER BY cr.created_at DESC`, userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []CollaborationRequest
	for rows.Next() {
		var cr CollaborationRequest
		if err := rows.Scan(&cr.ID, &cr.FromUserID, &cr.ToUserID, &cr.Message, &cr.Status, &cr.CreatedAt, &cr.UpdatedAt,
			&cr.FromName, &cr.FromUsername, &cr.FromImage); err != nil {
			return nil, err
		}
		results = append(results, cr)
	}
	if results == nil {
		results = []CollaborationRequest{}
	}
	return results, nil
}

// OutgoingCollaborations returns requests sent BY this user
func (s *Store) OutgoingCollaborations(ctx context.Context, userID string) ([]CollaborationRequest, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT cr.id, cr.from_user_id, cr.to_user_id, cr.message, cr.status, cr.created_at, cr.updated_at,
		        t.name, t.username, t.image
		 FROM collaboration_requests cr
		 JOIN users t ON t.id = cr.to_user_id
		 WHERE cr.from_user_id = $1
		 ORDER BY cr.created_at DESC`, userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []CollaborationRequest
	for rows.Next() {
		var cr CollaborationRequest
		if err := rows.Scan(&cr.ID, &cr.FromUserID, &cr.ToUserID, &cr.Message, &cr.Status, &cr.CreatedAt, &cr.UpdatedAt,
			&cr.ToName, &cr.ToUsername, &cr.ToImage); err != nil {
			return nil, err
		}
		results = append(results, cr)
	}
	if results == nil {
		results = []CollaborationRequest{}
	}
	return results, nil
}

type DiscoverUser struct {
	ID           string  `json:"id"`
	Name         *string `json:"name"`
	Username     *string `json:"username"`
	Image        *string `json:"image"`
	Bio          *string `json:"bio"`
	CreatorScore *int    `json:"creatorScore"`
	IsVerified   bool    `json:"isVerified"`
	Connections  int     `json:"connections"`
}

func (s *Store) DiscoverCreators(ctx context.Context, minScore int, platform, search string, limit, offset int) ([]DiscoverUser, int, error) {
	// Build dynamic query based on filters
	baseWhere := `WHERE u.onboarded = true AND u.role = 'CREATOR'`
	args := []interface{}{}
	argIdx := 1

	if minScore > 0 {
		baseWhere += ` AND u.creator_score >= $` + itoa(argIdx)
		args = append(args, minScore)
		argIdx++
	}

	if platform != "" {
		baseWhere += ` AND EXISTS (SELECT 1 FROM connections c WHERE c.user_id = u.id AND c.platform = $` + itoa(argIdx) + `)`
		args = append(args, platform)
		argIdx++
	}

	if search != "" {
		baseWhere += ` AND (u.name ILIKE '%' || $` + itoa(argIdx) + ` || '%' OR u.username ILIKE '%' || $` + itoa(argIdx) + ` || '%' OR u.bio ILIKE '%' || $` + itoa(argIdx) + ` || '%')`
		args = append(args, search)
		argIdx++
	}

	// Count
	var total int
	countQuery := `SELECT COUNT(*) FROM users u ` + baseWhere
	err := s.pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Fetch
	selectQuery := `SELECT u.id, u.name, u.username, u.image, u.bio, u.creator_score, u.is_verified,
	                       COALESCE((SELECT COUNT(*) FROM connections WHERE user_id = u.id), 0)
	                FROM users u ` + baseWhere +
		` ORDER BY u.creator_score DESC NULLS LAST LIMIT $` + itoa(argIdx) + ` OFFSET $` + itoa(argIdx+1)
	args = append(args, limit, offset)

	rows, err := s.pool.Query(ctx, selectQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []DiscoverUser
	for rows.Next() {
		var u DiscoverUser
		if err := rows.Scan(&u.ID, &u.Name, &u.Username, &u.Image, &u.Bio, &u.CreatorScore, &u.IsVerified, &u.Connections); err != nil {
			return nil, 0, err
		}
		users = append(users, u)
	}
	if users == nil {
		users = []DiscoverUser{}
	}
	return users, total, nil
}

func itoa(i int) string {
	return fmt.Sprintf("%d", i)
}
