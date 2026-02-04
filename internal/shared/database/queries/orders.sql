-- Sqlc Query
-- name: CreateOrder :exec
INSERT INTO
    orders (
        id,
        customer_id,
        total_quantity,
        total_price,
        created_at
    )
VALUES
    (?, ?, ?, ?, ?);

-- name: CreateOrderItem :exec
INSERT INTO
    order_items (id, order_id, product_id, quantity, unit_price)
VALUES
    (?, ?, ?, ?, ?);

-- name: GetOrders :many
SELECT
    o.id,
    o.total_quantity,
    o.total_price,
    o.created_at,
    o.customer_id,
    c.name AS customer_name,
    c.email AS customer_email,
    CAST(
        JSON_ARRAYAGG(
            JSON_OBJECT(
                'id',
                oi.id,
                'product_id',
                p.id,
                'product_name',
                p.name,
                'quantity',
                oi.quantity,
                'unit_price',
                oi.unit_price
            )
        ) AS JSON
    ) AS items
FROM
    orders o
    JOIN customers c ON o.customer_id = c.id
    JOIN order_items oi ON o.id = oi.order_id
    JOIN products p ON oi.product_id = p.id
GROUP BY
    o.id,
    c.id
ORDER BY
    o.created_at DESC
LIMIT
    ?
OFFSET
    ?;

-- name: GetOrderByID :one
SELECT
    o.id,
    o.total_quantity,
    o.total_price,
    o.created_at,
    o.customer_id,
    c.name AS customer_name,
    c.email AS customer_email,
    CAST(
        JSON_ARRAYAGG(
            JSON_OBJECT(
                'id',
                oi.id,
                'product_id',
                p.id,
                'product_name',
                p.name,
                'quantity',
                oi.quantity,
                'unit_price',
                oi.unit_price
            )
        ) AS JSON
    ) AS items
FROM
    orders o
    JOIN customers c ON o.customer_id = c.id
    LEFT JOIN order_items oi ON o.id = oi.order_id
    LEFT JOIN products p ON oi.product_id = p.id
WHERE
    o.id = ?
GROUP BY
    o.id,
    c.id
LIMIT
    1;

-- name: GetOrderItemsByOrderID :many
SELECT
    id,
    order_id,
    product_id,
    quantity,
    unit_price
FROM
    order_items
WHERE
    order_id = ?;

-- name: DeleteOrder :exec
DELETE FROM orders
WHERE
    id = ?;