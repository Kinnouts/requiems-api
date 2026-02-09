class CreateApiKeys < ActiveRecord::Migration[8.1]
  def change
    create_table :api_keys do |t|
      t.references :user, null: false, foreign_key: true
      t.string :key_prefix
      t.string :key_hash
      t.string :name
      t.string :environment
      t.datetime :last_used_at
      t.string :last_used_ip
      t.boolean :active
      t.datetime :revoked_at
      t.string :revoked_reason

      t.timestamps
    end
  end
end
