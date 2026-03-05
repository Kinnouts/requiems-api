# frozen_string_literal: true
class CreateSubscriptions < ActiveRecord::Migration[8.1]
  def change
    create_table :subscriptions do |t|
      t.references :user, null: false, foreign_key: true
      t.string :plan
      t.string :stripe_subscription_id
      t.string :stripe_customer_id
      t.string :status
      t.integer :credit_limit
      t.datetime :current_period_start
      t.datetime :current_period_end
      t.datetime :trial_ends_at
      t.datetime :canceled_at

      t.timestamps
    end
  end
end
