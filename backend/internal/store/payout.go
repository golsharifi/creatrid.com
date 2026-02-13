package store

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

type CreatorPayout struct {
	ID               string     `json:"id"`
	UserID           string     `json:"userId"`
	PurchaseID       string     `json:"purchaseId"`
	StripeTransferID *string    `json:"stripeTransferId"`
	AmountCents      int        `json:"amountCents"`
	Currency         string     `json:"currency"`
	Status           string     `json:"status"`
	ErrorMessage     *string    `json:"errorMessage"`
	CreatedAt        time.Time  `json:"createdAt"`
	CompletedAt      *time.Time `json:"completedAt"`
}

type PayoutDashboard struct {
	TotalEarnedCents int `json:"totalEarnedCents"`
	TotalPaidCents   int `json:"totalPaidCents"`
	PendingCents     int `json:"pendingCents"`
}

func (s *Store) CreatePayout(ctx context.Context, payout *CreatorPayout) error {
	_, err := s.pool.Exec(ctx,
		`INSERT INTO creator_payouts (id, user_id, purchase_id, stripe_transfer_id, amount_cents, currency, status, error_message, created_at, completed_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		payout.ID, payout.UserID, payout.PurchaseID, payout.StripeTransferID,
		payout.AmountCents, payout.Currency, payout.Status, payout.ErrorMessage,
		payout.CreatedAt, payout.CompletedAt,
	)
	return err
}

func (s *Store) UpdatePayoutStatus(ctx context.Context, id, status string, transferID *string, errorMsg *string) error {
	if status == "completed" {
		_, err := s.pool.Exec(ctx,
			`UPDATE creator_payouts SET status = $1, stripe_transfer_id = $2, error_message = $3, completed_at = NOW() WHERE id = $4`,
			status, transferID, errorMsg, id,
		)
		return err
	}
	_, err := s.pool.Exec(ctx,
		`UPDATE creator_payouts SET status = $1, stripe_transfer_id = $2, error_message = $3 WHERE id = $4`,
		status, transferID, errorMsg, id,
	)
	return err
}

func (s *Store) ListPayoutsByUser(ctx context.Context, userID string, limit, offset int) ([]*CreatorPayout, int, error) {
	var total int
	err := s.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM creator_payouts WHERE user_id = $1`, userID,
	).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := s.pool.Query(ctx,
		`SELECT id, user_id, purchase_id, stripe_transfer_id, amount_cents, currency, status, error_message, created_at, completed_at
		 FROM creator_payouts
		 WHERE user_id = $1
		 ORDER BY created_at DESC
		 LIMIT $2 OFFSET $3`,
		userID, limit, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var payouts []*CreatorPayout
	for rows.Next() {
		var p CreatorPayout
		if err := rows.Scan(&p.ID, &p.UserID, &p.PurchaseID, &p.StripeTransferID, &p.AmountCents, &p.Currency, &p.Status, &p.ErrorMessage, &p.CreatedAt, &p.CompletedAt); err != nil {
			return nil, 0, err
		}
		payouts = append(payouts, &p)
	}
	if payouts == nil {
		payouts = []*CreatorPayout{}
	}
	return payouts, total, nil
}

func (s *Store) GetPayoutDashboard(ctx context.Context, userID string) (*PayoutDashboard, error) {
	dashboard := &PayoutDashboard{}

	_ = s.pool.QueryRow(ctx,
		`SELECT COALESCE(SUM(lp.creator_payout_cents), 0)
		 FROM license_purchases lp
		 JOIN content_items ci ON ci.id = lp.content_id
		 WHERE ci.user_id = $1 AND lp.status = 'completed'`,
		userID,
	).Scan(&dashboard.TotalEarnedCents)

	_ = s.pool.QueryRow(ctx,
		`SELECT COALESCE(SUM(amount_cents), 0) FROM creator_payouts WHERE user_id = $1 AND status = 'completed'`,
		userID,
	).Scan(&dashboard.TotalPaidCents)

	_ = s.pool.QueryRow(ctx,
		`SELECT COALESCE(SUM(amount_cents), 0) FROM creator_payouts WHERE user_id = $1 AND status = 'pending'`,
		userID,
	).Scan(&dashboard.PendingCents)

	return dashboard, nil
}

func (s *Store) UpdateUserStripeConnect(ctx context.Context, userID, accountID string, onboarded bool) error {
	_, err := s.pool.Exec(ctx,
		`UPDATE users SET stripe_connect_account_id = $1, stripe_connect_onboarded = $2 WHERE id = $3`,
		accountID, onboarded, userID,
	)
	return err
}

func (s *Store) GetUserStripeConnectID(ctx context.Context, userID string) (*string, bool, error) {
	var accountID *string
	var onboarded bool
	err := s.pool.QueryRow(ctx,
		`SELECT stripe_connect_account_id, COALESCE(stripe_connect_onboarded, false) FROM users WHERE id = $1`,
		userID,
	).Scan(&accountID, &onboarded)
	if err == pgx.ErrNoRows {
		return nil, false, nil
	}
	return accountID, onboarded, err
}
