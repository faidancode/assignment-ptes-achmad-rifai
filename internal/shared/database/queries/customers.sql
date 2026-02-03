-- name: CreateCustomer :exec
INSERT INTO customers (id, name, email, created_at)
VALUES (?, ?, ?, ?);

-- name: GetCustomers :many
SELECT id, name, email, created_at 
FROM customers 
ORDER BY created_at DESC;

-- name: GetCustomerByID :one
SELECT id, name, email, created_at 
FROM customers 
WHERE id = ? LIMIT 1;

-- name: UpdateCustomer :exec
UPDATE customers
SET name = ?, email = ?
WHERE id = ?;

-- name: DeleteCustomer :exec
DELETE FROM customers 
WHERE id = ?;