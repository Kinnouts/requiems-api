# frozen_string_literal: true

class Admin::AnalyticsController < ApplicationController
  before_action :authenticate_user!
  before_action :require_admin!
  layout "admin"

  def usage
    @date_range = params[:date_range] || "30"
    @start_date = @date_range.to_i.days.ago.beginning_of_day
    @end_date = Time.current.end_of_day

    # Total requests in period
    @total_requests = UsageLog.where(used_at: @start_date..@end_date).count

    # Requests per day (for chart)
    @requests_by_day = UsageLog
      .where(used_at: @start_date..@end_date)
      .group("DATE(used_at)")
      .count
      .transform_keys { |date| date.to_s }

    # Requests by endpoint (top 10)
    @requests_by_endpoint = UsageLog
      .where(used_at: @start_date..@end_date)
      .group(:endpoint)
      .count
      .sort_by { |_, count| -count }
      .first(10)
      .to_h

    # Requests by plan
    @requests_by_plan = UsageLog
      .where(used_at: @start_date..@end_date)
      .joins(user: :subscription)
      .group("subscriptions.plan_name")
      .count

    # Add free plan users
    free_requests = UsageLog
      .where(used_at: @start_date..@end_date)
      .joins(:user)
      .where.missing(:subscription)
      .count
    @requests_by_plan["free"] = free_requests if free_requests > 0

    # Average response times by endpoint (top 10)
    @avg_response_by_endpoint = UsageLog
      .where(used_at: @start_date..@end_date)
      .where.not(response_time_ms: nil)
      .group(:endpoint)
      .average(:response_time_ms)
      .sort_by { |_, avg| -avg }
      .first(10)
      .transform_values { |v| v.round(2) }

    # Top users by usage
    @top_users_by_usage = UsageLog
      .where(used_at: @start_date..@end_date)
      .joins(:user)
      .group("users.id", "users.email")
      .select("users.id, users.email, COUNT(*) as request_count, SUM(credits_used) as total_requests")
      .order("request_count DESC")
      .limit(10)
  end

  def revenue
    # Monthly Recurring Revenue
    @mrr = calculate_mrr

    # Annual Recurring Revenue
    @arr = @mrr * 12

    # Revenue by plan
    plan_prices = {
      "developer" => { monthly: 30, yearly: 25 },
      "business" => { monthly: 75, yearly: 62.5 },
      "professional" => { monthly: 150, yearly: 125 }
    }

    @revenue_by_plan = Subscription
      .where.not(plan_name: "free")
      .where(cancel_at_period_end: [ false, nil ])
      .group(:plan_name, :billing_cycle)
      .count
      .transform_keys { |plan, cycle| "#{plan.titleize} (#{cycle&.titleize || 'Monthly'})" }
      .transform_values do |(plan, cycle), count|
        price = plan_prices[plan]&.fetch(cycle&.to_sym || :monthly, 0) || 0
        price * count
      end

    # Revenue trend (last 12 months)
    @revenue_trend = {}
    (0..11).each do |i|
      month_start = i.months.ago.beginning_of_month
      month_end = i.months.ago.end_of_month
      month_label = month_start.strftime("%b %Y")

      # Calculate revenue for that month (subscriptions active during that period)
      month_revenue = Subscription
        .where.not(plan_name: "free")
        .where("billing_cycle_start <= ?", month_end)
        .where("billing_cycle_end IS NULL OR billing_cycle_end >= ?", month_start)
        .sum do |sub|
          plan_prices[sub.plan_name]&.fetch(sub.billing_cycle&.to_sym || :monthly, 0) || 0
        end

      @revenue_trend[month_label] = month_revenue
    end
    @revenue_trend = @revenue_trend.reverse_each.to_h

    # Active subscriptions
    @active_subscriptions = Subscription
      .where.not(plan_name: "free")
      .where(cancel_at_period_end: [ false, nil ])
      .count

    # Subscriptions by plan
    @subscriptions_by_plan = Subscription
      .where.not(plan_name: "free")
      .group(:plan_name)
      .count
      .transform_keys(&:titleize)

    # Churn rate (last 30 days)
    start_of_period = 30.days.ago
    subscriptions_at_start = Subscription
      .where("created_at < ?", start_of_period)
      .where.not(plan_name: "free")
      .count

    canceled_in_period = Subscription
      .where(cancel_at_period_end: true)
      .where("canceled_at >= ?", start_of_period)
      .where.not(plan_name: "free")
      .count

    @churn_rate = subscriptions_at_start > 0 ? ((canceled_in_period.to_f / subscriptions_at_start) * 100).round(2) : 0

    # New vs Canceled subscriptions (last 12 months)
    @new_vs_canceled = {}
    (0..11).each do |i|
      month_start = i.months.ago.beginning_of_month
      month_end = i.months.ago.end_of_month
      month_label = month_start.strftime("%b %Y")

      new_subs = Subscription
        .where.not(plan_name: "free")
        .where(created_at: month_start..month_end)
        .count

      canceled_subs = Subscription
        .where(cancel_at_period_end: true)
        .where(canceled_at: month_start..month_end)
        .where.not(plan_name: "free")
        .count

      @new_vs_canceled[month_label] = { new: new_subs, canceled: canceled_subs }
    end
    @new_vs_canceled = @new_vs_canceled.reverse_each.to_h
  end

  def system_health
    # Time range for health metrics
    @time_range = params[:time_range] || "24h"
    @start_time = case @time_range
    when "1h" then 1.hour.ago
    when "24h" then 24.hours.ago
    when "7d" then 7.days.ago
    when "30d" then 30.days.ago
    else 24.hours.ago
    end

    # API Uptime (percentage of successful requests)
    total_requests = UsageLog.where(used_at: @start_time..Time.current).count
    successful_requests = UsageLog.where(used_at: @start_time..Time.current).where(status_code: 200..299).count
    @uptime_percentage = total_requests > 0 ? ((successful_requests.to_f / total_requests) * 100).round(2) : 100.0

    # Average response times (P50, P95, P99)
    response_times = UsageLog
      .where(used_at: @start_time..Time.current)
      .where.not(response_time_ms: nil)
      .order(:response_time_ms)
      .pluck(:response_time_ms)

    if response_times.any?
      @p50_response_time = percentile(response_times, 50).round(2)
      @p95_response_time = percentile(response_times, 95).round(2)
      @p99_response_time = percentile(response_times, 99).round(2)
    else
      @p50_response_time = @p95_response_time = @p99_response_time = 0
    end

    # Error rate trend (by hour or day)
    group_interval = @time_range == "1h" ? 5.minutes : (@time_range == "24h" ? 1.hour : 1.day)
    @error_rate_trend = {}

    current = @start_time
    while current < Time.current
      next_time = current + group_interval
      period_total = UsageLog.where(used_at: current..next_time).count
      period_errors = UsageLog.where(used_at: current..next_time).where("status_code >= ?", 400).count

      error_rate = period_total > 0 ? ((period_errors.to_f / period_total) * 100).round(2) : 0
      label = current.strftime(@time_range == "1h" ? "%H:%M" : (@time_range == "24h" ? "%H:00" : "%b %d"))
      @error_rate_trend[label] = error_rate

      current = next_time
    end

    # Rate limit hits (last 24h)
    # Note: This would need to be tracked in the database
    # For now, we'll show 0 as a placeholder
    @rate_limit_hits = 0

    # Failed authentication attempts (last 24h)
    @failed_auth_attempts = UsageLog
      .where(used_at: 24.hours.ago..Time.current)
      .where(status_code: 401)
      .count

    # Most common errors (top 10)
    @common_errors = UsageLog
      .where(used_at: @start_time..Time.current)
      .where("status_code >= ?", 400)
      .group(:status_code)
      .count
      .sort_by { |_, count| -count }
      .first(10)
      .to_h

    # Requests per minute (last hour) for real-time monitoring
    if @time_range == "1h"
      @requests_per_minute = UsageLog
        .where(used_at: 1.hour.ago..Time.current)
        .group("DATE_TRUNC('minute', used_at)")
        .count
        .transform_keys { |time| time.strftime("%H:%M") }
    end
  end

  private

  def require_admin!
    unless current_user.admin?
      redirect_to root_path, alert: "Access denied. Admin privileges required."
    end
  end

  def calculate_mrr
    plan_prices = {
      "developer" => { monthly: 30, yearly: 25 },
      "business" => { monthly: 75, yearly: 62.5 },
      "professional" => { monthly: 150, yearly: 125 }
    }

    Subscription
      .where.not(plan_name: "free")
      .where(cancel_at_period_end: [ false, nil ])
      .sum do |sub|
        price = plan_prices[sub.plan_name]&.fetch(sub.billing_cycle&.to_sym || :monthly, 0) || 0
        sub.billing_cycle == "yearly" ? (price * 12 / 12.0) : price
      end
  end

  def percentile(sorted_array, percentile)
    return 0 if sorted_array.empty?

    index = (percentile / 100.0) * (sorted_array.length - 1)
    lower = sorted_array[index.floor]
    upper = sorted_array[index.ceil]

    lower + (upper - lower) * (index - index.floor)
  end
end
