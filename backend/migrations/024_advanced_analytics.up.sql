-- Advanced analytics: device, browser, OS, and geo data on profile views
ALTER TABLE profile_views ADD COLUMN IF NOT EXISTS browser TEXT;
ALTER TABLE profile_views ADD COLUMN IF NOT EXISTS os TEXT;
ALTER TABLE profile_views ADD COLUMN IF NOT EXISTS device_type TEXT;
ALTER TABLE profile_views ADD COLUMN IF NOT EXISTS country TEXT;
ALTER TABLE profile_views ADD COLUMN IF NOT EXISTS city TEXT;

CREATE INDEX IF NOT EXISTS idx_profile_views_device ON profile_views (device_type);
CREATE INDEX IF NOT EXISTS idx_profile_views_country ON profile_views (country);
