create table "users"
(
    id       uuid primary key default gen_random_uuid(),
    username varchar(100) not null,
    email    varchar(50)  unique not null,
    password varchar(100) not null
);