# frozen_string_literal: true

class AccountDeletionMailer < ApplicationMailer
  # Send a time-limited confirmation link for account deletion.
  #
  # @param user [User] the user requesting deletion
  def confirmation(user)
    @user = user
    @confirm_url = confirm_deletion_dashboard_settings_url(token: user.deletion_token)
    @expires_at = user.deletion_token_sent_at + 1.hour
    @reason = user.deletion_reason

    mail(
      to: @user.email,
      subject: "Confirm account deletion — Requiems API"
    )
  end
end
