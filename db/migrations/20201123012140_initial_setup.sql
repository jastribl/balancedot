-- +goose Up
-- SQL in this section is executed when the migration is applied.
-- Transaction Date,Post Date,Description,Category,Type,Amount,Memo
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE bank_names AS ENUM ('chase', 'bofa');

CREATE TABLE IF NOT EXISTS accounts (
  uuid        UUID       PRIMARY KEY DEFAULT uuid_generate_v4(),
  last_four   VARCHAR(4) NOT NULL,
  description TEXT       NOT NULL,
  bank_name   bank_names NOT NULL,

  CONSTRAINT accounts_last_four_unique UNIQUE (last_four)
);

CREATE TABLE IF NOT EXISTS account_activities (
  uuid             UUID             PRIMARY KEY DEFAULT uuid_generate_v4(),
  account_uuid     UUID             NOT NULL,
  details          TEXT             NOT NULL,
  posting_date     TIMESTAMP        NOT NULL,
  description      TEXT             NOT NULL,
  amount           NUMERIC(10, 2)   NOT NULL,
  type             TEXT             NOT NULL,
  CONSTRAINT
    fk_account
      FOREIGN KEY(account_uuid)
      REFERENCES accounts(uuid)
);

CREATE TABLE IF NOT EXISTS cards (
  uuid        UUID       PRIMARY KEY DEFAULT uuid_generate_v4(),
  last_four   VARCHAR(4) NOT NULL,
  description TEXT       NOT NULL,
  bank_name   bank_names NOT NULL,

  CONSTRAINT cards_last_four_unique UNIQUE (last_four)
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
  uuid                 UUID           PRIMARY KEY DEFAULT uuid_generate_v4(),
  splitwise_id         int            NOT NULL,
  description          TEXT           NOT NULL,
  details              TEXT           NOT NULL,
  currency_code        TEXT           NOT NULL,
  amount               NUMERIC(10, 2) NOT NULL,
  amount_paid          NUMERIC(10, 2) NOT NULL,
  date                 TIMESTAMP      NOT NULL,
  splitwise_created_at TIMESTAMP      NOT NULL,
  splitwise_updated_at TIMESTAMP,
  splitwise_deleted_at TIMESTAMP,
  category             TEXT,
  creation_method      TEXT,

  CONSTRAINT splitwise_expenses_splitwise_id_unique UNIQUE (splitwise_id)
);

CREATE TABLE IF NOT EXISTS expense_links (
  card_activity_uuid     UUID NOT NULL REFERENCES card_activities(uuid),
  splitwise_expense_uuid UUID NOT NULL REFERENCES splitwise_expenses(uuid),

  CONSTRAINT fk_card_activity_splitwise_expense
    PRIMARY KEY(card_activity_uuid, splitwise_expense_uuid)
);

CREATE TABLE IF NOT EXISTS account_activity_links (
  account_activity_uuid  UUID NOT NULL REFERENCES account_activities(uuid),
  splitwise_expense_uuid UUID NOT NULL REFERENCES splitwise_expenses(uuid),

  CONSTRAINT fk_account_activity_splitwise_expense
    PRIMARY KEY(account_activity_uuid, splitwise_expense_uuid)
);

CREATE TABLE IF NOT EXISTS splitwise_links (
  splitwise_expense_uuid UUID NOT NULL REFERENCES splitwise_expenses(uuid),
  linked_splitwise_expense_uuid UUID NOT NULL REFERENCES splitwise_expenses(uuid),

  CONSTRAINT fk_splitwise_expense_splitwise_expense
    PRIMARY KEY(splitwise_expense_uuid, linked_splitwise_expense_uuid)
);

CREATE TABLE IF NOT EXISTS account_card_links (
  card_activity_uuid     UUID NOT NULL REFERENCES card_activities(uuid),
  account_activity_uuid  UUID NOT NULL REFERENCES account_activities(uuid),

  CONSTRAINT fk_card_activity_account_activity
    PRIMARY KEY(card_activity_uuid, account_activity_uuid)
);

-- Example data
-- TODO: remove eventually
INSERT INTO accounts (last_four, description) VALUES (3682, 'Chase Chequing Account');
INSERT INTO cards (last_four, description) VALUES (2427, 'Chase Freedom Unlimited');
INSERT INTO cards (last_four, description) VALUES (9307, 'Chase Sapphire Reserve');

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS splitwise_links;
DROP TABLE IF EXISTS account_activity_links;
DROP TABLE IF EXISTS expense_links;
DROP TABLE IF EXISTS splitwise_expenses;
DROP TABLE IF EXISTS card_activities;
DROP TABLE IF EXISTS cards;
DROP TABLE IF EXISTS account_activities;
DROP TABLE IF EXISTS accounts;
