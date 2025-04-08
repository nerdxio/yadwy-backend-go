CREATE TABLE IF NOT EXISTS users
(
    id         serial PRIMARY KEY,
    name       VARCHAR(50)         NOT NULL,
    password   VARCHAR(300)        NOT NULL,
    email      VARCHAR(300) UNIQUE NOT NULL,
    phone      VARCHAR(20),
    role       VARCHAR(20)         NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
