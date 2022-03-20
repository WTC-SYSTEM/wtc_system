create table public.user
(
    id       uuid primary key default gen_random_uuid(),
    username varchar(100) not null,
    email    varchar(50)  unique not null,
    password varchar(100) not null
);

alter table public.user
add constraint UQ_user_username UNIQUE(username)