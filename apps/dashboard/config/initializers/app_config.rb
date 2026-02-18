# frozen_string_literal: true

# Load and validate application configuration at startup
# This ensures the app fails fast if required environment variables are missing
Rails.application.config.to_prepare do
  # Skip config validation for rake tasks that don't need the full app
  # (db:create, db:migrate, assets:precompile, etc.)
  next if defined?(Rake) && Rake.application.top_level_tasks.any? { |task|
    task.start_with?("db:", "assets:", "tmp:", "log:", "about")
  }

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
