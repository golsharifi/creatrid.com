-- Add email notification preferences to users
ALTER TABLE users ADD COLUMN IF NOT EXISTS email_prefs JSONB NOT NULL DEFAULT '{"welcome":true,"connectionAlert":true,"weeklyDigest":true,"collaborations":true}'::jsonb;
