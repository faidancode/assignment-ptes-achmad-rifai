CREATE TABLE
    products (
        id CHAR(36) PRIMARY KEY,
        name VARCHAR(150) NOT NULL,
        description TEXT,
        price DECIMAL(15, 2) NOT NULL,
        category_id CHAR(36) NOT NULL,
        stock_quantity INT NOT NULL DEFAULT 0,
        is_active BOOLEAN NOT NULL DEFAULT TRUE,
        created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        CONSTRAINT fk_products_category FOREIGN KEY (category_id) REFERENCES categories (id) ON UPDATE CASCADE ON DELETE RESTRICT
    ) ENGINE = InnoDB;

-- Indexing for performance
CREATE INDEX idx_products_name ON products (name);

CREATE INDEX idx_products_category_id ON products (category_id);

CREATE INDEX idx_products_price ON products (price);

CREATE INDEX idx_products_stock_quantity ON products (stock_quantity);

CREATE INDEX idx_products_is_active ON products (is_active);

CREATE INDEX idx_products_created_at ON products (created_at);