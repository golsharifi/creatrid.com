CREATE TABLE IF NOT EXISTS moderation_flags (
    id TEXT PRIMARY KEY,
    content_id TEXT NOT NULL,
    reason TEXT NOT NULL,
    details TEXT,
    status TEXT NOT NULL DEFAULT 'pending',
    resolved_by TEXT,
    resolved_at TIMESTAMPTZ,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_moderation_status ON moderation_flags (status);
CREATE INDEX IF NOT EXISTS idx_moderation_content ON moderation_flags (content_id);
