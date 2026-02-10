# frozen_string_literal: true

class Dashboard::BillingController < ApplicationController
  before_action :authenticate_user!
  layout "dashboard"

  def show
    @subscription = current_user.subscription
    @current_plan = @subscription&.plan_name || "free"

    # Calculate usage for current billing period
    if @subscription
      period_start = @subscription.current_period_start || Time.current.beginning_of_month
      @usage_this_period = current_user.usage_logs
        .where("used_at >= ?", period_start)
        .count
    else
      @usage_this_period = current_user.usage_logs
        .where("used_at >= ?", Time.current.beginning_of_month)
        .count
    end

    # Get plan limits and pricing
    @plan_info = get_plan_info(@current_plan)
    @available_plans = get_all_plans

    # Calculate quota usage percentage
    quota_limit = @plan_info[:requests_per_month]
    @quota_used_percent = quota_limit > 0 ? ((@usage_this_period.to_f / quota_limit) * 100).round(1) : 0

    # Handle success/cancel callbacks from LemonSqueezy
    if params[:success] == "true"
      flash.now[:notice] = "Subscription activated successfully! Your new plan is now active."
    elsif params[:canceled] == "true"
      flash.now[:alert] = "Checkout was canceled. Your subscription was not changed."
    end
  end

  def checkout
    # Create LemonSqueezy checkout session for plan upgrade
    plan = params[:plan]

    unless valid_plan?(plan)
      redirect_to dashboard_billing_path, alert: "Invalid plan selected"
      return
    end

    billing_cycle = params[:billing_cycle] || "monthly"
    checkout_url = create_lemonsqueezy_checkout_url(plan, billing_cycle)

    redirect_to checkout_url, allow_other_host: true
  end

  def portal
    # Redirect to LemonSqueezy customer portal
    unless current_user.subscription&.lemonsqueezy_subscription_id
      redirect_to dashboard_billing_path, alert: "No active subscription found"
      return
    end

    # LemonSqueezy customer portal URL
    # Format: https://[store-slug].lemonsqueezy.com/billing
    store_slug = ENV["LEMONSQUEEZY_STORE_SLUG"] || "requiems"
    portal_url = "https://#{store_slug}.lemonsqueezy.com/billing"

    redirect_to portal_url, allow_other_host: true
  end

  def cancel_subscription
    unless current_user.subscription
      redirect_to dashboard_billing_path, alert: "No active subscription to cancel"
      return
    end

    # For LemonSqueezy, we'll direct users to the customer portal
    # where they can manage their subscription
    redirect_to portal_dashboard_billing_path
  end

  private

  def get_plan_info(plan_name)
    plans = {
      "free" => {
        name: "Free",
        price_monthly: 0,
        price_yearly: 0,
        requests_per_month: 500,
        rate_limit_per_minute: 30,
        features: [
          "500 requests/month",
          "30 requests/minute",
          "Community support",
          "US data centers"
        ]
      },
      "developer" => {
        name: "Developer",
        price_monthly: 59,
        price_yearly: 468, # $39/month billed yearly
        requests_per_month: 100_000,
        rate_limit_per_minute: 5_000,
        lemonsqueezy_variant_id_monthly: ENV["LEMONSQUEEZY_DEVELOPER_MONTHLY_VARIANT_ID"],
        lemonsqueezy_variant_id_yearly: ENV["LEMONSQUEEZY_DEVELOPER_YEARLY_VARIANT_ID"],
        features: [
          "100,000 requests/month",
          "5,000 requests/minute",
          "Email support",
          "US data centers"
        ]
      },
      "business" => {
        name: "Business",
        price_monthly: 149,
        price_yearly: 1188, # $99/month billed yearly
        requests_per_month: 1_000_000,
        rate_limit_per_minute: 10_000,
        lemonsqueezy_variant_id_monthly: ENV["LEMONSQUEEZY_BUSINESS_MONTHLY_VARIANT_ID"],
        lemonsqueezy_variant_id_yearly: ENV["LEMONSQUEEZY_BUSINESS_YEARLY_VARIANT_ID"],
        features: [
          "1M requests/month",
          "10,000 requests/minute",
          "Priority email support",
          "US & EU data centers",
          "99.9% SLA"
        ]
      },
      "professional" => {
        name: "Professional",
        price_monthly: 299,
        price_yearly: 2388, # $199/month billed yearly
        requests_per_month: 10_000_000,
        rate_limit_per_minute: 50_000,
        lemonsqueezy_variant_id_monthly: ENV["LEMONSQUEEZY_PROFESSIONAL_MONTHLY_VARIANT_ID"],
        lemonsqueezy_variant_id_yearly: ENV["LEMONSQUEEZY_PROFESSIONAL_YEARLY_VARIANT_ID"],
        features: [
          "10M requests/month",
          "50,000 requests/minute",
          "24/7 priority support",
          "US & EU data centers",
          "99.99% SLA",
          "Dedicated support engineer"
        ]
      }
    }

    plans[plan_name] || plans["free"]
  end

  def get_all_plans
    ["free", "developer", "business", "professional"].map { |plan| get_plan_info(plan).merge(id: plan) }
  end

  def valid_plan?(plan)
    %w[developer business professional].include?(plan)
  end

  def create_lemonsqueezy_checkout_url(plan, billing_cycle)
    plan_info = get_plan_info(plan)

    variant_id = if billing_cycle == "yearly"
                   plan_info[:lemonsqueezy_variant_id_yearly]
                 else
                   plan_info[:lemonsqueezy_variant_id_monthly]
                 end

    store_id = ENV["LEMONSQUEEZY_STORE_ID"]

    # LemonSqueezy checkout URL format
    # https://[store-slug].lemonsqueezy.com/checkout/buy/[variant-id]?checkout[email]=user@example.com&checkout[custom][user_id]=123
    store_slug = ENV["LEMONSQUEEZY_STORE_SLUG"] || "requiems"
    checkout_url = "https://#{store_slug}.lemonsqueezy.com/checkout/buy/#{variant_id}"

    # Add query parameters
    params = {
      "checkout[email]" => current_user.email,
      "checkout[custom][user_id]" => current_user.id,
      "checkout[custom][plan]" => plan
    }

    "#{checkout_url}?#{params.to_query}"
  end
end
