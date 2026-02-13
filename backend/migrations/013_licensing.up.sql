CREATE TABLE license_offerings (
    id TEXT PRIMARY KEY,
    content_id TEXT NOT NULL REFERENCES content_items(id) ON DELETE CASCADE,
    license_type TEXT NOT NULL,
    price_cents INT NOT NULL,
    currency TEXT NOT NULL DEFAULT 'usd',
    is_active BOOLEAN NOT NULL DEFAULT true,
    terms_text TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(content_id, license_type)
);

CREATE TABLE license_purchases (
    id TEXT PRIMARY KEY,
    offering_id TEXT NOT NULL REFERENCES license_offerings(id),
    content_id TEXT NOT NULL REFERENCES content_items(id),
    buyer_user_id TEXT REFERENCES users(id),
    buyer_email TEXT NOT NULL,
    buyer_company TEXT,
    stripe_session_id TEXT,
    amount_cents INT NOT NULL,
    platform_fee_cents INT NOT NULL,
    creator_payout_cents INT NOT NULL,
    status TEXT NOT NULL DEFAULT 'completed',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_purchases_buyer ON license_purchases(buyer_user_id);
CREATE INDEX idx_purchases_content ON license_purchases(content_id);

CREATE TABLE takedown_requests (
    id TEXT PRIMARY KEY,
    reporter_email TEXT NOT NULL,
    reporter_name TEXT NOT NULL,
    content_id TEXT NOT NULL REFERENCES content_items(id),
    reason TEXT NOT NULL,
    evidence_url TEXT,
    status TEXT NOT NULL DEFAULT 'pending',
    resolved_by TEXT REFERENCES users(id),
    resolution_notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    resolved_at TIMESTAMPTZ
);
