# frozen_string_literal: true

class EncryptTenantSecretOnPrivateDeploymentRequests < ActiveRecord::Migration[8.1]
  def up
    # Encrypted values are significantly longer than 255 chars — widen to text.
    change_column :private_deployment_requests, :tenant_secret, :text
  end

  def down
    change_column :private_deployment_requests, :tenant_secret, :string
  end
end
