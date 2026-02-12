package store

import (
	"context"

	"github.com/creatrid/creatrid/internal/model"
	"github.com/jackc/pgx/v5"
)

func (s *Store) FindAccountByProvider(ctx context.Context, provider, providerAccountID string) (*model.Account, error) {
	var a model.Account
	err := s.pool.QueryRow(ctx,
		`SELECT user_id, type, provider, provider_account_id,
		        refresh_token, access_token, expires_at, token_type, scope, id_token
		 FROM accounts WHERE provider = $1 AND provider_account_id = $2`,
		provider, providerAccountID,
	).Scan(
		&a.UserID, &a.Type, &a.Provider, &a.ProviderAccountID,
		&a.RefreshToken, &a.AccessToken, &a.ExpiresAt, &a.TokenType, &a.Scope, &a.IDToken,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return &a, err
}

func (s *Store) UpsertAccount(ctx context.Context, a *model.Account) error {
	_, err := s.pool.Exec(ctx,
		`INSERT INTO accounts (user_id, type, provider, provider_account_id,
		                       refresh_token, access_token, expires_at, token_type, scope, id_token)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		 ON CONFLICT (provider, provider_account_id)
		 DO UPDATE SET
			access_token = EXCLUDED.access_token,
			refresh_token = COALESCE(EXCLUDED.refresh_token, accounts.refresh_token),
			expires_at = EXCLUDED.expires_at,
			id_token = EXCLUDED.id_token`,
		a.UserID, a.Type, a.Provider, a.ProviderAccountID,
		a.RefreshToken, a.AccessToken, a.ExpiresAt, a.TokenType, a.Scope, a.IDToken,
	)
	return err
}
