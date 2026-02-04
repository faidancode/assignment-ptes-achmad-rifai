-- Menghapus index pada kolom created_at di tabel orders
DROP INDEX idx_orders_created_at ON orders;

-- Menghapus composite index pada tabel order_items
DROP INDEX idx_order_items_order_product ON order_items;
