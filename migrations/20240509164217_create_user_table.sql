-- +goose Up
create type user_role as enum('user', 'admin');

create table "user" (
    id serial primary key,
    name text not null,
    email text not null,
    role user_role not null default 'user',
    password_hash text not null,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null,
    unique (email)
);

-- +goose Down
drop table "user";

drop type user_role;
