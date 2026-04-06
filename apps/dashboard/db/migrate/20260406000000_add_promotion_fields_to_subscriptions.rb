# frozen_string_literal: true

class AddPromotionFieldsToSubscriptions < ActiveRecord::Migration[8.1]
  def change
    add_column :subscriptions, :promoted_by_id, :bigint, null: true
    add_column :subscriptions, :promotion_reason, :text, null: true
    add_column :subscriptions, :promotion_expires_at, :datetime, null: true

    add_index :subscriptions, :promotion_expires_at
    add_foreign_key :subscriptions, :users, column: :promoted_by_id
  end
end
