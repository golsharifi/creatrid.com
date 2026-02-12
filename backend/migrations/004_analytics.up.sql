-- Profile view and link click analytics
CREATE TABLE profile_views (
    id         BIGSERIAL PRIMARY KEY,
    user_id    TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    viewer_ip  TEXT,
    referrer   TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_profile_views_user_id ON profile_views (user_id);
CREATE INDEX idx_profile_views_created_at ON profile_views (created_at);

CREATE TABLE link_clicks (
    id         BIGSERIAL PRIMARY KEY,
    user_id    TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    link_type  TEXT NOT NULL,  -- 'connection', 'custom', 'share'
    link_value TEXT NOT NULL,  -- platform name, link URL, or share type
    clicker_ip TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_link_clicks_user_id ON link_clicks (user_id);
CREATE INDEX idx_link_clicks_created_at ON link_clicks (created_at);
