-- Agency accounts
CREATE TABLE IF NOT EXISTS agencies (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    website TEXT,
    logo_url TEXT,
    description TEXT,
    is_verified BOOLEAN DEFAULT false,
    max_creators INT DEFAULT 10,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(user_id)
);

-- Agency-creator relationships
CREATE TABLE IF NOT EXISTS agency_creators (
    id TEXT PRIMARY KEY,
    agency_id TEXT NOT NULL REFERENCES agencies(id) ON DELETE CASCADE,
    creator_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status TEXT NOT NULL DEFAULT 'pending', -- pending, active, removed
    invited_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    joined_at TIMESTAMPTZ,
    UNIQUE(agency_id, creator_id)
);
CREATE INDEX IF NOT EXISTS idx_agency_creators_agency ON agency_creators(agency_id);
CREATE INDEX IF NOT EXISTS idx_agency_creators_creator ON agency_creators(creator_id);

-- API usage tracking
CREATE TABLE IF NOT EXISTS api_usage (
    id BIGSERIAL PRIMARY KEY,
    api_key_id TEXT NOT NULL,
    endpoint TEXT NOT NULL,
    method TEXT NOT NULL,
    status_code INT,
    response_time_ms INT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_api_usage_key ON api_usage(api_key_id);
CREATE INDEX IF NOT EXISTS idx_api_usage_created ON api_usage(created_at);
