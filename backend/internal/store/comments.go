package store

import (
	"context"
	"time"
)

type ContentComment struct {
	ID        string    `json:"id"`
	ContentID string    `json:"contentId"`
	UserID    string    `json:"userId"`
	Username  *string   `json:"username"`
	Name      *string   `json:"name"`
	Image     *string   `json:"image"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"createdAt"`
}

func (s *Store) CreateContentComment(ctx context.Context, id, contentID, userID, body string) error {
	_, err := s.pool.Exec(ctx,
		`INSERT INTO content_comments (id, content_id, user_id, body) VALUES ($1, $2, $3, $4)`,
		id, contentID, userID, body,
	)
	return err
}

func (s *Store) ListContentComments(ctx context.Context, contentID string, limit, offset int) ([]*ContentComment, int, error) {
	var total int
	if err := s.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM content_comments WHERE content_id = $1`, contentID,
	).Scan(&total); err != nil {
		return nil, 0, err
	}

	rows, err := s.pool.Query(ctx,
		`SELECT c.id, c.content_id, c.user_id, u.username, u.name, u.image, c.body, c.created_at
		 FROM content_comments c JOIN users u ON u.id = c.user_id
		 WHERE c.content_id = $1
		 ORDER BY c.created_at DESC
		 LIMIT $2 OFFSET $3`,
		contentID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var out []*ContentComment
	for rows.Next() {
		var c ContentComment
		if err := rows.Scan(&c.ID, &c.ContentID, &c.UserID, &c.Username, &c.Name, &c.Image, &c.Body, &c.CreatedAt); err != nil {
			return nil, 0, err
		}
		out = append(out, &c)
	}
	if out == nil {
		out = []*ContentComment{}
	}
	return out, total, nil
}

// DeleteContentComment removes a comment if requester is its author, the
// content owner, or an admin (checked by the caller for admin).
func (s *Store) DeleteContentComment(ctx context.Context, commentID, requesterID string, isAdmin bool) (bool, error) {
	var res int64
	if isAdmin {
		ct, err := s.pool.Exec(ctx, `DELETE FROM content_comments WHERE id = $1`, commentID)
		if err != nil {
			return false, err
		}
		res = ct.RowsAffected()
	} else {
		ct, err := s.pool.Exec(ctx,
			`DELETE FROM content_comments c
			 WHERE c.id = $1 AND (c.user_id = $2 OR EXISTS (
			   SELECT 1 FROM content_items ci WHERE ci.id = c.content_id AND ci.user_id = $2))`,
			commentID, requesterID)
		if err != nil {
			return false, err
		}
		res = ct.RowsAffected()
	}
	return res > 0, nil
}
