ALTER TABLE users ADD COLUMN IF NOT EXISTS referral_code TEXT UNIQUE;
ALTER TABLE users ADD COLUMN IF NOT EXISTS referred_by TEXT REFERENCES users(id);

CREATE TABLE referral_rewards (
    id TEXT PRIMARY KEY,
    referrer_user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    referred_user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    reward_type TEXT NOT NULL DEFAULT 'score_bonus',
    reward_value INT NOT NULL DEFAULT 5,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_referral_rewards_referrer ON referral_rewards(referrer_user_id);
