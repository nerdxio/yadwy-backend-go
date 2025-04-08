CREATE TABLE IF NOT EXISTS categories
(
    id          serial PRIMARY KEY,
    name        VARCHAR(50) UNIQUE NOT NULL,
    description TEXT,
    image_url    VARCHAR(300),
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);