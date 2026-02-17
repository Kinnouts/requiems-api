require "test_helper"

class UserTest < ActiveSupport::TestCase
  def setup
    @user = create_user(
      email: "test@example.com",
      password: "password123",
      password_confirmation: "password123"
    )
  end

  test "valid user with email and password" do
    assert @user.valid?
    assert @user.persisted?
  end

  test "requires email" do
    user = User.new(password: "password123")
    assert_not user.valid?
    assert_includes user.errors[:email], "can't be blank"
  end

  test "requires valid email format" do
    user = User.new(email: "invalid", password: "password123")
    assert_not user.valid?
  end

  test "requires unique email" do
    duplicate_user = User.new(
      email: @user.email,
      password: "password123"
    )
    assert_not duplicate_user.valid?
    assert_includes duplicate_user.errors[:email], "has already been taken"
  end

  test "admin? returns admin status" do
    assert_not @user.admin?

    @user.update(admin: true)
    assert @user.admin?
  end

  test "suspended? returns suspension status" do
    assert_not @user.suspended?

    @user.update(status: "suspended")
    assert @user.suspended?
  end

  test "banned? returns ban status" do
    assert_not @user.banned?

    @user.update(status: "banned", banned_at: Time.current)
    assert @user.banned?
  end

  test "has many api_keys" do
    assert_respond_to @user, :api_keys

    api_key = @user.api_keys.create!(
      name: "Test Key",
      environment: "test"
    )

    assert_includes @user.api_keys, api_key
  end

  test "has one subscription" do
    assert_respond_to @user, :subscription
  end

  test "has many usage_logs" do
    assert_respond_to @user, :usage_logs
  end

  test "destroys dependent api_keys when destroyed" do
    api_key = @user.api_keys.create!(
      name: "Test Key",
      environment: "test"
    )

    assert_difference "ApiKey.count", -1 do
      @user.destroy
    end
  end
end
