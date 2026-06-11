-- VR Community Arena: stickers, points, leaderboard, comments. All closed-loop:
-- points and stickers have in-platform utility only and are NOT cashable or
-- externally tradeable (deliberate securities posture).

ALTER TABLE users ADD COLUMN IF NOT EXISTS arena_points INT NOT NULL DEFAULT 0;
ALTER TABLE users ADD COLUMN IF NOT EXISTS arena_last_explore TIMESTAMPTZ;

-- Sticker catalog (seeded, farm-themed starter set)
CREATE TABLE IF NOT EXISTS stickers (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    emoji TEXT NOT NULL,
    rarity TEXT NOT NULL DEFAULT 'common', -- common | uncommon | rare | legendary
    description TEXT NOT NULL DEFAULT ''
);

INSERT INTO stickers (id, name, emoji, rarity, description) VALUES
    ('rooster',   'Rooster',        '🐓', 'common',    'First on the farm, first on the feed.'),
    ('sheep',     'Sheep',          '🐑', 'common',    'Follows the flock — but collects them all.'),
    ('pig',       'Pig',            '🐖', 'common',    'Lives for the mud and the grind.'),
    ('cow',       'Cow',            '🐄', 'common',    'Steady output, every single day.'),
    ('duck',      'Duck',           '🦆', 'common',    'Calm on the surface, paddling underneath.'),
    ('goat',      'Goat',           '🐐', 'uncommon',  'The greatest of all time, obviously.'),
    ('horse',     'Horse',          '🐎', 'uncommon',  'Built for the long run.'),
    ('owl',       'Night Owl',      '🦉', 'uncommon',  'Creates while the farm sleeps.'),
    ('fox',       'Fox',            '🦊', 'rare',      'Clever finds, rare drops.'),
    ('peacock',   'Peacock',        '🦚', 'rare',      'Style is a strategy.'),
    ('dragon',    'Barn Dragon',    '🐉', 'legendary', 'Legend says it guards the vault.'),
    ('unicorn',   'Unicorn',        '🦄', 'legendary', 'A one-of-a-kind creator spirit.')
ON CONFLICT (id) DO NOTHING;

-- Stickers a user has collected
CREATE TABLE IF NOT EXISTS user_stickers (
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    sticker_id TEXT NOT NULL REFERENCES stickers(id) ON DELETE CASCADE,
    count INT NOT NULL DEFAULT 1,
    first_earned_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, sticker_id)
);

-- Point ledger (audit trail for arena_points)
CREATE TABLE IF NOT EXISTS point_events (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    delta INT NOT NULL,
    reason TEXT NOT NULL, -- explore | sticker_bonus | gift_sent | gift_received | comment_received
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_point_events_user ON point_events(user_id, created_at);

-- Comments on public vault content
CREATE TABLE IF NOT EXISTS content_comments (
    id TEXT PRIMARY KEY,
    content_id TEXT NOT NULL REFERENCES content_items(id) ON DELETE CASCADE,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    body TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_content_comments_content ON content_comments(content_id, created_at);
