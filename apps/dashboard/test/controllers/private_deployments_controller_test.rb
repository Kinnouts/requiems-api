# frozen_string_literal: true

require "test_helper"

class PrivateDeploymentsControllerTest < ActionDispatch::IntegrationTest
  def setup
    @user = create_user(
      email: "owner@example.com",
      name: "Morgan",
      company: "Acme Labs"
    )
    sign_in @user
  end

  test "new preselects all services" do
    get new_private_deployment_path

    assert_response :success
    PrivateDeploymentRequest::VALID_SERVICES.each do |service|
      assert_match(/private_deployment_request_selected_services_#{service}/, response.body)
    end
  end

  test "create persists request and redirects to lemonsqueezy checkout" do
    post private_deployments_path, params: {
      private_deployment_request: {
        server_tier: "growth",
        billing_cycle: "monthly",
        selected_services: %w[email tech],
        admin_notes: "Need SSO"
      }
    }

    request = PrivateDeploymentRequest.order(:id).last

    assert_redirected_to(/lemonsqueezy\.com\/checkout\/buy\//)
    assert_equal "Acme Labs", request.company
    assert_equal "Morgan", request.contact_name
    assert_equal "owner@example.com", request.contact_email
    assert_equal "pending_payment", request.status
    assert_equal 30_000, request.monthly_price_cents
    assert_equal %w[email tech], request.selected_services

    redirect_uri = URI.parse(response.redirect_url)
    params = Rack::Utils.parse_nested_query(redirect_uri.query)

    assert_equal "00000000-0000-0000-0000-000000000023", redirect_uri.path.split("/").last
    assert_equal @user.email, params.dig("checkout", "email")
    assert_equal @user.id.to_s, params.dig("checkout", "custom", "user_id")
    assert_equal request.id.to_s, params.dig("checkout", "custom", "private_deployment_request_id")
  end

  test "create derives effective monthly price for yearly billing" do
    post private_deployments_path, params: {
      private_deployment_request: {
        server_tier: "scale",
        billing_cycle: "yearly",
        selected_services: %w[email]
      }
    }

    assert_equal 42_500, PrivateDeploymentRequest.order(:id).last.monthly_price_cents
  end

  test "create falls back to talk to sales when checkout config is missing" do
    config = AppConfig.instance
    original = config.instance_variable_get(:@lemonsqueezy_private_starter_monthly_checkout_uuid)
    config.instance_variable_set(:@lemonsqueezy_private_starter_monthly_checkout_uuid, nil)

    post private_deployments_path, params: {
      private_deployment_request: {
        server_tier: "starter",
        billing_cycle: "monthly",
        selected_services: %w[email]
      }
    }

    assert_redirected_to talk_to_sales_url(locale: I18n.default_locale)
  ensure
    config.instance_variable_set(:@lemonsqueezy_private_starter_monthly_checkout_uuid, original)
  end

  test "create re-renders form for invalid request" do
    assert_no_difference "PrivateDeploymentRequest.count" do
      post private_deployments_path, params: {
        private_deployment_request: {
          server_tier: "starter",
          billing_cycle: "monthly",
          selected_services: []
        }
      }
    end

    assert_response :unprocessable_entity
  end
end
