# frozen_string_literal: true

module Cloudflare
  class KvSyncService
    CLOUDFLARE_API_BASE = "https://api.cloudflare.com/client/v4"

    def initialize(api_key_record)
      @api_key = api_key_record
      @account_id = ENV.fetch("CLOUDFLARE_ACCOUNT_ID")
      @namespace_id = ENV.fetch("CLOUDFLARE_KV_NAMESPACE_ID")
      @api_token = ENV.fetch("CLOUDFLARE_API_TOKEN")
    end

    # Sync API key to Cloudflare KV when created
    def sync_create
      key_name = "key:#{@api_key.full_key}"
      value = {
        userId: @api_key.user_id.to_s,
        plan: @api_key.user.current_plan,
        createdAt: @api_key.created_at.iso8601
      }

      response = connection.put(values_path(key_name)) do |req|
        req.headers["Content-Type"] = "application/json"
        req.body = value.to_json
      end

      handle_response(response, "create")
    end

    # Remove API key from Cloudflare KV when deleted
    def sync_delete
      key_name = "key:#{@api_key.full_key}"

      response = connection.delete(values_path(key_name))

      handle_response(response, "delete")
    end

    private

    def connection
      @connection ||= Faraday.new(
        url: "#{CLOUDFLARE_API_BASE}/accounts/#{@account_id}/storage/kv/namespaces/#{@namespace_id}",
        headers: {
          "Authorization" => "Bearer #{@api_token}",
          "Content-Type" => "application/json"
        }
      ) do |f|
        f.request :retry, max: 3, interval: 0.5
        f.adapter Faraday.default_adapter
      end
    end

    def values_path(key_name)
      "/values/#{key_name}"
    end

    def handle_response(response, action)
      if response.success?
        Rails.logger.info("Cloudflare KV #{action} succeeded for API key #{@api_key.key_prefix}")
        true
      else
        Rails.logger.error("Cloudflare KV #{action} failed for API key #{@api_key.key_prefix}: #{response.body}")
        false
      end
    rescue Faraday::Error => e
      Rails.logger.error("Cloudflare KV #{action} error for API key #{@api_key.key_prefix}: #{e.message}")
      false
    end
  end
end
