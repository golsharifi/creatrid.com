CREATE TABLE error_log (
    id BIGSERIAL PRIMARY KEY,
    source TEXT NOT NULL DEFAULT 'frontend',
    level TEXT NOT NULL DEFAULT 'error',
    message TEXT NOT NULL,
    stack TEXT,
    url TEXT,
    user_agent TEXT,
    user_id TEXT,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_error_log_created ON error_log(created_at DESC);
CREATE INDEX idx_error_log_source ON error_log(source, created_at DESC);
