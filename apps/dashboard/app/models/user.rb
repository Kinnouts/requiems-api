# frozen_string_literal: true

class User < ApplicationRecord
  # Include default devise modules. Others available are:
  # :confirmable, :lockable, :timeoutable, :trackable and :omniauthable
  devise :database_authenticatable, :registerable,
         :recoverable, :rememberable, :validatable,
         :confirmable, :trackable

  # Associations
  has_many :api_keys, dependent: :destroy
  has_one :subscription, dependent: :destroy
  has_many :usage_logs, dependent: :destroy
  has_many :daily_usage_summaries, dependent: :destroy
  has_many :credit_adjustments, dependent: :destroy
  has_many :audit_logs, dependent: :destroy
  has_many :abuse_reports, dependent: :destroy

  # Scopes
  scope :admins, -> { where(admin: true) }
  scope :active_users, -> { where(status: "active") }
  scope :suspended, -> { where(status: "suspended") }
  scope :banned, -> { where(status: "banned") }

  # Admin methods
  def admin?
    admin == true
  end

  def active_user?
    status == "active" && !banned_at
  end

  def suspended?
    status == "suspended"
  end

  def banned?
    status == "banned" || banned_at.present?
  end

  def current_plan
    subscription&.plan_name || "free"
  end

  def credit_limit
    subscription&.credit_limit || 50 # Free tier default
  end

  # Ban/suspend methods
  def ban!(reason:, admin_user:)
    transaction do
      update!(
        status: "banned",
        banned_at: Time.current,
        banned_reason: reason,
        active: false
      )

      api_keys.update_all(active: false, revoked_at: Time.current)

      AuditLog.create!(
        user: self,
        admin_user: admin_user,
        action: "ban_user",
        details: { reason: reason }.to_json
      )
    end
  end

  def suspend!(admin_user:)
    update!(status: "suspended", active: false)
    AuditLog.create!(user: self, admin_user: admin_user, action: "suspend_user")
  end

  def unsuspend!(admin_user:)
    update!(status: "active", active: true)
    AuditLog.create!(user: self, admin_user: admin_user, action: "unsuspend_user")
  end

  # Account deletion request methods
  def request_account_deletion!(reason)
    update!(
      deletion_token: SecureRandom.urlsafe_base64(32),
      deletion_token_sent_at: Time.current,
      deletion_reason: reason
    )
  end

  def deletion_token_valid?(token)
    deletion_token.present? &&
      deletion_token_sent_at.present? &&
      deletion_token_sent_at > 1.hour.ago &&
      ActiveSupport::SecurityUtils.secure_compare(deletion_token, token)
  end

  def clear_deletion_token!
    update_columns(deletion_token: nil, deletion_token_sent_at: nil, deletion_reason: nil)
  end
end
