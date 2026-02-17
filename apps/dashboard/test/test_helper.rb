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

    # Helper to create confirmed users for tests (Devise confirmable)
    def create_user(email: "test@example.com", password: "password123", **attributes)
      User.create!(
        email: email,
        password: password,
        password_confirmation: password,
        confirmed_at: Time.current,
        **attributes
      )
    end
  end
end

# Include Devise test helpers for integration tests
class ActionDispatch::IntegrationTest
  include Devise::Test::IntegrationHelpers
end
