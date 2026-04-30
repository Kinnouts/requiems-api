# frozen_string_literal: true

Rails.application.configure do
  config.enable_reloading = false
  config.eager_load = ENV["CI"].present?
  config.public_file_server.headers = { "cache-control" => "public, max-age=3600" }
  config.consider_all_requests_local = true
  config.cache_store = :null_store
  config.action_dispatch.show_exceptions = :rescuable
  config.action_controller.allow_forgery_protection = false
  config.action_mailer.delivery_method = :test
  config.action_mailer.default_url_options = { host: "example.com", locale: "en" }
  config.active_support.deprecation = :stderr
  config.action_controller.raise_on_missing_callback_actions = true

  # Fixed test keys for ActiveRecord encryption (tenant_secret on PrivateDeploymentRequest).
  # These are test-only values and must never be used in production.
  config.active_record.encryption.primary_key         = "test-primary-key-for-rails-testing!!"
  config.active_record.encryption.deterministic_key   = "test-deterministic-key-rails-tests"
  config.active_record.encryption.key_derivation_salt = "test-key-derivation-salt-for-rails"
end
