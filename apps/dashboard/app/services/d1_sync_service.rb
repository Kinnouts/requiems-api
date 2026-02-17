# frozen_string_literal: true

# Service to sync usage data from Cloudflare D1 to PostgreSQL
#
# Usage:
#   service = D1SyncService.new
#   result = service.fetch_usage(since: 1.hour.ago)
#
class D1SyncService
  class Error < StandardError; end
  class UnauthorizedError < Error; end
  class InvalidResponseError < Error; end

  API_MANAGEMENT_URL = ENV.fetch("API_MANAGEMENT_URL", "https://api-management.requiems.xyz")
  TIMEOUT = 30 # seconds
  MAX_RETRIES = 3

  def self.api_management_key
    ENV.fetch("API_MANAGEMENT_API_KEY")
  end

  # Fetch usage data from Cloudflare D1
  #
  # @param since [Time, String] Fetch records after this timestamp
  # @param limit [Integer] Max records per request (default: 1000)
  # @return [Hash] Response with :usage array and pagination info
  def fetch_usage(since:, limit: 1000)
    since_iso = since.is_a?(String) ? since : since.iso8601

    all_records = []
    cursor = nil
    has_more = true

    while has_more
      response = fetch_page(since: since_iso, limit: limit, cursor: cursor)

      all_records.concat(response[:usage])

      has_more = response[:has_more]
      cursor = response[:next_cursor]

      # Safety: stop after 100 pages (100k records max per call)
      break if all_records.size >= 100_000
    end

    {
      usage: all_records,
      total: all_records.size
    }
  rescue Faraday::Error => e
    Rails.logger.error("D1 sync failed: #{e.class} - #{e.message}")
    raise Error, "Failed to fetch usage from Cloudflare: #{e.message}"
  end

  # Bulk insert usage records into PostgreSQL
  #
  # @param records [Array<Hash>] Usage records from D1
  # @return [Integer] Number of records inserted
  def bulk_insert(records)
    return 0 if records.empty?

    # Map D1 records to UsageLog format
    values = records.map do |record|
      {
        api_key_id: resolve_api_key_id(record[:api_key]),
        user_id: resolve_user_id(record[:api_key]),
        endpoint: record[:endpoint],
        credits_used: record[:credits_used],
        status_code: 200, # D1 only records successful requests
        response_time_ms: nil, # Not tracked in D1
        used_at: Time.parse(record[:used_at]),
        usage_date: Date.parse(record[:used_at]),
        created_at: Time.current,
        updated_at: Time.current
      }
    end.compact

    # Bulk insert using ActiveRecord
    UsageLog.insert_all(values, unique_by: [ :api_key_id, :used_at, :endpoint ])

    values.size
  end

  private

  # Fetch a single page of usage data
  def fetch_page(since:, limit:, cursor: nil)
    url = "#{API_MANAGEMENT_URL}/usage/export"

    params = {
      since: since,
      limit: limit
    }
    params[:cursor] = cursor if cursor

    response = connection.get(url, params)

    unless response.success?
      handle_error_response(response)
    end

    parse_response(response.body)
  end

  # Create Faraday connection with retries
  def connection
    @connection ||= Faraday.new do |conn|
      conn.request :url_encoded
      conn.request :retry, {
        max: MAX_RETRIES,
        interval: 0.5,
        backoff_factor: 2,
        retry_statuses: [ 500, 502, 503, 504 ],
        methods: [ :get ]
      }
      conn.response :json, content_type: /\bjson$/
      conn.adapter Faraday.default_adapter
      conn.options.timeout = TIMEOUT
      conn.headers["X-API-Management-Key"] = self.class.api_management_key
      conn.headers["Content-Type"] = "application/json"
    end
  end

  def parse_response(body)
    {
      usage: body["usage"] || [],
      total: body["total"] || 0,
      has_more: body["hasMore"] || false,
      next_cursor: body["nextCursor"]
    }
  rescue => e
    raise InvalidResponseError, "Invalid response format: #{e.message}"
  end

  def handle_error_response(response)
    case response.status
    when 401
      raise UnauthorizedError, "Invalid API management key"
    when 400
      raise Error, "Bad request: #{response.body}"
    else
      raise Error, "HTTP #{response.status}: #{response.body}"
    end
  end

  # Resolve API key string to database ID
  # Uses in-memory cache to avoid repeated lookups
  def resolve_api_key_id(key_string)
    @api_key_cache ||= {}
    @api_key_cache[key_string] ||= begin
      # Extract prefix from full key (first 12 chars)
      prefix = key_string[0...12]
      api_key = ApiKey.find_by(key_prefix: prefix)
      api_key&.id
    end
  end

  # Resolve API key to user ID
  def resolve_user_id(key_string)
    key_id = resolve_api_key_id(key_string)
    return nil unless key_id

    @user_id_cache ||= {}
    @user_id_cache[key_id] ||= ApiKey.find(key_id).user_id
  end
end
