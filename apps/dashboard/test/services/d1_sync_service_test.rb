# frozen_string_literal: true

require "test_helper"

class D1SyncServiceTest < ActiveSupport::TestCase
  test "bulk_insert maps request method and telemetry fields" do
    user = create_user(email: "d1-sync-user@example.com")
    api_key = user.api_keys.create!(name: "Sync Key", environment: "test")

    record = {
      api_key: api_key.key_prefix + "rest",
      endpoint: "/v1/text/advice",
      credits_used: 2,
      request_method: "PATCH",
      status_code: 503,
      response_time_ms: 128,
      used_at: "2026-04-19T12:00:00Z"
    }

    inserted_values = nil

    relation = Object.new
    relation.define_singleton_method(:pluck) do |*|
      [ [ api_key.key_prefix, api_key.id, user.id ] ]
    end

    UsageLog.stub(:insert_all, ->(values, unique_by:) {
      inserted_values = values
      assert_equal [ :api_key_id, :used_at, :endpoint ], unique_by
    }) do
      ApiKey.stub(:where, relation) do
        result = D1SyncService.new.bulk_insert([ record ])

        assert_equal 1, result
      end
    end

    assert_equal 1, inserted_values.size
    assert_equal "PATCH", inserted_values.first[:request_method]
    assert_equal 503, inserted_values.first[:status_code]
    assert_equal 128, inserted_values.first[:response_time_ms]
  end
end