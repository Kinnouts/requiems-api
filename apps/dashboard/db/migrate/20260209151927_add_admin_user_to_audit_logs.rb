class AddAdminUserToAuditLogs < ActiveRecord::Migration[8.1]
  def change
    add_column :audit_logs, :admin_user_id, :integer
  end
end
