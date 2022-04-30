CREATE TABLE IF NOT EXISTS "users"
(
    id       BIGSERIAL PRIMARY KEY,
    username VARCHAR(100)       NOT NULL,
    email    VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(100)       NOT NULL
);