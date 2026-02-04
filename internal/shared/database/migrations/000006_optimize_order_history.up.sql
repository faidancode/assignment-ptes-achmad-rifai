-- Ini mempercepat join sekaligus pengambilan data qty dan price
CREATE INDEX idx_order_items_order_product ON order_items(order_id, product_id);

-- Index pada created_at untuk mempercepat sorting riwayat terbaru
CREATE INDEX idx_orders_created_at ON orders(created_at);