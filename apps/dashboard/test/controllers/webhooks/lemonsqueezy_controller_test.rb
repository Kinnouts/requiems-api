# frozen_string_literal: true

require "test_helper"

class Webhooks::LemonsqueezyControllerTest < ActionDispatch::IntegrationTest
  include ActionMailer::TestHelper
  include ActiveJob::TestHelper

  def setup
    @user = create_user(email: "payer@example.com")
    @request = PrivateDeploymentRequest.create!(
      user: @user,
      company: "Acme Corp",
      contact_name: "Sam",
      contact_email: @user.email,
      server_tier: "starter",
      billing_cycle: "monthly",
      status: "pending_payment",
      selected_services: %w[email text],
      monthly_price_cents: 20_000
    )
    clear_enqueued_jobs
    ActionMailer::Base.deliveries.clear
  end

  test "rejects webhook without signature" do
    post webhooks_lemonsqueezy_path, params: webhook_payload.to_json, headers: {
      "CONTENT_TYPE" => "application/json"
    }

    assert_response :unauthorized
  end

  test "marks private deployment request pending after subscription is created" do
    assert_enqueued_emails 2 do
      post webhooks_lemonsqueezy_path, params: webhook_payload.to_json, headers: signed_headers(webhook_payload)
    end

    assert_response :ok

    @request.reload
    assert_equal "pending", @request.status
    assert_equal "sub_123", @request.lemonsqueezy_subscription_id
  end

  test "returns ok for unknown private deployment request id" do
    payload = webhook_payload(private_deployment_request_id: "999999")

    assert_no_enqueued_emails do
      post webhooks_lemonsqueezy_path, params: payload.to_json, headers: signed_headers(payload)
    end

    assert_response :ok
  end

  private

  def webhook_payload(private_deployment_request_id: @request.id.to_s)
    {
      meta: {
        event_name: "subscription_created",
        custom_data: {
          private_deployment_request_id: private_deployment_request_id,
          user_id: @user.id.to_s
        }
      },
      data: {
        id: "sub_123",
        attributes: {
          customer_id: "cust_123",
          status: "active",
          variant_id: AppConfig.lemonsqueezy_developer_monthly_variant_id
        }
      }
    }
  end

  def signed_headers(payload)
    raw_payload = payload.to_json
    signature = OpenSSL::HMAC.hexdigest(
      OpenSSL::Digest.new("sha256"),
      AppConfig.lemonsqueezy_signing_secret,
      raw_payload
    )

    {
      "CONTENT_TYPE" => "application/json",
      "X-Signature" => signature
    }
  end
end
