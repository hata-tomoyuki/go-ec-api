-- +goose Up
CREATE TYPE status AS ENUM ('pending', 'completed', 'cancelled');
ALTER TABLE orders
  ADD COLUMN status status NOT NULL DEFAULT 'pending';

-- +goose Down
DROP TYPE status;
