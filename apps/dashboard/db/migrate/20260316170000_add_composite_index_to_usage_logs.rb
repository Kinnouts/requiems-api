# frozen_string_literal: true

class AddCompositeIndexToUsageLogs < ActiveRecord::Migration[8.1]
  def change
    add_index :usage_logs, [ :user_id, :used_at ]
  end
end
