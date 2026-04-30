# frozen_string_literal: true

# This file should ensure the existence of records required to run the application in every environment (production,
# development, test). The code here should be idempotent so that it can be executed at any point in every environment.
# The data can then be loaded with the bin/rails db:seed command (or created alongside the database with db:setup).

# Development: Create test user if it doesn't exist
if Rails.env.development?
  test_user = User.find_or_initialize_by(email: "eliaz.bobadilladeva@gmail.com")
  if test_user.new_record?
    test_user.password = SecureRandom.hex(16)
    test_user.admin = true
    test_user.save!
    puts "✓ Created test user: #{test_user.email} (admin: true)"
  else
    puts "✓ Test user already exists: #{test_user.email}"
  end

  # Create additional test user
  user2 = User.find_or_initialize_by(email: "test@example.com")
  if user2.new_record?
    user2.password = "password123!"
    user2.save!
    puts "✓ Created additional test user: #{user2.email}"
  end
end
