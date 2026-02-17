-- Migration: Create api_keys table for metadata and audit trail
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
