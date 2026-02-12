-- Audit log for account deletions (GDPR compliance)
CREATE TABLE IF NOT EXISTS deletion_log (
    id           TEXT PRIMARY KEY,
    user_email   TEXT NOT NULL,
    requested_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
