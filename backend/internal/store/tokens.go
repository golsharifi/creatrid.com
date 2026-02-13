package store

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/nrednav/cuid2"
)

// CreatorToken represents a creator's social token.
type CreatorToken struct {
	ID          string    `json:"id"`
	UserID      string    `json:"userId"`
	Name        string    `json:"name"`
	Symbol      string    `json:"symbol"`
	Description *string   `json:"description"`
	TotalSupply int       `json:"totalSupply"`
	PriceCents  int       `json:"priceCents"`
	IsActive    bool      `json:"isActive"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// TokenBalance represents a user's balance of a creator token.
type TokenBalance struct {
	ID        string    `json:"id"`
	TokenID   string    `json:"tokenId"`
	UserID    string    `json:"userId"`
	Balance   int       `json:"balance"`
	UpdatedAt time.Time `json:"updatedAt"`
	// Display fields
	UserName     *string `json:"userName,omitempty"`
	UserUsername *string `json:"userUsername,omitempty"`
	UserImage    *string `json:"userImage,omitempty"`
}

// Tip represents a one-time tip from one user to another.
type Tip struct {
	ID              string    `json:"id"`
	FromUserID      string    `json:"fromUserId"`
	ToUserID        string    `json:"toUserId"`
	AmountCents     int       `json:"amountCents"`
	Message         *string   `json:"message"`
	StripePaymentID *string   `json:"-"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"createdAt"`
	// Display fields
	FromName     *string `json:"fromName,omitempty"`
	FromUsername *string `json:"fromUsername,omitempty"`
	FromImage    *string `json:"fromImage,omitempty"`
	ToName       *string `json:"toName,omitempty"`
	ToUsername   *string `json:"toUsername,omitempty"`
}

// FanSubscription represents a recurring subscription from a fan to a creator.
type FanSubscription struct {
	ID                   string     `json:"id"`
	FanUserID            string     `json:"fanUserId"`
	CreatorUserID        string     `json:"creatorUserId"`
	Tier                 string     `json:"tier"`
	PriceCents           int        `json:"priceCents"`
	StripeSubscriptionID *string    `json:"-"`
	Status               string     `json:"status"`
	StartedAt            time.Time  `json:"startedAt"`
	CanceledAt           *time.Time `json:"canceledAt"`
	// Display fields
	CreatorName     *string `json:"creatorName,omitempty"`
	CreatorUsername *string `json:"creatorUsername,omitempty"`
	FanName         *string `json:"fanName,omitempty"`
	FanUsername     *string `json:"fanUsername,omitempty"`
}

// GatedContent represents token/subscription gating on a content item.
type GatedContent struct {
	ID               string    `json:"id"`
	ContentID        string    `json:"contentId"`
	TokenID          *string   `json:"tokenId"`
	MinTokens        int       `json:"minTokens"`
	SubscriptionTier *string   `json:"subscriptionTier"`
	CreatedAt        time.Time `json:"createdAt"`
}

// TokenTransaction represents a transaction in the token ledger.
type TokenTransaction struct {
	ID          string    `json:"id"`
	TokenID     string    `json:"tokenId"`
	FromUserID  *string   `json:"fromUserId"`
	ToUserID    *string   `json:"toUserId"`
	Amount      int       `json:"amount"`
	TxType      string    `json:"txType"`
	ReferenceID *string   `json:"referenceId"`
	CreatedAt   time.Time `json:"createdAt"`
}

// TipStats holds aggregate tip statistics for a user.
type TipStats struct {
	TotalReceived int `json:"totalReceived"`
	TotalSent     int `json:"totalSent"`
	TipCount      int `json:"tipCount"`
}

// --- Creator Tokens ---

func (s *Store) CreateCreatorToken(ctx context.Context, token *CreatorToken) error {
	_, err := s.pool.Exec(ctx,
		`INSERT INTO creator_tokens (id, user_id, name, symbol, description, total_supply, price_cents, is_active, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		token.ID, token.UserID, token.Name, token.Symbol, token.Description,
		token.TotalSupply, token.PriceCents, token.IsActive, token.CreatedAt, token.UpdatedAt,
	)
	return err
}

func (s *Store) FindTokenByUserID(ctx context.Context, userID string) (*CreatorToken, error) {
	var t CreatorToken
	err := s.pool.QueryRow(ctx,
		`SELECT id, user_id, name, symbol, description, total_supply, price_cents, is_active, created_at, updated_at
		 FROM creator_tokens WHERE user_id = $1`, userID,
	).Scan(&t.ID, &t.UserID, &t.Name, &t.Symbol, &t.Description,
		&t.TotalSupply, &t.PriceCents, &t.IsActive, &t.CreatedAt, &t.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return &t, err
}

func (s *Store) FindTokenByID(ctx context.Context, id string) (*CreatorToken, error) {
	var t CreatorToken
	err := s.pool.QueryRow(ctx,
		`SELECT id, user_id, name, symbol, description, total_supply, price_cents, is_active, created_at, updated_at
		 FROM creator_tokens WHERE id = $1`, id,
	).Scan(&t.ID, &t.UserID, &t.Name, &t.Symbol, &t.Description,
		&t.TotalSupply, &t.PriceCents, &t.IsActive, &t.CreatedAt, &t.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return &t, err
}

func (s *Store) FindTokenBySymbol(ctx context.Context, symbol string) (*CreatorToken, error) {
	var t CreatorToken
	err := s.pool.QueryRow(ctx,
		`SELECT id, user_id, name, symbol, description, total_supply, price_cents, is_active, created_at, updated_at
		 FROM creator_tokens WHERE symbol = $1`, symbol,
	).Scan(&t.ID, &t.UserID, &t.Name, &t.Symbol, &t.Description,
		&t.TotalSupply, &t.PriceCents, &t.IsActive, &t.CreatedAt, &t.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return &t, err
}

func (s *Store) UpdateCreatorToken(ctx context.Context, id, name string, description *string, priceCents int, isActive bool) error {
	_, err := s.pool.Exec(ctx,
		`UPDATE creator_tokens SET name = $2, description = $3, price_cents = $4, is_active = $5, updated_at = NOW()
		 WHERE id = $1`,
		id, name, description, priceCents, isActive,
	)
	return err
}

// --- Token Balances ---

func (s *Store) GetTokenBalance(ctx context.Context, tokenID, userID string) (int, error) {
	var balance int
	err := s.pool.QueryRow(ctx,
		`SELECT COALESCE(balance, 0) FROM token_balances WHERE token_id = $1 AND user_id = $2`,
		tokenID, userID,
	).Scan(&balance)
	if err == pgx.ErrNoRows {
		return 0, nil
	}
	return balance, err
}

func (s *Store) MintTokens(ctx context.Context, tokenID, toUserID string, amount int, txType, referenceID string) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Upsert balance
	_, err = tx.Exec(ctx,
		`INSERT INTO token_balances (id, token_id, user_id, balance, updated_at)
		 VALUES ($1, $2, $3, $4, NOW())
		 ON CONFLICT (token_id, user_id) DO UPDATE SET balance = token_balances.balance + $4, updated_at = NOW()`,
		cuid2.Generate(), tokenID, toUserID, amount,
	)
	if err != nil {
		return err
	}

	// Update total supply
	_, err = tx.Exec(ctx,
		`UPDATE creator_tokens SET total_supply = total_supply + $2, updated_at = NOW() WHERE id = $1`,
		tokenID, amount,
	)
	if err != nil {
		return err
	}

	// Record transaction
	var refID *string
	if referenceID != "" {
		refID = &referenceID
	}
	_, err = tx.Exec(ctx,
		`INSERT INTO token_transactions (id, token_id, from_user_id, to_user_id, amount, tx_type, reference_id, created_at)
		 VALUES ($1, $2, NULL, $3, $4, $5, $6, NOW())`,
		cuid2.Generate(), tokenID, toUserID, amount, txType, refID,
	)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (s *Store) TransferTokens(ctx context.Context, tokenID, fromUserID, toUserID string, amount int) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Check balance
	var balance int
	err = tx.QueryRow(ctx,
		`SELECT balance FROM token_balances WHERE token_id = $1 AND user_id = $2 FOR UPDATE`,
		tokenID, fromUserID,
	).Scan(&balance)
	if err != nil {
		return err
	}
	if balance < amount {
		return pgx.ErrNoRows // insufficient balance
	}

	// Deduct from sender
	_, err = tx.Exec(ctx,
		`UPDATE token_balances SET balance = balance - $3, updated_at = NOW()
		 WHERE token_id = $1 AND user_id = $2`,
		tokenID, fromUserID, amount,
	)
	if err != nil {
		return err
	}

	// Credit to receiver
	_, err = tx.Exec(ctx,
		`INSERT INTO token_balances (id, token_id, user_id, balance, updated_at)
		 VALUES ($1, $2, $3, $4, NOW())
		 ON CONFLICT (token_id, user_id) DO UPDATE SET balance = token_balances.balance + $4, updated_at = NOW()`,
		cuid2.Generate(), tokenID, toUserID, amount,
	)
	if err != nil {
		return err
	}

	// Record transaction
	_, err = tx.Exec(ctx,
		`INSERT INTO token_transactions (id, token_id, from_user_id, to_user_id, amount, tx_type, reference_id, created_at)
		 VALUES ($1, $2, $3, $4, $5, 'transfer', NULL, NOW())`,
		cuid2.Generate(), tokenID, fromUserID, toUserID, amount,
	)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (s *Store) ListTokenHolders(ctx context.Context, tokenID string, limit, offset int) ([]TokenBalance, int, error) {
	var total int
	err := s.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM token_balances WHERE token_id = $1 AND balance > 0`, tokenID,
	).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := s.pool.Query(ctx,
		`SELECT tb.id, tb.token_id, tb.user_id, tb.balance, tb.updated_at, u.name, u.username, u.image
		 FROM token_balances tb
		 JOIN users u ON u.id = tb.user_id
		 WHERE tb.token_id = $1 AND tb.balance > 0
		 ORDER BY tb.balance DESC
		 LIMIT $2 OFFSET $3`, tokenID, limit, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var holders []TokenBalance
	for rows.Next() {
		var h TokenBalance
		if err := rows.Scan(&h.ID, &h.TokenID, &h.UserID, &h.Balance, &h.UpdatedAt,
			&h.UserName, &h.UserUsername, &h.UserImage); err != nil {
			return nil, 0, err
		}
		holders = append(holders, h)
	}
	return holders, total, nil
}

// --- Tips ---

func (s *Store) CreateTip(ctx context.Context, tip *Tip) error {
	_, err := s.pool.Exec(ctx,
		`INSERT INTO tips (id, from_user_id, to_user_id, amount_cents, message, stripe_payment_id, status, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		tip.ID, tip.FromUserID, tip.ToUserID, tip.AmountCents, tip.Message,
		tip.StripePaymentID, tip.Status, tip.CreatedAt,
	)
	return err
}

func (s *Store) UpdateTipStatus(ctx context.Context, id, status string, stripePaymentID *string) error {
	_, err := s.pool.Exec(ctx,
		`UPDATE tips SET status = $2, stripe_payment_id = COALESCE($3, stripe_payment_id) WHERE id = $1`,
		id, status, stripePaymentID,
	)
	return err
}

func (s *Store) ListTipsReceived(ctx context.Context, userID string, limit, offset int) ([]Tip, int, error) {
	var total int
	err := s.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM tips WHERE to_user_id = $1 AND status = 'completed'`, userID,
	).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := s.pool.Query(ctx,
		`SELECT t.id, t.from_user_id, t.to_user_id, t.amount_cents, t.message, t.status, t.created_at,
		        fu.name, fu.username, fu.image
		 FROM tips t
		 JOIN users fu ON fu.id = t.from_user_id
		 WHERE t.to_user_id = $1 AND t.status = 'completed'
		 ORDER BY t.created_at DESC
		 LIMIT $2 OFFSET $3`, userID, limit, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var tips []Tip
	for rows.Next() {
		var tip Tip
		if err := rows.Scan(&tip.ID, &tip.FromUserID, &tip.ToUserID, &tip.AmountCents,
			&tip.Message, &tip.Status, &tip.CreatedAt,
			&tip.FromName, &tip.FromUsername, &tip.FromImage); err != nil {
			return nil, 0, err
		}
		tips = append(tips, tip)
	}
	return tips, total, nil
}

func (s *Store) ListTipsSent(ctx context.Context, userID string, limit, offset int) ([]Tip, int, error) {
	var total int
	err := s.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM tips WHERE from_user_id = $1 AND status = 'completed'`, userID,
	).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := s.pool.Query(ctx,
		`SELECT t.id, t.from_user_id, t.to_user_id, t.amount_cents, t.message, t.status, t.created_at,
		        tu.name, tu.username
		 FROM tips t
		 JOIN users tu ON tu.id = t.to_user_id
		 WHERE t.from_user_id = $1 AND t.status = 'completed'
		 ORDER BY t.created_at DESC
		 LIMIT $2 OFFSET $3`, userID, limit, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var tips []Tip
	for rows.Next() {
		var tip Tip
		if err := rows.Scan(&tip.ID, &tip.FromUserID, &tip.ToUserID, &tip.AmountCents,
			&tip.Message, &tip.Status, &tip.CreatedAt,
			&tip.ToName, &tip.ToUsername); err != nil {
			return nil, 0, err
		}
		tips = append(tips, tip)
	}
	return tips, total, nil
}

func (s *Store) GetTipStats(ctx context.Context, userID string) (*TipStats, error) {
	stats := &TipStats{}

	// Total received
	err := s.pool.QueryRow(ctx,
		`SELECT COALESCE(SUM(amount_cents), 0) FROM tips WHERE to_user_id = $1 AND status = 'completed'`,
		userID,
	).Scan(&stats.TotalReceived)
	if err != nil {
		return nil, err
	}

	// Total sent
	err = s.pool.QueryRow(ctx,
		`SELECT COALESCE(SUM(amount_cents), 0) FROM tips WHERE from_user_id = $1 AND status = 'completed'`,
		userID,
	).Scan(&stats.TotalSent)
	if err != nil {
		return nil, err
	}

	// Tip count (received)
	err = s.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM tips WHERE to_user_id = $1 AND status = 'completed'`,
		userID,
	).Scan(&stats.TipCount)
	if err != nil {
		return nil, err
	}

	return stats, nil
}

// --- Fan Subscriptions ---

func (s *Store) CreateFanSubscription(ctx context.Context, sub *FanSubscription) error {
	_, err := s.pool.Exec(ctx,
		`INSERT INTO fan_subscriptions (id, fan_user_id, creator_user_id, tier, price_cents, stripe_subscription_id, status, started_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		sub.ID, sub.FanUserID, sub.CreatorUserID, sub.Tier, sub.PriceCents,
		sub.StripeSubscriptionID, sub.Status, sub.StartedAt,
	)
	return err
}

func (s *Store) FindFanSubscription(ctx context.Context, fanUserID, creatorUserID string) (*FanSubscription, error) {
	var sub FanSubscription
	err := s.pool.QueryRow(ctx,
		`SELECT id, fan_user_id, creator_user_id, tier, price_cents, stripe_subscription_id, status, started_at, canceled_at
		 FROM fan_subscriptions WHERE fan_user_id = $1 AND creator_user_id = $2`,
		fanUserID, creatorUserID,
	).Scan(&sub.ID, &sub.FanUserID, &sub.CreatorUserID, &sub.Tier, &sub.PriceCents,
		&sub.StripeSubscriptionID, &sub.Status, &sub.StartedAt, &sub.CanceledAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return &sub, err
}

func (s *Store) UpdateFanSubscriptionStatus(ctx context.Context, id, status string) error {
	q := `UPDATE fan_subscriptions SET status = $2 WHERE id = $1`
	if status == "canceled" {
		q = `UPDATE fan_subscriptions SET status = $2, canceled_at = NOW() WHERE id = $1`
	}
	_, err := s.pool.Exec(ctx, q, id, status)
	return err
}

func (s *Store) ListFansByCreator(ctx context.Context, creatorUserID string, limit, offset int) ([]FanSubscription, int, error) {
	var total int
	err := s.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM fan_subscriptions WHERE creator_user_id = $1 AND status = 'active'`, creatorUserID,
	).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := s.pool.Query(ctx,
		`SELECT fs.id, fs.fan_user_id, fs.creator_user_id, fs.tier, fs.price_cents, fs.status, fs.started_at, fs.canceled_at,
		        u.name, u.username
		 FROM fan_subscriptions fs
		 JOIN users u ON u.id = fs.fan_user_id
		 WHERE fs.creator_user_id = $1 AND fs.status = 'active'
		 ORDER BY fs.started_at DESC
		 LIMIT $2 OFFSET $3`, creatorUserID, limit, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var subs []FanSubscription
	for rows.Next() {
		var sub FanSubscription
		if err := rows.Scan(&sub.ID, &sub.FanUserID, &sub.CreatorUserID, &sub.Tier, &sub.PriceCents,
			&sub.Status, &sub.StartedAt, &sub.CanceledAt,
			&sub.FanName, &sub.FanUsername); err != nil {
			return nil, 0, err
		}
		subs = append(subs, sub)
	}
	return subs, total, nil
}

func (s *Store) ListSubscriptionsByFan(ctx context.Context, fanUserID string) ([]FanSubscription, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT fs.id, fs.fan_user_id, fs.creator_user_id, fs.tier, fs.price_cents, fs.status, fs.started_at, fs.canceled_at,
		        u.name, u.username
		 FROM fan_subscriptions fs
		 JOIN users u ON u.id = fs.creator_user_id
		 WHERE fs.fan_user_id = $1
		 ORDER BY fs.started_at DESC`, fanUserID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subs []FanSubscription
	for rows.Next() {
		var sub FanSubscription
		if err := rows.Scan(&sub.ID, &sub.FanUserID, &sub.CreatorUserID, &sub.Tier, &sub.PriceCents,
			&sub.Status, &sub.StartedAt, &sub.CanceledAt,
			&sub.CreatorName, &sub.CreatorUsername); err != nil {
			return nil, err
		}
		subs = append(subs, sub)
	}
	return subs, nil
}

// --- Gated Content ---

func (s *Store) SetGatedContent(ctx context.Context, gated *GatedContent) error {
	_, err := s.pool.Exec(ctx,
		`INSERT INTO gated_content (id, content_id, token_id, min_tokens, subscription_tier, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 ON CONFLICT (content_id) DO UPDATE SET
		   token_id = EXCLUDED.token_id,
		   min_tokens = EXCLUDED.min_tokens,
		   subscription_tier = EXCLUDED.subscription_tier`,
		gated.ID, gated.ContentID, gated.TokenID, gated.MinTokens, gated.SubscriptionTier, gated.CreatedAt,
	)
	return err
}

func (s *Store) FindGatedContent(ctx context.Context, contentID string) (*GatedContent, error) {
	var g GatedContent
	err := s.pool.QueryRow(ctx,
		`SELECT id, content_id, token_id, min_tokens, subscription_tier, created_at
		 FROM gated_content WHERE content_id = $1`, contentID,
	).Scan(&g.ID, &g.ContentID, &g.TokenID, &g.MinTokens, &g.SubscriptionTier, &g.CreatedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return &g, err
}

func (s *Store) RemoveGatedContent(ctx context.Context, contentID string) error {
	_, err := s.pool.Exec(ctx, `DELETE FROM gated_content WHERE content_id = $1`, contentID)
	return err
}

func (s *Store) CheckContentAccess(ctx context.Context, contentID, userID string) (bool, error) {
	// Find gated content
	gated, err := s.FindGatedContent(ctx, contentID)
	if err != nil {
		return false, err
	}
	if gated == nil {
		// Not gated — everyone has access
		return true, nil
	}

	// Check token balance if token gating is set
	if gated.TokenID != nil {
		balance, err := s.GetTokenBalance(ctx, *gated.TokenID, userID)
		if err != nil {
			return false, err
		}
		if balance >= gated.MinTokens {
			return true, nil
		}
	}

	// Check subscription tier if set
	if gated.SubscriptionTier != nil {
		// Find the content owner
		var ownerID string
		err := s.pool.QueryRow(ctx,
			`SELECT user_id FROM content_items WHERE id = $1`, contentID,
		).Scan(&ownerID)
		if err != nil {
			return false, err
		}

		sub, err := s.FindFanSubscription(ctx, userID, ownerID)
		if err != nil {
			return false, err
		}
		if sub != nil && sub.Status == "active" {
			// Check tier hierarchy: patron > superfan > supporter
			tierRank := map[string]int{"supporter": 1, "superfan": 2, "patron": 3}
			requiredRank := tierRank[*gated.SubscriptionTier]
			userRank := tierRank[sub.Tier]
			if userRank >= requiredRank {
				return true, nil
			}
		}
	}

	return false, nil
}

// --- Token Transactions ---

func (s *Store) RecordTokenTransaction(ctx context.Context, tx *TokenTransaction) error {
	_, err := s.pool.Exec(ctx,
		`INSERT INTO token_transactions (id, token_id, from_user_id, to_user_id, amount, tx_type, reference_id, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		tx.ID, tx.TokenID, tx.FromUserID, tx.ToUserID, tx.Amount, tx.TxType, tx.ReferenceID, tx.CreatedAt,
	)
	return err
}

func (s *Store) ListTokenTransactions(ctx context.Context, tokenID string, limit, offset int) ([]TokenTransaction, int, error) {
	var total int
	err := s.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM token_transactions WHERE token_id = $1`, tokenID,
	).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := s.pool.Query(ctx,
		`SELECT id, token_id, from_user_id, to_user_id, amount, tx_type, reference_id, created_at
		 FROM token_transactions
		 WHERE token_id = $1
		 ORDER BY created_at DESC
		 LIMIT $2 OFFSET $3`, tokenID, limit, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var txs []TokenTransaction
	for rows.Next() {
		var t TokenTransaction
		if err := rows.Scan(&t.ID, &t.TokenID, &t.FromUserID, &t.ToUserID, &t.Amount,
			&t.TxType, &t.ReferenceID, &t.CreatedAt); err != nil {
			return nil, 0, err
		}
		txs = append(txs, t)
	}
	return txs, total, nil
}
