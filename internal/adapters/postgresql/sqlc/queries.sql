-- name: ListProducts :many
SELECT
 *
FROM
    products;

-- name: FindProductById :one
SELECT
 *
FROM
    products
WHERE
    id = $1;

-- name: CreateProduct :one
INSERT INTO products (name, price_in_cents) VALUES ($1, $2) RETURNING *;

-- name: UpdateProduct :one
UPDATE products
SET name = $2, price_in_cents = $3
WHERE id = $1
RETURNING *;

-- name: CreateOrder :one
INSERT INTO orders (customer_id) VALUES ($1) RETURNING *;

-- name: CreateOrderItem :one
INSERT INTO order_items (order_id, product_id, quantity, price_in_cents)
VALUES ($1, $2, $3, $4) RETURNING *;
