# frozen_string_literal: true

require "test_helper"

class Admin::PromotionsControllerTest < ActionDispatch::IntegrationTest
  def setup
    @admin = create_user(email: "admin@example.com", admin: true)
    @user  = create_user(email: "user@example.com")
    sign_in @admin
  end

  # ── create ──────────────────────────────────────────────────────────────────

  test "create grants promotion and creates subscription" do
    assert_nil @user.subscription

    post admin_user_promotion_path(@user), params: {
      promotion: {
        plan_name: "developer",
        expires_at: 1.month.from_now.to_date.to_s,
        reason: "YouTube video collaboration"
      }
    }

    assert_redirected_to admin_user_path(@user)
    assert_match /upgraded to the developer/i, flash[:notice]

    @user.reload
    sub = @user.subscription
    assert_not_nil sub
    assert_equal "developer", sub.plan_name
    assert_equal "active",    sub.status
    assert_equal @admin,      sub.promoted_by
    assert_equal "YouTube video collaboration", sub.promotion_reason
    assert sub.promoted?
    assert sub.promotion_expires_at > Time.current
  end

  test "create updates existing free subscription" do
    sub = Subscription.create!(user: @user, plan_name: "free", status: "active")

    post admin_user_promotion_path(@user), params: {
      promotion: {
        plan_name: "business",
        expires_at: 3.months.from_now.to_date.to_s,
        reason: "Blog post partner"
      }
    }

    sub.reload
    assert_equal "business", sub.plan_name
    assert sub.promoted?
  end

  test "create writes an audit log entry" do
    assert_difference "AuditLog.count", 1 do
      post admin_user_promotion_path(@user), params: {
        promotion: {
          plan_name: "developer",
          expires_at: 1.month.from_now.to_date.to_s,
          reason: "Test"
        }
      }
    end

    log = AuditLog.last
    assert_equal "promotion_granted", log.action
    assert_equal @user.id,  log.user_id
    assert_equal @admin.id, log.admin_user_id
  end

  test "create rejects missing reason" do
    assert_no_difference "Subscription.count" do
      post admin_user_promotion_path(@user), params: {
        promotion: {
          plan_name: "developer",
          expires_at: 1.month.from_now.to_date.to_s,
          reason: ""
        }
      }
    end

    assert_redirected_to admin_user_path(@user)
    assert_match /reason is required/i, flash[:alert]
  end

  test "create rejects missing expiry date" do
    assert_no_difference "Subscription.count" do
      post admin_user_promotion_path(@user), params: {
        promotion: {
          plan_name: "developer",
          expires_at: "",
          reason: "Test"
        }
      }
    end

    assert_redirected_to admin_user_path(@user)
    assert_match /expiry date must be in the future/i, flash[:alert]
  end

  test "create rejects past expiry date" do
    assert_no_difference "Subscription.count" do
      post admin_user_promotion_path(@user), params: {
        promotion: {
          plan_name: "developer",
          expires_at: 1.day.ago.to_date.to_s,
          reason: "Test"
        }
      }
    end

    assert_redirected_to admin_user_path(@user)
    assert_match /expiry date must be in the future/i, flash[:alert]
  end

  test "create rejects invalid plan" do
    assert_no_difference "Subscription.count" do
      post admin_user_promotion_path(@user), params: {
        promotion: {
          plan_name: "free",
          expires_at: 1.month.from_now.to_date.to_s,
          reason: "Test"
        }
      }
    end

    assert_redirected_to admin_user_path(@user)
    assert_match /invalid plan/i, flash[:alert]
  end

  test "non-admin cannot create promotion" do
    sign_out @admin
    sign_in @user

    post admin_user_promotion_path(@user), params: {
      promotion: { plan_name: "developer", expires_at: 1.month.from_now.to_date.to_s, reason: "Test" }
    }

    assert_response :not_found
  end

  # ── destroy ─────────────────────────────────────────────────────────────────

  test "destroy revokes active promotion" do
    sub = Subscription.create!(
      user: @user,
      plan_name: "developer",
      status: "active",
      promoted_by: @admin,
      promotion_reason: "Test",
      promotion_expires_at: 1.month.from_now
    )

    delete admin_user_promotion_path(@user)

    assert_redirected_to admin_user_path(@user)
    assert_match /revoked/i, flash[:notice]

    sub.reload
    assert_equal "free", sub.plan_name
    assert_not sub.promoted?
    assert_nil sub.promotion_expires_at
    assert_nil sub.promotion_reason
  end

  test "destroy writes an audit log entry" do
    Subscription.create!(
      user: @user,
      plan_name: "developer",
      status: "active",
      promoted_by: @admin,
      promotion_reason: "Test",
      promotion_expires_at: 1.month.from_now
    )

    assert_difference "AuditLog.count", 1 do
      delete admin_user_promotion_path(@user)
    end

    assert_equal "promotion_revoked", AuditLog.last.action
  end

  test "destroy returns alert when no active promotion" do
    delete admin_user_promotion_path(@user)

    assert_redirected_to admin_user_path(@user)
    assert_match /no active promotion/i, flash[:alert]
  end
end
