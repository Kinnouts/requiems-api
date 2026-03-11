# frozen_string_literal: true

require_relative "boot"

require "rails"

require "active_model/railtie"
require "active_job/railtie"
require "active_record/railtie"
require "action_controller/railtie"
require "action_mailer/railtie"
require "action_view/railtie"

Bundler.require(*Rails.groups)

module Dashboard
  class Application < Rails::Application
    config.load_defaults 8.1

    config.autoload_lib(ignore: %w[assets tasks])

    config.generators.system_tests = nil

    config.middleware.use Rack::Attack

    config.i18n.available_locales = %i[en es]
    config.i18n.default_locale = :en
    config.i18n.fallbacks = true
    config.i18n.load_path += Dir[Rails.root.join("config/locales/**/*.yml")]

    # Fixed effective date for legal documents (Privacy Policy, Terms of Service).
    # Update this when the documents are materially revised.
    config.x.legal_effective_date = Date.new(2026, 2, 17)
  end
end
