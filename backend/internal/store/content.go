package store

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

type ContentItem struct {
	ID           string    `json:"id"`
	UserID       string    `json:"userId"`
	Title        string    `json:"title"`
	Description  *string   `json:"description"`
	ContentType  string    `json:"contentType"`
	MimeType     string    `json:"mimeType"`
	FileSize     int64     `json:"fileSize"`
	FileURL      string    `json:"fileUrl"`
	ThumbnailURL *string   `json:"thumbnailUrl"`
	HashSHA256   string    `json:"hashSha256"`
	IsPublic     bool      `json:"isPublic"`
	Tags         []string  `json:"tags"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

func (s *Store) CreateContentItem(ctx context.Context, item *ContentItem) error {
	_, err := s.pool.Exec(ctx,
		`INSERT INTO content_items (id, user_id, title, description, content_type, mime_type, file_size, file_url, thumbnail_url, hash_sha256, is_public, tags, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`,
		item.ID, item.UserID, item.Title, item.Description, item.ContentType,
		item.MimeType, item.FileSize, item.FileURL, item.ThumbnailURL,
		item.HashSHA256, item.IsPublic, item.Tags, item.CreatedAt, item.UpdatedAt,
	)
	return err
}

func (s *Store) FindContentItemByID(ctx context.Context, id string) (*ContentItem, error) {
	var item ContentItem
	err := s.pool.QueryRow(ctx,
		`SELECT id, user_id, title, description, content_type, mime_type, file_size, file_url, thumbnail_url, hash_sha256, is_public, tags, created_at, updated_at
		 FROM content_items WHERE id = $1`, id,
	).Scan(
		&item.ID, &item.UserID, &item.Title, &item.Description, &item.ContentType,
		&item.MimeType, &item.FileSize, &item.FileURL, &item.ThumbnailURL,
		&item.HashSHA256, &item.IsPublic, &item.Tags, &item.CreatedAt, &item.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return &item, err
}

func (s *Store) ListContentItemsByUser(ctx context.Context, userID string, limit, offset int) ([]*ContentItem, int, error) {
	var total int
	err := s.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM content_items WHERE user_id = $1`, userID,
	).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := s.pool.Query(ctx,
		`SELECT id, user_id, title, description, content_type, mime_type, file_size, file_url, thumbnail_url, hash_sha256, is_public, tags, created_at, updated_at
		 FROM content_items
		 WHERE user_id = $1
		 ORDER BY created_at DESC
		 LIMIT $2 OFFSET $3`,
		userID, limit, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var items []*ContentItem
	for rows.Next() {
		var item ContentItem
		if err := rows.Scan(
			&item.ID, &item.UserID, &item.Title, &item.Description, &item.ContentType,
			&item.MimeType, &item.FileSize, &item.FileURL, &item.ThumbnailURL,
			&item.HashSHA256, &item.IsPublic, &item.Tags, &item.CreatedAt, &item.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		items = append(items, &item)
	}
	if items == nil {
		items = []*ContentItem{}
	}
	return items, total, nil
}

func (s *Store) UpdateContentItem(ctx context.Context, id, title string, description *string, tags []string, isPublic bool) error {
	_, err := s.pool.Exec(ctx,
		`UPDATE content_items SET title = $1, description = $2, tags = $3, is_public = $4, updated_at = NOW() WHERE id = $5`,
		title, description, tags, isPublic, id,
	)
	return err
}

func (s *Store) DeleteContentItem(ctx context.Context, id string) error {
	_, err := s.pool.Exec(ctx, `DELETE FROM content_items WHERE id = $1`, id)
	return err
}

func (s *Store) FindContentByHash(ctx context.Context, hash string) (*ContentItem, error) {
	var item ContentItem
	err := s.pool.QueryRow(ctx,
		`SELECT id, user_id, title, description, content_type, mime_type, file_size, file_url, thumbnail_url, hash_sha256, is_public, tags, created_at, updated_at
		 FROM content_items WHERE hash_sha256 = $1`, hash,
	).Scan(
		&item.ID, &item.UserID, &item.Title, &item.Description, &item.ContentType,
		&item.MimeType, &item.FileSize, &item.FileURL, &item.ThumbnailURL,
		&item.HashSHA256, &item.IsPublic, &item.Tags, &item.CreatedAt, &item.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return &item, err
}

func (s *Store) ListPublicContent(ctx context.Context, contentType, query string, limit, offset int) ([]*ContentItem, int, error) {
	baseWhere := `WHERE is_public = true`
	args := []interface{}{}
	argIdx := 1

	if contentType != "" {
		baseWhere += ` AND content_type = $` + fmt.Sprintf("%d", argIdx)
		args = append(args, contentType)
		argIdx++
	}

	if query != "" {
		baseWhere += ` AND (title ILIKE '%' || $` + fmt.Sprintf("%d", argIdx) + ` || '%' OR description ILIKE '%' || $` + fmt.Sprintf("%d", argIdx) + ` || '%')`
		args = append(args, query)
		argIdx++
	}

	var total int
	countQuery := `SELECT COUNT(*) FROM content_items ` + baseWhere
	err := s.pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	selectQuery := `SELECT id, user_id, title, description, content_type, mime_type, file_size, file_url, thumbnail_url, hash_sha256, is_public, tags, created_at, updated_at
		 FROM content_items ` + baseWhere +
		` ORDER BY created_at DESC LIMIT $` + fmt.Sprintf("%d", argIdx) + ` OFFSET $` + fmt.Sprintf("%d", argIdx+1)
	args = append(args, limit, offset)

	rows, err := s.pool.Query(ctx, selectQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var items []*ContentItem
	for rows.Next() {
		var item ContentItem
		if err := rows.Scan(
			&item.ID, &item.UserID, &item.Title, &item.Description, &item.ContentType,
			&item.MimeType, &item.FileSize, &item.FileURL, &item.ThumbnailURL,
			&item.HashSHA256, &item.IsPublic, &item.Tags, &item.CreatedAt, &item.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		items = append(items, &item)
	}
	if items == nil {
		items = []*ContentItem{}
	}
	return items, total, nil
}

func (s *Store) CountContentByUser(ctx context.Context, userID string) (int, error) {
	var count int
	err := s.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM content_items WHERE user_id = $1`, userID,
	).Scan(&count)
	return count, err
}

