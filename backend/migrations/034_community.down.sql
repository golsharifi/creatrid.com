DROP TABLE IF EXISTS content_comments;
DROP TABLE IF EXISTS point_events;
DROP TABLE IF EXISTS user_stickers;
DROP TABLE IF EXISTS stickers;
ALTER TABLE users DROP COLUMN IF EXISTS arena_points;
ALTER TABLE users DROP COLUMN IF EXISTS arena_last_explore;
