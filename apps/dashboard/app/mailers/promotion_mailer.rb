# frozen_string_literal: true

class PromotionMailer < ApplicationMailer
  # Notify a user that an admin has upgraded their plan.
  #
  # @param user [User] the user who was promoted
  # @param plan_name [String] the plan they were upgraded to (e.g. "developer")
  # @param expires_at [Time] when the promotion expires
  # @param reason [String] the reason given by the admin
  def upgrade_notification(user, plan_name, expires_at, reason)
    @user = user
    @plan_name = plan_name
    @expires_at = expires_at
    @reason = reason
    @dashboard_url = dashboard_root_url

    mail(
      to: @user.email,
      subject: t("dashboard.mailers.promotion.upgrade_notification.subject",
                 plan: plan_name.titleize)
    )
  end
end
