class CreateDailyUsageSummaries < ActiveRecord::Migration[8.1]
  def change
    create_table :daily_usage_summaries do |t|
      t.references :user, null: false, foreign_key: true
      t.date :date
      t.integer :total_requests
      t.integer :total_credits

      t.timestamps
    end
  end
end
