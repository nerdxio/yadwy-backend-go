CREATE TABLE IF NOT EXISTS banners
(
    id         serial PRIMARY KEY,
    name       VARCHAR(50) UNIQUE NOT NULL,
    index      INT,
    image_url  VARCHAR(300),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);