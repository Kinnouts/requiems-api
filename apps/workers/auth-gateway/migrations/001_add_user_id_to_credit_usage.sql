-- Migration: Add user_id to credit_usage table
-- This allows checking quotas at the user level (all API keys share quota)

-- Add user_id column
ALTER TABLE credit_usage ADD COLUMN user_id TEXT;

-- Create index for fast user-level usage queries
CREATE INDEX IF NOT EXISTS idx_credit_usage_user_lookup
ON credit_usage (user_id, used_at);

-- Note: Existing records will have NULL user_id
-- They will be gradually replaced as new requests come in
