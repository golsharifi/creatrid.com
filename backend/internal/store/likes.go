package store

import "context"

// LikeContent records a like; returns true if it was newly created.
func (s *Store) LikeContent(ctx context.Context, contentID, userID string) (bool, error) {
	res, err := s.pool.Exec(ctx,
		`INSERT INTO content_likes (content_id, user_id) VALUES ($1, $2)
		 ON CONFLICT (content_id, user_id) DO NOTHING`,
		contentID, userID)
	if err != nil {
		return false, err
	}
	return res.RowsAffected() > 0, nil
}

// UnlikeContent removes a like; returns true if one existed.
func (s *Store) UnlikeContent(ctx context.Context, contentID, userID string) (bool, error) {
	res, err := s.pool.Exec(ctx,
		`DELETE FROM content_likes WHERE content_id = $1 AND user_id = $2`,
		contentID, userID)
	if err != nil {
		return false, err
	}
	return res.RowsAffected() > 0, nil
}

// ContentLikeStats returns the like count and whether userID liked it
// (userID may be empty for anonymous visitors).
func (s *Store) ContentLikeStats(ctx context.Context, contentID, userID string) (int, bool, error) {
	var count int
	var liked bool
	err := s.pool.QueryRow(ctx,
		`SELECT COUNT(*),
		        COALESCE(BOOL_OR(user_id = $2), false)
		 FROM content_likes WHERE content_id = $1`,
		contentID, userID,
	).Scan(&count, &liked)
	return count, liked, err
}
