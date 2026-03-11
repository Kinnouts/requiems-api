# frozen_string_literal: true

class Admin::ApiKeysController < ApplicationController
  before_action :authenticate_user!
  before_action :require_admin!
  before_action :set_api_key, only: [ :show, :revoke ]
  layout "admin"

  def index
    @api_keys = ApiKey.includes(:user).order(created_at: :desc)

    if params[:search].present?
      search_term = "%#{params[:search]}%"
      @api_keys = @api_keys.joins(:user)
        .where("api_keys.key_prefix ILIKE ? OR users.email ILIKE ? OR api_keys.name ILIKE ?",
               search_term, search_term, search_term)
    end

    @api_keys = @api_keys.where(active: params[:active] == "true") if params[:active].present?

    @pagy, @api_keys = pagy(@api_keys, items: 25)
  end

  def show
    @user = @api_key.user
    @recent_usage = @api_key.usage_logs.order(used_at: :desc).limit(20)
  end

  def revoke
    if @api_key.revoke!(reason: "Revoked by admin #{current_user.email}")
      redirect_to admin_api_key_path(@api_key), notice: "API key revoked successfully."
    else
      redirect_to admin_api_key_path(@api_key), alert: "Failed to revoke API key."
    end
  rescue StandardError => e
    redirect_to admin_api_key_path(@api_key), alert: "Failed to revoke API key: #{e.message}"
  end

  private

  def require_admin!
    unless current_user.admin?
      redirect_to root_path, alert: "Access denied. Admin privileges required."
    end
  end

  def set_api_key
    @api_key = ApiKey.find(params[:id])
  rescue ActiveRecord::RecordNotFound
    redirect_to admin_api_keys_path, alert: "API key not found."
  end
end
