# frozen_string_literal: true

# Background job to aggregate daily usage summaries
#
# This job runs once per day at 00:05 UTC via Solid Queue recurring tasks.
# It aggregates usage_logs data into daily_usage_summaries for faster analytics queries.
#
# Usage:
#   AggregateDailyUsageJob.perform_later(date: Date.yesterday)
#
class AggregateDailyUsageJob < ApplicationJob
  queue_as :default

  def perform(date: Date.yesterday)
    Rails.logger.info("Starting daily usage aggregation for #{date}")

    start_time = Time.current

    # Aggregate by user and date
    # IMPORTANT: Group by user_id only (not api_key_id) because all API keys
    # for a user share the same quota
    aggregated = UsageLog
      .where(usage_date: date)
      .group(:user_id, :usage_date)
      .select(
        "user_id",
        "usage_date as date",
        "COUNT(*) as total_requests",
        "SUM(credits_used) as total_credits",
        "AVG(response_time_ms) as avg_response_time_ms",
        "COUNT(CASE WHEN status_code >= 400 THEN 1 END) as error_count"
      )

    inserted_count = 0

    aggregated.each do |summary|
      DailyUsageSummary.upsert(
        {
          user_id: summary.user_id,
          date: summary.date,
          total_requests: summary.total_requests,
          total_credits: summary.total_credits,
          updated_at: Time.current
        },
        unique_by: [ :user_id, :date ]
      )
      inserted_count += 1
    end

    duration = (Time.current - start_time).round(2)

    Rails.logger.info(
      "Daily aggregation completed",
      date: date,
      summaries_created: inserted_count,
      duration_seconds: duration
    )

    inserted_count
  end
end
