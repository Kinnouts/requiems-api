# frozen_string_literal: true

class Dashboard::SettingsController < ApplicationController
  before_action :authenticate_user!
  layout "dashboard"

  def show
    # Settings page - view current settings
  end

  def update
    # Update account settings
    if current_user.update(user_params)
      # If email was changed and Devise confirmable is enabled, send confirmation
      if user_params[:email].present? && current_user.email != current_user.email_was
        flash[:notice] = "Settings updated. Please check your email to confirm your new address."
      else
        flash[:notice] = "Settings updated successfully."
      end

      redirect_to dashboard_settings_path
    else
      flash.now[:alert] = "Failed to update settings. Please check the form for errors."
      render :show, status: :unprocessable_entity
    end
  end

  def request_deletion
    reason = params[:deletion_reason].to_s.strip

    if reason.length < 10
      redirect_to dashboard_settings_path, alert: "Please provide a reason (at least 10 characters)."
      return
    end

    current_user.request_account_deletion!(reason)
    AccountDeletionMailer.confirmation(current_user).deliver_later

    redirect_to dashboard_settings_path,
      notice: "Check your email — we sent a confirmation link. It expires in 1 hour."
  end

  def confirm_deletion
    token = params[:token].to_s

    unless current_user.deletion_token_valid?(token)
      redirect_to dashboard_settings_path, alert: "This link is invalid or has expired. Please request a new one."
      return
    end

    @token = token
    @reason = current_user.deletion_reason
  end

  def execute_deletion
    token = params[:token].to_s

    unless current_user.deletion_token_valid?(token)
      redirect_to dashboard_settings_path, alert: "This link is invalid or has expired. Please request a new one."
      return
    end

    # Revoke all API keys
    current_user.api_keys.each do |key|
      key.revoke!(reason: "Account deleted by user")
    end

    # Cancel subscription if exists
    if current_user.subscription
      current_user.subscription.update(
        cancel_at_period_end: true,
        canceled_at: Time.current
      )
    end

    current_user.destroy
    sign_out current_user

    redirect_to root_path, notice: "Your account has been permanently deleted. We're sorry to see you go."
  end

  private

  def user_params
    p = params.require(:user).permit(:email, :name, :company, :locale, :email_notifications, :usage_alerts, :weekly_reports)
    p[:locale] = p[:locale].presence # convert "" (auto-detect) to nil
    p
  end
end
