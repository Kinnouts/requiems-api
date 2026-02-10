# frozen_string_literal: true

class Rack::Attack
  ### Configure Cache ###

  # If you don't want to use Rails.cache (Solid Cache), you can set it to Redis
  # Rack::Attack.cache.store = ActiveSupport::Cache::RedisStore.new

  # Use Rails cache (Solid Cache in Rails 8)
  Rack::Attack.cache.store = Rails.cache

  ### Throttle Configuration ###

  # Throttle API playground requests for anonymous users
  # Allow 10 requests per minute per IP
  throttle("api_proxy/ip", limit: 10, period: 1.minute) do |req|
    if req.path == "/api/proxy" && req.post?
      # Return IP to throttle by IP address
      # Don't throttle authenticated users
      req.ip unless req.env["warden"]&.user
    end
  end

  # Optional: General login throttling
  # Throttle login attempts by IP address
  throttle("logins/ip", limit: 5, period: 20.seconds) do |req|
    if req.path == "/users/sign_in" && req.post?
      req.ip
    end
  end

  # Optional: Throttle login attempts by email
  throttle("logins/email", limit: 5, period: 20.seconds) do |req|
    if req.path == "/users/sign_in" && req.post?
      # Normalize email
      req.params["user"]&.dig("email")&.to_s&.downcase&.gsub(/\s+/, "")
    end
  end

  ### Custom Throttle Response ###

  # Customize the response when throttled
  self.throttled_responder = lambda do |request|
    match_data = request.env["rack.attack.match_data"]
    now = match_data[:epoch_time]

    headers = {
      "RateLimit-Limit" => match_data[:limit].to_s,
      "RateLimit-Remaining" => "0",
      "RateLimit-Reset" => (now + (match_data[:period] - (now % match_data[:period]))).to_s,
      "Content-Type" => "application/json"
    }

    body = {
      error: "Rate limit exceeded",
      message: "Too many requests. Please try again later.",
      retry_after: match_data[:period] - (now % match_data[:period])
    }.to_json

    [429, headers, [body]]
  end

  ### Logging ###

  # Log blocked requests
  ActiveSupport::Notifications.subscribe("throttle.rack_attack") do |_name, _start, _finish, _request_id, payload|
    req = payload[:request]
    Rails.logger.warn("[Rack::Attack][Throttled] IP: #{req.ip} | Path: #{req.path}")
  end
end
