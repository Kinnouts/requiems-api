# frozen_string_literal: true

class PrivateDeploymentRequest < ApplicationRecord
  belongs_to :user

  VALID_STATUSES = %w[pending_payment pending deploying active cancelled].freeze
  VALID_TIERS = %w[starter growth scale enterprise].freeze
  VALID_SERVICES = %w[email text tech places finance entertainment ai convert fitness misc].freeze

  SERVICE_META = {
    "email"         => "Validation, disposable check, normalization",
    "text"          => "Quotes, dictionary, lorem ipsum, advice, profanity, spellcheck",
    "tech"          => "IP/VPN info, phone, QR codes, domain, WHOIS, MX lookup",
    "places"        => "Cities, geocoding, timezone, holidays, postal codes",
    "finance"       => "BIN lookup, crypto, exchange rates, IBAN, SWIFT, mortgage",
    "entertainment" => "Jokes, facts, horoscope, trivia, emoji, sudoku",
    "ai"            => "Similarity, sentiment analysis, language detection",
    "convert"       => "Base64, color, markdown, number base, format conversion",
    "fitness"       => "Exercise database by body part, equipment, muscle",
    "misc"          => "Counter, random user, unit conversion"
  }.freeze

  TIER_PRICES_MONTHLY = {
    "starter"    => 20000,
    "growth"     => 30000,
    "scale"      => 50000,
    "enterprise" => 100000
  }.freeze

  # 15% yearly discount, billed as one annual charge
  TIER_PRICES_YEARLY = {
    "starter"    => 204000,   # $2,040/yr  ($170/mo)
    "growth"     => 306000,   # $3,060/yr  ($255/mo)
    "scale"      => 510000,   # $5,100/yr  ($425/mo)
    "enterprise" => 1020000   # $10,200/yr ($850/mo)
  }.freeze

  TIER_SPECS = {
    "starter"    => { hetzner: "CPX21", vcpu: 3, ram: "4 GB",  ssd: "80 GB"  },
    "growth"     => { hetzner: "CPX31", vcpu: 4, ram: "8 GB",  ssd: "160 GB" },
    "scale"      => { hetzner: "CPX41", vcpu: 8, ram: "16 GB", ssd: "240 GB" },
    "enterprise" => { hetzner: "CPX51", vcpu: 16, ram: "32 GB", ssd: "360 GB" }
  }.freeze

  validates :contact_name, :contact_email, :server_tier, :billing_cycle, presence: true
  validates :server_tier, inclusion: { in: VALID_TIERS }
  validates :billing_cycle, inclusion: { in: %w[monthly yearly] }
  validates :status, inclusion: { in: VALID_STATUSES }
  validates :lemonsqueezy_subscription_id, uniqueness: true, allow_nil: true
  validate :at_least_one_service_selected
  validates :subdomain_slug, uniqueness: true, allow_nil: true,
            format: { with: /\A[a-z0-9\-]{2,40}\z/, message: "must be lowercase letters, numbers, and hyphens only (2–40 chars)" }
  validates :contact_email, format: { with: URI::MailTo::EMAIL_REGEXP }

  scope :pending,    -> { where(status: "pending") }
  scope :deploying,  -> { where(status: "deploying") }
  scope :active,     -> { where(status: "active") }
  scope :cancelled,  -> { where(status: "cancelled") }

  def monthly_price_dollars
    monthly_prices = TIER_PRICES_MONTHLY.fetch(server_tier, 0) / 100.0
    return monthly_prices if billing_cycle == "monthly"
    # For yearly, show the effective per-month rate
    TIER_PRICES_YEARLY.fetch(server_tier, 0) / 100.0 / 12.0
  end

  def total_price_dollars
    price_table = billing_cycle == "yearly" ? TIER_PRICES_YEARLY : TIER_PRICES_MONTHLY
    price_table.fetch(server_tier, 0) / 100.0
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
