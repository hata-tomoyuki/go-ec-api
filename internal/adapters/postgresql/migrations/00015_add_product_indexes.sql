-- +goose Up

-- pg_trgm GIN index for ILIKE '%keyword%' search on product name
CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE INDEX idx_products_name_trgm ON products USING gin (name gin_trgm_ops);

-- category_id index for category filtering (PK is (product_id, category_id) so category_id alone needs its own index)
CREATE INDEX idx_product_categories_category_id ON product_categories (category_id);

-- created_at DESC index for default sort order (newest first)
CREATE INDEX idx_products_created_at ON products (created_at DESC);

-- price_in_cents index for price sort
CREATE INDEX idx_products_price_in_cents ON products (price_in_cents);

-- +goose Down
DROP INDEX IF EXISTS idx_products_price_in_cents;
DROP INDEX IF EXISTS idx_products_created_at;
DROP INDEX IF EXISTS idx_product_categories_category_id;
DROP INDEX IF EXISTS idx_products_name_trgm;
DROP EXTENSION IF EXISTS pg_trgm;
