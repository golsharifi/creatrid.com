package store

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type AuditEntry struct {
	ID          int64           `json:"id"`
	AdminUserID string          `json:"adminUserId"`
	Action      string          `json:"action"`
	TargetType  string          `json:"targetType"`
	TargetID    *string         `json:"targetId"`
	Details     json.RawMessage `json:"details"`
	IPAddress   *string         `json:"ipAddress"`
	CreatedAt   time.Time       `json:"createdAt"`
	AdminName   *string         `json:"adminName,omitempty"`
}

func (s *Store) CreateAuditEntry(ctx context.Context, adminUserID, action, targetType string, targetID *string, details json.RawMessage, ipAddress *string) error {
	_, err := s.pool.Exec(ctx,
		`INSERT INTO audit_log (admin_user_id, action, target_type, target_id, details, ip_address)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		adminUserID, action, targetType, targetID, details, ipAddress,
	)
	return err
}

func (s *Store) ListAuditLog(ctx context.Context, limit, offset int) ([]*AuditEntry, int, error) {
	var total int
	err := s.pool.QueryRow(ctx, `SELECT COUNT(*) FROM audit_log`).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := s.pool.Query(ctx,
		`SELECT a.id, a.admin_user_id, a.action, a.target_type, a.target_id, a.details, a.ip_address, a.created_at, u.name
		 FROM audit_log a
		 LEFT JOIN users u ON u.id = a.admin_user_id
		 ORDER BY a.created_at DESC
		 LIMIT $1 OFFSET $2`,
		limit, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	entries := []*AuditEntry{}
	for rows.Next() {
		var e AuditEntry
		if err := rows.Scan(&e.ID, &e.AdminUserID, &e.Action, &e.TargetType, &e.TargetID, &e.Details, &e.IPAddress, &e.CreatedAt, &e.AdminName); err != nil {
			return nil, 0, fmt.Errorf("scan audit entry: %w", err)
		}
		entries = append(entries, &e)
	}
	return entries, total, nil
}
