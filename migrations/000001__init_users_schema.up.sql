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

CREATE TABLE IF NOT EXISTS customers
(
    id         serial PRIMARY KEY,
    user_id    int UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE TABLE IF NOT EXISTS sellers
(
    id         serial PRIMARY KEY,
    user_id    int UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE TABLE IF NOT EXISTS admins
(
    id         serial PRIMARY KEY,
    user_id    int UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id)
);
