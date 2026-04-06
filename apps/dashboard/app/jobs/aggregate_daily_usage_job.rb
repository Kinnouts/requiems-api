# frozen_string_literal: true

# Background job to aggregate daily usage summaries
#
# This job runs once per day at 00:05 UTC via Sidekiq Cron.
# It aggregates usage_logs data into daily_usage_summaries for faster analytics queries.
#
# Usage:
#   AggregateDailyUsageJob.perform_later(date: Date.yesterday)
#
class AggregateDailyUsageJob < ApplicationJob
  queue_as :default

  retry_on StandardError, attempts: 3, wait: :polynomially_longer do |job, error|
    Rails.logger.error(
      "AggregateDailyUsageJob permanently failed after retries: #{error.message}",
      job_id: job.job_id
    )
  end

  BATCH_SIZE = 1000
  private_constant :BATCH_SIZE

  def perform(date: Date.yesterday)
    Rails.logger.info("Starting daily usage aggregation for #{date}")

    start_time = Time.current
    inserted_count = 0
    now = Time.current

    # Pluck distinct user IDs (integers only) to keep initial memory low,
    # then aggregate and upsert in bounded chunks.
    # IMPORTANT: Group by user_id only (not api_key_id) because all API keys
    # for a user share the same quota
    UsageLog.where(usage_date: date).distinct.pluck(:user_id).each_slice(BATCH_SIZE) do |user_ids|
      records = UsageLog
        .where(usage_date: date, user_id: user_ids)
        .group(:user_id, :usage_date)
        .select(
          "user_id",
          "usage_date as date",
          "COUNT(*) as total_requests",
          "SUM(credits_used) as total_credits",
          "AVG(response_time_ms) as avg_response_time_ms",
          "COUNT(CASE WHEN status_code >= 400 THEN 1 END) as error_count"
        )
        .map do |summary|
          {
            user_id: summary.user_id,
            date: summary.date,
            total_requests: summary.total_requests,
            total_credits: summary.total_credits,
            updated_at: now
          }
        end

      DailyUsageSummary.upsert_all(records, unique_by: [ :user_id, :date ]) if records.any?
      inserted_count += records.size
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
