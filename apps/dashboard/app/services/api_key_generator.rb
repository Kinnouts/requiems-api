# frozen_string_literal: true

require "bcrypt"

class ApiKeyGenerator
  # Generate a new API key in the format: rq_live_<24_random_chars>
  # Returns the full key (which should be shown to the user once)
  def self.generate(environment: :live)
    prefix = environment == :test ? "rq_test" : "rq_live"
    random_part = Nanoid.generate(size: 24, alphabet: "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

    "#{prefix}_#{random_part}"
  end

  # Extract the prefix for display (first 12 characters)
  def self.extract_prefix(full_key)
    full_key[0..11] if full_key
  end

  # Hash the full key for secure storage
  def self.hash_key(full_key)
    BCrypt::Password.create(full_key)
  end

  # Verify a key against a hash
  def self.verify_key(full_key, key_hash)
    BCrypt::Password.new(key_hash) == full_key
  rescue BCrypt::Errors::InvalidHash
    false
  end
end
