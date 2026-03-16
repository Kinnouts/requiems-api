# frozen_string_literal: true

class UsageLog < ApplicationRecord
  belongs_to :user
  belongs_to :api_key

  scope :this_month, -> { where("used_at >= ?", Time.current.beginning_of_month) }
  scope :last_7_days, -> { where("used_at >= ?", 7.days.ago) }
  scope :with_response_time, -> { where.not(response_time_ms: nil) }
  scope :recent, ->(limit = 10) { order(used_at: :desc).limit(limit).includes(:api_key) }
end
