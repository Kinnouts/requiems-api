class CreateAbuseReports < ActiveRecord::Migration[8.1]
  def change
    create_table :abuse_reports do |t|
      t.references :user, null: false, foreign_key: true
      t.references :api_key, null: false, foreign_key: true
      t.string :report_type
      t.text :description
      t.string :status
      t.datetime :resolved_at

      t.timestamps
    end
  end
end
