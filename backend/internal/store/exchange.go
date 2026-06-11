package store

import (
	"context"
	"errors"
	"time"
)

// Closed-loop sticker exchange: stickers trade for arena points, entirely
// in-platform. Listings escrow the sticker (decremented from the seller on
// listing, returned on cancel, granted to the buyer on sale).

type StickerListing struct {
	ID          string     `json:"id"`
	SellerID    string     `json:"sellerId"`
	Seller      *string    `json:"seller"` // username
	StickerID   string     `json:"stickerId"`
	Name        string     `json:"name"`
	Emoji       string     `json:"emoji"`
	Rarity      string     `json:"rarity"`
	PricePoints int        `json:"pricePoints"`
	Status      string     `json:"status"`
	CreatedAt   time.Time  `json:"createdAt"`
	SoldAt      *time.Time `json:"soldAt"`
}

type PointEvent struct {
	Delta     int       `json:"delta"`
	Reason    string    `json:"reason"`
	CreatedAt time.Time `json:"createdAt"`
}

var (
	ErrListingUnavailable = errors.New("listing unavailable")
	ErrInsufficientPoints = errors.New("not enough points")
)

// ListPointEvents returns the user's recent point ledger (the wallet history).
func (s *Store) ListPointEvents(ctx context.Context, userID string, limit int) ([]*PointEvent, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT delta, reason, created_at FROM point_events
		 WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2`,
		userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*PointEvent
	for rows.Next() {
		var e PointEvent
		if err := rows.Scan(&e.Delta, &e.Reason, &e.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, &e)
	}
	if out == nil {
		out = []*PointEvent{}
	}
	return out, nil
}

// CreateListing escrows one sticker from the seller and opens a listing.
func (s *Store) CreateListing(ctx context.Context, id, sellerID, stickerID string, pricePoints int) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	res, err := tx.Exec(ctx,
		`UPDATE user_stickers SET count = count - 1
		 WHERE user_id = $1 AND sticker_id = $2 AND count > 0`,
		sellerID, stickerID)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return ErrInsufficientStickers
	}
	if _, err := tx.Exec(ctx,
		`DELETE FROM user_stickers WHERE user_id = $1 AND sticker_id = $2 AND count <= 0`,
		sellerID, stickerID); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx,
		`INSERT INTO sticker_listings (id, seller_id, sticker_id, price_points) VALUES ($1, $2, $3, $4)`,
		id, sellerID, stickerID, pricePoints); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// CancelListing returns the escrowed sticker to the seller.
func (s *Store) CancelListing(ctx context.Context, listingID, sellerID string) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var stickerID string
	err = tx.QueryRow(ctx,
		`UPDATE sticker_listings SET status = 'cancelled'
		 WHERE id = $1 AND seller_id = $2 AND status = 'open'
		 RETURNING sticker_id`,
		listingID, sellerID).Scan(&stickerID)
	if err != nil {
		return ErrListingUnavailable
	}
	if _, err := tx.Exec(ctx,
		`INSERT INTO user_stickers (user_id, sticker_id, count) VALUES ($1, $2, 1)
		 ON CONFLICT (user_id, sticker_id) DO UPDATE SET count = user_stickers.count + 1`,
		sellerID, stickerID); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// BuyListing atomically: checks buyer points, moves points buyer→seller,
// grants the sticker, marks the listing sold, and writes both ledger entries.
func (s *Store) BuyListing(ctx context.Context, eventID1, eventID2, listingID, buyerID string) (*StickerListing, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var l StickerListing
	err = tx.QueryRow(ctx,
		`UPDATE sticker_listings SET status = 'sold', buyer_id = $2, sold_at = NOW()
		 WHERE id = $1 AND status = 'open' AND seller_id <> $2
		 RETURNING id, seller_id, sticker_id, price_points`,
		listingID, buyerID).Scan(&l.ID, &l.SellerID, &l.StickerID, &l.PricePoints)
	if err != nil {
		return nil, ErrListingUnavailable
	}

	// Deduct buyer points only if they have enough.
	res, err := tx.Exec(ctx,
		`UPDATE users SET arena_points = arena_points - $2, updated_at = NOW()
		 WHERE id = $1 AND arena_points >= $2`,
		buyerID, l.PricePoints)
	if err != nil {
		return nil, err
	}
	if res.RowsAffected() == 0 {
		return nil, ErrInsufficientPoints
	}
	if _, err := tx.Exec(ctx,
		`UPDATE users SET arena_points = arena_points + $2, updated_at = NOW() WHERE id = $1`,
		l.SellerID, l.PricePoints); err != nil {
		return nil, err
	}
	if _, err := tx.Exec(ctx,
		`INSERT INTO user_stickers (user_id, sticker_id, count) VALUES ($1, $2, 1)
		 ON CONFLICT (user_id, sticker_id) DO UPDATE SET count = user_stickers.count + 1`,
		buyerID, l.StickerID); err != nil {
		return nil, err
	}
	if _, err := tx.Exec(ctx,
		`INSERT INTO point_events (id, user_id, delta, reason) VALUES ($1, $2, $3, 'trade_buy'), ($4, $5, $6, 'trade_sale')`,
		eventID1, buyerID, -l.PricePoints, eventID2, l.SellerID, l.PricePoints); err != nil {
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	l.Status = "sold"
	return &l, nil
}

// ListOpenListings returns the exchange floor, newest first.
func (s *Store) ListOpenListings(ctx context.Context, limit int) ([]*StickerListing, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT l.id, l.seller_id, u.username, l.sticker_id, st.name, st.emoji, st.rarity,
		        l.price_points, l.status, l.created_at, l.sold_at
		 FROM sticker_listings l
		 JOIN users u ON u.id = l.seller_id
		 JOIN stickers st ON st.id = l.sticker_id
		 WHERE l.status = 'open'
		 ORDER BY l.created_at DESC
		 LIMIT $1`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanListings(rows)
}

// ListUserListings returns a seller's own listings (any status), newest first.
func (s *Store) ListUserListings(ctx context.Context, sellerID string, limit int) ([]*StickerListing, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT l.id, l.seller_id, u.username, l.sticker_id, st.name, st.emoji, st.rarity,
		        l.price_points, l.status, l.created_at, l.sold_at
		 FROM sticker_listings l
		 JOIN users u ON u.id = l.seller_id
		 JOIN stickers st ON st.id = l.sticker_id
		 WHERE l.seller_id = $1
		 ORDER BY l.created_at DESC
		 LIMIT $2`, sellerID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanListings(rows)
}

type pgxRows interface {
	Next() bool
	Scan(dest ...any) error
}

func scanListings(rows pgxRows) ([]*StickerListing, error) {
	var out []*StickerListing
	for rows.Next() {
		var l StickerListing
		if err := rows.Scan(&l.ID, &l.SellerID, &l.Seller, &l.StickerID, &l.Name, &l.Emoji,
			&l.Rarity, &l.PricePoints, &l.Status, &l.CreatedAt, &l.SoldAt); err != nil {
			return nil, err
		}
		out = append(out, &l)
	}
	if out == nil {
		out = []*StickerListing{}
	}
	return out, nil
}
