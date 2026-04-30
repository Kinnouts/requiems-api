# frozen_string_literal: true

redis_config = { url: ENV.fetch("REDIS_URL", "redis://localhost:6379") }

Sidekiq.configure_server do |config|
  config.redis = redis_config

  config.on(:startup) do
    schedule_file = Rails.root.join("config/sidekiq_schedule.yml")
    Sidekiq::Cron::Job.load_from_hash(YAML.load_file(schedule_file)) if schedule_file.exist?
  end
end

Sidekiq.configure_client do |config|
  config.redis = redis_config
end
