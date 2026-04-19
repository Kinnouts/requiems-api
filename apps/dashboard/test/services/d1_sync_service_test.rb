# frozen_string_literal: true

require "test_helper"

class D1SyncServiceTest < ActiveSupport::TestCase
  test "bulk_insert persists telemetry fields from D1 export format" do
    user = create_user(email: "d1-sync-user@example.com")
    api_key = user.api_keys.create!(name: "Sync Key", environment: "test")

    d1_record = {
      api_key: "#{api_key.key_prefix}rest",
      endpoint: "/v1/text/advice",
      credits_used: 2,
      request_method: "PATCH",
      status_code: 503,
      response_time_ms: 128,
      used_at: Time.current.iso8601
    }

    service = D1SyncService.new
    result = service.bulk_insert([ d1_record ])

    assert_equal 1, result

    # Retrieve and verify all telemetry fields are persisted
    log = UsageLog.find_by(endpoint: "/v1/text/advice", api_key_id: api_key.id)
    assert_not_nil log
    assert_equal "PATCH", log.request_method
    assert_equal 503, log.status_code
    assert_equal 128, log.response_time_ms
    assert_equal 2, log.credits_used
  end
end
