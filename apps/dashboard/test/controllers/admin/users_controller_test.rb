require "test_helper"

class Admin::UsersControllerTest < ActionDispatch::IntegrationTest
  def setup
    @admin = User.create!(
      email: "admin@example.com",
      password: "password123",
      password_confirmation: "password123",
      admin: true
    )

    @regular_user = User.create!(
      email: "user@example.com",
      password: "password123",
      password_confirmation: "password123"
    )

    sign_in @admin
  end

  test "index requires admin authentication" do
    sign_out @admin
    sign_in @regular_user

    get admin_users_path

    assert_redirected_to root_path
    assert_equal "Access denied. Admin privileges required.", flash[:alert]
  end

  test "index shows all users" do
    get admin_users_path

    assert_response :success
    assert_select "h1", text: /User Management/i
    assert_match @admin.email, response.body
    assert_match @regular_user.email, response.body
  end

  test "index filters by search term" do
    get admin_users_path, params: { search: @regular_user.email }

    assert_response :success
    assert_match @regular_user.email, response.body
  end

  test "index filters by plan" do
    subscription = Subscription.create!(
      user: @regular_user,
      plan_name: "developer",
      billing_cycle: "monthly"
    )

    get admin_users_path, params: { plan: "developer" }

    assert_response :success
    assert_match @regular_user.email, response.body
  end

  test "index filters by status" do
    @regular_user.update(admin: true)

    get admin_users_path, params: { status: "admin" }

    assert_response :success
    assert_match @regular_user.email, response.body
  end

  test "show displays user details" do
    get admin_user_path(@regular_user)

    assert_response :success
    assert_select "h1", text: /User Details/i
    assert_match @regular_user.email, response.body
  end

  test "suspend suspends user and revokes api keys" do
    api_key = @regular_user.api_keys.create!(
      name: "Test Key",
      key: "rq_test_" + SecureRandom.hex(32),
      prefix: "rq_test_abc123",
      environment: "test"
    )

    post suspend_admin_user_path(@regular_user)

    @regular_user.reload
    assert @regular_user.suspended?
    assert_not_nil @regular_user.suspended_at

    api_key.reload
    assert_not_nil api_key.revoked_at
    assert_equal "User suspended by admin", api_key.revoked_reason

    assert_redirected_to admin_user_path(@regular_user)
    assert_match /suspended successfully/i, flash[:notice]
  end

  test "unsuspend removes suspension" do
    @regular_user.update(suspended: true, suspended_at: Time.current)

    post unsuspend_admin_user_path(@regular_user)

    @regular_user.reload
    assert_not @regular_user.suspended?
    assert_nil @regular_user.suspended_at

    assert_redirected_to admin_user_path(@regular_user)
    assert_match /unsuspended successfully/i, flash[:notice]
  end

  test "ban bans user, revokes keys, and cancels subscription" do
    subscription = Subscription.create!(
      user: @regular_user,
      plan_name: "developer",
      billing_cycle: "monthly"
    )

    api_key = @regular_user.api_keys.create!(
      name: "Test Key",
      key: "rq_test_" + SecureRandom.hex(32),
      prefix: "rq_test_abc123",
      environment: "test"
    )

    post ban_admin_user_path(@regular_user)

    @regular_user.reload
    assert @regular_user.banned?
    assert @regular_user.suspended?
    assert_not_nil @regular_user.banned_at

    api_key.reload
    assert_not_nil api_key.revoked_at

    subscription.reload
    assert subscription.cancel_at_period_end?

    assert_redirected_to admin_user_path(@regular_user)
    assert_match /banned successfully/i, flash[:notice]
  end

  test "make_admin grants admin privileges" do
    assert_not @regular_user.admin?

    post make_admin_admin_user_path(@regular_user)

    @regular_user.reload
    assert @regular_user.admin?

    assert_redirected_to admin_user_path(@regular_user)
    assert_match /granted admin privileges/i, flash[:notice]
  end

  test "remove_admin removes admin privileges" do
    @regular_user.update(admin: true)

    post remove_admin_admin_user_path(@regular_user)

    @regular_user.reload
    assert_not @regular_user.admin?

    assert_redirected_to admin_user_path(@regular_user)
    assert_match /removed from user/i, flash[:notice]
  end

  test "cannot remove own admin privileges" do
    post remove_admin_admin_user_path(@admin)

    @admin.reload
    assert @admin.admin?

    assert_redirected_to admin_user_path(@admin)
    assert_match /cannot remove your own admin privileges/i, flash[:alert]
  end

  test "non-admin cannot access admin actions" do
    sign_out @admin
    sign_in @regular_user

    post suspend_admin_user_path(@regular_user)

    assert_redirected_to root_path
    assert_equal "Access denied. Admin privileges required.", flash[:alert]
  end
end
