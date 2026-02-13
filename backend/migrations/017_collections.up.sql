CREATE TABLE content_collections (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    description TEXT,
    is_public BOOLEAN NOT NULL DEFAULT TRUE,
    cover_image_url TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_collections_user ON content_collections(user_id);
CREATE INDEX idx_collections_public ON content_collections(is_public, created_at DESC);

CREATE TABLE collection_items (
    collection_id TEXT NOT NULL REFERENCES content_collections(id) ON DELETE CASCADE,
    content_id TEXT NOT NULL REFERENCES content_items(id) ON DELETE CASCADE,
    position INT NOT NULL DEFAULT 0,
    added_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (collection_id, content_id)
);

CREATE INDEX idx_collection_items_content ON collection_items(content_id);
