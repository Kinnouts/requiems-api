# frozen_string_literal: true

class PrivateDeploymentsController < ApplicationController
  before_action :authenticate_user!

  def new
    @deployment_request = PrivateDeploymentRequest.new(billing_cycle: "monthly")
  end

  def create
    @deployment_request = current_user.private_deployment_requests.build(deployment_params)
    @deployment_request.company       = current_user.company.presence
    @deployment_request.contact_name  = current_user.name.presence || current_user.email
    @deployment_request.contact_email = current_user.email
    @deployment_request.monthly_price_cents = price_cents_for(@deployment_request)
    @deployment_request.status = "pending_payment"

    if @deployment_request.save
      checkout_url = build_checkout_url(@deployment_request)
      redirect_to checkout_url, allow_other_host: true
    else
      render :new, status: :unprocessable_entity
    end
  end

  private

  def deployment_params
    params.require(:private_deployment_request).permit(
      :server_tier, :billing_cycle, :admin_notes,
      selected_services: []
    )
  end

  def price_cents_for(request)
    if request.billing_cycle == "yearly"
      PrivateDeploymentRequest::TIER_PRICES_YEARLY[request.server_tier] || 0
    else
      PrivateDeploymentRequest::TIER_PRICES_MONTHLY[request.server_tier] || 0
    end
  end

  def build_checkout_url(deployment_request)
    tier = deployment_request.server_tier
    cycle = deployment_request.billing_cycle

    checkout_uuid = AppConfig.private_deployment_checkout_uuid_for(tier: tier, billing_cycle: cycle)

    store_slug = AppConfig.lemonsqueezy_store_slug
    base_url = "https://#{store_slug}.lemonsqueezy.com/checkout/buy/#{checkout_uuid}"

    query = {
      "checkout[email]"                                        => deployment_request.contact_email,
      "checkout[custom][user_id]"                              => deployment_request.user_id,
      "checkout[custom][private_deployment_request_id]"        => deployment_request.id,
      "checkout[redirect_url]"                                 => dashboard_root_url(private_deployment: "success")
    }

    "#{base_url}?#{query.to_query}"
  rescue AppConfig::InvalidConfigError => e
    Rails.logger.error "[PrivateDeploymentsController] Checkout UUID not configured: #{e.message}"
    # Fallback: redirect to sales inquiry with a notice
    talk_to_sales_url
  end
end
