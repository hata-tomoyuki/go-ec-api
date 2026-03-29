-- +goose Up
CREATE TABLE IF NOT EXISTS product_categories (
    product_id BIGSERIAL NOT NULL,
    category_id BIGSERIAL NOT NULL,
    PRIMARY KEY (product_id, category_id),
    CONSTRAINT fk_product
        FOREIGN KEY(product_id)
            REFERENCES products(id)
            ON DELETE CASCADE,
    CONSTRAINT fk_category
        FOREIGN KEY(category_id)
            REFERENCES categories(id)
            ON DELETE CASCADE
);

-- +goose Down
DROP TABLE IF EXISTS product_categories;
