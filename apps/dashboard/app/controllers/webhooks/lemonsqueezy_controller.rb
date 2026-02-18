# frozen_string_literal: true

class Webhooks::LemonsqueezyController < ApplicationController
  skip_before_action :verify_authenticity_token
  before_action :verify_signature

  def create
    event_name = params[:meta][:event_name]

    Rails.logger.info "[LemonSqueezy Webhook] Received: #{event_name}"

    case event_name
    when "subscription_created"
      handle_subscription_created
    when "subscription_updated"
      handle_subscription_updated
    when "subscription_cancelled", "subscription_expired"
      handle_subscription_cancelled
    when "subscription_resumed"
      handle_subscription_resumed
    when "subscription_payment_success"
      handle_payment_success
    else
      Rails.logger.warn "[LemonSqueezy Webhook] Unhandled event: #{event_name}"
    end

    head :ok
  rescue => e
    Rails.logger.error "[LemonSqueezy Webhook] Error: #{e.message}"
    Rails.logger.error e.backtrace.join("\n")
    head :bad_request
  end

  private

  def verify_signature
    signature = request.headers["X-Signature"]

    unless signature
      Rails.logger.error "[LemonSqueezy Webhook] Missing signature"
      head :unauthorized
      return
    end

    secret = ENV["LEMONSQUEEZY_SIGNING_SECRET"]
    payload = request.body.read

    expected_signature = OpenSSL::HMAC.hexdigest(
      OpenSSL::Digest.new("sha256"),
      secret,
      payload
    )

    unless Rack::Utils.secure_compare(signature, expected_signature)
      Rails.logger.error "[LemonSqueezy Webhook] Invalid signature"
      head :unauthorized
      return
    end

    # Reset request body for controller to read
    request.body.rewind
  end

  def handle_subscription_created
    data = params[:data][:attributes]
    custom_data = data[:custom_data] || {}

    user = User.find_by(id: custom_data[:user_id])
    unless user
      Rails.logger.error "[LemonSqueezy] User not found: #{custom_data[:user_id]}"
      return
    end

    plan_name = determine_plan_name(data[:variant_id])

    subscription = user.subscription || user.build_subscription
    subscription.update!(
      lemonsqueezy_subscription_id: data[:id],
      lemonsqueezy_customer_id: data[:customer_id],
      plan_name: plan_name,
      status: data[:status],
      current_period_start: data[:renews_at] ? Time.zone.parse(data[:renews_at]) - 1.month : Time.current,
      current_period_end: data[:renews_at],
      trial_ends_at: data[:trial_ends_at],
      cancel_at_period_end: false
    )

    # Sync to Cloudflare KV
    Cloudflare::ApiManagementService.new.sync_user_plan(user, plan_name)

    Rails.logger.info "[LemonSqueezy] Subscription created for user #{user.id}: #{plan_name}"
  end

  def handle_subscription_updated
    data = params[:data][:attributes]

    subscription = Subscription.find_by(lemonsqueezy_subscription_id: data[:id])
    unless subscription
      Rails.logger.error "[LemonSqueezy] Subscription not found: #{data[:id]}"
      return
    end

    plan_name = determine_plan_name(data[:variant_id])

    subscription.update!(
      status: data[:status],
      plan_name: plan_name,
      current_period_end: data[:renews_at],
      cancel_at_period_end: data[:ends_at].present?
    )

    # Sync to Cloudflare KV
    Cloudflare::ApiManagementService.new.sync_user_plan(subscription.user, plan_name)

    Rails.logger.info "[LemonSqueezy] Subscription updated: #{subscription.id}"
  end

  def handle_subscription_cancelled
    data = params[:data][:attributes]

    subscription = Subscription.find_by(lemonsqueezy_subscription_id: data[:id])
    unless subscription
      Rails.logger.error "[LemonSqueezy] Subscription not found: #{data[:id]}"
      return
    end

    subscription.update!(
      status: "cancelled",
      cancel_at_period_end: true,
      plan_name: "free"
    )

    # Downgrade to free plan in Cloudflare KV
    Cloudflare::ApiManagementService.new.sync_user_plan(subscription.user, "free")

    Rails.logger.info "[LemonSqueezy] Subscription cancelled: #{subscription.id}"
  end

  def handle_subscription_resumed
    data = params[:data][:attributes]

    subscription = Subscription.find_by(lemonsqueezy_subscription_id: data[:id])
    unless subscription
      Rails.logger.error "[LemonSqueezy] Subscription not found: #{data[:id]}"
      return
    end

    plan_name = determine_plan_name(data[:variant_id])

    subscription.update!(
      status: data[:status],
      plan_name: plan_name,
      cancel_at_period_end: false
    )

    # Restore plan in Cloudflare KV
    Cloudflare::ApiManagementService.new.sync_user_plan(subscription.user, plan_name)

    Rails.logger.info "[LemonSqueezy] Subscription resumed: #{subscription.id}"
  end

  def handle_payment_success
    data = params[:data][:attributes]

    subscription = Subscription.find_by(lemonsqueezy_subscription_id: data[:subscription_id])
    return unless subscription

    subscription.update!(
      current_period_start: Time.current,
      current_period_end: data[:renews_at]
    )

    Rails.logger.info "[LemonSqueezy] Payment success for subscription: #{subscription.id}"
  end

  def determine_plan_name(variant_id)
    variant_id = variant_id.to_s

    case variant_id
    when ENV["LEMONSQUEEZY_DEVELOPER_MONTHLY_VARIANT_ID"], ENV["LEMONSQUEEZY_DEVELOPER_YEARLY_VARIANT_ID"]
      "developer"
    when ENV["LEMONSQUEEZY_BUSINESS_MONTHLY_VARIANT_ID"], ENV["LEMONSQUEEZY_BUSINESS_YEARLY_VARIANT_ID"]
      "business"
    when ENV["LEMONSQUEEZY_PROFESSIONAL_MONTHLY_VARIANT_ID"], ENV["LEMONSQUEEZY_PROFESSIONAL_YEARLY_VARIANT_ID"]
      "professional"
    else
      Rails.logger.warn "[LemonSqueezy] Unknown variant_id: #{variant_id}"
      "free"
    end
  end
end
