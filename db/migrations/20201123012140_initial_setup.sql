-- +goose Up
-- SQL in this section is executed when the migration is applied.
-- Transaction Date,Post Date,Description,Category,Type,Amount,Memo
CREATE TABLE IF NOT EXISTS cards (
  uuid UUID PRIMARY KEY NOT NULL,
  last_four VARCHAR(4)  NOT NULL,
  bank_name VARCHAR     NOT NULL
);

CREATE TABLE IF NOT EXISTS card_activities (
  uuid             UUID PRIMARY KEY NOT NULL,
  card_uuid        UUID             NOT NULL,
  transaction_date TIMESTAMP        NOT NULL,
  post_date        TIMESTAMP        NOT NULL,
  description      VARCHAR          NOT NULL,
  category         VARCHAR          NOT NULL,
  type             VARCHAR          NOT NULL,
  amount           NUMERIC(6, 2)    NOT NULL,
  CONSTRAINT
    fk_card
      FOREIGN KEY(card_uuid)
      REFERENCES cards(uuid)
);

CREATE TABLE IF NOT EXISTS splitwise_expenses (
  uuid         UUID PRIMARY KEY NOT NULL,
  splitwise_id int              NOT NULL,
  description  VARCHAR          NOT NULL,
  details      VARCHAR          NOT NULL,
  amount       NUMERIC(6, 2)    NOT NULL,
  total_amount NUMERIC(6, 2)    NOT NULL,
  date         TIMESTAMP        NOT NULL,
  created_at   TIMESTAMP        NOT NULL,
  updated_at   TIMESTAMP,
  deleted_at   TIMESTAMP,
  category     VARCHAR
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS cards, card_activities, splitwise_expenses;
