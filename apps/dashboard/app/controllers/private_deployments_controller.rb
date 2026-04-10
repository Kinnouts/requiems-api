# frozen_string_literal: true

class PrivateDeploymentsController < ApplicationController
  before_action :authenticate_user!

  def new
    @deployment_request = PrivateDeploymentRequest.new(
      contact_name: current_user.name,
      contact_email: current_user.email
    )
  end

  def create
    @deployment_request = current_user.private_deployment_requests.build(deployment_params)
    @deployment_request.monthly_price_cents = PrivateDeploymentRequest::TIER_PRICES[@deployment_request.server_tier] || 0

    if @deployment_request.save
      PrivateDeploymentMailer.request_received(@deployment_request).deliver_later
      PrivateDeploymentMailer.admin_notification(@deployment_request).deliver_later

      redirect_to dashboard_root_path,
        notice: "Your private deployment request was received! We'll have it ready within 24–48 hours."
    else
      render :new, status: :unprocessable_entity
    end
  end

  private

  def deployment_params
    params.require(:private_deployment_request).permit(
      :company, :contact_name, :contact_email, :server_tier, :admin_notes,
      selected_services: []
    )
  end
end
