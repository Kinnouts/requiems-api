class CreateCreditAdjustments < ActiveRecord::Migration[8.1]
  def change
    create_table :credit_adjustments do |t|
      t.references :user, null: false, foreign_key: true
      t.integer :amount
      t.text :reason
      t.string :adjustment_type

      t.timestamps
    end
  end
end
