# frozen_string_literal: true

class AddCompositeIndexToUsageLogs < ActiveRecord::Migration[8.1]
  def change
    add_index :usage_logs, [ :api_key_id, :used_at, :endpoint ],
              unique: true,
              name: "index_usage_logs_dedup"
  end
end
