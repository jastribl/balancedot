-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE cards
ADD CONSTRAINT cards_last_four_unique UNIQUE (last_four);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE cards
DROP CONSTRAINT IF EXISTS cards_last_four_unique;
