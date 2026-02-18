# frozen_string_literal: true

# Load and validate application configuration at startup
# This ensures the app fails fast if required environment variables are missing
Rails.application.config.to_prepare do
  begin
    AppConfig.instance
    Rails.logger.info("AppConfig initialized successfully")
  rescue AppConfig::MissingConfigError, AppConfig::InvalidConfigError => e
    Rails.logger.error("Configuration Error: #{e.message}")
    raise e if Rails.env.production? || Rails.env.test?
    # In development, we'll allow the app to start but log the warning
    Rails.logger.warn("Starting in development mode despite config errors")
  end
end
