DROP TRIGGER IF EXISTS content_search_vector_trigger ON content_items;
DROP FUNCTION IF EXISTS content_search_vector_update();
DROP INDEX IF EXISTS idx_content_search;
DROP INDEX IF EXISTS idx_content_title_trgm;
DROP INDEX IF EXISTS idx_content_desc_trgm;
ALTER TABLE content_items DROP COLUMN IF EXISTS search_vector;
