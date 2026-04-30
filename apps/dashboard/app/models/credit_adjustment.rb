# frozen_string_literal: true

class CreditAdjustment < ApplicationRecord
  belongs_to :user

  validates :amount, :adjustment_type, presence: true
  validates :amount, numericality: true
end
