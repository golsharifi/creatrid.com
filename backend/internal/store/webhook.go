package store

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jackc/pgx/v5"
)

type WebhookEndpoint struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	URL       string    `json:"url"`
	Secret    string    `json:"-"`
	Events    []string  `json:"events"`
	IsActive  bool      `json:"isActive"`
	CreatedAt time.Time `json:"createdAt"`
}

type WebhookDelivery struct {
	ID             int64            `json:"id"`
	EndpointID     string           `json:"endpointId"`
	EventType      string           `json:"eventType"`
	Payload        json.RawMessage  `json:"payload"`
	ResponseStatus *int             `json:"responseStatus"`
	ResponseBody   *string          `json:"responseBody"`
	DeliveredAt    *time.Time       `json:"deliveredAt"`
	CreatedAt      time.Time        `json:"createdAt"`
}

func (s *Store) CreateWebhookEndpoint(ctx context.Context, ep *WebhookEndpoint) error {
	_, err := s.pool.Exec(ctx,
		`INSERT INTO webhook_endpoints (id, user_id, url, secret, events, is_active, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		ep.ID, ep.UserID, ep.URL, ep.Secret, ep.Events, ep.IsActive, ep.CreatedAt,
	)
	return err
}

func (s *Store) ListWebhookEndpointsByUser(ctx context.Context, userID string) ([]*WebhookEndpoint, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT id, user_id, url, secret, events, is_active, created_at
		 FROM webhook_endpoints WHERE user_id = $1 ORDER BY created_at DESC`, userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var endpoints []*WebhookEndpoint
	for rows.Next() {
		var ep WebhookEndpoint
		if err := rows.Scan(&ep.ID, &ep.UserID, &ep.URL, &ep.Secret, &ep.Events, &ep.IsActive, &ep.CreatedAt); err != nil {
			return nil, err
		}
		endpoints = append(endpoints, &ep)
	}
	if endpoints == nil {
		endpoints = []*WebhookEndpoint{}
	}
	return endpoints, nil
}

func (s *Store) FindWebhookEndpointByID(ctx context.Context, id string) (*WebhookEndpoint, error) {
	var ep WebhookEndpoint
	err := s.pool.QueryRow(ctx,
		`SELECT id, user_id, url, secret, events, is_active, created_at
		 FROM webhook_endpoints WHERE id = $1`, id,
	).Scan(&ep.ID, &ep.UserID, &ep.URL, &ep.Secret, &ep.Events, &ep.IsActive, &ep.CreatedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return &ep, err
}

func (s *Store) UpdateWebhookEndpoint(ctx context.Context, id, url string, events []string, isActive bool) error {
	_, err := s.pool.Exec(ctx,
		`UPDATE webhook_endpoints SET url = $1, events = $2, is_active = $3 WHERE id = $4`,
		url, events, isActive, id,
	)
	return err
}

func (s *Store) DeleteWebhookEndpoint(ctx context.Context, id string) error {
	_, err := s.pool.Exec(ctx, `DELETE FROM webhook_endpoints WHERE id = $1`, id)
	return err
}

func (s *Store) ListActiveWebhookEndpointsForEvent(ctx context.Context, userID, eventType string) ([]*WebhookEndpoint, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT id, user_id, url, secret, events, is_active, created_at
		 FROM webhook_endpoints
		 WHERE user_id = $1 AND is_active = true AND $2 = ANY(events)`,
		userID, eventType,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var endpoints []*WebhookEndpoint
	for rows.Next() {
		var ep WebhookEndpoint
		if err := rows.Scan(&ep.ID, &ep.UserID, &ep.URL, &ep.Secret, &ep.Events, &ep.IsActive, &ep.CreatedAt); err != nil {
			return nil, err
		}
		endpoints = append(endpoints, &ep)
	}
	if endpoints == nil {
		endpoints = []*WebhookEndpoint{}
	}
	return endpoints, nil
}

func (s *Store) CreateWebhookDelivery(ctx context.Context, endpointID, eventType string, payload json.RawMessage) (int64, error) {
	var id int64
	err := s.pool.QueryRow(ctx,
		`INSERT INTO webhook_deliveries (endpoint_id, event_type, payload, created_at)
		 VALUES ($1, $2, $3, $4) RETURNING id`,
		endpointID, eventType, payload, time.Now(),
	).Scan(&id)
	return id, err
}

func (s *Store) UpdateWebhookDelivery(ctx context.Context, id int64, status int, body string) error {
	now := time.Now()
	_, err := s.pool.Exec(ctx,
		`UPDATE webhook_deliveries SET response_status = $1, response_body = $2, delivered_at = $3 WHERE id = $4`,
		status, body, now, id,
	)
	return err
}

func (s *Store) ListWebhookDeliveries(ctx context.Context, endpointID string, limit, offset int) ([]*WebhookDelivery, int, error) {
	var total int
	err := s.pool.QueryRow(ctx, `SELECT COUNT(*) FROM webhook_deliveries WHERE endpoint_id = $1`, endpointID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := s.pool.Query(ctx,
		`SELECT id, endpoint_id, event_type, payload, response_status, response_body, delivered_at, created_at
		 FROM webhook_deliveries WHERE endpoint_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		endpointID, limit, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var deliveries []*WebhookDelivery
	for rows.Next() {
		var d WebhookDelivery
		if err := rows.Scan(&d.ID, &d.EndpointID, &d.EventType, &d.Payload, &d.ResponseStatus, &d.ResponseBody, &d.DeliveredAt, &d.CreatedAt); err != nil {
			return nil, 0, err
		}
		deliveries = append(deliveries, &d)
	}
	if deliveries == nil {
		deliveries = []*WebhookDelivery{}
	}
	return deliveries, total, nil
}
