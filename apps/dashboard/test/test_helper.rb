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
    # Ensure locale is not injected as first positional arg in URL helpers.
    # With scope "(:locale)" routes, omitting this causes positional args like
    # path(@model) to land in the :locale slot instead of :id.
    self.default_url_options = { locale: nil }
  end
end
