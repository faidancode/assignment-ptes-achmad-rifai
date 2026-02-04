-- name: CreateCustomer :exec
INSERT INTO
    customers (id, name, email, created_at)
VALUES
    (?, ?, ?, ?);

-- name: GetCustomers :many
SELECT
    id,
    name,
    email,
    created_at
FROM
    customers
ORDER BY
    created_at DESC
LIMIT
    ?
OFFSET
    ?;

-- name: GetCustomerByID :one
SELECT
    id,
    name,
    email,
    created_at
FROM
    customers
WHERE
    id = ?
LIMIT
    1;

-- name: UpdateCustomer :exec
UPDATE customers
SET
    name = ?,
    email = ?
WHERE
    id = ?;

-- name: DeleteCustomer :exec
DELETE FROM customers
WHERE
    id = ?;

-- name: GetTopCustomers :many
SELECT
    c.id,
    c.name,
    c.email,
    CAST(SUM(o.total_price) AS DECIMAL(10, 2)) as total_spent,
    COUNT(o.id) as total_orders
FROM
    customers c
    JOIN orders o ON c.id = o.customer_id
GROUP BY
    c.id
ORDER BY
    total_spent DESC
LIMIT
    ?;