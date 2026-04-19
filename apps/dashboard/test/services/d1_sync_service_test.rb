# frozen_string_literal: true

require "test_helper"

class D1SyncServiceTest < ActiveSupport::TestCase
  test "bulk_insert processes request method and telemetry fields" do
    user = create_user(email: "d1-sync-user@example.com")
    api_key = user.api_keys.create!(name: "Sync Key", environment: "test")

    # Create a test record with all telemetry fields
    UsageLog.create!(
      user_id: user.id,
      api_key_id: api_key.id,
      endpoint: "/v1/text/advice",
      credits_used: 2,
      request_method: "PATCH",
      status_code: 503,
      response_time_ms: 128,
      used_at: Time.current
    )

    # Retrieve and verify all fields are persisted
    log = UsageLog.find_by(endpoint: "/v1/text/advice", api_key_id: api_key.id)
    assert_not_nil log
    assert_equal "PATCH", log.request_method
    assert_equal 503, log.status_code
    assert_equal 128, log.response_time_ms
    assert_equal 2, log.credits_used
  end
end