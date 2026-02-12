package store

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

type Subscription struct {
	ID                   string     `json:"id"`
	UserID               string     `json:"userId"`
	StripeCustomerID     string     `json:"stripeCustomerId"`
	StripeSubscriptionID *string    `json:"stripeSubscriptionId"`
	Plan                 string     `json:"plan"`
	Status               string     `json:"status"`
	CurrentPeriodEnd     *time.Time `json:"currentPeriodEnd"`
	CreatedAt            time.Time  `json:"createdAt"`
	UpdatedAt            time.Time  `json:"updatedAt"`
}

func (s *Store) UpsertSubscription(ctx context.Context, sub *Subscription) error {
	_, err := s.pool.Exec(ctx,
		`INSERT INTO subscriptions (id, user_id, stripe_customer_id, stripe_subscription_id, plan, status, current_period_end, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		 ON CONFLICT (user_id) DO UPDATE SET
		   stripe_customer_id = EXCLUDED.stripe_customer_id,
		   stripe_subscription_id = EXCLUDED.stripe_subscription_id,
		   plan = EXCLUDED.plan,
		   status = EXCLUDED.status,
		   current_period_end = EXCLUDED.current_period_end,
		   updated_at = NOW()`,
		sub.ID, sub.UserID, sub.StripeCustomerID, sub.StripeSubscriptionID,
		sub.Plan, sub.Status, sub.CurrentPeriodEnd, sub.CreatedAt, sub.UpdatedAt,
	)
	return err
}

func (s *Store) FindSubscriptionByUserID(ctx context.Context, userID string) (*Subscription, error) {
	var sub Subscription
	err := s.pool.QueryRow(ctx,
		`SELECT id, user_id, stripe_customer_id, stripe_subscription_id, plan, status, current_period_end, created_at, updated_at
		 FROM subscriptions WHERE user_id = $1`, userID,
	).Scan(
		&sub.ID, &sub.UserID, &sub.StripeCustomerID, &sub.StripeSubscriptionID,
		&sub.Plan, &sub.Status, &sub.CurrentPeriodEnd, &sub.CreatedAt, &sub.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return &sub, err
}

func (s *Store) FindSubscriptionByStripeCustomerID(ctx context.Context, customerID string) (*Subscription, error) {
	var sub Subscription
	err := s.pool.QueryRow(ctx,
		`SELECT id, user_id, stripe_customer_id, stripe_subscription_id, plan, status, current_period_end, created_at, updated_at
		 FROM subscriptions WHERE stripe_customer_id = $1`, customerID,
	).Scan(
		&sub.ID, &sub.UserID, &sub.StripeCustomerID, &sub.StripeSubscriptionID,
		&sub.Plan, &sub.Status, &sub.CurrentPeriodEnd, &sub.CreatedAt, &sub.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return &sub, err
}

func (s *Store) UpdateSubscriptionPlan(ctx context.Context, userID, plan, status string, periodEnd *time.Time) error {
	_, err := s.pool.Exec(ctx,
		`UPDATE subscriptions SET plan = $2, status = $3, current_period_end = $4, updated_at = NOW() WHERE user_id = $1`,
		userID, plan, status, periodEnd,
	)
	return err
}

func (s *Store) CreatePaymentEvent(ctx context.Context, id, userID, stripeEventID, eventType string, data []byte) error {
	_, err := s.pool.Exec(ctx,
		`INSERT INTO payment_events (id, user_id, stripe_event_id, event_type, data) VALUES ($1, $2, $3, $4, $5)`,
		id, userID, stripeEventID, eventType, data,
	)
	return err
}

func (s *Store) FindPaymentEventByStripeID(ctx context.Context, stripeEventID string) (bool, error) {
	var exists bool
	err := s.pool.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM payment_events WHERE stripe_event_id = $1)`, stripeEventID,
	).Scan(&exists)
	return exists, err
}
