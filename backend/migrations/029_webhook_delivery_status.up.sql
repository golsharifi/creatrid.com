ALTER TABLE webhook_deliveries ADD COLUMN IF NOT EXISTS attempts INT DEFAULT 0;
ALTER TABLE webhook_deliveries ADD COLUMN IF NOT EXISTS max_attempts INT DEFAULT 5;
ALTER TABLE webhook_deliveries ADD COLUMN IF NOT EXISTS next_retry_at TIMESTAMPTZ;
ALTER TABLE webhook_deliveries ADD COLUMN IF NOT EXISTS status TEXT DEFAULT 'pending';
CREATE INDEX IF NOT EXISTS idx_webhook_deliveries_status ON webhook_deliveries(status, next_retry_at);
