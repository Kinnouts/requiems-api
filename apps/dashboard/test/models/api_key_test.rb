require "test_helper"

class ApiKeyTest < ActiveSupport::TestCase
  def setup
    @user = User.create!(
      email: "test@example.com",
      password: "password123",
      password_confirmation: "password123"
    )

    @api_key = @user.api_keys.create!(
      name: "Test Key",
      key: "rq_test_" + SecureRandom.hex(32),
      prefix: "rq_test_abc123",
      environment: "test"
    )
  end

  test "valid api key with required attributes" do
    assert @api_key.valid?
    assert @api_key.persisted?
  end

  test "requires name" do
    api_key = @user.api_keys.build(
      key: "rq_test_" + SecureRandom.hex(32),
      prefix: "rq_test_xyz",
      environment: "test"
    )
    api_key.name = nil

    assert_not api_key.valid?
    assert_includes api_key.errors[:name], "can't be blank"
  end

  test "requires key" do
    api_key = @user.api_keys.build(
      name: "Test",
      prefix: "rq_test_xyz",
      environment: "test"
    )
    api_key.key = nil

    assert_not api_key.valid?
    assert_includes api_key.errors[:key], "can't be blank"
  end

  test "requires unique key" do
    duplicate_key = @user.api_keys.build(
      name: "Duplicate",
      key: @api_key.key,
      prefix: "rq_test_xyz",
      environment: "test"
    )

    assert_not duplicate_key.valid?
    assert_includes duplicate_key.errors[:key], "has already been taken"
  end

  test "requires environment" do
    api_key = @user.api_keys.build(
      name: "Test",
      key: "rq_test_" + SecureRandom.hex(32),
      prefix: "rq_test_xyz"
    )
    api_key.environment = nil

    assert_not api_key.valid?
    assert_includes api_key.errors[:environment], "can't be blank"
  end

  test "validates environment is test or live" do
    @api_key.environment = "invalid"
    assert_not @api_key.valid?

    @api_key.environment = "test"
    assert @api_key.valid?

    @api_key.environment = "live"
    assert @api_key.valid?
  end

  test "belongs to user" do
    assert_equal @user, @api_key.user
  end

  test "has many usage_logs" do
    assert_respond_to @api_key, :usage_logs
  end

  test "active_keys scope returns non-revoked keys" do
    active_key = @user.api_keys.create!(
      name: "Active",
      key: "rq_test_" + SecureRandom.hex(32),
      prefix: "rq_test_active",
      environment: "test"
    )

    revoked_key = @user.api_keys.create!(
      name: "Revoked",
      key: "rq_test_" + SecureRandom.hex(32),
      prefix: "rq_test_revoked",
      environment: "test",
      revoked_at: Time.current
    )

    active_keys = ApiKey.active_keys

    assert_includes active_keys, active_key
    assert_not_includes active_keys, revoked_key
  end

  test "revoked scope returns revoked keys" do
    revoked_key = @user.api_keys.create!(
      name: "Revoked",
      key: "rq_test_" + SecureRandom.hex(32),
      prefix: "rq_test_revoked",
      environment: "test",
      revoked_at: Time.current
    )

    revoked_keys = ApiKey.revoked

    assert_includes revoked_keys, revoked_key
    assert_not_includes revoked_keys, @api_key
  end

  test "revoke! sets revoked_at and reason" do
    assert_nil @api_key.revoked_at
    assert_nil @api_key.revoked_reason

    @api_key.revoke!(reason: "User requested")

    @api_key.reload
    assert_not_nil @api_key.revoked_at
    assert_equal "User requested", @api_key.revoked_reason
  end

  test "generates prefix from key on creation" do
    new_key = @user.api_keys.create!(
      name: "New Key",
      key: "rq_test_abc123def456",
      environment: "test"
    )

    assert_equal "rq_test_abc123", new_key.prefix
  end
end
