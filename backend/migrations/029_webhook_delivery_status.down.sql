DROP INDEX IF EXISTS idx_webhook_deliveries_status;
ALTER TABLE webhook_deliveries DROP COLUMN IF EXISTS status;
ALTER TABLE webhook_deliveries DROP COLUMN IF EXISTS next_retry_at;
ALTER TABLE webhook_deliveries DROP COLUMN IF EXISTS max_attempts;
ALTER TABLE webhook_deliveries DROP COLUMN IF EXISTS attempts;
