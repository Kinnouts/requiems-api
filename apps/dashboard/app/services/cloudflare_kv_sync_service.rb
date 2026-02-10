# frozen_string_literal: true

class CloudflareKvSyncService
  class << self
    def sync_user_plan(user, plan_name)
      # Update all active API keys for this user in Cloudflare KV
      user.api_keys.active_keys.find_each do |api_key|
        sync_api_key(api_key, plan_name)
      end
    end

    def sync_api_key(api_key, plan_name = nil)
      plan_name ||= api_key.user.subscription&.plan_name || "free"

      key_data = {
        userId: api_key.user_id.to_s,
        plan: plan_name,
        billingCycleStart: billing_cycle_start(api_key.user).iso8601,
        createdAt: api_key.created_at.iso8601
      }

      update_kv(api_key.key_value, key_data)

      Rails.logger.info "[CloudflareKV] Synced API key #{api_key.id} to plan: #{plan_name}"
    end

    private

    def billing_cycle_start(user)
      subscription = user.subscription
      return Time.current.beginning_of_month unless subscription

      subscription.current_period_start || Time.current.beginning_of_month
    end

    def update_kv(api_key, data)
      account_id = ENV["CLOUDFLARE_ACCOUNT_ID"]
      namespace_id = ENV["CLOUDFLARE_KV_NAMESPACE_ID"]
      api_token = ENV["CLOUDFLARE_API_TOKEN"]

      url = "https://api.cloudflare.com/client/v4/accounts/#{account_id}/storage/kv/namespaces/#{namespace_id}/values/key:#{api_key}"

      response = HTTP
        .auth("Bearer #{api_token}")
        .headers("Content-Type" => "application/json")
        .put(url, json: data)

      unless response.status.success?
        Rails.logger.error "[CloudflareKV] Failed to update key: #{response.body}"
        raise "Failed to sync to Cloudflare KV: #{response.status}"
      end
    end
  end
end
