Sentry.init do |config|
  config.dsn = ENV.fetch("SENTRY_DSN", "https://8cef11165f1846f39d700ecb0cf781cf@issues.bobadilla.tech/2")
  config.environment = Rails.env
  config.enabled_environments = %w[production staging]
  config.traces_sample_rate = 0.01
  config.breadcrumbs_logger = [ :active_support_logger, :http_logger ]
end
