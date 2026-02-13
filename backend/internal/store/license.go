package store

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

type LicenseOffering struct {
	ID          string    `json:"id"`
	ContentID   string    `json:"contentId"`
	LicenseType string    `json:"licenseType"`
	PriceCents  int       `json:"priceCents"`
	Currency    string    `json:"currency"`
	IsActive    bool      `json:"isActive"`
	TermsText   *string   `json:"termsText"`
	CreatedAt   time.Time `json:"createdAt"`
}

type LicensePurchase struct {
	ID                 string    `json:"id"`
	OfferingID         string    `json:"offeringId"`
	ContentID          string    `json:"contentId"`
	BuyerUserID        *string   `json:"buyerUserId"`
	BuyerEmail         string    `json:"buyerEmail"`
	BuyerCompany       *string   `json:"buyerCompany"`
	StripeSessionID    *string   `json:"stripeSessionId"`
	AmountCents        int       `json:"amountCents"`
	PlatformFeeCents   int       `json:"platformFeeCents"`
	CreatorPayoutCents int       `json:"creatorPayoutCents"`
	Status             string    `json:"status"`
	CreatedAt          time.Time `json:"createdAt"`
}

type TakedownRequest struct {
	ID              string     `json:"id"`
	ReporterEmail   string     `json:"reporterEmail"`
	ReporterName    string     `json:"reporterName"`
	ContentID       string     `json:"contentId"`
	Reason          string     `json:"reason"`
	EvidenceURL     *string    `json:"evidenceUrl"`
	Status          string     `json:"status"`
	ResolvedBy      *string    `json:"resolvedBy"`
	ResolutionNotes *string    `json:"resolutionNotes"`
	CreatedAt       time.Time  `json:"createdAt"`
	ResolvedAt      *time.Time `json:"resolvedAt"`
}

type MarketplaceItem struct {
	ContentItem
	CreatorName      *string `json:"creatorName"`
	CreatorUsername   *string `json:"creatorUsername"`
	CreatorImage     *string `json:"creatorImage"`
	LowestPriceCents int     `json:"lowestPriceCents"`
}

// --- License Offerings ---

func (s *Store) CreateLicenseOffering(ctx context.Context, offering *LicenseOffering) error {
	_, err := s.pool.Exec(ctx,
		`INSERT INTO license_offerings (id, content_id, license_type, price_cents, currency, is_active, terms_text, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		offering.ID, offering.ContentID, offering.LicenseType, offering.PriceCents,
		offering.Currency, offering.IsActive, offering.TermsText, offering.CreatedAt,
	)
	return err
}

func (s *Store) ListOfferingsByContent(ctx context.Context, contentID string) ([]*LicenseOffering, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT id, content_id, license_type, price_cents, currency, is_active, terms_text, created_at
		 FROM license_offerings
		 WHERE content_id = $1 AND is_active = true
		 ORDER BY price_cents ASC`, contentID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var offerings []*LicenseOffering
	for rows.Next() {
		var o LicenseOffering
		if err := rows.Scan(&o.ID, &o.ContentID, &o.LicenseType, &o.PriceCents, &o.Currency, &o.IsActive, &o.TermsText, &o.CreatedAt); err != nil {
			return nil, err
		}
		offerings = append(offerings, &o)
	}
	if offerings == nil {
		offerings = []*LicenseOffering{}
	}
	return offerings, nil
}

func (s *Store) FindOfferingByID(ctx context.Context, id string) (*LicenseOffering, error) {
	var o LicenseOffering
	err := s.pool.QueryRow(ctx,
		`SELECT id, content_id, license_type, price_cents, currency, is_active, terms_text, created_at
		 FROM license_offerings WHERE id = $1`, id,
	).Scan(&o.ID, &o.ContentID, &o.LicenseType, &o.PriceCents, &o.Currency, &o.IsActive, &o.TermsText, &o.CreatedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return &o, err
}

func (s *Store) UpdateOffering(ctx context.Context, id string, priceCents int, isActive bool, terms *string) error {
	_, err := s.pool.Exec(ctx,
		`UPDATE license_offerings SET price_cents = $1, is_active = $2, terms_text = $3 WHERE id = $4`,
		priceCents, isActive, terms, id,
	)
	return err
}

func (s *Store) DeleteOffering(ctx context.Context, id string) error {
	_, err := s.pool.Exec(ctx, `DELETE FROM license_offerings WHERE id = $1`, id)
	return err
}

// --- License Purchases ---

func (s *Store) CreateLicensePurchase(ctx context.Context, purchase *LicensePurchase) error {
	_, err := s.pool.Exec(ctx,
		`INSERT INTO license_purchases (id, offering_id, content_id, buyer_user_id, buyer_email, buyer_company, stripe_session_id, amount_cents, platform_fee_cents, creator_payout_cents, status, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
		purchase.ID, purchase.OfferingID, purchase.ContentID, purchase.BuyerUserID,
		purchase.BuyerEmail, purchase.BuyerCompany, purchase.StripeSessionID,
		purchase.AmountCents, purchase.PlatformFeeCents, purchase.CreatorPayoutCents,
		purchase.Status, purchase.CreatedAt,
	)
	return err
}

func (s *Store) ListPurchasesByBuyer(ctx context.Context, userID string) ([]*LicensePurchase, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT id, offering_id, content_id, buyer_user_id, buyer_email, buyer_company, stripe_session_id, amount_cents, platform_fee_cents, creator_payout_cents, status, created_at
		 FROM license_purchases
		 WHERE buyer_user_id = $1
		 ORDER BY created_at DESC`, userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var purchases []*LicensePurchase
	for rows.Next() {
		var p LicensePurchase
		if err := rows.Scan(&p.ID, &p.OfferingID, &p.ContentID, &p.BuyerUserID, &p.BuyerEmail, &p.BuyerCompany, &p.StripeSessionID, &p.AmountCents, &p.PlatformFeeCents, &p.CreatorPayoutCents, &p.Status, &p.CreatedAt); err != nil {
			return nil, err
		}
		purchases = append(purchases, &p)
	}
	if purchases == nil {
		purchases = []*LicensePurchase{}
	}
	return purchases, nil
}

func (s *Store) ListSalesByCreator(ctx context.Context, userID string) ([]*LicensePurchase, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT lp.id, lp.offering_id, lp.content_id, lp.buyer_user_id, lp.buyer_email, lp.buyer_company, lp.stripe_session_id, lp.amount_cents, lp.platform_fee_cents, lp.creator_payout_cents, lp.status, lp.created_at
		 FROM license_purchases lp
		 JOIN content_items ci ON ci.id = lp.content_id
		 WHERE ci.user_id = $1
		 ORDER BY lp.created_at DESC`, userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sales []*LicensePurchase
	for rows.Next() {
		var p LicensePurchase
		if err := rows.Scan(&p.ID, &p.OfferingID, &p.ContentID, &p.BuyerUserID, &p.BuyerEmail, &p.BuyerCompany, &p.StripeSessionID, &p.AmountCents, &p.PlatformFeeCents, &p.CreatorPayoutCents, &p.Status, &p.CreatedAt); err != nil {
			return nil, err
		}
		sales = append(sales, &p)
	}
	if sales == nil {
		sales = []*LicensePurchase{}
	}
	return sales, nil
}

func (s *Store) HasLicense(ctx context.Context, userID, contentID string) (bool, error) {
	var exists bool
	err := s.pool.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM license_purchases WHERE buyer_user_id = $1 AND content_id = $2 AND status = 'completed')`,
		userID, contentID,
	).Scan(&exists)
	return exists, err
}

// --- Takedown Requests ---

func (s *Store) CreateTakedownRequest(ctx context.Context, req *TakedownRequest) error {
	_, err := s.pool.Exec(ctx,
		`INSERT INTO takedown_requests (id, reporter_email, reporter_name, content_id, reason, evidence_url, status, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		req.ID, req.ReporterEmail, req.ReporterName, req.ContentID,
		req.Reason, req.EvidenceURL, req.Status, req.CreatedAt,
	)
	return err
}

func (s *Store) ListTakedownRequests(ctx context.Context, status string, limit, offset int) ([]*TakedownRequest, int, error) {
	baseWhere := ""
	args := []interface{}{}
	argIdx := 1

	if status != "" {
		baseWhere = ` WHERE status = $` + fmt.Sprintf("%d", argIdx)
		args = append(args, status)
		argIdx++
	}

	var total int
	countQuery := `SELECT COUNT(*) FROM takedown_requests` + baseWhere
	err := s.pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	selectQuery := `SELECT id, reporter_email, reporter_name, content_id, reason, evidence_url, status, resolved_by, resolution_notes, created_at, resolved_at
		 FROM takedown_requests` + baseWhere +
		` ORDER BY created_at DESC LIMIT $` + fmt.Sprintf("%d", argIdx) + ` OFFSET $` + fmt.Sprintf("%d", argIdx+1)
	args = append(args, limit, offset)

	rows, err := s.pool.Query(ctx, selectQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var requests []*TakedownRequest
	for rows.Next() {
		var t TakedownRequest
		if err := rows.Scan(&t.ID, &t.ReporterEmail, &t.ReporterName, &t.ContentID, &t.Reason, &t.EvidenceURL, &t.Status, &t.ResolvedBy, &t.ResolutionNotes, &t.CreatedAt, &t.ResolvedAt); err != nil {
			return nil, 0, err
		}
		requests = append(requests, &t)
	}
	if requests == nil {
		requests = []*TakedownRequest{}
	}
	return requests, total, nil
}

func (s *Store) ResolveTakedown(ctx context.Context, id, resolvedBy, status, notes string) error {
	_, err := s.pool.Exec(ctx,
		`UPDATE takedown_requests SET status = $1, resolved_by = $2, resolution_notes = $3, resolved_at = NOW() WHERE id = $4`,
		status, resolvedBy, notes, id,
	)
	return err
}

// --- Marketplace ---

func (s *Store) ListMarketplaceContent(ctx context.Context, contentType, query, sort string, limit, offset int) ([]MarketplaceItem, int, error) {
	baseWhere := `WHERE ci.is_public = true AND EXISTS (SELECT 1 FROM license_offerings lo WHERE lo.content_id = ci.id AND lo.is_active = true)`
	args := []interface{}{}
	argIdx := 1

	if contentType != "" {
		baseWhere += ` AND ci.content_type = $` + fmt.Sprintf("%d", argIdx)
		args = append(args, contentType)
		argIdx++
	}

	if query != "" {
		baseWhere += ` AND ci.search_vector @@ plainto_tsquery('english', $` + fmt.Sprintf("%d", argIdx) + `)`
		args = append(args, query)
		argIdx++
	}

	// Count
	var total int
	countQuery := `SELECT COUNT(*) FROM content_items ci ` + baseWhere
	err := s.pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Determine ORDER BY
	orderBy := `ci.created_at DESC`
	if sort == "popular" {
		orderBy = `lowest_price ASC, ci.created_at DESC`
	}

	selectQuery := `SELECT ci.id, ci.user_id, ci.title, ci.description, ci.content_type, ci.mime_type, ci.file_size, ci.file_url, ci.thumbnail_url, ci.hash_sha256, ci.is_public, ci.tags, ci.created_at, ci.updated_at,
	                       u.name, u.username, u.image,
	                       COALESCE((SELECT MIN(lo.price_cents) FROM license_offerings lo WHERE lo.content_id = ci.id AND lo.is_active = true), 0) AS lowest_price
	                FROM content_items ci
	                JOIN users u ON u.id = ci.user_id
	                ` + baseWhere +
		` ORDER BY ` + orderBy + ` LIMIT $` + fmt.Sprintf("%d", argIdx) + ` OFFSET $` + fmt.Sprintf("%d", argIdx+1)
	args = append(args, limit, offset)

	rows, err := s.pool.Query(ctx, selectQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var items []MarketplaceItem
	for rows.Next() {
		var m MarketplaceItem
		if err := rows.Scan(
			&m.ID, &m.UserID, &m.Title, &m.Description, &m.ContentType,
			&m.MimeType, &m.FileSize, &m.FileURL, &m.ThumbnailURL,
			&m.HashSHA256, &m.IsPublic, &m.Tags, &m.CreatedAt, &m.UpdatedAt,
			&m.CreatorName, &m.CreatorUsername, &m.CreatorImage,
			&m.LowestPriceCents,
		); err != nil {
			return nil, 0, err
		}
		items = append(items, m)
	}
	if items == nil {
		items = []MarketplaceItem{}
	}
	return items, total, nil
}
