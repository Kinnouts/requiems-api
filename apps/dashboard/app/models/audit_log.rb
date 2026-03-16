# frozen_string_literal: true

class AuditLog < ApplicationRecord
  belongs_to :user
  belongs_to :admin_user, class_name: "User", optional: true
end
