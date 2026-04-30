# frozen_string_literal: true

class AddLocaleConstraintToUsers < ActiveRecord::Migration[8.1]
  def up
    # Hard-codes the currently supported locales ('en', 'es') to match User::SUPPORTED_LOCALES.
    # If a new locale is added, update User::SUPPORTED_LOCALES AND add a new migration to
    # drop and recreate this constraint (or use ALTER TABLE ... DROP CONSTRAINT / ADD CONSTRAINT).
    execute <<~SQL.squish
      ALTER TABLE users
        ADD CONSTRAINT locale_valid_values
        CHECK (locale IS NULL OR locale IN ('en', 'es'));
    SQL
  end

  def down
    execute "ALTER TABLE users DROP CONSTRAINT locale_valid_values;"
  end
end
