CREATE TABLE IF NOT EXISTS product_labels
(
    product_id BIGINT      NOT NULL,
    label_name VARCHAR(50) NOT NULL,
    PRIMARY KEY (product_id, label_name),
    FOREIGN KEY (product_id) REFERENCES products (id) ON DELETE CASCADE
);