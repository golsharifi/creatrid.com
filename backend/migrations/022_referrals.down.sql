DROP TABLE IF EXISTS referral_rewards;
ALTER TABLE users DROP COLUMN IF EXISTS referred_by;
ALTER TABLE users DROP COLUMN IF EXISTS referral_code;
