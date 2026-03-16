# frozen_string_literal: true

require "test_helper"

class UsageLogTest < ActiveSupport::TestCase
  def setup
    @user = create_user(email: "usage_log_test@example.com")
    @api_key = @user.api_keys.create!(name: "Test Key", environment: "test")
  end

  def create_log(used_at:, status_code:)
    UsageLog.create!(
      user: @user,
      api_key: @api_key,
      endpoint: "/v1/test",
      status_code: status_code,
      used_at: used_at
    )
  end

  test "in_range returns logs within the date range" do
    inside  = create_log(used_at: 5.days.ago,  status_code: 200)
    outside = create_log(used_at: 40.days.ago, status_code: 200)

    results = UsageLog.in_range(30.days.ago, Time.current)

    assert_includes     results, inside
    assert_not_includes results, outside
  end

  test "successful returns logs with 2xx status codes" do
    ok      = create_log(used_at: 1.day.ago, status_code: 200)
    created = create_log(used_at: 1.day.ago, status_code: 201)
    error   = create_log(used_at: 1.day.ago, status_code: 422)

    results = UsageLog.successful

    assert_includes     results, ok
    assert_includes     results, created
    assert_not_includes results, error
  end

  test "with_errors returns logs with status_code >= 400" do
    ok        = create_log(used_at: 1.day.ago, status_code: 200)
    not_found = create_log(used_at: 1.day.ago, status_code: 404)
    server    = create_log(used_at: 1.day.ago, status_code: 500)

    results = UsageLog.with_errors

    assert_not_includes results, ok
    assert_includes     results, not_found
    assert_includes     results, server
  end

  test "scopes can be chained" do
    inside_ok    = create_log(used_at: 5.days.ago,  status_code: 200)
    inside_error = create_log(used_at: 5.days.ago,  status_code: 500)
    outside_ok   = create_log(used_at: 40.days.ago, status_code: 200)

    results = UsageLog.in_range(30.days.ago, Time.current).with_errors

    assert_includes     results, inside_error
    assert_not_includes results, inside_ok
    assert_not_includes results, outside_ok
  end
end
