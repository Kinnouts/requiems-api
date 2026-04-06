# frozen_string_literal: true

require "test_helper"

class SubscriptionTest < ActiveSupport::TestCase
  def setup
    @user  = create_user(email: "user@example.com")
    @admin = create_user(email: "admin@example.com", admin: true)
  end

  test "promoted? returns false when promoted_by_id is nil" do
    sub = Subscription.create!(user: @user, plan_name: "developer", status: "active")
    assert_not sub.promoted?
  end

  test "promoted? returns true when promoted_by_id is set" do
    sub = Subscription.create!(
      user: @user,
      plan_name: "developer",
      status: "active",
      promoted_by: @admin,
      promotion_reason: "Blog post",
      promotion_expires_at: 1.month.from_now
    )
    assert sub.promoted?
  end

  test "promotional scope returns only promoted subscriptions" do
    paid_sub = Subscription.create!(
      user: @user,
      plan_name: "developer",
      status: "active",
      lemonsqueezy_subscription_id: "ls_123"
    )

    other_user = create_user(email: "other@example.com")
    promo_sub = Subscription.create!(
      user: other_user,
      plan_name: "business",
      status: "active",
      promoted_by: @admin,
      promotion_reason: "Test",
      promotion_expires_at: 1.month.from_now
    )

    promotional = Subscription.promotional
    assert_includes promotional, promo_sub
    assert_not_includes promotional, paid_sub
  end
end
