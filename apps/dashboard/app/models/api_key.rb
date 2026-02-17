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
  before_validation :generate_key, on: :create
  after_create :sync_to_cloudflare
  after_destroy :remove_from_cloudflare
  after_update :sync_revocation_to_cloudflare, if: :saved_change_to_active?

  # Generate a new API key
  def generate_key
    env = environment&.to_sym || :live
    self.full_key = ApiKeyGenerator.generate(environment: env)
    self.key_prefix = ApiKeyGenerator.extract_prefix(full_key)
    self.key_hash = ApiKeyGenerator.hash_key(full_key)
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

  def sync_to_cloudflare
    return unless Rails.env.production? || ENV["SYNC_TO_CLOUDFLARE"] == "true"

    Cloudflare::KvSyncService.new(self).sync_create
  rescue StandardError => e
    Rails.logger.error("Failed to sync API key to Cloudflare: #{e.message}")
    # Don't fail the creation if Cloudflare sync fails
  end

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
