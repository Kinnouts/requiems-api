# frozen_string_literal: true

class SubscriptionMailer < ApplicationMailer
  # Notify a user that their paid subscription has been activated or upgraded.
  #
  # @param user [User] the user whose plan changed
  # @param plan_name [String] the new plan name (e.g. "developer")
  def upgrade_notification(user, plan_name)
    @user = user
    @plan_name = plan_name
    @dashboard_url = dashboard_root_url

    mail(
      to: @user.email,
      subject: t("dashboard.mailers.subscription.upgrade_notification.subject",
                 plan: plan_name.titleize)
    )
  end
end
