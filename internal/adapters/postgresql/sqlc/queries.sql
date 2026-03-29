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

-- name: FindUserById :one
SELECT * FROM users WHERE id = $1;

-- name: UpdateUser :one
UPDATE users
SET name = $2, email = $3, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: UpdateUserPassword :one
UPDATE users
SET password_hash = $2, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: RevokeToken :exec
INSERT INTO revoked_tokens (jti, expired_at) VALUES ($1, $2);

-- name: IsTokenRevoked :one
SELECT EXISTS (SELECT 1 FROM revoked_tokens WHERE jti = $1);

-- name: InsertRefreshToken :one
INSERT INTO refresh_tokens (user_id, token_hash, expires_at) VALUES ($1, $2, $3) RETURNING *;

-- name: ConsumeRefreshToken :one
DELETE FROM refresh_tokens
WHERE token_hash = $1 AND expires_at > now()
RETURNING *;

-- name: DeleteRefreshToken :exec
DELETE FROM refresh_tokens WHERE id = $1;

-- name: DeleteRefreshTokensByUserId :exec
DELETE FROM refresh_tokens WHERE user_id = $1;

-- name: CreateCart :one
INSERT INTO carts (user_id) VALUES ($1) RETURNING *;

-- name: AddItemToCart :one
INSERT INTO cart_items (cart_id, product_id, quantity) VALUES ($1, $2, $3) RETURNING *;

-- name: ListCartItemsByUserId :many
SELECT
    ci.id,
    ci.cart_id,
    ci.product_id,
    ci.quantity,
    p.name AS product_name,
    p.price_in_cents AS product_price_in_cents
FROM
    carts c
JOIN
    cart_items ci ON ci.cart_id = c.id
JOIN
    products p ON ci.product_id = p.id
WHERE
    c.user_id = $1;

-- name: UpdateCartItemQuantity :one
UPDATE cart_items
SET quantity = $2
WHERE id = $1
RETURNING *;

-- name: RemoveItemFromCart :one
DELETE FROM cart_items
WHERE id = $1
RETURNING *;

-- name: ClearCart :exec
DELETE FROM cart_items
WHERE cart_id = (SELECT id FROM carts WHERE user_id = $1);
