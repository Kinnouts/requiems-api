# frozen_string_literal: true

class Admin::PrivateDeploymentsController < ApplicationController
  before_action :authenticate_user!
  before_action :require_admin!
  before_action :set_deployment_request, only: [:show, :activate, :cancel]
  layout "admin"

  def index
    @deployment_requests = PrivateDeploymentRequest.includes(:user)
                                                   .order(created_at: :desc)

    @deployment_requests = case params[:status]
    when "pending"   then @deployment_requests.pending
    when "deploying" then @deployment_requests.deploying
    when "active"    then @deployment_requests.active
    when "cancelled" then @deployment_requests.cancelled
    else @deployment_requests
    end
  end

  def show
  end

  def activate
    subdomain_slug = params[:subdomain_slug].to_s.strip.downcase
    tenant_secret  = params[:tenant_secret].to_s.strip
    admin_notes    = params[:admin_notes].to_s.strip

    if subdomain_slug.blank?
      redirect_to admin_private_deployment_path(@deployment_request), alert: "Subdomain slug is required." and return
    end

    if tenant_secret.blank?
      redirect_to admin_private_deployment_path(@deployment_request), alert: "Tenant secret is required." and return
    end

    @deployment_request.update!(
      subdomain_slug: subdomain_slug,
      tenant_secret: tenant_secret,
      admin_notes: admin_notes.presence,
      status: "active",
      deployed_at: Time.current
    )

    begin
      PrivateDeploymentMailer.deployment_ready(@deployment_request).deliver_later
    rescue StandardError => e
      Rails.logger.error "[Admin::PrivateDeploymentsController] Failed to enqueue deployment_ready email for request #{@deployment_request.id}: #{e.message}"
    end

    redirect_to admin_private_deployment_path(@deployment_request),
      notice: "Deployment marked as active. Confirmation email sent to #{@deployment_request.contact_email}."
  rescue ActiveRecord::RecordInvalid => e
    redirect_to admin_private_deployment_path(@deployment_request), alert: "Failed to activate: #{e.message}"
  end

  def cancel
    @deployment_request.update!(status: "cancelled")
    redirect_to admin_private_deployments_path, notice: "Deployment request cancelled."
  rescue ActiveRecord::RecordInvalid => e
    redirect_to admin_private_deployment_path(@deployment_request), alert: "Failed to cancel: #{e.message}"
  end

  private

  def require_admin!
    unless current_user.admin?
      redirect_to root_path, alert: "Access denied. Admin privileges required."
    end
  end

  def set_deployment_request
    @deployment_request = PrivateDeploymentRequest.find(params[:id])
  rescue ActiveRecord::RecordNotFound
    redirect_to admin_private_deployments_path, alert: "Deployment request not found."
  end
end
