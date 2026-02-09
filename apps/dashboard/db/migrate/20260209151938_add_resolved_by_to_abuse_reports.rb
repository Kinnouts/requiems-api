class AddResolvedByToAbuseReports < ActiveRecord::Migration[8.1]
  def change
    add_column :abuse_reports, :resolved_by_id, :integer
  end
end
