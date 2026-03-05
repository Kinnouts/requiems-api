# frozen_string_literal: true
class AddAdminUserToCreditAdjustments < ActiveRecord::Migration[8.1]
  def change
    add_column :credit_adjustments, :admin_user_id, :integer
  end
end
