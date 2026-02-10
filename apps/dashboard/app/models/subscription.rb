class Subscription < ApplicationRecord
  belongs_to :user

  # Validations
  validates :plan_name, presence: true, inclusion: { in: %w[free developer business professional] }
  validates :lemonsqueezy_subscription_id, uniqueness: true, allow_nil: true
  validates :status, presence: true

  # Scopes
  scope :active, -> { where(status: %w[active trialing]) }
  scope :cancelled, -> { where(status: "cancelled") }

  # Callbacks
  after_create :sync_to_cloudflare
  after_update :sync_to_cloudflare, if: :saved_change_to_plan_name?

  private

  def sync_to_cloudflare
    CloudflareKvSyncService.sync_user_plan(user, plan_name)
  rescue StandardError => e
    Rails.logger.error "[Subscription] Failed to sync to Cloudflare KV: #{e.message}"
    # Don't raise - subscription should still be saved even if KV sync fails
  end
end
