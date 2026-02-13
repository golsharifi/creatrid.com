-- Creator tokens (each creator can have one token)
CREATE TABLE IF NOT EXISTS creator_tokens (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    symbol TEXT NOT NULL,
    description TEXT,
    total_supply INT NOT NULL DEFAULT 0,
    price_cents INT NOT NULL DEFAULT 100,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(user_id),
    UNIQUE(symbol)
);

-- Token balances
CREATE TABLE IF NOT EXISTS token_balances (
    id TEXT PRIMARY KEY,
    token_id TEXT NOT NULL REFERENCES creator_tokens(id) ON DELETE CASCADE,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    balance INT NOT NULL DEFAULT 0,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(token_id, user_id)
);
CREATE INDEX IF NOT EXISTS idx_token_balances_user ON token_balances(user_id);

-- Tips (one-time payments)
CREATE TABLE IF NOT EXISTS tips (
    id TEXT PRIMARY KEY,
    from_user_id TEXT NOT NULL REFERENCES users(id),
    to_user_id TEXT NOT NULL REFERENCES users(id),
    amount_cents INT NOT NULL,
    message TEXT,
    stripe_payment_id TEXT,
    status TEXT NOT NULL DEFAULT 'pending',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_tips_to_user ON tips(to_user_id);
CREATE INDEX IF NOT EXISTS idx_tips_from_user ON tips(from_user_id);

-- Fan subscriptions (recurring)
CREATE TABLE IF NOT EXISTS fan_subscriptions (
    id TEXT PRIMARY KEY,
    fan_user_id TEXT NOT NULL REFERENCES users(id),
    creator_user_id TEXT NOT NULL REFERENCES users(id),
    tier TEXT NOT NULL DEFAULT 'supporter',
    price_cents INT NOT NULL,
    stripe_subscription_id TEXT,
    status TEXT NOT NULL DEFAULT 'active',
    started_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    canceled_at TIMESTAMPTZ,
    UNIQUE(fan_user_id, creator_user_id)
);
CREATE INDEX IF NOT EXISTS idx_fan_subs_creator ON fan_subscriptions(creator_user_id);
CREATE INDEX IF NOT EXISTS idx_fan_subs_fan ON fan_subscriptions(fan_user_id);

-- Token-gated content
CREATE TABLE IF NOT EXISTS gated_content (
    id TEXT PRIMARY KEY,
    content_id TEXT NOT NULL REFERENCES content_items(id) ON DELETE CASCADE,
    token_id TEXT REFERENCES creator_tokens(id),
    min_tokens INT DEFAULT 1,
    subscription_tier TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(content_id)
);

-- Token transactions log
CREATE TABLE IF NOT EXISTS token_transactions (
    id TEXT PRIMARY KEY,
    token_id TEXT NOT NULL REFERENCES creator_tokens(id),
    from_user_id TEXT REFERENCES users(id),
    to_user_id TEXT REFERENCES users(id),
    amount INT NOT NULL,
    tx_type TEXT NOT NULL,
    reference_id TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_token_tx_token ON token_transactions(token_id);
