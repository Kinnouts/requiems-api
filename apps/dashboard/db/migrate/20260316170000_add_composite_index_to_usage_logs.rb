# frozen_string_literal: true

class AddCompositeIndexToUsageLogs < ActiveRecord::Migration[8.1]
  def change
    # Index on (user_id, used_at) already exists as index_usage_logs_on_user_and_time
    # added in 20260210174437_add_indexes_to_usage_tables.rb
  end
end
