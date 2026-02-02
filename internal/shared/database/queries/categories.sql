-- name: CreateCategory :exec
INSERT INTO categories (
    id,
    name,
    description
) VALUES (
    ?, ?, ?
);

-- name: GetCategories :many
SELECT
    id,
    name,
    description
FROM categories
ORDER BY name ASC;

-- name: GetCategoryByID :one
SELECT
    id,
    name,
    description
FROM categories
WHERE id = ?
LIMIT 1;

-- name: UpdateCategory :exec
UPDATE categories
SET
    name = ?,
    description = ?
WHERE id = ?;

-- name: DeleteCategory :exec
DELETE FROM categories
WHERE id = ?;
