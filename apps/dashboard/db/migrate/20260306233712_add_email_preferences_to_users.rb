class AddEmailPreferencesToUsers < ActiveRecord::Migration[8.1]
  def change
    add_column :users, :email_notifications, :boolean, default: true, null: false
    add_column :users, :usage_alerts, :boolean, default: true, null: false
    add_column :users, :weekly_reports, :boolean, default: false, null: false
  end
end
