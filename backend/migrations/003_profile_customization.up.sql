-- Profile customization: theme and custom links
ALTER TABLE users ADD COLUMN IF NOT EXISTS theme TEXT NOT NULL DEFAULT 'default';
ALTER TABLE users ADD COLUMN IF NOT EXISTS custom_links JSONB NOT NULL DEFAULT '[]';
