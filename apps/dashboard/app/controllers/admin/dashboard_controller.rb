# frozen_string_literal: true

class Admin::DashboardController < ApplicationController
  before_action :authenticate_user!
  before_action :require_admin!
  layout "admin"

  def index
    # System Statistics
    @total_users = User.count
    @active_users = User.where("last_sign_in_at >= ?", 30.days.ago).count
    @total_api_keys = ApiKey.count
    @active_api_keys = ApiKey.active_keys.count

    # Usage Statistics (last 30 days)
    @total_requests_30d = UsageLog.where("used_at >= ?", 30.days.ago).count
    @total_requests_today = UsageLog.where("used_at >= ?", Time.current.beginning_of_day).count

    # Revenue Statistics
    @active_subscriptions = Subscription.where.not(plan_name: "free").count
    @mrr = calculate_mrr

    # System Health
    @avg_response_time = calculate_avg_response_time
    @error_rate = calculate_error_rate

    # Recent Activity
    @recent_users = User.order(created_at: :desc).limit(5)
    @recent_api_keys = ApiKey.order(created_at: :desc).limit(5)
    @recent_subscriptions = Subscription.where.not(plan_name: "free").order(created_at: :desc).limit(5)

    # Chart data (for AJAX loading)
    # Will be loaded via separate endpoints
  end

  private

  def require_admin!
    unless current_user.admin?
      redirect_to root_path, alert: "Access denied. Admin privileges required."
    end
  end

  def calculate_mrr
    # Calculate Monthly Recurring Revenue
    plan_prices = {
      "developer" => 30,
      "business" => 75,
      "professional" => 150
    }

    Subscription.where.not(plan_name: "free")
      .where(cancel_at_period_end: [false, nil])
      .group(:plan_name)
      .count
      .sum { |plan, count| (plan_prices[plan] || 0) * count }
  end

  def calculate_avg_response_time
    avg = UsageLog.where("used_at >= ?", 7.days.ago)
      .where.not(response_time_ms: nil)
      .average(:response_time_ms)

    avg&.round || 0
  end

  def calculate_error_rate
    total = UsageLog.where("used_at >= ?", 7.days.ago).count
    return 0 if total.zero?

    errors = UsageLog.where("used_at >= ?", 7.days.ago)
      .where("status_code >= ?", 400)
      .count

    ((errors.to_f / total) * 100).round(2)
  end
end
