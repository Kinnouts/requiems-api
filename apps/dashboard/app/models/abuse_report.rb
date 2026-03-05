# frozen_string_literal: true
class AbuseReport < ApplicationRecord
  belongs_to :user
  belongs_to :api_key
  belongs_to :resolved_by, class_name: "User", optional: true

  # Validations
  validates :report_type, presence: true
  validates :status, presence: true, inclusion: { in: %w[pending investigating resolved] }
  validates :description, length: { maximum: 5000 }

  # Scopes
  scope :pending, -> { where(status: "pending") }
  scope :investigating, -> { where(status: "investigating") }
  scope :resolved, -> { where(status: "resolved") }
  scope :recent, -> { order(created_at: :desc) }

  # Status helpers
  def pending?
    status == "pending"
  end

  def investigating?
    status == "investigating"
  end

  def resolved?
    status == "resolved"
  end
end
