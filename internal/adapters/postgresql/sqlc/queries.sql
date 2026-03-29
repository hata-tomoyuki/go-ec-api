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

-- name: DeleteProduct :one
DELETE FROM products
WHERE id = $1
RETURNING *;

-- name: CreateOrder :one
INSERT INTO orders (customer_id) VALUES ($1) RETURNING *;

-- name: ListCategories :many
SELECT * FROM categories;

-- name: FindCategoryById :one
SELECT * FROM categories WHERE id = $1;

-- name: CreateCategory :one
INSERT INTO categories (name, description) VALUES ($1, $2) RETURNING *;

-- name: UpdateCategory :one
UPDATE categories
SET name = $2, description = $3, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteCategory :one
DELETE FROM categories
WHERE id = $1
RETURNING *;

-- name: ListProductsByCategory :many
SELECT
    p.id,
    p.name,
    p.price_in_cents,
    p.quantity,
    p.created_at
FROM
    products p
JOIN
    product_categories pc ON p.id = pc.product_id
WHERE
    pc.category_id = $1;

-- name: AddProductToCategory :exec
INSERT INTO product_categories (product_id, category_id) VALUES ($1, $2);

-- name: RemoveProductFromCategory :exec
DELETE FROM product_categories
WHERE product_id = $1 AND category_id = $2;

-- name: CreateOrderItem :one
INSERT INTO order_items (order_id, product_id, quantity, price_in_cents)
VALUES ($1, $2, $3, $4) RETURNING *;

-- name: CancelOrder :one
UPDATE orders
SET status = 'cancelled', updated_at = now()
WHERE id = $1
RETURNING *;

-- name: ListOrdersByCustomerID :many
SELECT
    o.id,
    o.customer_id,
    o.status,
    o.created_at,
    o.updated_at,
    oi.product_id,
    oi.quantity,
    oi.price_in_cents
FROM
    orders o
JOIN
    order_items oi ON o.id = oi.order_id
WHERE
    o.customer_id = $1;

-- name: ListAllOrders :many
SELECT
    o.id,
    o.customer_id,
    o.status,
    o.created_at,
    o.updated_at,
    oi.product_id,
    oi.quantity,
    oi.price_in_cents
FROM
    orders o
JOIN
    order_items oi ON o.id = oi.order_id;

-- name: FindOrderById :one
SELECT
    o.id,
    o.customer_id,
    o.status,
    o.created_at,
    o.updated_at,
    oi.product_id,
    oi.quantity,
    oi.price_in_cents
FROM
    orders o
JOIN
    order_items oi ON o.id = oi.order_id
WHERE
    o.id = $1;

-- name: UpdateOrderStatus :one
UPDATE orders
SET status = $2, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: CreateUser :one
INSERT INTO users (name, email, password_hash, role) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: FindUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: RevokeToken :exec
INSERT INTO revoked_tokens (jti, expired_at) VALUES ($1, $2);

-- name: IsTokenRevoked :one
SELECT EXISTS (SELECT 1 FROM revoked_tokens WHERE jti = $1);
