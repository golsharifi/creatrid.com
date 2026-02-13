package store

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

type ContentAnchor struct {
	ID              string     `json:"id"`
	ContentID       string     `json:"contentId"`
	UserID          string     `json:"userId"`
	ContentHash     string     `json:"contentHash"`
	TxHash          *string    `json:"txHash"`
	Chain           string     `json:"chain"`
	BlockNumber     *int64     `json:"blockNumber"`
	ContractAddress *string    `json:"contractAddress"`
	AnchorStatus    string     `json:"anchorStatus"`
	ErrorMessage    *string    `json:"errorMessage,omitempty"`
	CreatedAt       time.Time  `json:"createdAt"`
	ConfirmedAt     *time.Time `json:"confirmedAt,omitempty"`
}

func (s *Store) CreateContentAnchor(ctx context.Context, anchor *ContentAnchor) error {
	_, err := s.pool.Exec(ctx,
		`INSERT INTO content_anchors (id, content_id, user_id, content_hash, tx_hash, chain, block_number, contract_address, anchor_status, error_message, created_at, confirmed_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
		anchor.ID, anchor.ContentID, anchor.UserID, anchor.ContentHash,
		anchor.TxHash, anchor.Chain, anchor.BlockNumber, anchor.ContractAddress,
		anchor.AnchorStatus, anchor.ErrorMessage, anchor.CreatedAt, anchor.ConfirmedAt,
	)
	return err
}

func (s *Store) FindAnchorByContentID(ctx context.Context, contentID string) (*ContentAnchor, error) {
	var a ContentAnchor
	err := s.pool.QueryRow(ctx,
		`SELECT id, content_id, user_id, content_hash, tx_hash, chain, block_number, contract_address, anchor_status, error_message, created_at, confirmed_at
		 FROM content_anchors WHERE content_id = $1`, contentID,
	).Scan(
		&a.ID, &a.ContentID, &a.UserID, &a.ContentHash, &a.TxHash,
		&a.Chain, &a.BlockNumber, &a.ContractAddress, &a.AnchorStatus,
		&a.ErrorMessage, &a.CreatedAt, &a.ConfirmedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return &a, err
}

func (s *Store) FindAnchorByHash(ctx context.Context, hash string) (*ContentAnchor, error) {
	var a ContentAnchor
	err := s.pool.QueryRow(ctx,
		`SELECT id, content_id, user_id, content_hash, tx_hash, chain, block_number, contract_address, anchor_status, error_message, created_at, confirmed_at
		 FROM content_anchors WHERE content_hash = $1`, hash,
	).Scan(
		&a.ID, &a.ContentID, &a.UserID, &a.ContentHash, &a.TxHash,
		&a.Chain, &a.BlockNumber, &a.ContractAddress, &a.AnchorStatus,
		&a.ErrorMessage, &a.CreatedAt, &a.ConfirmedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return &a, err
}

func (s *Store) FindAnchorByTxHash(ctx context.Context, txHash string) (*ContentAnchor, error) {
	var a ContentAnchor
	err := s.pool.QueryRow(ctx,
		`SELECT id, content_id, user_id, content_hash, tx_hash, chain, block_number, contract_address, anchor_status, error_message, created_at, confirmed_at
		 FROM content_anchors WHERE tx_hash = $1`, txHash,
	).Scan(
		&a.ID, &a.ContentID, &a.UserID, &a.ContentHash, &a.TxHash,
		&a.Chain, &a.BlockNumber, &a.ContractAddress, &a.AnchorStatus,
		&a.ErrorMessage, &a.CreatedAt, &a.ConfirmedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return &a, err
}

func (s *Store) UpdateAnchorStatus(ctx context.Context, id, status string, errorMsg string, blockNumber int64, confirmedAt *time.Time) error {
	var errPtr *string
	if errorMsg != "" {
		errPtr = &errorMsg
	}
	var blkPtr *int64
	if blockNumber > 0 {
		blkPtr = &blockNumber
	}
	_, err := s.pool.Exec(ctx,
		`UPDATE content_anchors
		 SET anchor_status = $1, block_number = $2, confirmed_at = $3, error_message = $4
		 WHERE id = $5`,
		status, blkPtr, confirmedAt, errPtr, id,
	)
	return err
}

func (s *Store) ListPendingAnchors(ctx context.Context) ([]*ContentAnchor, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT id, content_id, user_id, content_hash, tx_hash, chain, block_number, contract_address, anchor_status, error_message, created_at, confirmed_at
		 FROM content_anchors
		 WHERE anchor_status = 'pending'
		 ORDER BY created_at ASC
		 LIMIT 50`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var anchors []*ContentAnchor
	for rows.Next() {
		var a ContentAnchor
		if err := rows.Scan(
			&a.ID, &a.ContentID, &a.UserID, &a.ContentHash, &a.TxHash,
			&a.Chain, &a.BlockNumber, &a.ContractAddress, &a.AnchorStatus,
			&a.ErrorMessage, &a.CreatedAt, &a.ConfirmedAt,
		); err != nil {
			return nil, err
		}
		anchors = append(anchors, &a)
	}
	return anchors, nil
}

func (s *Store) ListAnchorsByUser(ctx context.Context, userID string, limit, offset int) ([]*ContentAnchor, int, error) {
	var total int
	err := s.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM content_anchors WHERE user_id = $1`, userID,
	).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := s.pool.Query(ctx,
		`SELECT id, content_id, user_id, content_hash, tx_hash, chain, block_number, contract_address, anchor_status, error_message, created_at, confirmed_at
		 FROM content_anchors
		 WHERE user_id = $1
		 ORDER BY created_at DESC
		 LIMIT $2 OFFSET $3`,
		userID, limit, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var anchors []*ContentAnchor
	for rows.Next() {
		var a ContentAnchor
		if err := rows.Scan(
			&a.ID, &a.ContentID, &a.UserID, &a.ContentHash, &a.TxHash,
			&a.Chain, &a.BlockNumber, &a.ContractAddress, &a.AnchorStatus,
			&a.ErrorMessage, &a.CreatedAt, &a.ConfirmedAt,
		); err != nil {
			return nil, 0, err
		}
		anchors = append(anchors, &a)
	}
	if anchors == nil {
		anchors = []*ContentAnchor{}
	}
	return anchors, total, nil
}
