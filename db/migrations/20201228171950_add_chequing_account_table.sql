-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS accounts (
  uuid        UUID       PRIMARY KEY DEFAULT uuid_generate_v4(),
  last_four   VARCHAR(4) NOT NULL,
  description TEXT       NOT NULL
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

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS account_activities, accounts;
