class AddTrgmIndexToApiKeysKeyPrefix < ActiveRecord::Migration[8.1]
  def up
    enable_extension "pg_trgm" unless extension_enabled?("pg_trgm")
    add_index :api_keys, :key_prefix, using: :gin, opclass: :gin_trgm_ops,
              name: "index_api_keys_on_key_prefix_trgm"
  end

  def down
    remove_index :api_keys, name: "index_api_keys_on_key_prefix_trgm"
  end
end
