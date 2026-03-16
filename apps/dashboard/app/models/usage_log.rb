# frozen_string_literal: true

class UsageLog < ApplicationRecord
  belongs_to :user
  belongs_to :api_key

  scope :this_month, -> { where("used_at >= ?", Time.current.beginning_of_month) }
  scope :last_7_days, -> { where("used_at >= ?", 7.days.ago) }
  scope :with_response_time, -> { where.not(response_time_ms: nil) }
  scope :recent, ->(limit = 10) { order(used_at: :desc).limit(limit).includes(:api_key) }
  validates :user_id, :used_at, :endpoint, presence: true
  validates :user_id, :used_at, :endpoint, presence: true

  def self.error_rate_for(scope)
    total = scope.count
    return 0 if total.zero?

    errors = scope.where("status_code >= ?", 400).count
    ((errors.to_f / total) * 100).round(2)
  end
end
