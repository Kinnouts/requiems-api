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
    limit = PlanConfig.requests_per_month(@current_plan)
    [ limit - @usage_this_month, 0 ].max
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
end
