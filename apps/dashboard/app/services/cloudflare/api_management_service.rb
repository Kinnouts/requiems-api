# frozen_string_literal: true

require "cgi"

module Cloudflare
  # Service to interact with the API Management Cloudflare Worker.
  # All API key CRUD operations go through this service so that both
  # Cloudflare KV (auth validation) and D1 (usage tracking) stay in sync.
  class ApiManagementService
    def initialize
      @base_url = ENV.fetch("API_MANAGEMENT_URL", "https://api-management.requiems.xyz")
      @management_key = ENV["API_MANAGEMENT_API_KEY"]
    end

    # Create a new API key via the api-management worker.
    # Returns the full API key string (shown to the user only once) or nil on failure.
    def create_key(user_id:, plan:, name:, billing_cycle_start: nil)
      payload = {
        userId: user_id.to_s,
        plan: plan,
        name: name,
        billingCycleStart: billing_cycle_start || Time.current.iso8601
      }

      response = connection.post("/api-keys") do |req|
        req.body = payload.to_json
      end

      if response.success?
        body = JSON.parse(response.body)
        Rails.logger.info "[ApiManagement] API key created: #{body['keyPrefix']}"
        body["apiKey"]
      else
        Rails.logger.error "[ApiManagement] Failed to create API key (#{response.status}): #{response.body}"
        nil
      end
    rescue Faraday::Error => e
      Rails.logger.error "[ApiManagement] Create key network error: #{e.message}"
      nil
    end

    # Revoke an API key by prefix via the api-management worker.
    def revoke_key(key_prefix)
      response = connection.delete("/api-keys/#{CGI.escape(key_prefix.to_s)}")
      handle_response(response, "revoke", key_prefix)
    rescue Faraday::Error => e
      Rails.logger.error "[ApiManagement] Revoke key network error for #{key_prefix}: #{e.message}"
      false
    end

    # Update an API key's plan and/or billing cycle start.
    def update_key(key_prefix, plan: nil, billing_cycle_start: nil)
      payload = {}
      payload[:plan] = plan if plan
      payload[:billingCycleStart] = billing_cycle_start if billing_cycle_start
      return false if payload.empty?

      response = connection.patch("/api-keys/#{CGI.escape(key_prefix.to_s)}") do |req|
        req.body = payload.to_json
      end

      handle_response(response, "update", key_prefix)
    rescue Faraday::Error => e
      Rails.logger.error "[ApiManagement] Update key network error for #{key_prefix}: #{e.message}"
      false
    end

    # Sync all active API keys for a user to a new plan.
    # Called when a user upgrades/downgrades their subscription.
    def sync_user_plan(user, plan_name)
      billing_start = resolve_billing_start(user)
      user.api_keys.active_keys.find_each do |api_key|
        update_key(api_key.key_prefix, plan: plan_name, billing_cycle_start: billing_start)
      end
    end

    private

    def resolve_billing_start(user)
      subscription = user.subscription
      start_date = subscription&.current_period_start || subscription&.created_at || Time.current
      start_date.iso8601
    end

    def connection
      @connection ||= Faraday.new(url: @base_url) do |f|
        f.headers["Content-Type"] = "application/json"
        f.headers["X-API-Management-Key"] = @management_key.to_s
        f.request :retry, max: 3, interval: 0.5
        f.adapter Faraday.default_adapter
      end
    end

    def handle_response(response, action, key_prefix)
      if response.success?
        Rails.logger.info "[ApiManagement] Key #{action} succeeded for #{key_prefix}"
        true
      else
        Rails.logger.error "[ApiManagement] Key #{action} failed for #{key_prefix} (#{response.status}): #{response.body}"
        false
      end
    end
  end
end
