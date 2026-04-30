# frozen_string_literal: true

class CreateUsageLogs < ActiveRecord::Migration[8.1]
  def change
    create_table :usage_logs do |t|
      t.references :user, null: false, foreign_key: true
      t.references :api_key, null: false, foreign_key: true
      t.string :endpoint
      t.integer :credits_used
      t.integer :response_time_ms
      t.integer :status_code
      t.datetime :used_at
      t.date :usage_date

      t.timestamps
    end
  end
end
