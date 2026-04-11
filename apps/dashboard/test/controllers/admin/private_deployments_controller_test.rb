# frozen_string_literal: true

require "test_helper"

class Admin::PrivateDeploymentsControllerTest < ActionDispatch::IntegrationTest
  include ActionMailer::TestHelper
  include ActiveJob::TestHelper

  def setup
    @admin = create_user(email: "admin@example.com", admin: true)
    @user = create_user(email: "customer@example.com", company: "Acme Corp")
    @request = PrivateDeploymentRequest.create!(
      user: @user,
      company: "Acme Corp",
      contact_name: "Jamie",
      contact_email: @user.email,
      server_tier: "starter",
      billing_cycle: "monthly",
      status: "pending",
      selected_services: %w[email tech],
      monthly_price_cents: 20_000,
      admin_notes: "Existing note"
    )
    clear_enqueued_jobs
    sign_in @admin
  end

  test "index shows requests and filters by status" do
    active_request = PrivateDeploymentRequest.create!(
      user: @user,
      company: "Acme Corp",
      contact_name: "Jamie",
      contact_email: @user.email,
      server_tier: "growth",
      billing_cycle: "monthly",
      status: "active",
      selected_services: %w[email],
      monthly_price_cents: 30_000,
      tenant_secret: "x" * 32,
      subdomain_slug: "acme-growth"
    )

    get admin_private_deployments_path, params: { status: "active" }

    assert_response :success
    assert_match active_request.subdomain_slug, response.body
    assert_no_match @request.contact_email, response.body
  end

  test "show renders request details" do
    get admin_private_deployment_path(@request)

    assert_response :success
    assert_match @request.contact_email, response.body
    assert_match @request.company, response.body
  end

  test "activate marks request active and enqueues deployment email" do
    assert_enqueued_emails 1 do
      patch activate_admin_private_deployment_path(@request), params: {
        subdomain_slug: "Acme-Prod",
        tenant_secret: "s" * 32,
        admin_notes: "Tenant provisioned"
      }
    end

    @request.reload
    assert_equal "active", @request.status
    assert_equal "acme-prod", @request.subdomain_slug
    assert_equal "Existing note\n\n---\n\nTenant provisioned", @request.admin_notes
    assert_not_nil @request.deployed_at
    assert_redirected_to admin_private_deployment_path(@request)
  end

  test "activate rejects missing subdomain" do
    patch activate_admin_private_deployment_path(@request), params: {
      subdomain_slug: "",
      tenant_secret: "s" * 32
    }

    assert_redirected_to admin_private_deployment_path(@request)
    assert_match(/Subdomain slug is required/i, flash[:alert])
  end

  test "cancel marks request cancelled" do
    patch cancel_admin_private_deployment_path(@request)

    @request.reload
    assert_equal "cancelled", @request.status
    assert_redirected_to admin_private_deployments_path
  end

  test "non admin cannot access admin private deployments" do
    sign_out @admin
    sign_in @user

    get admin_private_deployments_path

    assert_response :not_found
  end
end
