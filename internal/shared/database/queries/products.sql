-- name: CreateProduct :exec
INSERT INTO
    products (
        id,
        name,
        description,
        price,
        category_id,
        stock_quantity,
        is_active
    )
VALUES
    (?, ?, ?, ?, ?, ?, ?);

-- name: GetProductByID :one
SELECT
    p.id,
    p.name,
    p.description,
    p.price,
    p.stock_quantity,
    p.is_active,
    p.created_at,
    p.updated_at,
    c.id AS category_id,
    c.name AS category_name,
    c.description AS category_description
FROM
    products p
    JOIN categories c ON c.id = p.category_id
WHERE
    p.id = ?
LIMIT
    1;

-- name: ListProducts :many
SELECT
    p.id,
    p.name,
    p.description,
    p.price,
    p.stock_quantity,
    p.is_active,
    p.created_at,
    c.id AS category_id,
    c.name AS category_name,
    CAST(IFNULL (SUM(oi.quantity), 0) AS UNSIGNED) AS total_sold
FROM
    products p
    JOIN categories c ON c.id = p.category_id
    LEFT JOIN order_items oi ON oi.product_id = p.id
WHERE
    (
        sqlc.arg ('search_name') = ''
        OR p.name LIKE CONCAT ('%', sqlc.arg ('search_name'), '%')
    )
    AND (
        sqlc.arg ('category_id') = ''
        OR p.category_id = sqlc.arg ('category_id')
    )
    AND (
        sqlc.arg ('min_price') = 0
        OR p.price >= sqlc.arg ('min_price')
    )
    AND (
        sqlc.arg ('max_price') = 0
        OR p.price <= sqlc.arg ('max_price')
    )
    AND (
        sqlc.arg ('min_stock') = 0
        OR p.stock_quantity >= sqlc.arg ('min_stock')
    )
    AND (
        sqlc.arg ('max_stock') = 0
        OR p.stock_quantity <= sqlc.arg ('max_stock')
    )
GROUP BY
    p.id,
    c.id
ORDER BY
    CASE
        WHEN sqlc.arg ('order_by') = 'sold_desc' THEN IFNULL (SUM(oi.quantity), 0)
    END DESC,
    p.created_at DESC
LIMIT
    ?
OFFSET
    ?;

-- name: CountProducts :one
SELECT
    COUNT(DISTINCT p.id) AS total
FROM
    products p
WHERE
    (
        sqlc.arg ('search_name') = ''
        OR p.name LIKE CONCAT ('%', sqlc.arg ('search_name'), '%')
    )
    AND (
        sqlc.arg ('category_id') = ''
        OR p.category_id = sqlc.arg ('category_id')
    )
    AND (
        sqlc.arg ('min_price') = 0
        OR p.price >= sqlc.arg ('min_price')
    )
    AND (
        sqlc.arg ('max_price') = 0
        OR p.price <= sqlc.arg ('max_price')
    )
    AND (
        sqlc.arg ('min_stock') = 0
        OR p.stock_quantity >= sqlc.arg ('min_stock')
    )
    AND (
        sqlc.arg ('max_stock') = 0
        OR p.stock_quantity <= sqlc.arg ('max_stock')
    );

-- name: UpdateProduct :exec
UPDATE products
SET
    name = ?,
    description = ?,
    price = ?,
    category_id = ?,
    stock_quantity = ?,
    is_active = ?
WHERE
    id = ?;

-- name: DeleteProduct :exec
DELETE FROM products
WHERE
    id = ?;