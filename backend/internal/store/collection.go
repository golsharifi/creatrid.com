package store

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

type Collection struct {
	ID            string    `json:"id"`
	UserID        string    `json:"userId"`
	Title         string    `json:"title"`
	Description   *string   `json:"description"`
	IsPublic      bool      `json:"isPublic"`
	CoverImageURL *string   `json:"coverImageUrl"`
	ItemCount     int       `json:"itemCount"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

type CollectionItemDetail struct {
	ContentID    string    `json:"contentId"`
	Title        string    `json:"title"`
	ContentType  string    `json:"contentType"`
	ThumbnailURL *string   `json:"thumbnailUrl"`
	Position     int       `json:"position"`
	AddedAt      time.Time `json:"addedAt"`
}

func (s *Store) CreateCollection(ctx context.Context, coll *Collection) error {
	_, err := s.pool.Exec(ctx,
		`INSERT INTO content_collections (id, user_id, title, description, is_public, cover_image_url, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		coll.ID, coll.UserID, coll.Title, coll.Description, coll.IsPublic, coll.CoverImageURL, coll.CreatedAt, coll.UpdatedAt,
	)
	return err
}

func (s *Store) FindCollectionByID(ctx context.Context, id string) (*Collection, error) {
	var coll Collection
	err := s.pool.QueryRow(ctx,
		`SELECT c.id, c.user_id, c.title, c.description, c.is_public, c.cover_image_url, c.created_at, c.updated_at,
		        COALESCE((SELECT COUNT(*) FROM collection_items ci WHERE ci.collection_id = c.id), 0) AS item_count
		 FROM content_collections c
		 WHERE c.id = $1`, id,
	).Scan(
		&coll.ID, &coll.UserID, &coll.Title, &coll.Description, &coll.IsPublic, &coll.CoverImageURL,
		&coll.CreatedAt, &coll.UpdatedAt, &coll.ItemCount,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return &coll, err
}

func (s *Store) ListCollectionsByUser(ctx context.Context, userID string, limit, offset int) ([]*Collection, int, error) {
	var total int
	err := s.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM content_collections WHERE user_id = $1`, userID,
	).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := s.pool.Query(ctx,
		`SELECT c.id, c.user_id, c.title, c.description, c.is_public, c.cover_image_url, c.created_at, c.updated_at,
		        COALESCE((SELECT COUNT(*) FROM collection_items ci WHERE ci.collection_id = c.id), 0) AS item_count
		 FROM content_collections c
		 WHERE c.user_id = $1
		 ORDER BY c.created_at DESC
		 LIMIT $2 OFFSET $3`,
		userID, limit, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var collections []*Collection
	for rows.Next() {
		var coll Collection
		if err := rows.Scan(
			&coll.ID, &coll.UserID, &coll.Title, &coll.Description, &coll.IsPublic, &coll.CoverImageURL,
			&coll.CreatedAt, &coll.UpdatedAt, &coll.ItemCount,
		); err != nil {
			return nil, 0, err
		}
		collections = append(collections, &coll)
	}
	if collections == nil {
		collections = []*Collection{}
	}
	return collections, total, nil
}

func (s *Store) ListPublicCollectionsByUser(ctx context.Context, userID string, limit, offset int) ([]*Collection, int, error) {
	var total int
	err := s.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM content_collections WHERE user_id = $1 AND is_public = true`, userID,
	).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := s.pool.Query(ctx,
		`SELECT c.id, c.user_id, c.title, c.description, c.is_public, c.cover_image_url, c.created_at, c.updated_at,
		        COALESCE((SELECT COUNT(*) FROM collection_items ci WHERE ci.collection_id = c.id), 0) AS item_count
		 FROM content_collections c
		 WHERE c.user_id = $1 AND c.is_public = true
		 ORDER BY c.created_at DESC
		 LIMIT $2 OFFSET $3`,
		userID, limit, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var collections []*Collection
	for rows.Next() {
		var coll Collection
		if err := rows.Scan(
			&coll.ID, &coll.UserID, &coll.Title, &coll.Description, &coll.IsPublic, &coll.CoverImageURL,
			&coll.CreatedAt, &coll.UpdatedAt, &coll.ItemCount,
		); err != nil {
			return nil, 0, err
		}
		collections = append(collections, &coll)
	}
	if collections == nil {
		collections = []*Collection{}
	}
	return collections, total, nil
}

func (s *Store) UpdateCollection(ctx context.Context, id, title string, description *string, isPublic bool) error {
	_, err := s.pool.Exec(ctx,
		`UPDATE content_collections SET title = $1, description = $2, is_public = $3, updated_at = NOW() WHERE id = $4`,
		title, description, isPublic, id,
	)
	return err
}

func (s *Store) DeleteCollection(ctx context.Context, id string) error {
	_, err := s.pool.Exec(ctx, `DELETE FROM content_collections WHERE id = $1`, id)
	return err
}

func (s *Store) AddItemToCollection(ctx context.Context, collectionID, contentID string, position int) error {
	_, err := s.pool.Exec(ctx,
		`INSERT INTO collection_items (collection_id, content_id, position) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`,
		collectionID, contentID, position,
	)
	return err
}

func (s *Store) RemoveItemFromCollection(ctx context.Context, collectionID, contentID string) error {
	_, err := s.pool.Exec(ctx,
		`DELETE FROM collection_items WHERE collection_id = $1 AND content_id = $2`,
		collectionID, contentID,
	)
	return err
}

func (s *Store) ListCollectionItems(ctx context.Context, collectionID string) ([]*CollectionItemDetail, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT ci_item.content_id, ct.title, ct.content_type, ct.thumbnail_url, ci_item.position, ci_item.added_at
		 FROM collection_items ci_item
		 JOIN content_items ct ON ct.id = ci_item.content_id
		 WHERE ci_item.collection_id = $1
		 ORDER BY ci_item.position`,
		collectionID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*CollectionItemDetail
	for rows.Next() {
		var item CollectionItemDetail
		if err := rows.Scan(&item.ContentID, &item.Title, &item.ContentType, &item.ThumbnailURL, &item.Position, &item.AddedAt); err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	if items == nil {
		items = []*CollectionItemDetail{}
	}
	return items, nil
}
