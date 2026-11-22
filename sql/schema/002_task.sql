-- +goose Up
create table if not exists tasks(
    id uuid primary key,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now(),
    task_name text not null default 'untitled',
    task_desc text not null,
    deadline timestamp,
    user_id uuid not null references users(id) on delete cascade
);

-- +goose Down
drop table tasks;