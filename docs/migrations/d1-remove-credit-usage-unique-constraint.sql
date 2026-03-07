-- Migration: Remove UNIQUE(api_key, used_at, endpoint) from credit_usage
--
-- The original constraint silently drops inserts when two requests for the same
-- api_key + endpoint arrive within the same second (datetime('now') resolution).
-- The table already has id INTEGER PRIMARY KEY AUTOINCREMENT which is sufficient
-- as a unique key. Deduplication during the Rails D1 sync is handled at the
-- PostgreSQL layer via insert_all unique_by: [:api_key_id, :used_at, :endpoint].
--
-- SQLite does not support DROP CONSTRAINT, so we recreate the table.
-- This file is idempotent and safe to re-run.
--
-- Canonical runnable file: apps/workers/auth-gateway/migrations/001_remove_credit_usage_unique_constraint.sql
--
-- Dev (local) — run from apps/workers/auth-gateway/:
--   pnpm run db:migrate:001
-- Production:
--   pnpm run db:migrate:001:prod

-- Clean up any leftover temp table from a failed previous run
DROP TABLE IF EXISTS credit_usage_new;

CREATE TABLE credit_usage_new (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  api_key TEXT NOT NULL,
  user_id TEXT NOT NULL,
  endpoint TEXT NOT NULL,
  credits_used INTEGER NOT NULL,
  used_at TEXT NOT NULL DEFAULT (datetime('now'))
);

INSERT OR IGNORE INTO credit_usage_new (id, api_key, user_id, endpoint, credits_used, used_at)
SELECT id, api_key, user_id, endpoint, credits_used, used_at
FROM credit_usage;

DROP TABLE credit_usage;

ALTER TABLE credit_usage_new RENAME TO credit_usage;

CREATE INDEX IF NOT EXISTS idx_credit_usage_user_lookup
ON credit_usage (user_id, used_at);

CREATE INDEX IF NOT EXISTS idx_credit_usage_lookup
ON credit_usage (api_key, used_at);

CREATE INDEX IF NOT EXISTS idx_credit_usage_endpoint
ON credit_usage (endpoint, used_at);
