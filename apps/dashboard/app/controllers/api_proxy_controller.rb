# frozen_string_literal: true

require "base64"

class ApiProxyController < ApplicationController
  # Rate limiting handled by Rack::Attack (anonymous: 10/min, authenticated: 30/min)

  def create
    endpoint = params[:endpoint]
    method = params[:method]&.upcase || "GET"
    request_params = params[:params] || {}

    unless valid_endpoint?(endpoint)
      return render json: {
        error: "Invalid endpoint",
        message: "The endpoint must start with /v1/"
      }, status: :bad_request
    end

    start_time = Time.current
    result = make_api_request(endpoint, method, request_params)
    response_time = ((Time.current - start_time) * 1000).round

    render json: {
      status_code: result[:status_code],
      response_time_ms: response_time,
      data: result[:data],
      error: result[:error]
    }, status: result[:status_code]
  rescue StandardError => e
    Rails.logger.error("API Proxy Error: #{e.message}")
    Rails.logger.error(e.backtrace.join("\n"))

    render json: {
      error: "Proxy error",
      message: "Failed to connect to API: #{e.message}"
    }, status: :internal_server_error
  end

  private

  def valid_endpoint?(endpoint)
    return false if endpoint.blank?

    endpoint.match?(/\A\/v1\/[a-zA-Z0-9\/\-_.]+\z/)
  end

  def make_api_request(endpoint, method, request_params)
    # Convert ActionController::Parameters to a plain hash
    request_params = request_params.to_unsafe_h if request_params.respond_to?(:to_unsafe_h)

    internal_url = ::AppConfig.internal_api_url
    Rails.logger.debug { "Playground proxy: #{method} #{internal_url}#{endpoint} params=#{request_params.inspect}" }

    uri = URI(internal_url)
    uri.path = endpoint
    uri.query = nil

    # CF-Connecting-IP is set by Cloudflare with the real client IP.
    # request.remote_ip alone returns Cloudflare's edge node IP because the
    # dashboard sits behind Cloudflare and Rails has no trusted-proxy config.
    # Fall back to remote_ip for local/direct traffic (dev, internal).
    headers = {
      "X-Backend-Secret" => ::AppConfig.backend_secret,
      "Content-Type" => "application/json",
      "User-Agent" => "Requiems-Playground/1.0",
      "X-Forwarded-For" => request.headers["CF-Connecting-IP"] || request.remote_ip
    }

    http = Net::HTTP.new(uri.host, uri.port)
    http.use_ssl = uri.scheme == "https"
    http.open_timeout = 10
    http.read_timeout = 30

    request = case method
    when "GET"
                uri.query = URI.encode_www_form(request_params) if request_params.any?
                Net::HTTP::Get.new(uri)
    when "POST"
                req = Net::HTTP::Post.new(uri)
                req.body = request_params.to_json
                req
    else
                raise "Unsupported HTTP method: #{method}"
    end

    headers.each { |key, value| request[key] = value }

    response = http.request(request)

    {
      status_code: response.code.to_i,
      data: parse_response_body(response),
      error: response.is_a?(Net::HTTPSuccess) ? nil : response.message
    }
  rescue Net::OpenTimeout, Net::ReadTimeout => e
    {
      status_code: 504,
      data: nil,
      error: "Request timeout: #{e.message}"
    }
  rescue StandardError => e
    Rails.logger.error("Playground proxy error: #{e.class}: #{e.message}")
    Rails.logger.error(e.backtrace.first(3).join("\n"))
    {
      status_code: 500,
      data: nil,
      error: "Request failed: #{e.message}"
    }
  end

  def parse_response_body(response)
    content_type = response["Content-Type"].to_s
    body = response.body

    return nil if body.blank?

    if content_type.start_with?("image/", "application/octet-stream")
      return { "type" => "image", "content_type" => content_type, "base64" => Base64.strict_encode64(body) }
    end

    JSON.parse(body)
  rescue JSON::ParserError
    body.encode("UTF-8", invalid: :replace, undef: :replace)
  end
end
