# frozen_string_literal: true

class AddLocaleConstraintToUsers < ActiveRecord::Migration[8.1]
  def up
    execute <<~SQL
      ALTER TABLE users
        ADD CONSTRAINT locale_valid_values
        CHECK (locale IS NULL OR locale IN ('en', 'es'));
    SQL
  end

  def down
    execute "ALTER TABLE users DROP CONSTRAINT locale_valid_values;"
  end
end
