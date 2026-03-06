-- name: GetTasksByScheduleID :many
SELECT * FROM tasks
WHERE schedule_id = $1
ORDER BY id ASC;

-- name: GetTaskByID :one
SELECT * FROM tasks
WHERE id = $1;

-- name: UpdateTaskStatus :one
UPDATE tasks
SET status = $2,
    reason = $3,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: CreateTask :one
INSERT INTO tasks (
    schedule_id, title, status
) VALUES (
    $1, $2, 'pending'
)
RETURNING *;
