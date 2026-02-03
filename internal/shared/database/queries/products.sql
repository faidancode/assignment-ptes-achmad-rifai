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
    c.name AS category_name
FROM
    products p
    JOIN categories c ON c.id = p.category_id
WHERE
    (
        sqlc.narg ('search_name') IS NULL
        OR p.name LIKE CONCAT ('%', sqlc.narg ('search_name'), '%')
    )
    AND (
        sqlc.narg ('category_id') IS NULL
        OR p.category_id = sqlc.narg ('category_id')
    )
    AND p.price >= IFNULL (
        CAST(sqlc.narg ('min_price') AS DECIMAL(10, 2)),
        0
    )
    AND p.price <= IFNULL (
        CAST(sqlc.narg ('max_price') AS DECIMAL(10, 2)),
        999999999.99
    )
    AND (
        sqlc.narg ('min_stock') IS NULL
        OR p.stock_quantity >= sqlc.narg ('min_stock')
    )
    AND (
        sqlc.narg ('max_stock') IS NULL
        OR p.stock_quantity <= sqlc.narg ('max_stock')
    )
ORDER BY
    CASE
        WHEN sqlc.arg ('order_by') = 'name_asc' THEN p.name
    END ASC,
    CASE
        WHEN sqlc.arg ('order_by') = 'name_desc' THEN p.name
    END DESC,
    CASE
        WHEN sqlc.arg ('order_by') = 'price_asc' THEN p.price
    END ASC,
    CASE
        WHEN sqlc.arg ('order_by') = 'price_desc' THEN p.price
    END DESC,
    CASE
        WHEN sqlc.arg ('order_by') = 'stock_asc' THEN p.stock_quantity
    END ASC,
    CASE
        WHEN sqlc.arg ('order_by') = 'stock_desc' THEN p.stock_quantity
    END DESC,
    p.created_at DESC
LIMIT
    ?
OFFSET
    ?;

-- name: CountProducts :one
SELECT
    COUNT(*) AS total
FROM
    products p
WHERE
    (
        sqlc.narg ('search_name') IS NULL
        OR p.name LIKE CONCAT ('%', sqlc.narg ('search_name'), '%')
    )
    AND (
        sqlc.narg ('category_id') IS NULL
        OR p.category_id = sqlc.narg ('category_id')
    )
    AND (
        sqlc.narg ('min_price') IS NULL
        OR p.price >= sqlc.narg ('min_price')
    )
    AND (
        sqlc.narg ('max_price') IS NULL
        OR p.price <= sqlc.narg ('max_price')
    )
    AND (
        sqlc.narg ('min_stock') IS NULL
        OR p.stock_quantity >= sqlc.narg ('min_stock')
    )
    AND (
        sqlc.narg ('max_stock') IS NULL
        OR p.stock_quantity <= sqlc.narg ('max_stock')
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