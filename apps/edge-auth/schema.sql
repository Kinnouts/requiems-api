-- D1 Schema for Requiem API Gateway
-- This runs in Cloudflare D1 (edge SQLite)

-- Usage tracking table
CREATE TABLE IF NOT EXISTS credit_usage (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  api_key TEXT NOT NULL,
  endpoint TEXT NOT NULL,
  credits_used INTEGER NOT NULL,
  used_at TEXT NOT NULL DEFAULT (datetime('now')),
  
  -- Index for fast usage queries
  CONSTRAINT idx_api_key_date UNIQUE (api_key, used_at, endpoint)
);

-- Index for querying usage by API key and date range
CREATE INDEX IF NOT EXISTS idx_credit_usage_lookup 
ON credit_usage (api_key, used_at);

-- Index for querying by endpoint (for analytics)
CREATE INDEX IF NOT EXISTS idx_credit_usage_endpoint 
ON credit_usage (endpoint, used_at);
