-- +goose Up
create table if not exists users(
    id uuid primary key,
    useraname text not null,
    FirstName text not null,
    LastName text not null,
    email text not null,
    password text not null,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now()
);

-- +goose Down
drop table users;