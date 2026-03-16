# frozen_string_literal: true

class Dashboard::OverviewController < ApplicationController
  before_action :authenticate_user!
  layout "dashboard"

  def index
    @current_plan = current_user.subscription&.plan_name || "free"
    @usage_this_month = calculate_usage_this_month
    @total_requests = calculate_total_requests
    @requests_remaining = calculate_requests_remaining
    @avg_response_time = calculate_avg_response_time
    @recent_activity = fetch_recent_activity
    @api_keys_count = current_user.api_keys.active_keys.count
    @next_billing_date = current_user.subscription&.current_period_end
    @usage_percentage = calculate_usage_percentage
    @bar_color = calculate_bar_color
  end

  private

  def calculate_usage_this_month
    # Get usage from current month
    start_of_month = Time.current.beginning_of_month

    current_user.usage_logs
      .where("used_at >= ?", start_of_month)
      .count
  end

  def calculate_total_requests
    current_user.usage_logs.count
  end

  def calculate_requests_remaining
    # Plan limits (should match pricing config)
    plan_limits = {
      "free" => 500,
      "developer" => 100_000,
      "business" => 1_000_000,
      "professional" => 10_000_000
    }

    limit = plan_limits[@current_plan] || 500
    used = @usage_this_month

    [ limit - used, 0 ].max
  end

  def calculate_avg_response_time
    # Calculate average response time from recent requests
    recent_logs = current_user.usage_logs
      .where("used_at >= ?", 7.days.ago)
      .where.not(response_time_ms: nil)

    return 0 if recent_logs.empty?

    (recent_logs.average(:response_time_ms) || 0).round
  end

  def fetch_recent_activity
    current_user.usage_logs
      .order(used_at: :desc)
      .limit(10)
      .includes(:api_key)
  end

  def calculate_usage_percentage
    total_limit = @usage_this_month + @requests_remaining
    return 0 if total_limit <= 0

    ((@usage_this_month.to_f / total_limit) * 100).round
  end

  def calculate_bar_color
    if @usage_percentage >= 90
      "bg-red-500"
    elsif @usage_percentage >= 75
      "bg-yellow-500"
    else
      "bg-blue-500"
    end
  end
end
