-- Watermark image (PNG composited onto image uploads) + client-side encrypted vault files
ALTER TABLE users ADD COLUMN IF NOT EXISTS watermark_url TEXT;
ALTER TABLE content_items ADD COLUMN IF NOT EXISTS is_encrypted BOOLEAN NOT NULL DEFAULT false;
ALTER TABLE content_items ADD COLUMN IF NOT EXISTS is_watermarked BOOLEAN NOT NULL DEFAULT false;
