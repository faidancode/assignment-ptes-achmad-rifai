-- Sqlc Query
-- name: CreateOrder :exec
INSERT INTO orders (id, customer_id, total_quantity, total_price, created_at)
VALUES (?, ?, ?, ?, ?);

-- name: CreateOrderItem :exec
INSERT INTO order_items (id, order_id, product_id, quantity, unit_price)
VALUES (?, ?, ?, ?, ?);

-- name: GetOrders :many
SELECT id, customer_id, total_quantity, total_price, created_at 
FROM orders 
ORDER BY created_at DESC;

-- name: GetOrderByID :one
SELECT id, customer_id, total_quantity, total_price, created_at 
FROM orders 
WHERE id = ? LIMIT 1;

-- name: GetOrderItemsByOrderID :many
SELECT id, order_id, product_id, quantity, unit_price
FROM order_items
WHERE order_id = ?;

-- name: DeleteOrder :exec
DELETE FROM orders WHERE id = ?;