-- +goose Up
-- SQL in this section is executed when the migration is applied.
-- Transaction Date,Post Date,Description,Category,Type,Amount,Memo
CREATE TABLE IF NOT EXISTS cards (
  id SERIAL PRIMARY KEY NOT NULL,
  last_four INT NOT NULL,
  bank_name VARCHAR NOT NULL,
  UNIQUE (last_four)
);
CREATE TABLE IF NOT EXISTS card_activities (
  id SERIAL NOT NULL,
  uuid UUID PRIMARY KEY NOT NULL,
  card_id INT,
  transaction_date TIMESTAMP NOT NULL,
  post_date TIMESTAMP NOT NULL,
  description VARCHAR NOT NULL,
  category VARCHAR NOT NULL,
  type VARCHAR NOT NULL,
  amount NUMERIC(6, 2) NOT NULL,
  UNIQUE (uuid),
  CONSTRAINT
    fk_card
      FOREIGN KEY(card_id)
      REFERENCES cards(id)
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS cards, card_activities;
