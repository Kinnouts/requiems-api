# frozen_string_literal: true

# Background job to sync usage data from Cloudflare D1 to PostgreSQL
#
# This job runs every 5 minutes via Solid Queue recurring tasks.
# It fetches new usage records from D1 and inserts them into the usage_logs table.
#
# Usage:
#   SyncD1UsageJob.perform_later
#
class SyncD1UsageJob < ApplicationJob
  queue_as :default

  retry_on D1SyncService::Error, attempts: 5, wait: :polynomially_longer do |job, error|
    Rails.logger.error(
      "SyncD1UsageJob permanently failed after retries: #{error.message}",
      job_id: job.job_id
    )
  end

  LAST_SYNC_KEY = "d1_sync:last_sync_at"
  SYNC_INTERVAL = 5.minutes

  def perform
    start_time = Time.current
    last_sync_at = get_last_sync_time

    Rails.logger.info("Starting D1 usage sync from #{last_sync_at}")

    service = D1SyncService.new
    result = service.fetch_usage(since: last_sync_at)

    if result[:usage].empty?
      Rails.logger.info("No new usage records to sync")
      return
    end

    inserted = service.bulk_insert(result[:usage])

    # Update last sync timestamp
    set_last_sync_time(start_time)

    Rails.logger.info("D1 sync completed: #{inserted} records inserted")

    # Track metrics
    track_sync_metrics(inserted, start_time)
  rescue D1SyncService::Error => e
    Rails.logger.error("D1 sync failed: #{e.message}")
    # Don't update last_sync_time on failure - will retry from same point
    raise
  end

  private

  def get_last_sync_time
    timestamp = Rails.cache.read(LAST_SYNC_KEY)
    return Time.parse(timestamp) if timestamp

    # Fallback: get timestamp of most recent usage log
    last_log = UsageLog.order(used_at: :desc).first
    last_log ? last_log.used_at : 1.hour.ago
  end

  def set_last_sync_time(time)
    Rails.cache.write(LAST_SYNC_KEY, time.iso8601, expires_in: 7.days)
  end

  def track_sync_metrics(inserted, start_time)
    duration = (Time.current - start_time).round(2)

    Rails.logger.info(
      "D1 sync metrics",
      records_inserted: inserted,
      duration_seconds: duration,
      records_per_second: (inserted / duration).round(2)
    )
  end
end
