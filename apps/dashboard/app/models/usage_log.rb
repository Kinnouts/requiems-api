# frozen_string_literal: true

class UsageLog < ApplicationRecord
  belongs_to :user
  belongs_to :api_key

  scope :in_range, ->(start_date, end_date) { where(used_at: start_date..end_date) }
  scope :successful, -> { where(status_code: 200..299) }
  scope :with_errors, -> { where("status_code >= ?", 400) }
end
