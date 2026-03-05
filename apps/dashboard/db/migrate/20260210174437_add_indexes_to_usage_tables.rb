# frozen_string_literal: true
class AddIndexesToUsageTables < ActiveRecord::Migration[8.1]
  def change
    # Usage logs indexes for analytics queries
    add_index :usage_logs, [ :user_id, :used_at ], name: 'index_usage_logs_on_user_and_time'
    add_index :usage_logs, [ :api_key_id, :used_at ], name: 'index_usage_logs_on_api_key_and_time'
    add_index :usage_logs, [ :endpoint, :used_at ], name: 'index_usage_logs_on_endpoint_and_time'
    add_index :usage_logs, :status_code, name: 'index_usage_logs_on_status_code'
    add_index :usage_logs, :usage_date, name: 'index_usage_logs_on_usage_date'

    # Daily usage summaries indexes for dashboard queries
    # Note: Quotas are per-user (all API keys share the same quota)
    add_index :daily_usage_summaries, [ :user_id, :date ], name: 'index_daily_usage_on_user_and_date'
  end
end
