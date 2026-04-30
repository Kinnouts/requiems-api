# frozen_string_literal: true

class AddSubscriptionFieldsToPrivateDeploymentRequests < ActiveRecord::Migration[8.1]
  def change
    add_column :private_deployment_requests, :billing_cycle, :string, null: false, default: "monthly"
    add_column :private_deployment_requests, :lemonsqueezy_subscription_id, :string
    add_index :private_deployment_requests, :lemonsqueezy_subscription_id, unique: true,
              where: "lemonsqueezy_subscription_id IS NOT NULL"
  end
end
