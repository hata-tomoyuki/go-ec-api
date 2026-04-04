-- +goose Up
ALTER TABLE products ADD COLUMN IF NOT EXISTS description TEXT NOT NULL DEFAULT '';
ALTER TABLE products ADD COLUMN IF NOT EXISTS image_color TEXT NOT NULL DEFAULT 'from-gray-400 to-gray-600';
ALTER TABLE categories ADD COLUMN IF NOT EXISTS image_color TEXT NOT NULL DEFAULT 'from-gray-400 to-gray-600';

-- +goose Down
ALTER TABLE products DROP COLUMN IF EXISTS description;
ALTER TABLE products DROP COLUMN IF EXISTS image_color;
ALTER TABLE categories DROP COLUMN IF EXISTS image_color;
