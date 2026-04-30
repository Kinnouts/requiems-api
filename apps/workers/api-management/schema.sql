-- NOTE: existing production databases must be migrated to drop the old UNIQUE constraint.
-- See docs/migrations/d1-remove-credit-usage-unique-constraint.sql
DROP TABLE IF EXISTS credit_usage;
DROP TABLE IF EXISTS api_keys;

CREATE TABLE credit_usage (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  api_key TEXT NOT NULL,
  user_id TEXT NOT NULL,
  endpoint TEXT NOT NULL,
  credits_used INTEGER NOT NULL,
  request_method TEXT NOT NULL,
  status_code INTEGER NOT NULL,
  response_time_ms INTEGER NOT NULL,
  used_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE api_keys (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  key_prefix TEXT NOT NULL UNIQUE,
  user_id TEXT NOT NULL,
  plan TEXT NOT NULL,
  active INTEGER NOT NULL DEFAULT 1,
  created_at TEXT NOT NULL,
  updated_at TEXT,
  revoked_at TEXT,
  billing_cycle_start TEXT NOT NULL,

  CONSTRAINT idx_api_keys_user UNIQUE (user_id, key_prefix)
);

CREATE INDEX IF NOT EXISTS idx_credit_usage_user_lookup
ON credit_usage (user_id, used_at);

CREATE INDEX IF NOT EXISTS idx_credit_usage_lookup
ON credit_usage (api_key, used_at);

CREATE INDEX IF NOT EXISTS idx_credit_usage_endpoint
ON credit_usage (endpoint, used_at);

CREATE INDEX IF NOT EXISTS idx_api_keys_user_lookup
ON api_keys (user_id, active);

CREATE INDEX IF NOT EXISTS idx_api_keys_prefix_lookup
ON api_keys (key_prefix, active);
