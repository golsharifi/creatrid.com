-- Enable trigram extension for fast ILIKE searches
CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- Create GIN trigram index for user search
CREATE INDEX IF NOT EXISTS idx_users_search_name ON users USING gin (name gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_users_search_username ON users USING gin (username gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_users_search_bio ON users USING gin (bio gin_trgm_ops);
