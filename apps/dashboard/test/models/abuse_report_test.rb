require "test_helper"

class AbuseReportTest < ActiveSupport::TestCase
  def setup
    @user = User.create!(
      email: "reported@example.com",
      password: "password123",
      password_confirmation: "password123"
    )

    @api_key = @user.api_keys.create!(
      name: "Test Key",
      key: "rq_test_" + SecureRandom.hex(32),
      prefix: "rq_test_abc123",
      environment: "test"
    )

    @report = AbuseReport.create!(
      user: @user,
      api_key: @api_key,
      report_type: "spam",
      description: "This user is sending spam",
      status: "pending"
    )
  end

  test "valid abuse report with required attributes" do
    assert @report.valid?
    assert @report.persisted?
  end

  test "requires report_type" do
    report = AbuseReport.new(
      user: @user,
      api_key: @api_key,
      description: "Test",
      status: "pending"
    )
    report.report_type = nil

    assert_not report.valid?
    assert_includes report.errors[:report_type], "can't be blank"
  end

  test "requires status" do
    report = AbuseReport.new(
      user: @user,
      api_key: @api_key,
      report_type: "spam",
      description: "Test"
    )
    report.status = nil

    assert_not report.valid?
    assert_includes report.errors[:status], "can't be blank"
  end

  test "validates status inclusion" do
    @report.status = "invalid_status"
    assert_not @report.valid?

    @report.status = "pending"
    assert @report.valid?

    @report.status = "investigating"
    assert @report.valid?

    @report.status = "resolved"
    assert @report.valid?
  end

  test "validates description length" do
    @report.description = "a" * 5001
    assert_not @report.valid?
    assert_includes @report.errors[:description], "is too long (maximum is 5000 characters)"

    @report.description = "a" * 5000
    assert @report.valid?
  end

  test "belongs to user" do
    assert_equal @user, @report.user
  end

  test "belongs to api_key" do
    assert_equal @api_key, @report.api_key
  end

  test "optionally belongs to resolved_by user" do
    resolver = User.create!(
      email: "admin@example.com",
      password: "password123",
      password_confirmation: "password123",
      admin: true
    )

    @report.update(
      resolved_by: resolver,
      resolved_at: Time.current,
      status: "resolved"
    )

    assert_equal resolver, @report.resolved_by
  end

  test "pending scope returns pending reports" do
    pending_report = AbuseReport.create!(
      user: @user,
      api_key: @api_key,
      report_type: "abuse",
      status: "pending"
    )

    resolved_report = AbuseReport.create!(
      user: @user,
      api_key: @api_key,
      report_type: "spam",
      status: "resolved"
    )

    pending_reports = AbuseReport.pending

    assert_includes pending_reports, pending_report
    assert_not_includes pending_reports, resolved_report
  end

  test "investigating scope returns investigating reports" do
    investigating_report = AbuseReport.create!(
      user: @user,
      api_key: @api_key,
      report_type: "malicious",
      status: "investigating"
    )

    investigating_reports = AbuseReport.investigating

    assert_includes investigating_reports, investigating_report
    assert_not_includes investigating_reports, @report
  end

  test "resolved scope returns resolved reports" do
    resolved_report = AbuseReport.create!(
      user: @user,
      api_key: @api_key,
      report_type: "terms_violation",
      status: "resolved",
      resolved_at: Time.current
    )

    resolved_reports = AbuseReport.resolved

    assert_includes resolved_reports, resolved_report
    assert_not_includes resolved_reports, @report
  end

  test "pending? helper method" do
    assert @report.pending?

    @report.update(status: "investigating")
    assert_not @report.pending?
  end

  test "investigating? helper method" do
    assert_not @report.investigating?

    @report.update(status: "investigating")
    assert @report.investigating?
  end

  test "resolved? helper method" do
    assert_not @report.resolved?

    @report.update(status: "resolved", resolved_at: Time.current)
    assert @report.resolved?
  end

  test "recent scope orders by created_at desc" do
    old_report = AbuseReport.create!(
      user: @user,
      api_key: @api_key,
      report_type: "spam",
      status: "pending",
      created_at: 1.day.ago
    )

    recent_reports = AbuseReport.recent

    assert_equal @report, recent_reports.first
    assert_equal old_report, recent_reports.last
  end
end
