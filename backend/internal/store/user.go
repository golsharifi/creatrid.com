package store

import (
	"context"
	"time"

	"github.com/creatrid/creatrid/internal/model"
	"github.com/jackc/pgx/v5"
)

// EmailVerification represents a pending or completed email verification.
type EmailVerification struct {
	ID         string
	UserID     string
	Token      string
	ExpiresAt  time.Time
	VerifiedAt *time.Time
	CreatedAt  time.Time
}

func (s *Store) FindUserByID(ctx context.Context, id string) (*model.User, error) {
	var u model.User
	err := s.pool.QueryRow(ctx,
		`SELECT id, name, email, email_verified, image, username, bio, role,
		        creator_score, is_verified, onboarded, theme, custom_links, email_prefs, created_at, updated_at
		 FROM users WHERE id = $1`, id,
	).Scan(
		&u.ID, &u.Name, &u.Email, &u.EmailVerified, &u.Image,
		&u.Username, &u.Bio, &u.Role, &u.CreatorScore,
		&u.IsVerified, &u.Onboarded, &u.Theme, &u.CustomLinks, &u.EmailPrefsRaw, &u.CreatedAt, &u.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return &u, err
}

func (s *Store) FindUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var u model.User
	err := s.pool.QueryRow(ctx,
		`SELECT id, name, email, email_verified, image, username, bio, role,
		        creator_score, is_verified, onboarded, theme, custom_links, email_prefs, created_at, updated_at
		 FROM users WHERE email = $1`, email,
	).Scan(
		&u.ID, &u.Name, &u.Email, &u.EmailVerified, &u.Image,
		&u.Username, &u.Bio, &u.Role, &u.CreatorScore,
		&u.IsVerified, &u.Onboarded, &u.Theme, &u.CustomLinks, &u.EmailPrefsRaw, &u.CreatedAt, &u.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return &u, err
}

func (s *Store) FindUserByUsername(ctx context.Context, username string) (*model.User, error) {
	var u model.User
	err := s.pool.QueryRow(ctx,
		`SELECT id, name, email, email_verified, image, username, bio, role,
		        creator_score, is_verified, onboarded, theme, custom_links, email_prefs, created_at, updated_at
		 FROM users WHERE username = $1`, username,
	).Scan(
		&u.ID, &u.Name, &u.Email, &u.EmailVerified, &u.Image,
		&u.Username, &u.Bio, &u.Role, &u.CreatorScore,
		&u.IsVerified, &u.Onboarded, &u.Theme, &u.CustomLinks, &u.EmailPrefsRaw, &u.CreatedAt, &u.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return &u, err
}

func (s *Store) CreateUser(ctx context.Context, u *model.User) error {
	_, err := s.pool.Exec(ctx,
		`INSERT INTO users (id, name, email, email_verified, image, role)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		u.ID, u.Name, u.Email, u.EmailVerified, u.Image, u.Role,
	)
	return err
}

func (s *Store) UpdateUserOnboarding(ctx context.Context, id, username, name string) error {
	_, err := s.pool.Exec(ctx,
		`UPDATE users SET username = $1, name = $2, onboarded = true WHERE id = $3`,
		username, name, id,
	)
	return err
}

func (s *Store) UpdateUserProfile(ctx context.Context, id string, name, bio, username, theme *string) error {
	_, err := s.pool.Exec(ctx,
		`UPDATE users SET
			name = COALESCE($1, name),
			bio = COALESCE($2, bio),
			username = COALESCE($3, username),
			theme = COALESCE($4, theme)
		 WHERE id = $5`,
		name, bio, username, theme, id,
	)
	return err
}

func (s *Store) UpdateUserCustomLinks(ctx context.Context, id string, links []byte) error {
	_, err := s.pool.Exec(ctx,
		`UPDATE users SET custom_links = $1 WHERE id = $2`,
		links, id,
	)
	return err
}

func (s *Store) UpdateUserImage(ctx context.Context, id string, imageURL string) error {
	_, err := s.pool.Exec(ctx,
		`UPDATE users SET image = $1 WHERE id = $2`,
		imageURL, id,
	)
	return err
}

func (s *Store) UpdateUserScore(ctx context.Context, id string, score int) error {
	_, err := s.pool.Exec(ctx,
		`UPDATE users SET creator_score = $1 WHERE id = $2`,
		score, id,
	)
	return err
}

func (s *Store) UpdateUserEmailPrefs(ctx context.Context, id string, prefs []byte) error {
	_, err := s.pool.Exec(ctx,
		`UPDATE users SET email_prefs = $1 WHERE id = $2`,
		prefs, id,
	)
	return err
}

func (s *Store) DeleteUser(ctx context.Context, id string) error {
	_, err := s.pool.Exec(ctx, `DELETE FROM users WHERE id = $1`, id)
	return err
}

func (s *Store) LogDeletion(ctx context.Context, id, email string) error {
	_, err := s.pool.Exec(ctx,
		`INSERT INTO deletion_log (id, user_email) VALUES ($1, $2)`,
		id, email,
	)
	return err
}

func (s *Store) CreateEmailVerification(ctx context.Context, id, userID, token string, expiresAt time.Time) error {
	_, err := s.pool.Exec(ctx,
		`INSERT INTO email_verifications (id, user_id, token, expires_at) VALUES ($1, $2, $3, $4)`,
		id, userID, token, expiresAt,
	)
	return err
}

func (s *Store) FindEmailVerification(ctx context.Context, token string) (*EmailVerification, error) {
	var ev EmailVerification
	err := s.pool.QueryRow(ctx,
		`SELECT id, user_id, token, expires_at, verified_at, created_at
		 FROM email_verifications WHERE token = $1`, token,
	).Scan(&ev.ID, &ev.UserID, &ev.Token, &ev.ExpiresAt, &ev.VerifiedAt, &ev.CreatedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return &ev, err
}

func (s *Store) MarkEmailVerified(ctx context.Context, token string, verifiedAt time.Time) error {
	_, err := s.pool.Exec(ctx,
		`UPDATE email_verifications SET verified_at = $1 WHERE token = $2`,
		verifiedAt, token,
	)
	return err
}

func (s *Store) SetUserEmailVerified(ctx context.Context, userID string, verifiedAt time.Time) error {
	_, err := s.pool.Exec(ctx,
		`UPDATE users SET email_verified = $1 WHERE id = $2`,
		verifiedAt, userID,
	)
	return err
}
