-- AI generation usage tracking (for monthly per-user quotas / plan gating)
CREATE TABLE IF NOT EXISTS ai_generations (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    task TEXT NOT NULL, -- logos | copy | refine | legal
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_ai_generations_user_created ON ai_generations(user_id, created_at);
