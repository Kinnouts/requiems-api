# frozen_string_literal: true

namespace :usage do
  desc "Backfill usage data from Cloudflare D1 to PostgreSQL"
  task backfill_from_d1: :environment do
    puts "Starting D1 usage backfill..."
    puts "This will fetch ALL historical usage data from Cloudflare D1"
    puts ""

    # Get the earliest date we want to backfill from
    # Default: 90 days ago or from the beginning of D1 data
    since = ENV.fetch("SINCE", 90.days.ago.iso8601)

    puts "Fetching usage data since: #{since}"
    puts ""

    service = D1SyncService.new

    begin
      result = service.fetch_usage(since: since)

      puts "Fetched #{result[:total]} usage records from D1"
      puts ""

      if result[:usage].empty?
        puts "No usage data to backfill"
        exit 0
      end

      puts "Inserting records into PostgreSQL..."
      inserted = service.bulk_insert(result[:usage])

      puts ""
      puts "✓ Backfill completed successfully!"
      puts "  - Records fetched: #{result[:total]}"
      puts "  - Records inserted: #{inserted}"
      puts ""
      puts "Next steps:"
      puts "  1. Run: rake usage:aggregate_daily to create daily summaries"
      puts "  2. The sync job will continue running every 5 minutes automatically"
    rescue D1SyncService::Error => e
      puts ""
      puts "✗ Backfill failed: #{e.message}"
      puts ""
      puts "Please check:"
      puts "  - CLOUDFLARE_WORKER_URL environment variable is set"
      puts "  - BACKEND_SECRET environment variable is set"
      puts "  - Cloudflare Worker is deployed and accessible"
      exit 1
    end
  end

  desc "Aggregate daily usage summaries for a date range"
  task aggregate_daily: :environment do
    # Get date range from environment variables
    start_date = Date.parse(ENV.fetch("START_DATE", 30.days.ago.to_date.to_s))
    end_date = Date.parse(ENV.fetch("END_DATE", Date.yesterday.to_s))

    puts "Aggregating daily usage summaries"
    puts "Date range: #{start_date} to #{end_date}"
    puts ""

    total_summaries = 0

    (start_date..end_date).each do |date|
      print "Processing #{date}... "

      count = AggregateDailyUsageJob.new.perform(date: date)
      total_summaries += count

      puts "#{count} summaries created"
    end

    puts ""
    puts "✓ Daily aggregation completed!"
    puts "  - Days processed: #{(end_date - start_date).to_i + 1}"
    puts "  - Total summaries created: #{total_summaries}"
  end

  desc "Show usage sync status"
  task status: :environment do
    puts "Usage Sync Status"
    puts "=" * 60
    puts ""

    # Check last sync time
    last_sync_key = "d1_sync:last_sync_at"
    last_sync = Rails.cache.read(last_sync_key)

    if last_sync
      puts "Last D1 sync: #{Time.parse(last_sync).in_time_zone}"
    else
      puts "Last D1 sync: Never (or cache expired)"
    end
    puts ""

    # Check usage_logs table
    usage_logs_count = UsageLog.count
    puts "Usage logs in PostgreSQL: #{usage_logs_count.to_s.reverse.gsub(/(\d{3})(?=\d)/, '\\1,').reverse}"

    if usage_logs_count > 0
      oldest = UsageLog.order(:used_at).first&.used_at
      newest = UsageLog.order(:used_at).last&.used_at
      puts "  Oldest record: #{oldest}"
      puts "  Newest record: #{newest}"
    end
    puts ""

    # Check daily summaries
    summaries_count = DailyUsageSummary.count
    puts "Daily summaries: #{summaries_count}"

    if summaries_count > 0
      oldest_date = DailyUsageSummary.order(:date).first&.date
      newest_date = DailyUsageSummary.order(:date).last&.date
      puts "  Oldest date: #{oldest_date}"
      puts "  Newest date: #{newest_date}"
    end
    puts ""

    # Check if sync job is configured
    puts "Recurring jobs configured:"
    puts "  - D1 sync: Every 5 minutes"
    puts "  - Daily aggregation: Every day at 00:05 UTC"
    puts ""
    puts "To manually trigger:"
    puts "  - D1 sync: SyncD1UsageJob.perform_later"
    puts "  - Daily aggregation: AggregateDailyUsageJob.perform_later(date: Date.yesterday)"
  end
end
