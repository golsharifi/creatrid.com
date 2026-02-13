ALTER TABLE content_items ADD COLUMN IF NOT EXISTS search_vector tsvector;

CREATE INDEX IF NOT EXISTS idx_content_search ON content_items USING gin(search_vector);
CREATE INDEX IF NOT EXISTS idx_content_title_trgm ON content_items USING gin(title gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_content_desc_trgm ON content_items USING gin(description gin_trgm_ops);

CREATE OR REPLACE FUNCTION content_search_vector_update() RETURNS trigger AS $$
BEGIN
    NEW.search_vector :=
        setweight(to_tsvector('english', COALESCE(NEW.title, '')), 'A') ||
        setweight(to_tsvector('english', COALESCE(NEW.description, '')), 'B') ||
        setweight(to_tsvector('english', COALESCE(array_to_string(NEW.tags, ' '), '')), 'C');
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER content_search_vector_trigger
    BEFORE INSERT OR UPDATE OF title, description, tags
    ON content_items
    FOR EACH ROW EXECUTE FUNCTION content_search_vector_update();

UPDATE content_items SET search_vector =
    setweight(to_tsvector('english', COALESCE(title, '')), 'A') ||
    setweight(to_tsvector('english', COALESCE(description, '')), 'B') ||
    setweight(to_tsvector('english', COALESCE(array_to_string(tags, ' '), '')), 'C');
