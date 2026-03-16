class AddCompositeIndexToAbuseReports < ActiveRecord::Migration[8.1]
  def change
    add_index :abuse_reports, [ :status, :created_at ]
  end
end
