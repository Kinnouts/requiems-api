CREATE TABLE IF NOT EXISTS credit_usage (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  api_key TEXT NOT NULL,
  user_id TEXT NOT NULL,
  endpoint TEXT NOT NULL,
  credits_used INTEGER NOT NULL,
  used_at TEXT NOT NULL DEFAULT (datetime('now')),

  -- Index for fast usage queries
  CONSTRAINT idx_api_key_date UNIQUE (api_key, used_at, endpoint)
);

-- Index for querying usage by user (for quota checking - all keys share quota)
CREATE INDEX IF NOT EXISTS idx_credit_usage_user_lookup
ON credit_usage (user_id, used_at);

-- Index for querying usage by API key and date range (for analytics)
CREATE INDEX IF NOT EXISTS idx_credit_usage_lookup
ON credit_usage (api_key, used_at);

-- Index for querying by endpoint (for analytics)
CREATE INDEX IF NOT EXISTS idx_credit_usage_endpoint
ON credit_usage (endpoint, used_at);

-- API keys table for metadata and audit trail
-- Actual validation happens via KV, this is for tracking/analytics
CREATE TABLE IF NOT EXISTS api_keys (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  key_prefix TEXT NOT NULL UNIQUE,
  user_id TEXT NOT NULL,
  plan TEXT NOT NULL,
  active INTEGER NOT NULL DEFAULT 1,
  created_at TEXT NOT NULL,
  updated_at TEXT,
  revoked_at TEXT,
  billing_cycle_start TEXT NOT NULL,

  -- Index for user lookups
  CONSTRAINT idx_api_keys_user UNIQUE (user_id, key_prefix)
);

-- Index for looking up keys by user
CREATE INDEX IF NOT EXISTS idx_api_keys_user_lookup
ON api_keys (user_id, active);

-- Index for looking up by key prefix
CREATE INDEX IF NOT EXISTS idx_api_keys_prefix_lookup
ON api_keys (key_prefix, active);
