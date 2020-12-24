-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS users (
  uuid     UUID        PRIMARY KEY DEFAULT uuid_generate_v4(),
  username VARCHAR(20) NOT NULL,
  password CHAR(60),
  CONSTRAINT username_unique UNIQUE (username)
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS users