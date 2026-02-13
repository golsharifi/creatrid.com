package store

import (
	"context"
)

func (s *Store) SetTOTPSecret(ctx context.Context, userID, secret string) error {
	_, err := s.pool.Exec(ctx,
		`UPDATE users SET totp_secret = $1 WHERE id = $2`,
		secret, userID)
	return err
}

func (s *Store) EnableTOTP(ctx context.Context, userID string, backupCodes []string) error {
	_, err := s.pool.Exec(ctx,
		`UPDATE users SET totp_enabled = TRUE, totp_backup_codes = $1 WHERE id = $2`,
		backupCodes, userID)
	return err
}

func (s *Store) DisableTOTP(ctx context.Context, userID string) error {
	_, err := s.pool.Exec(ctx,
		`UPDATE users SET totp_enabled = FALSE, totp_secret = NULL, totp_backup_codes = NULL WHERE id = $1`,
		userID)
	return err
}

func (s *Store) GetTOTPSecret(ctx context.Context, userID string) (string, bool, error) {
	var secret *string
	var enabled bool
	err := s.pool.QueryRow(ctx,
		`SELECT totp_secret, COALESCE(totp_enabled, false) FROM users WHERE id = $1`, userID,
	).Scan(&secret, &enabled)
	if err != nil {
		return "", false, err
	}
	if secret == nil {
		return "", enabled, nil
	}
	return *secret, enabled, nil
}

func (s *Store) ConsumeTOTPBackupCode(ctx context.Context, userID, code string) (bool, error) {
	result, err := s.pool.Exec(ctx,
		`UPDATE users SET totp_backup_codes = array_remove(totp_backup_codes, $1)
		 WHERE id = $2 AND $1 = ANY(totp_backup_codes)`,
		code, userID)
	if err != nil {
		return false, err
	}
	return result.RowsAffected() > 0, nil
}
