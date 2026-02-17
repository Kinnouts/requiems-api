# frozen_string_literal: true

module Cloudflare
  class KvSyncService
    def initialize(api_key_record)
      @api_key = api_key_record
      @api_management_url = ENV.fetch("API_MANAGEMENT_URL", "https://api-management.requiems.xyz")
      @api_management_key = ENV.fetch("API_MANAGEMENT_API_KEY")
    end

    # Create API key via api-management service
    # Returns the full API key from the server (only shown once)
    def sync_create
      payload = {
        userId: @api_key.user_id.to_s,
        plan: @api_key.user.current_plan,
        name: @api_key.name,
        billingCycleStart: billing_cycle_start
      }

      response = connection.post("/api-keys") do |req|
        req.headers["Content-Type"] = "application/json"
        req.headers["X-API-Management-Key"] = @api_management_key
        req.body = payload.to_json
      end

      if response.success?
        body = JSON.parse(response.body)
        Rails.logger.info("API key created successfully via api-management: #{body['keyPrefix']}")

        # Return the full API key so the model can hash and store it
        body["apiKey"]
      else
        Rails.logger.error("API key creation failed: #{response.body}")
        nil
      end
    rescue Faraday::Error => e
      Rails.logger.error("API key creation error: #{e.message}")
      nil
    end

    # Revoke API key via api-management service
    def sync_delete
      response = connection.delete("/api-keys/#{@api_key.key_prefix}") do |req|
        req.headers["X-API-Management-Key"] = @api_management_key
      end

      handle_response(response, "revoke")
    rescue Faraday::Error => e
      Rails.logger.error("API key revocation error for #{@api_key.key_prefix}: #{e.message}")
      false
    end

    # Update API key plan via api-management service
    def sync_update
      payload = {
        plan: @api_key.user.current_plan,
        billingCycleStart: billing_cycle_start
      }

      response = connection.patch("/api-keys/#{@api_key.key_prefix}") do |req|
        req.headers["Content-Type"] = "application/json"
        req.headers["X-API-Management-Key"] = @api_management_key
        req.body = payload.to_json
      end

      handle_response(response, "update")
    rescue Faraday::Error => e
      Rails.logger.error("API key update error for #{@api_key.key_prefix}: #{e.message}")
      false
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
        Rails.logger.info("API key #{action} succeeded for #{@api_key.key_prefix}")
        true
      else
        Rails.logger.error("API key #{action} failed for #{@api_key.key_prefix}: #{response.body}")
        false
      end
    end
  end
end
