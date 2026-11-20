-- name: CreateTask :one
insert into tasks(id, created_at, updated_at, task_name,task_desc,deadline)
values(
    gen_random_uuid(),
    now(),
    now(),
    $1,
    $2,
    $3
    )
returning *;