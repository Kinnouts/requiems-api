class AbuseReport < ApplicationRecord
  belongs_to :user
  belongs_to :api_key
end
