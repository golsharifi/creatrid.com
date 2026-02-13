CREATE TABLE content_items (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    description TEXT,
    content_type TEXT NOT NULL,
    mime_type TEXT NOT NULL,
    file_size BIGINT NOT NULL,
    file_url TEXT NOT NULL,
    thumbnail_url TEXT,
    hash_sha256 TEXT NOT NULL,
    is_public BOOLEAN NOT NULL DEFAULT true,
    tags TEXT[],
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_content_user ON content_items(user_id);
CREATE INDEX idx_content_hash ON content_items(hash_sha256);
CREATE INDEX idx_content_type ON content_items(content_type);
CREATE INDEX idx_content_tags ON content_items USING gin(tags);
CREATE INDEX idx_content_public ON content_items(is_public, created_at DESC);
