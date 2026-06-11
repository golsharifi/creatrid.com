ALTER TABLE users DROP COLUMN IF EXISTS watermark_url;
ALTER TABLE content_items DROP COLUMN IF EXISTS is_encrypted;
ALTER TABLE content_items DROP COLUMN IF EXISTS is_watermarked;
