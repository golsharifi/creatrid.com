-- Community economy: likes on public works (community ratings feed reputation)
-- and the closed-loop sticker exchange (trade stickers for points, in-platform
-- only — nothing is cashable or leaves the platform).

CREATE TABLE IF NOT EXISTS content_likes (
    content_id TEXT NOT NULL REFERENCES content_items(id) ON DELETE CASCADE,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (content_id, user_id)
);
CREATE INDEX IF NOT EXISTS idx_content_likes_user ON content_likes(user_id);

CREATE TABLE IF NOT EXISTS sticker_listings (
    id TEXT PRIMARY KEY,
    seller_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    sticker_id TEXT NOT NULL REFERENCES stickers(id) ON DELETE CASCADE,
    price_points INT NOT NULL CHECK (price_points > 0),
    status TEXT NOT NULL DEFAULT 'open', -- open | sold | cancelled
    buyer_id TEXT REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    sold_at TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_sticker_listings_status ON sticker_listings(status, created_at);
CREATE INDEX IF NOT EXISTS idx_sticker_listings_seller ON sticker_listings(seller_id);
