-- name: ListSchedules :many
SELECT * FROM schedules
ORDER BY start_at ASC;

-- name: GetScheduleByID :one
SELECT * FROM schedules
WHERE id = $1;

-- name: GetTodaySchedules :many
SELECT * FROM schedules
WHERE start_at >= @today_start::timestamptz AND start_at < @today_end::timestamptz
ORDER BY start_at ASC;

-- name: GetScheduleStats :one
SELECT
    COUNT(*)::int AS total,
    COUNT(*) FILTER (WHERE clock_in_at IS NULL AND end_at < NOW())::int AS missed,
    COUNT(*) FILTER (WHERE clock_in_at IS NULL AND end_at > NOW() AND start_at >= @today_start::timestamptz AND start_at < @today_end::timestamptz)::int AS upcoming,
    COUNT(*) FILTER (WHERE clock_out_at IS NOT NULL AND start_at >= @today_start::timestamptz AND start_at < @today_end::timestamptz)::int AS completed
FROM schedules;

-- name: ClockIn :one
UPDATE schedules
SET clock_in_at = NOW(),
    clock_in_lat = @latitude,
    clock_in_lng = @longitude,
    clock_in_verified = @verified
WHERE id = $1
RETURNING *;

-- name: ClockOut :one
UPDATE schedules
SET clock_out_at = NOW(),
    clock_out_lat = @latitude,
    clock_out_lng = @longitude,
    clock_out_verified = @verified
WHERE id = $1
RETURNING *;

-- name: CreateSchedule :one
INSERT INTO schedules (
    client_name, start_at, end_at,
    location, latitude, longitude
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING *;
