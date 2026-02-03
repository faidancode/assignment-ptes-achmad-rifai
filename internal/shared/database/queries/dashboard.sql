-- name: GetProductDashboardReport :one
SELECT 
    COUNT(*) AS total_products,
    CAST(IFNULL(SUM(stock_quantity), 0) AS SIGNED) AS total_stock,
    CAST(IFNULL(AVG(price), 0) AS DECIMAL(10,2)) AS avg_price
FROM 
    products;

-- name: GetRecentProducts :many
SELECT 
    id, name, price, stock_quantity, created_at
FROM 
    products
ORDER BY 
    created_at DESC
LIMIT ?;