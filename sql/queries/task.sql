-- name: CreateTask :one
insert into tasks(id, created_at, updated_at, task_name,task_desc,deadline, user_id)
values(
    gen_random_uuid(),
    now(),
    now(),
    $1,
    $2,
    $3,
    $4
    )
returning *;

-- name: GetTaskByUserId :many
select * from tasks where user_id = $1;