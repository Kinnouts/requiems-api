# Use a separate migration tracking table to avoid conflicts with Go backend
# Go backend uses 'schema_migrations' with golang-migrate
# Rails will use 'rails_schema_migrations' instead
Rails.application.config.active_record.schema_migrations_table_name = "rails_schema_migrations"
Rails.application.config.active_record.internal_metadata_table_name = "rails_internal_metadata"
