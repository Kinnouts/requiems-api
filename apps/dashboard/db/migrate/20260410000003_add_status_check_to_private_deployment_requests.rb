# frozen_string_literal: true

class AddStatusCheckToPrivateDeploymentRequests < ActiveRecord::Migration[8.1]
  def up
    add_check_constraint :private_deployment_requests,
      "status IN ('pending_payment', 'pending', 'deploying', 'active', 'cancelled')",
      name: "private_deployment_requests_status_check"
  end

  def down
    remove_check_constraint :private_deployment_requests,
      name: "private_deployment_requests_status_check"
  end
end
