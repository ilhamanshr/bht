package usecase

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"

	"bht-test/internal/domain"
	"bht-test/internal/repository"
	"bht-test/internal/repository/sqlc"
)

// ScheduleUsecase handles business logic for schedules
type ScheduleUsecase struct {
	repository repository.IRepository
	db         *sql.DB
}

// NewScheduleUsecase creates a new schedule usecase
func NewScheduleUsecase(repo repository.IRepository, db *sql.DB) domain.IScheduleUsecase {
	return &ScheduleUsecase{
		repository: repo,
		db:         db,
	}
}

// List returns all schedules
func (u *ScheduleUsecase) List(ctx context.Context) ([]*domain.Schedule, error) {
	schedules, err := u.repository.ListSchedules(ctx)
	if err != nil {
		return nil, err
	}

	return domain.NewScheduleDomainsFromQueries(schedules), nil
}

// GetByID returns a schedule by ID including its tasks
func (u *ScheduleUsecase) GetByID(ctx context.Context, id int32) (*domain.Schedule, error) {
	schedule, err := u.repository.GetScheduleByID(context.Background(), id)
	if err != nil {
		return nil, err
	}

	tasks, err := u.repository.GetTasksByScheduleID(ctx, id)
	if err != nil {
		return nil, err
	}

	scheduleDomain := domain.NewScheduleDomainFromQuery(schedule)
	scheduleDomain.Tasks = domain.NewTaskDomainsFromQueries(tasks)

	return scheduleDomain, nil
}

// GetToday returns today's schedules based on user timezone
func (u *ScheduleUsecase) GetToday(ctx context.Context, timezone string) ([]*domain.Schedule, error) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		slog.Warn("Invalid timezone, falling back to UTC", slog.String("timezone", timezone))
		loc = time.UTC
	}

	now := time.Now().In(loc)
	// Scoping "today" from midnight to midnight in the user's timezone.
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	todayEnd := todayStart.Add(24 * time.Hour)

	todaySchedule, err := u.repository.GetTodaySchedules(ctx, sqlc.GetTodaySchedulesParams{
		TodayStart: todayStart,
		TodayEnd:   todayEnd,
	})
	if err != nil {
		return nil, err
	}

	return domain.NewScheduleDomainsFromQueries(todaySchedule), nil
}

// GetStats returns dashboard statistics relative to user timezone
func (u *ScheduleUsecase) GetStats(ctx context.Context, timezone string) (*domain.ScheduleStats, error) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		slog.Warn("Invalid timezone, falling back to UTC", slog.String("timezone", timezone))
		loc = time.UTC
	}

	now := time.Now().In(loc)
	// Define "today" for upcoming and completed tasks dashboard numbers relative to user location
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	todayEnd := todayStart.Add(24 * time.Hour)

	stats, err := u.repository.GetScheduleStats(ctx, sqlc.GetScheduleStatsParams{
		TodayStart: todayStart,
		TodayEnd:   todayEnd,
	})
	if err != nil {
		return nil, err
	}

	return &domain.ScheduleStats{
		Total:     int(stats.Total),
		Missed:    int(stats.Missed),
		Upcoming:  int(stats.Upcoming),
		Completed: int(stats.Completed),
	}, nil
}

// ClockIn starts a visit
func (u *ScheduleUsecase) ClockIn(ctx context.Context, id int32, req domain.ClockInRequest) (*domain.Schedule, error) {
	// Validate schedule exists and is in valid state
	schedule, err := u.repository.GetScheduleByID(ctx, id)
	if err != nil {
		return nil, err
	}

	scheduleDomain := domain.NewScheduleDomainFromQuery(schedule)
	if scheduleDomain.Status != "upcoming" {
		return nil, errors.New("can only clock in to an upcoming schedule")
	}

	verified := req.Latitude != nil && req.Longitude != nil
	params := sqlc.ClockInParams{
		ID:       id,
		Verified: sql.NullBool{Bool: verified, Valid: true},
	}
	if verified {
		params.Latitude = sql.NullFloat64{Float64: *req.Latitude, Valid: true}
		params.Longitude = sql.NullFloat64{Float64: *req.Longitude, Valid: true}
	}

	schedule, err = u.repository.ClockIn(ctx, params)
	if err != nil {
		return nil, err
	}

	return domain.NewScheduleDomainFromQuery(schedule), nil
}

// ClockOut ends a visit
func (u *ScheduleUsecase) ClockOut(ctx context.Context, id int32, req domain.ClockOutRequest) (*domain.Schedule, error) {
	// Validate schedule exists and is in valid state
	schedule, err := u.repository.GetScheduleByID(ctx, id)
	if err != nil {
		return nil, err
	}

	scheduleDomain := domain.NewScheduleDomainFromQuery(schedule)
	if scheduleDomain.Status != "in_progress" {
		return nil, errors.New("can only clock out of an in-progress schedule")
	}

	verified := req.Latitude != nil && req.Longitude != nil
	params := sqlc.ClockOutParams{
		ID:       id,
		Verified: sql.NullBool{Bool: verified, Valid: true},
	}
	if verified {
		params.Latitude = sql.NullFloat64{Float64: *req.Latitude, Valid: true}
		params.Longitude = sql.NullFloat64{Float64: *req.Longitude, Valid: true}
	}

	schedule, err = u.repository.ClockOut(ctx, params)
	if err != nil {
		return nil, err
	}

	return domain.NewScheduleDomainFromQuery(schedule), nil
}

// Create creates a new schedule and its requested tasks
func (u *ScheduleUsecase) Create(ctx context.Context, req domain.CreateScheduleRequest) (*domain.Schedule, error) {
	startAt, err := time.Parse(time.RFC3339, req.StartAt)
	if err != nil {
		return nil, errors.New("invalid start_at format, expected RFC3339 (e.g. 2026-03-06T10:00:00Z)")
	}

	endAt, err := time.Parse(time.RFC3339, req.EndAt)
	if err != nil {
		return nil, errors.New("invalid end_at format, expected RFC3339 (e.g. 2026-03-06T12:00:00Z)")
	}

	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	qtx := u.repository.WithTx(tx)

	params := sqlc.CreateScheduleParams{
		ClientName: req.ClientName,
		StartAt:    startAt,
		EndAt:      endAt,
		Location:   req.Location,
		Latitude:   sql.NullFloat64{Float64: req.Latitude, Valid: true},
		Longitude:  sql.NullFloat64{Float64: req.Longitude, Valid: true},
	}

	schedule, err := qtx.CreateSchedule(ctx, params)
	if err != nil {
		return nil, err
	}

	var createdTasks []*domain.Task
	for _, taskReq := range req.Tasks {
		task, err := qtx.CreateTask(ctx, sqlc.CreateTaskParams{
			ScheduleID: schedule.ID,
			Title:      taskReq.Title,
		})
		if err != nil {
			return nil, err
		}
		createdTasks = append(createdTasks, domain.NewTaskDomainFromQuery(task))
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	scheduleDomain := domain.NewScheduleDomainFromQuery(schedule)
	scheduleDomain.Tasks = createdTasks

	return scheduleDomain, nil
}
