CREATE TABLE IF NOT EXISTS "users"
(
    id       UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(100)       NOT NULL,
    email    VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(100)       NOT NULL
);