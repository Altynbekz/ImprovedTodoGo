-- name: CreateUser :one
insert into users(id, useraname, first_name, last_name, email, password, created_at, updated_at)
values(
    gen_random_uuid(),
    $1,
    $2,
    $3,
    $4,
    $5,
    now(),
    now()
)
returning *;

-- name: DeleteUserById :exec
Delete from users where id = $1;

-- name: GetUserById :one
select * from users where id = $1;

-- name: UpdateUser :exec
update users
set useraname = $1,
    first_name = $2,
    last_name = $3,
    email = $4,
    password = $5,
    updated_at = now()
where id = $6;

-- name: GetUsernameById :one
select useraname from users where id= $1;

-- name: FindUsesrById :one
select exists(
    select 1 from users where id = $1
);
-- name: FindUserByUsername :one
select exists(
    select 1 from users where useraname = $1
);
-- name: FindUserByEmail :one
select exists(
    select 1 from users where email = $1
);

-- name: GetUserByEmail :one
select * from users where email = $1;

