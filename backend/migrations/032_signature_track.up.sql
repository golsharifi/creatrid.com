-- Signature track: a public Vault audio item that plays on the creator's public profile
ALTER TABLE users ADD COLUMN IF NOT EXISTS signature_track_id TEXT REFERENCES content_items(id) ON DELETE SET NULL;
