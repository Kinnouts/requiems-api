# frozen_string_literal: true

class AddRequestMethodToUsageLogs < ActiveRecord::Migration[8.1]
  def change
    add_column :usage_logs, :request_method, :string
    add_index :usage_logs, [:api_key_id, :used_at, :endpoint], unique: true
    add_index :usage_logs, :request_method
    add_index :usage_logs, [:user_id, :request_method]
  end
end
