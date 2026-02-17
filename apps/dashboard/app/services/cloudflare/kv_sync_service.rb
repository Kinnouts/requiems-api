# frozen_string_literal: true

module Cloudflare
  class KvSyncService
    def initialize(api_key_record)
      @api_key = api_key_record
      @api_management_url = ENV.fetch("API_MANAGEMENT_URL", "https://api-management.requiems.xyz")
      @api_management_key = ENV.fetch("API_MANAGEMENT_API_KEY")
    end

    # Sync API key to Cloudflare Worker (which writes to KV + D1)
    def sync_create
      payload = {
        action: "create",
        key: @api_key.full_key,
        userId: @api_key.user_id.to_s,
        plan: @api_key.user.current_plan,
        billingCycleStart: billing_cycle_start
      }

      response = connection.post("/api-keys") do |req|
        req.headers["Content-Type"] = "application/json"
        req.headers["X-API-Management-Key"] = @api_management_key
        req.body = payload.to_json
      end

      handle_response(response, "create")
    end

    # Remove API key from Cloudflare Worker (deletes from KV + marks in D1)
    def sync_delete
      payload = {
        action: "revoke",
        key: @api_key.full_key,
        userId: @api_key.user_id.to_s,
        plan: @api_key.user.current_plan
      }

      response = connection.post("/api-keys") do |req|
        req.headers["Content-Type"] = "application/json"
        req.headers["X-API-Management-Key"] = @api_management_key
        req.body = payload.to_json
      end

      handle_response(response, "delete")
    end

    # Update API key plan (e.g., after subscription change)
    def sync_update
      payload = {
        action: "update",
        key: @api_key.full_key,
        userId: @api_key.user_id.to_s,
        plan: @api_key.user.current_plan,
        billingCycleStart: billing_cycle_start
      }

      response = connection.post("/api-keys") do |req|
        req.headers["Content-Type"] = "application/json"
        req.headers["X-API-Management-Key"] = @api_management_key
        req.body = payload.to_json
      end

      handle_response(response, "update")
    end

    private

    def billing_cycle_start
      # Use subscription created_at or current time as billing cycle start
      subscription = @api_key.user.subscription
      start_date = subscription&.created_at || Time.current
      start_date.iso8601
    end

    def connection
      @connection ||= Faraday.new(
        url: @api_management_url,
        headers: {
          "Content-Type" => "application/json"
        }
      ) do |f|
        f.request :retry, max: 3, interval: 0.5
        f.adapter Faraday.default_adapter
      end
    end

    def handle_response(response, action)
      if response.success?
        Rails.logger.info("Cloudflare Worker #{action} succeeded for API key #{@api_key.key_prefix}")
        true
      else
        Rails.logger.error("Cloudflare Worker #{action} failed for API key #{@api_key.key_prefix}: #{response.body}")
        false
      end
    rescue Faraday::Error => e
      Rails.logger.error("Cloudflare Worker #{action} error for API key #{@api_key.key_prefix}: #{e.message}")
      false
    end
  end
end
