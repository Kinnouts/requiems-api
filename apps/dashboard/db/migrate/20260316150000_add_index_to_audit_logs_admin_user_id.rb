class AddIndexToAuditLogsAdminUserId < ActiveRecord::Migration[8.1]
  def change
    add_index :audit_logs, :admin_user_id
  end
end
