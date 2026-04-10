# frozen_string_literal: true

class PrivateDeploymentRequest < ApplicationRecord
  belongs_to :user

  VALID_STATUSES = %w[pending_payment pending deploying active cancelled].freeze
  VALID_TIERS = %w[starter growth scale enterprise].freeze
  VALID_SERVICES = %w[email text tech places finance entertainment ai convert fitness misc].freeze

  TIER_PRICES = {
    "starter"    => 12000,
    "growth"     => 22000,
    "scale"      => 42000,
    "enterprise" => 85000
  }.freeze

  TIER_SPECS = {
    "starter"    => { hetzner: "CPX21", vcpu: 3, ram: "4 GB",  ssd: "80 GB"  },
    "growth"     => { hetzner: "CPX31", vcpu: 4, ram: "8 GB",  ssd: "160 GB" },
    "scale"      => { hetzner: "CPX41", vcpu: 8, ram: "16 GB", ssd: "240 GB" },
    "enterprise" => { hetzner: "CPX51", vcpu: 16, ram: "32 GB", ssd: "360 GB" }
  }.freeze

  validates :company, :contact_name, :contact_email, :server_tier, presence: true
  validates :server_tier, inclusion: { in: VALID_TIERS }
  validates :status, inclusion: { in: VALID_STATUSES }
  validate :at_least_one_service_selected
  validates :subdomain_slug, uniqueness: true, allow_nil: true,
            format: { with: /\A[a-z0-9\-]{2,40}\z/, message: "must be lowercase letters, numbers, and hyphens only (2–40 chars)" }
  validates :contact_email, format: { with: URI::MailTo::EMAIL_REGEXP }

  scope :pending,    -> { where(status: "pending") }
  scope :deploying,  -> { where(status: "deploying") }
  scope :active,     -> { where(status: "active") }
  scope :cancelled,  -> { where(status: "cancelled") }

  def monthly_price_dollars
    TIER_PRICES.fetch(server_tier, 0) / 100.0
  end

  def live_url
    "https://#{subdomain_slug}.requiems.xyz" if subdomain_slug.present?
  end

  def tier_specs
    TIER_SPECS.fetch(server_tier, {})
  end

  private

  def at_least_one_service_selected
    if selected_services.blank? || selected_services.empty?
      errors.add(:selected_services, "must include at least one service")
    end
  end
end
