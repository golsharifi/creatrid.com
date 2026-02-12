CREATE TABLE IF NOT EXISTS collaboration_requests (
    id              TEXT PRIMARY KEY,
    from_user_id    TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    to_user_id      TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    message         TEXT NOT NULL DEFAULT '',
    status          TEXT NOT NULL DEFAULT 'pending',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (from_user_id, to_user_id)
);

CREATE INDEX IF NOT EXISTS idx_collab_to_user ON collaboration_requests(to_user_id, status);
CREATE INDEX IF NOT EXISTS idx_collab_from_user ON collaboration_requests(from_user_id, status);
