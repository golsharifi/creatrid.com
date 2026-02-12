CREATE TYPE user_role AS ENUM ('CREATOR', 'BRAND', 'ADMIN');

CREATE TABLE users (
    id             TEXT PRIMARY KEY,
    name           TEXT,
    email          TEXT NOT NULL UNIQUE,
    email_verified TIMESTAMPTZ,
    image          TEXT,
    username       TEXT UNIQUE,
    bio            TEXT,
    role           user_role NOT NULL DEFAULT 'CREATOR',
    creator_score  INTEGER,
    is_verified    BOOLEAN NOT NULL DEFAULT FALSE,
    onboarded      BOOLEAN NOT NULL DEFAULT FALSE,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_users_username ON users (username) WHERE username IS NOT NULL;

CREATE TABLE accounts (
    user_id             TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type                TEXT NOT NULL,
    provider            TEXT NOT NULL,
    provider_account_id TEXT NOT NULL,
    refresh_token       TEXT,
    access_token        TEXT,
    expires_at          INTEGER,
    token_type          TEXT,
    scope               TEXT,
    id_token            TEXT,
    session_state       TEXT,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (provider, provider_account_id)
);

CREATE INDEX idx_accounts_user_id ON accounts (user_id);

CREATE OR REPLACE FUNCTION update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER accounts_updated_at
    BEFORE UPDATE ON accounts
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();
