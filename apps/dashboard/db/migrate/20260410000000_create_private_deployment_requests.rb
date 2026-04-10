# frozen_string_literal: true

class CreatePrivateDeploymentRequests < ActiveRecord::Migration[8.1]
  def change
    create_table :private_deployment_requests do |t|
      t.references :user, null: false, foreign_key: true
      t.string :company, null: false
      t.string :contact_name, null: false
      t.string :contact_email, null: false
      t.string :server_tier, null: false
      t.integer :monthly_price_cents, null: false
      t.jsonb :selected_services, null: false, default: []
      t.string :subdomain_slug
      t.string :tenant_secret
      t.string :status, null: false, default: "pending"
      t.text :admin_notes
      t.datetime :deployed_at
      t.timestamps
    end

    add_index :private_deployment_requests, :status
    add_index :private_deployment_requests, :subdomain_slug, unique: true,
              where: "subdomain_slug IS NOT NULL"
  end
end
