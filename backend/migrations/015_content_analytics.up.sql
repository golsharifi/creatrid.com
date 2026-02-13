CREATE TABLE content_views (
    id BIGSERIAL PRIMARY KEY,
    content_id TEXT NOT NULL REFERENCES content_items(id) ON DELETE CASCADE,
    viewer_ip TEXT,
    referrer TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_content_views_content ON content_views(content_id);
CREATE INDEX idx_content_views_created ON content_views(created_at);

CREATE TABLE content_downloads (
    id BIGSERIAL PRIMARY KEY,
    content_id TEXT NOT NULL REFERENCES content_items(id) ON DELETE CASCADE,
    downloader_user_id TEXT REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_content_downloads_content ON content_downloads(content_id);
