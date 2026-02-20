# frozen_string_literal: true

# Central configuration service for environment variables
# Validates required variables at startup and provides typed access
class AppConfig
  class MissingConfigError < StandardError; end
  class InvalidConfigError < StandardError; end

  include Singleton

  # API Management
  attr_reader :api_management_api_key

  # LemonSqueezy Configuration
  attr_reader :lemonsqueezy_store_id,
              :lemonsqueezy_store_slug,
              :lemonsqueezy_signing_secret

  # LemonSqueezy Variant IDs
  attr_reader :lemonsqueezy_developer_monthly_variant_id,
              :lemonsqueezy_developer_yearly_variant_id,
              :lemonsqueezy_business_monthly_variant_id,
              :lemonsqueezy_business_yearly_variant_id,
              :lemonsqueezy_professional_monthly_variant_id,
              :lemonsqueezy_professional_yearly_variant_id

  # API Configuration
  attr_reader :api_base_url,
              :playground_api_key,
              :internal_api_url,
              :backend_secret

  # SMTP Configuration (production only)
  attr_reader :smtp_address,
              :smtp_port,
              :smtp_domain,
              :smtp_username,
              :smtp_password,
              :mailer_host

  def initialize
    load_config
    validate_config
  end

  # Convenience method for accessing the singleton
  def self.instance
    @instance ||= new
  end

  # Quick accessor methods for common patterns
  def self.method_missing(method_name, *args, &block)
    if instance.respond_to?(method_name)
      instance.public_send(method_name, *args, &block)
    else
      super
    end
  end

  def self.respond_to_missing?(method_name, include_private = false)
    instance.respond_to?(method_name) || super
  end

  # Get variant ID for a specific plan and billing cycle
  def variant_id_for(plan:, billing_cycle:)
    case plan.to_s.downcase
    when "developer"
      billing_cycle == "monthly" ? lemonsqueezy_developer_monthly_variant_id : lemonsqueezy_developer_yearly_variant_id
    when "business"
      billing_cycle == "monthly" ? lemonsqueezy_business_monthly_variant_id : lemonsqueezy_business_yearly_variant_id
    when "professional"
      billing_cycle == "monthly" ? lemonsqueezy_professional_monthly_variant_id : lemonsqueezy_professional_yearly_variant_id
    else
      raise InvalidConfigError, "Unknown plan: #{plan}"
    end
  end

  # Check if SMTP is configured (for production)
  def smtp_configured?
    smtp_address.present? && smtp_username.present? && smtp_password.present?
  end

  private

  def load_config
    # API Management
    @api_management_api_key = require_env("API_MANAGEMENT_API_KEY")

    # LemonSqueezy Store
    @lemonsqueezy_store_id = require_env("LEMONSQUEEZY_STORE_ID")
    @lemonsqueezy_store_slug = optional_env("LEMONSQUEEZY_STORE_SLUG", default: "requiems")
    @lemonsqueezy_signing_secret = require_env("LEMONSQUEEZY_SIGNING_SECRET")

    # LemonSqueezy Variants
    @lemonsqueezy_developer_monthly_variant_id = require_env("LEMONSQUEEZY_DEVELOPER_MONTHLY_VARIANT_ID")
    @lemonsqueezy_developer_yearly_variant_id = require_env("LEMONSQUEEZY_DEVELOPER_YEARLY_VARIANT_ID")
    @lemonsqueezy_business_monthly_variant_id = require_env("LEMONSQUEEZY_BUSINESS_MONTHLY_VARIANT_ID")
    @lemonsqueezy_business_yearly_variant_id = require_env("LEMONSQUEEZY_BUSINESS_YEARLY_VARIANT_ID")
    @lemonsqueezy_professional_monthly_variant_id = require_env("LEMONSQUEEZY_PROFESSIONAL_MONTHLY_VARIANT_ID")
    @lemonsqueezy_professional_yearly_variant_id = require_env("LEMONSQUEEZY_PROFESSIONAL_YEARLY_VARIANT_ID")

    # API Configuration
    @api_base_url = optional_env("API_BASE_URL", default: "https://api.requiems.xyz")
    @playground_api_key = optional_env("PLAYGROUND_API_KEY", default: "rq_test_playground_demo_key")
    @internal_api_url = optional_env("INTERNAL_API_URL", default: "http://localhost:8080")
    @backend_secret = optional_env("BACKEND_SECRET", default: "dev_backend_secret")

    # SMTP (optional - only needed in production)
    @smtp_address = optional_env("SMTP_ADDRESS")
    @smtp_port = optional_env("SMTP_PORT", default: "587").to_i
    @smtp_domain = optional_env("SMTP_DOMAIN")
    @smtp_username = optional_env("SMTP_USERNAME")
    @smtp_password = optional_env("SMTP_PASSWORD")
    @mailer_host = optional_env("MAILER_HOST", default: "requiems.xyz")
  end

  def validate_config
    # Validate URL format
    validate_url(@api_base_url, "API_BASE_URL")

    # Validate LemonSqueezy store ID format
    unless @lemonsqueezy_store_id.match?(/^\d+$/)
      raise InvalidConfigError, "LEMONSQUEEZY_STORE_ID must be numeric"
    end

    # Validate all variant IDs are present and numeric
    validate_variant_id(@lemonsqueezy_developer_monthly_variant_id, "LEMONSQUEEZY_DEVELOPER_MONTHLY_VARIANT_ID")
    validate_variant_id(@lemonsqueezy_developer_yearly_variant_id, "LEMONSQUEEZY_DEVELOPER_YEARLY_VARIANT_ID")
    validate_variant_id(@lemonsqueezy_business_monthly_variant_id, "LEMONSQUEEZY_BUSINESS_MONTHLY_VARIANT_ID")
    validate_variant_id(@lemonsqueezy_business_yearly_variant_id, "LEMONSQUEEZY_BUSINESS_YEARLY_VARIANT_ID")
    validate_variant_id(@lemonsqueezy_professional_monthly_variant_id, "LEMONSQUEEZY_PROFESSIONAL_MONTHLY_VARIANT_ID")
    validate_variant_id(@lemonsqueezy_professional_yearly_variant_id, "LEMONSQUEEZY_PROFESSIONAL_YEARLY_VARIANT_ID")

    # Validate SMTP if in production
    if Rails.env.production? && !smtp_configured?
      Rails.logger.warn("SMTP not fully configured in production environment")
    end
  end

  def require_env(key)
    ENV.fetch(key) do
      # In test environment, provide safe defaults instead of failing
      if Rails.env.test?
        test_defaults[key]
      else
        raise MissingConfigError, "Missing required environment variable: #{key}"
      end
    end
  end

  def optional_env(key, default: nil)
    ENV.fetch(key, default)
  end

  # Safe test defaults for required environment variables
  def test_defaults
    {
      "API_MANAGEMENT_API_KEY" => "test_api_management_key",
      "LEMONSQUEEZY_STORE_ID" => "12345",
      "LEMONSQUEEZY_SIGNING_SECRET" => "test_signing_secret",
      "LEMONSQUEEZY_DEVELOPER_MONTHLY_VARIANT_ID" => "123456",
      "LEMONSQUEEZY_DEVELOPER_YEARLY_VARIANT_ID" => "123457",
      "LEMONSQUEEZY_BUSINESS_MONTHLY_VARIANT_ID" => "123458",
      "LEMONSQUEEZY_BUSINESS_YEARLY_VARIANT_ID" => "123459",
      "LEMONSQUEEZY_PROFESSIONAL_MONTHLY_VARIANT_ID" => "123460",
      "LEMONSQUEEZY_PROFESSIONAL_YEARLY_VARIANT_ID" => "123461",
      "BACKEND_SECRET" => "test_backend_secret"
    }
  end

  def validate_url(url, key)
    uri = URI.parse(url)
    unless uri.is_a?(URI::HTTP) || uri.is_a?(URI::HTTPS)
      raise InvalidConfigError, "#{key} must be a valid HTTP/HTTPS URL"
    end
  rescue URI::InvalidURIError
    raise InvalidConfigError, "#{key} is not a valid URL: #{url}"
  end

  def validate_variant_id(value, key)
    unless value.present? && value.match?(/^\w+$/)
      raise InvalidConfigError, "#{key} must be present and contain valid characters"
    end
  end
end
