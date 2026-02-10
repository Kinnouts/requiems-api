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

  def account
    # Delete account action
    if params[:confirm_email] == current_user.email
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

      # Mark user as deleted (soft delete or hard delete based on preference)
      current_user.destroy

      # Sign out
      sign_out current_user

      redirect_to root_path, notice: "Your account has been successfully deleted. We're sorry to see you go!"
    else
      redirect_to dashboard_settings_path, alert: "Email confirmation did not match. Account was not deleted."
    end
  end

  private

  def user_params
    params.require(:user).permit(:email, :name, :company, :email_notifications, :usage_alerts, :weekly_reports)
  end
end
