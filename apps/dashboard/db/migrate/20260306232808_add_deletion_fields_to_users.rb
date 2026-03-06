class AddDeletionFieldsToUsers < ActiveRecord::Migration[8.1]
  def change
    add_column :users, :deletion_token, :string
    add_column :users, :deletion_token_sent_at, :datetime
    add_column :users, :deletion_reason, :text
    add_index :users, :deletion_token, unique: true
  end
end
