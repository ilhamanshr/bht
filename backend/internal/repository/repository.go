package repository

import (
	"bht-test/internal/repository/sqlc"
	"context"
	"database/sql"
)

//go:generate mockgen -source=repository.go -destination=mock/repository.go -package=mock
type IRepository interface {
	// Schedule
	ClockIn(ctx context.Context, arg sqlc.ClockInParams) (sqlc.Schedule, error)
	ClockOut(ctx context.Context, arg sqlc.ClockOutParams) (sqlc.Schedule, error)
	CreateSchedule(ctx context.Context, arg sqlc.CreateScheduleParams) (sqlc.Schedule, error)
	GetScheduleByID(ctx context.Context, id int32) (sqlc.Schedule, error)
	GetScheduleStats(ctx context.Context, arg sqlc.GetScheduleStatsParams) (sqlc.GetScheduleStatsRow, error)
	GetTodaySchedules(ctx context.Context, arg sqlc.GetTodaySchedulesParams) ([]sqlc.Schedule, error)
	ListSchedules(ctx context.Context) ([]sqlc.Schedule, error)

	// Task
	CreateTask(ctx context.Context, arg sqlc.CreateTaskParams) (sqlc.Task, error)
	GetTaskByID(ctx context.Context, id int32) (sqlc.Task, error)
	GetTasksByScheduleID(ctx context.Context, scheduleID int32) ([]sqlc.Task, error)
	UpdateTaskStatus(ctx context.Context, arg sqlc.UpdateTaskStatusParams) (sqlc.Task, error)

	// Tx
	WithTx(tx *sql.Tx) *sqlc.Queries
}
