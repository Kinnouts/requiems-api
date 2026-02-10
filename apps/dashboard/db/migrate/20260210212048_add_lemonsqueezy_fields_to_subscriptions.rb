class AddLemonsqueezyFieldsToSubscriptions < ActiveRecord::Migration[8.1]
  def change
    add_column :subscriptions, :lemonsqueezy_subscription_id, :string
    add_column :subscriptions, :lemonsqueezy_customer_id, :string
    add_column :subscriptions, :plan_name, :string
    add_column :subscriptions, :cancel_at_period_end, :boolean, default: false

    add_index :subscriptions, :lemonsqueezy_subscription_id, unique: true
    add_index :subscriptions, :lemonsqueezy_customer_id
    add_index :subscriptions, :plan_name
  end
end
