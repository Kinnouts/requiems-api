# frozen_string_literal: true

ENV["RAILS_ENV"] ||= "test"
require_relative "../config/environment"
require "rails/test_help"

# Ensure Devise mappings are loaded
Rails.application.reload_routes!

module ActiveSupport
  class TestCase
    # Run tests in parallel with specified workers
    parallelize(workers: :number_of_processors)

    # Setup all fixtures in test/fixtures/*.yml for all tests in alphabetical order.
    fixtures :all

    # Password used by create_user for tests that need to sign in via the login form.
    TEST_USER_PASSWORD = "password123!"

    # Helper to create confirmed users for tests (Devise confirmable).
    def create_user(email: "test@example.com", **attributes)
      User.create!(
        email: email,
        password: TEST_USER_PASSWORD,
        password_confirmation: TEST_USER_PASSWORD,
        confirmed_at: Time.current,
        **attributes
      )
    end
  end
end

# Include Devise test helpers for integration tests
class ActionDispatch::IntegrationTest
  include Devise::Test::IntegrationHelpers

  setup do
    # Use the default locale explicitly so requests include the locale prefix
    # (e.g. /en/admin/users) and bypass the set_locale redirect that enforces
    # canonical locale-prefixed URLs. Using a symbol/string avoids the locale
    # slot accidentally consuming a positional model argument in path helpers.
    self.default_url_options = { locale: I18n.default_locale }
  end
end
