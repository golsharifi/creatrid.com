DROP TABLE IF EXISTS creator_payouts;
ALTER TABLE users DROP COLUMN IF EXISTS stripe_connect_onboarded;
ALTER TABLE users DROP COLUMN IF EXISTS stripe_connect_account_id;
