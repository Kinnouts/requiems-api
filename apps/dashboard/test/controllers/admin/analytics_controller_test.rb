# frozen_string_literal: true

require "test_helper"

class Admin::AnalyticsControllerTest < ActionDispatch::IntegrationTest
  def setup
    @admin = create_user(email: "analytics_admin@example.com", admin: true)
    @user  = create_user(email: "analytics_user@example.com")
    @api_key = @admin.api_keys.create!(name: "Test Key", environment: "test")

    sign_in @admin
  end

  def create_log(endpoint: "/v1/test", status_code: 200, used_at: 1.day.ago, response_time_ms: 50, credits_used: 1)
    UsageLog.create!(
      user: @admin,
      api_key: @api_key,
      endpoint: endpoint,
      status_code: status_code,
      used_at: used_at,
      response_time_ms: response_time_ms,
      credits_used: credits_used
    )
  end

  # ── usage ──────────────────────────────────────────────────────────────────

  test "usage requires authentication" do
    sign_out @admin
    get admin_analytics_usage_path
    assert_redirected_to new_user_session_path
  end

  test "usage requires admin" do
    sign_out @admin
    sign_in @user
    get admin_analytics_usage_path
    assert_response :not_found
  end

  test "usage renders successfully with no data" do
    get admin_analytics_usage_path
    assert_response :success
  end

  test "usage renders successfully with usage data" do
    create_log(endpoint: "/v1/email/disposable", status_code: 200)
    create_log(endpoint: "/v1/text/advice",      status_code: 200)
    create_log(endpoint: "/v1/email/disposable", status_code: 429, used_at: 2.days.ago)

    get admin_analytics_usage_path
    assert_response :success
  end

  test "usage accepts date_range param" do
    create_log(used_at: 5.days.ago)
    create_log(used_at: 40.days.ago)

    get admin_analytics_usage_path, params: { date_range: "7" }
    assert_response :success
  end

  test "usage groups by endpoint without sql errors" do
    5.times { |i| create_log(endpoint: "/v1/endpoint/#{i}") }

    # This hit COUNT(*) DESC without Arel.sql and raised ActiveRecord::UnknownAttributeReference in prod
    assert_nothing_raised { get admin_analytics_usage_path }
    assert_response :success
  end

  # ── revenue ────────────────────────────────────────────────────────────────

  test "revenue requires authentication" do
    sign_out @admin
    get admin_analytics_revenue_path
    assert_redirected_to new_user_session_path
  end

  test "revenue requires admin" do
    sign_out @admin
    sign_in @user
    get admin_analytics_revenue_path
    assert_response :not_found
  end

  test "revenue renders successfully" do
    get admin_analytics_revenue_path
    assert_response :success
  end

  # ── system_health ──────────────────────────────────────────────────────────

  test "system_health requires authentication" do
    sign_out @admin
    get admin_analytics_system_health_path
    assert_redirected_to new_user_session_path
  end

  test "system_health requires admin" do
    sign_out @admin
    sign_in @user
    get admin_analytics_system_health_path
    assert_response :not_found
  end

  test "system_health renders successfully with no data" do
    get admin_analytics_system_health_path
    assert_response :success
  end

  test "system_health renders with usage data across status codes" do
    create_log(status_code: 200, response_time_ms: 40)
    create_log(status_code: 200, response_time_ms: 80)
    create_log(status_code: 401, response_time_ms: 10)
    create_log(status_code: 500, response_time_ms: 200)

    # Exercises percentile_cont queries — would raise without Arel.sql wrapping
    assert_nothing_raised { get admin_analytics_system_health_path }
    assert_response :success
  end

  test "system_health accepts time_range param" do
    %w[1h 24h 7d 30d].each do |range|
      get admin_analytics_system_health_path, params: { time_range: range }
      assert_response :success, "Failed for time_range=#{range}"
    end
  end

  test "system_health 1h shows per-minute breakdown" do
    create_log(used_at: 30.minutes.ago)

    get admin_analytics_system_health_path, params: { time_range: "1h" }
    assert_response :success
  end
end
