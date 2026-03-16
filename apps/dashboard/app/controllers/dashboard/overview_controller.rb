# frozen_string_literal: true

class Dashboard::OverviewController < ApplicationController
  before_action :authenticate_user!
  layout "dashboard"

  def index
    @current_plan = current_user.current_plan
    @usage_this_month = current_user.usage_this_month
    @total_requests = current_user.total_requests
    @requests_remaining = current_user.requests_remaining
    @avg_response_time = current_user.avg_response_time_ms
    @recent_activity = current_user.recent_activity
    @api_keys_count = current_user.api_keys.active_keys.count
    @next_billing_date = current_user.subscription&.current_period_end
  end
end
