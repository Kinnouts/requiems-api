# frozen_string_literal: true

class AddRequestMethodToUsageLogs < ActiveRecord::Migration[8.1]
  def change
    add_column :usage_logs, :request_method, :string
  end
end