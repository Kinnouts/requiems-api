class ApiKey < ApplicationRecord
  belongs_to :user
  has_many :usage_logs, dependent: :destroy
  has_many :abuse_reports, dependent: :destroy

  # Virtual attribute to store the full key temporarily (only shown once to user)
  attr_accessor :full_key

  # Validations
  validates :key_prefix, presence: true, uniqueness: true
  validates :key_hash, presence: true
  validates :name, presence: true, length: { maximum: 100 }
  validates :environment, inclusion: { in: %w[test live] }, allow_nil: true

  # Scopes
  scope :active_keys, -> { where(active: true, revoked_at: nil) }
  scope :revoked, -> { where.not(revoked_at: nil) }
  scope :for_environment, ->(env) { where(environment: env) }

  # Callbacks
  before_validation :request_key_from_server, on: :create
  after_destroy :remove_from_cloudflare
  after_update :sync_revocation_to_cloudflare, if: :saved_change_to_active?

  # Request a new API key from api-management server
  # The server generates the key and returns it (only once)
  def request_key_from_server
    return if key_prefix.present? # Skip if already generated

    # In test/development, generate keys locally without external API calls
    if Rails.env.test? || Rails.env.development?
      generate_key_locally
    else
      # In production, request from api-management service
      service = Cloudflare::KvSyncService.new(self)
      generated_key = service.sync_create

      if generated_key
        # Store the returned key
        self.full_key = generated_key
        self.key_prefix = ApiKeyGenerator.extract_prefix(generated_key)
        self.key_hash = ApiKeyGenerator.hash_key(generated_key)
        self.active = true if active.nil?
      else
        # Key generation failed
        errors.add(:base, "Failed to generate API key")
        throw :abort
      end
    end
  end

  def generate_key_locally
    # Generate a local API key for test/development
    env = environment&.to_sym || :live
    generated_key = ApiKeyGenerator.generate(environment: env)
    self.full_key = generated_key
    self.key_prefix = ApiKeyGenerator.extract_prefix(generated_key)
    self.key_hash = ApiKeyGenerator.hash_key(generated_key)
    self.active = true if active.nil?
  end

  # Verify a key matches this record
  def verify_key(key_to_verify)
    ApiKeyGenerator.verify_key(key_to_verify, key_hash)
  end

  # Revoke the API key
  def revoke!(reason: nil)
    update!(
      active: false,
      revoked_at: Time.current,
      revoked_reason: reason
    )
  end

  # Display-friendly key (shows prefix + masked)
  def masked_key
    "#{key_prefix}••••••••••••"
  end

  private

  def remove_from_cloudflare
    return unless Rails.env.production? || ENV["SYNC_TO_CLOUDFLARE"] == "true"

    Cloudflare::KvSyncService.new(self).sync_delete
  rescue StandardError => e
    Rails.logger.error("Failed to remove API key from Cloudflare: #{e.message}")
  end

  def sync_revocation_to_cloudflare
    return unless !active && revoked_at

    remove_from_cloudflare
  end
end
