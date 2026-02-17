# frozen_string_literal: true

class Admin::AbuseReportsController < ApplicationController
  before_action :authenticate_user!
  before_action :require_admin!
  before_action :set_abuse_report, only: [ :show, :resolve, :investigate ]
  layout "admin"

  def index
    @abuse_reports = AbuseReport.includes(:user, :api_key).order(created_at: :desc)

    # Filter by status
    if params[:status].present? && params[:status] != "all"
      @abuse_reports = @abuse_reports.where(status: params[:status])
    end

    # Filter by type
    if params[:report_type].present? && params[:report_type] != "all"
      @abuse_reports = @abuse_reports.where(report_type: params[:report_type])
    end

    # Search by user email or API key
    if params[:search].present?
      search_term = "%#{params[:search]}%"
      @abuse_reports = @abuse_reports.joins(:user)
        .where("users.email ILIKE ? OR abuse_reports.description ILIKE ?", search_term, search_term)
    end

    # Paginate
    @pagy, @abuse_reports = pagy(@abuse_reports, items: 20)

    # Statistics
    @total_reports = AbuseReport.count
    @pending_reports = AbuseReport.where(status: "pending").count
    @investigating_reports = AbuseReport.where(status: "investigating").count
    @resolved_reports = AbuseReport.where(status: "resolved").count
  end

  def show
    @user = @abuse_report.user
    @api_key = @abuse_report.api_key
    @resolver = User.find(@abuse_report.resolved_by_id) if @abuse_report.resolved_by_id.present?

    # Get user's other abuse reports
    @other_reports = @user.abuse_reports.where.not(id: @abuse_report.id).order(created_at: :desc).limit(5)

    # Get user's API usage statistics
    @usage_stats = {
      total_requests: @user.usage_logs.count,
      requests_this_month: @user.usage_logs.where("used_at >= ?", Time.current.beginning_of_month).count,
      error_rate: calculate_error_rate(@user),
      last_request_at: @user.usage_logs.maximum(:used_at)
    }
  end

  def investigate
    if @abuse_report.update(status: "investigating")
      redirect_to admin_abuse_report_path(@abuse_report), notice: "Report marked as investigating."
    else
      redirect_to admin_abuse_report_path(@abuse_report), alert: "Failed to update report status."
    end
  end

  def resolve
    if @abuse_report.update(
      status: "resolved",
      resolved_at: Time.current,
      resolved_by_id: current_user.id
    )
      redirect_to admin_abuse_report_path(@abuse_report), notice: "Report resolved successfully."
    else
      redirect_to admin_abuse_report_path(@abuse_report), alert: "Failed to resolve report."
    end
  end

  private

  def require_admin!
    unless current_user.admin?
      redirect_to root_path, alert: "Access denied. Admin privileges required."
    end
  end

  def set_abuse_report
    @abuse_report = AbuseReport.find(params[:id])
  rescue ActiveRecord::RecordNotFound
    redirect_to admin_abuse_reports_path, alert: "Abuse report not found."
  end

  def calculate_error_rate(user)
    total = user.usage_logs.count
    return 0 if total.zero?

    errors = user.usage_logs.where("status_code >= ?", 400).count
    ((errors.to_f / total) * 100).round(2)
  end
end
