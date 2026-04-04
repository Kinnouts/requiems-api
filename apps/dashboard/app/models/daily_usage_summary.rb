# frozen_string_literal: true

class DailyUsageSummary < ApplicationRecord
  belongs_to :user

  validates :date, presence: true
  validates :date, uniqueness: { scope: :user_id }
end
