CREATE TABLE IF NOT EXISTS content_anchors (
    id TEXT PRIMARY KEY,
    content_id TEXT NOT NULL REFERENCES content_items(id) ON DELETE CASCADE,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content_hash TEXT NOT NULL,
    tx_hash TEXT,
    chain TEXT NOT NULL DEFAULT 'polygon',
    block_number BIGINT,
    contract_address TEXT,
    anchor_status TEXT NOT NULL DEFAULT 'pending',
    error_message TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    confirmed_at TIMESTAMPTZ,
    UNIQUE(content_id)
);
CREATE INDEX IF NOT EXISTS idx_content_anchors_user ON content_anchors(user_id);
CREATE INDEX IF NOT EXISTS idx_content_anchors_hash ON content_anchors(content_hash);
CREATE INDEX IF NOT EXISTS idx_content_anchors_tx ON content_anchors(tx_hash);
