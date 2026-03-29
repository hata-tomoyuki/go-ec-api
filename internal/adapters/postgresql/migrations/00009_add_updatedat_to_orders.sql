-- +goose Up
ALTER TABLE orders
  ADD COLUMN updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

-- +goose Down
ALTER TABLE orders
  DROP COLUMN updated_at;
