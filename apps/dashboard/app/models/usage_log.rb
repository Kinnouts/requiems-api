# frozen_string_literal: true

class UsageLog < ApplicationRecord
  belongs_to :user
  belongs_to :api_key

  def self.error_rate_for(scope)
    total = scope.count
    return 0 if total.zero?

    errors = scope.where("status_code >= ?", 400).count
    ((errors.to_f / total) * 100).round(2)
  end
end
