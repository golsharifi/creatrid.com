DROP INDEX IF EXISTS idx_profile_views_country;
DROP INDEX IF EXISTS idx_profile_views_device;
ALTER TABLE profile_views DROP COLUMN IF EXISTS city;
ALTER TABLE profile_views DROP COLUMN IF EXISTS country;
ALTER TABLE profile_views DROP COLUMN IF EXISTS device_type;
ALTER TABLE profile_views DROP COLUMN IF EXISTS os;
ALTER TABLE profile_views DROP COLUMN IF EXISTS browser;
