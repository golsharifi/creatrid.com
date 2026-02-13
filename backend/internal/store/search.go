package store

import (
	"context"
	"fmt"
)

type SearchContentResult struct {
	ID           string   `json:"id"`
	Title        string   `json:"title"`
	Description  *string  `json:"description"`
	ContentType  string   `json:"contentType"`
	ThumbnailURL *string  `json:"thumbnailUrl"`
	Tags         []string `json:"tags"`
	Rank         float64  `json:"rank"`
	CreatorName  *string  `json:"creatorName"`
	CreatorUsername *string `json:"creatorUsername"`
}

type SearchUserResult struct {
	ID           string  `json:"id"`
	Name         *string `json:"name"`
	Username     *string `json:"username"`
	Bio          *string `json:"bio"`
	Image        *string `json:"image"`
	CreatorScore *int    `json:"creatorScore"`
	IsVerified   bool    `json:"isVerified"`
}

func (s *Store) SearchContent(ctx context.Context, query string, limit, offset int) ([]*SearchContentResult, int, error) {
	var total int
	err := s.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM content_items WHERE is_public = true AND search_vector @@ plainto_tsquery('english', $1)`,
		query,
	).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := s.pool.Query(ctx,
		`SELECT ci.id, ci.title, ci.description, ci.content_type, ci.thumbnail_url, ci.tags,
		        ts_rank(ci.search_vector, plainto_tsquery('english', $1)) AS rank,
		        u.name, u.username
		 FROM content_items ci
		 JOIN users u ON u.id = ci.user_id
		 WHERE ci.is_public = true AND ci.search_vector @@ plainto_tsquery('english', $1)
		 ORDER BY rank DESC
		 LIMIT $2 OFFSET $3`,
		query, limit, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var results []*SearchContentResult
	for rows.Next() {
		var r SearchContentResult
		if err := rows.Scan(&r.ID, &r.Title, &r.Description, &r.ContentType, &r.ThumbnailURL, &r.Tags, &r.Rank, &r.CreatorName, &r.CreatorUsername); err != nil {
			return nil, 0, err
		}
		results = append(results, &r)
	}
	if results == nil {
		results = []*SearchContentResult{}
	}
	return results, total, nil
}

func (s *Store) SearchUsers(ctx context.Context, query string, limit, offset int) ([]*SearchUserResult, int, error) {
	pattern := "%" + query + "%"

	var total int
	err := s.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM users WHERE onboarded = true AND (name ILIKE $1 OR username ILIKE $1 OR bio ILIKE $1)`,
		pattern,
	).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := s.pool.Query(ctx,
		`SELECT id, name, username, bio, image, creator_score, is_verified
		 FROM users
		 WHERE onboarded = true AND (name ILIKE $1 OR username ILIKE $1 OR bio ILIKE $1)
		 ORDER BY creator_score DESC NULLS LAST
		 LIMIT $2 OFFSET $3`,
		pattern, limit, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var results []*SearchUserResult
	for rows.Next() {
		var r SearchUserResult
		if err := rows.Scan(&r.ID, &r.Name, &r.Username, &r.Bio, &r.Image, &r.CreatorScore, &r.IsVerified); err != nil {
			return nil, 0, err
		}
		results = append(results, &r)
	}
	if results == nil {
		results = []*SearchUserResult{}
	}
	return results, total, nil
}

func (s *Store) SearchSuggestions(ctx context.Context, query string, limit int) ([]string, error) {
	pattern := query + "%"

	rows, err := s.pool.Query(ctx,
		fmt.Sprintf(`SELECT DISTINCT unnest(tags) AS tag FROM content_items WHERE EXISTS (SELECT 1 FROM unnest(tags) t WHERE t ILIKE $1) LIMIT $2`),
		pattern, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var suggestions []string
	for rows.Next() {
		var tag string
		if err := rows.Scan(&tag); err != nil {
			return nil, err
		}
		suggestions = append(suggestions, tag)
	}
	if suggestions == nil {
		suggestions = []string{}
	}
	return suggestions, nil
}
