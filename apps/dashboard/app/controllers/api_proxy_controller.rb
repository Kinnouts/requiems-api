# frozen_string_literal: true

class ApiProxyController < ApplicationController
  # Rate limiting handled by Rack::Attack for anonymous users
  # No rate limit for authenticated users (handled by API gateway)

  def create
    # Extract request parameters
    endpoint = params[:endpoint]
    method = params[:method]&.upcase || "POST"
    request_params = params[:params] || {}
    request_headers = params[:headers] || {}

    # Validate endpoint
    unless valid_endpoint?(endpoint)
      return render json: {
        error: "Invalid endpoint",
        message: "The endpoint must start with /v1/"
      }, status: :bad_request
    end

    # Get API key
    api_key = get_api_key

    unless api_key
      return render json: {
        error: "Authentication required",
        message: "Please sign in or use the test playground key"
      }, status: :unauthorized
    end

    # Make the API request
    start_time = Time.current
    response = make_api_request(endpoint, method, request_params, api_key, request_headers)
    response_time = ((Time.current - start_time) * 1000).round

    # Return the response with metadata
    render json: {
      success: response[:success],
      status_code: response[:status_code],
      response_time_ms: response_time,
      data: response[:data],
      error: response[:error],
      request: {
        endpoint: endpoint,
        method: method,
        params: request_params,
        headers: sanitize_headers(request_headers)
      }
    }, status: response[:status_code]
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
    return false unless endpoint.start_with?("/v1/")
    # Reject anything that could manipulate URI parsing (schemes, fragments, newlines, path traversal)
    return false if endpoint.match?(/[#\r\n]|:\/\//)
    return false if endpoint.include?("..")

    true
  end

  def get_api_key
    if user_signed_in?
      # Use user's first active API key
      current_user.api_keys.active_keys.first&.full_key
    else
      # Use test/demo API key for anonymous users
      # This should be a special key with limited access and rate limits
      AppConfig.playground_api_key
    end
  end

  def make_api_request(endpoint, method, request_params, api_key, request_headers)
    api_base_url = AppConfig.api_base_url

    base_uri = URI(api_base_url)
    uri = base_uri.dup
    uri.path = endpoint
    uri.query = nil

    # Prepare headers
    headers = {
      "Authorization" => "Bearer #{api_key}",
      "Content-Type" => "application/json",
      "User-Agent" => "Requiems-Playground/1.0"
    }

    # Add custom headers (sanitized)
    request_headers.each do |key, value|
      headers[key] = value if allowed_header?(key)
    end

    http = Net::HTTP.new(uri.host, uri.port)
    http.use_ssl = uri.scheme == "https"
    http.open_timeout = 10
    http.read_timeout = 30

    # Create request based on method
    request = case method
    when "GET"
                # For GET, add params to query string
                uri.query = URI.encode_www_form(request_params) if request_params.any?
                Net::HTTP::Get.new(uri)
    when "POST"
                req = Net::HTTP::Post.new(uri)
                req.body = request_params.to_json
                req
    when "PUT"
                req = Net::HTTP::Put.new(uri)
                req.body = request_params.to_json
                req
    when "DELETE"
                Net::HTTP::Delete.new(uri)
    else
                raise "Unsupported HTTP method: #{method}"
    end

    # Set headers
    headers.each { |key, value| request[key] = value }

    # Execute request
    response = http.request(request)

    # Parse response
    {
      success: response.is_a?(Net::HTTPSuccess),
      status_code: response.code.to_i,
      data: parse_response_body(response.body),
      error: response.is_a?(Net::HTTPSuccess) ? nil : response.message
    }
  rescue Net::OpenTimeout, Net::ReadTimeout => e
    {
      success: false,
      status_code: 504,
      data: nil,
      error: "Request timeout: #{e.message}"
    }
  rescue StandardError => e
    {
      success: false,
      status_code: 500,
      data: nil,
      error: "Request failed: #{e.message}"
    }
  end

  def parse_response_body(body)
    return nil if body.blank?

    JSON.parse(body)
  rescue JSON::ParserError
    # Return raw body if not JSON
    body
  end

  def allowed_header?(header_name)
    # Only allow safe headers
    safe_headers = %w[
      Content-Type
      Accept
      Accept-Language
      X-Request-ID
    ]

    safe_headers.any? { |h| h.casecmp(header_name).zero? }
  end

  def sanitize_headers(headers)
    # Remove sensitive headers from response
    headers.except("Authorization", "Cookie", "X-Api-Key")
  end
end
