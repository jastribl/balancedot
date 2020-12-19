-- +goose Up
-- SQL in this section is executed when the migration is applied.
-- Transaction Date,Post Date,Description,Category,Type,Amount,Memo
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS cards (
  uuid        UUID        PRIMARY KEY DEFAULT uuid_generate_v4(),
  last_four   VARCHAR(4)  NOT NULL,
  description TEXT        NOT NULL
);

CREATE TABLE IF NOT EXISTS card_activities (
  uuid             UUID             PRIMARY KEY DEFAULT uuid_generate_v4(),
  card_uuid        UUID             NOT NULL,
  transaction_date TIMESTAMP        NOT NULL,
  post_date        TIMESTAMP        NOT NULL,
  description      TEXT             NOT NULL,
  category         TEXT             NOT NULL,
  type             TEXT             NOT NULL,
  amount           NUMERIC(10, 2)   NOT NULL,
  CONSTRAINT
    fk_card
      FOREIGN KEY(card_uuid)
      REFERENCES cards(uuid)
);

CREATE TABLE IF NOT EXISTS splitwise_expenses (
  uuid          UUID             PRIMARY KEY DEFAULT uuid_generate_v4(),
  splitwise_id  int              NOT NULL UNIQUE,
  description   TEXT             NOT NULL,
  details       TEXT             NOT NULL,
  currency_code TEXT             NOT NULL,
  amount        NUMERIC(10, 2)   NOT NULL,
  amount_paid   NUMERIC(10, 2)   NOT NULL,
  date          TIMESTAMP        NOT NULL,
  created_at    TIMESTAMP        NOT NULL,
  updated_at    TIMESTAMP,
  deleted_at    TIMESTAMP,
  category      TEXT   
);

CREATE TABLE IF NOT EXISTS expense_links (
  card_activity_uuid     UUID NOT NULL REFERENCES card_activities(uuid),
  splitwise_expense_uuid UUID NOT NULL REFERENCES splitwise_expenses(uuid),
  CONSTRAINT fk_card_activitiy_splitwise_expense
    PRIMARY KEY(card_activity_uuid, splitwise_expense_uuid)
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS cards, card_activities, splitwise_expenses, expense_links;
