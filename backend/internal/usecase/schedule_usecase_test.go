package usecase

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"
	"time"

	"bht-test/internal/domain"
	repoMock "bht-test/internal/repository/mock"
	"bht-test/internal/repository/sqlc"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestScheduleUsecase_ClockIn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repoMock.NewMockIRepository(ctrl)
	u := NewScheduleUsecase(mockRepo, nil)

	ctx := context.Background()
	scheduleID := int32(1)
	now := time.Now()
	lat := 35.6895
	lng := 139.6917

	testCases := []struct {
		name         string
		req          domain.ClockInRequest
		mock         func()
		expectErr    bool
		expectErrMsg string
	}{
		{
			name: "success verified",
			req: domain.ClockInRequest{
				Latitude:  &lat,
				Longitude: &lng,
			},
			mock: func() {
				mockRepo.EXPECT().GetScheduleByID(gomock.Any(), scheduleID).Return(sqlc.Schedule{
					ID:    scheduleID,
					EndAt: now.Add(time.Hour), // upcoming
				}, nil)

				params := sqlc.ClockInParams{
					ID:        scheduleID,
					Verified:  sql.NullBool{Bool: true, Valid: true},
					Latitude:  sql.NullFloat64{Float64: lat, Valid: true},
					Longitude: sql.NullFloat64{Float64: lng, Valid: true},
				}

				mockRepo.EXPECT().ClockIn(gomock.Any(), params).Return(sqlc.Schedule{
					ID:              scheduleID,
					ClockInAt:       sql.NullTime{Time: now, Valid: true},
					ClockInLat:      sql.NullFloat64{Float64: lat, Valid: true},
					ClockInLng:      sql.NullFloat64{Float64: lng, Valid: true},
					ClockInVerified: sql.NullBool{Bool: true, Valid: true},
				}, nil)
			},
			expectErr: false,
		},
		{
			name: "success unverified (fallback)",
			req: domain.ClockInRequest{
				Latitude:  nil,
				Longitude: nil,
			},
			mock: func() {
				mockRepo.EXPECT().GetScheduleByID(gomock.Any(), scheduleID).Return(sqlc.Schedule{
					ID:    scheduleID,
					EndAt: now.Add(time.Hour),
				}, nil)

				params := sqlc.ClockInParams{
					ID:       scheduleID,
					Verified: sql.NullBool{Bool: false, Valid: true},
				}

				mockRepo.EXPECT().ClockIn(gomock.Any(), params).Return(sqlc.Schedule{
					ID:              scheduleID,
					ClockInAt:       sql.NullTime{Time: now, Valid: true},
					ClockInVerified: sql.NullBool{Bool: false, Valid: true},
				}, nil)
			},
			expectErr: false,
		},
		{
			name: "failure - not upcoming",
			req: domain.ClockInRequest{
				Latitude:  &lat,
				Longitude: &lng,
			},
			mock: func() {
				mockRepo.EXPECT().GetScheduleByID(gomock.Any(), scheduleID).Return(sqlc.Schedule{
					ID:    scheduleID,
					EndAt: now.Add(-time.Hour), // not upcoming
				}, nil)
			},
			expectErr:    true,
			expectErrMsg: "can only clock in to an upcoming schedule",
		},
		{
			name: "failure - not found",
			req: domain.ClockInRequest{
				Latitude:  &lat,
				Longitude: &lng,
			},
			mock: func() {
				mockRepo.EXPECT().GetScheduleByID(gomock.Any(), scheduleID).Return(sqlc.Schedule{}, errors.New("not found"))
			},
			expectErr:    true,
			expectErrMsg: "not found",
		},
		{
			name: "failure - clock in failed",
			req: domain.ClockInRequest{
				Latitude:  &lat,
				Longitude: &lng,
			},
			mock: func() {
				mockRepo.EXPECT().GetScheduleByID(gomock.Any(), scheduleID).Return(sqlc.Schedule{
					ID:    scheduleID,
					EndAt: now.Add(time.Hour),
				}, nil)

				params := sqlc.ClockInParams{
					ID:        scheduleID,
					Verified:  sql.NullBool{Bool: true, Valid: true},
					Latitude:  sql.NullFloat64{Float64: lat, Valid: true},
					Longitude: sql.NullFloat64{Float64: lng, Valid: true},
				}

				mockRepo.EXPECT().ClockIn(gomock.Any(), params).Return(sqlc.Schedule{}, errors.New("clock in failed"))
			},
			expectErr:    true,
			expectErrMsg: "clock in failed",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mock()

			res, err := u.ClockIn(ctx, scheduleID, tc.req)
			if tc.expectErr {
				assert.Error(t, err)
				if tc.expectErrMsg != "" {
					assert.Equal(t, tc.expectErrMsg, err.Error())
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
			}
		})
	}
}

func TestScheduleUsecase_ClockOut(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repoMock.NewMockIRepository(ctrl)
	u := NewScheduleUsecase(mockRepo, nil)

	ctx := context.Background()
	scheduleID := int32(1)
	now := time.Now()

	lat := 35.6895
	lng := 139.6917

	testCases := []struct {
		name         string
		req          domain.ClockOutRequest
		mock         func()
		expectErr    bool
		expectErrMsg string
	}{
		{
			name: "success verified",
			req: domain.ClockOutRequest{
				Latitude:  &lat,
				Longitude: &lng,
			},
			mock: func() {
				mockRepo.EXPECT().GetScheduleByID(gomock.Any(), scheduleID).Return(sqlc.Schedule{
					ID:        scheduleID,
					ClockInAt: sql.NullTime{Time: now.Add(-time.Hour), Valid: true}, // in progress
				}, nil)

				params := sqlc.ClockOutParams{
					ID:        scheduleID,
					Verified:  sql.NullBool{Bool: true, Valid: true},
					Latitude:  sql.NullFloat64{Float64: lat, Valid: true},
					Longitude: sql.NullFloat64{Float64: lng, Valid: true},
				}

				mockRepo.EXPECT().ClockOut(gomock.Any(), params).Return(sqlc.Schedule{
					ID:               scheduleID,
					ClockInAt:        sql.NullTime{Time: now.Add(-time.Hour), Valid: true},
					ClockOutAt:       sql.NullTime{Time: now, Valid: true},
					ClockOutLat:      sql.NullFloat64{Float64: lat, Valid: true},
					ClockOutLng:      sql.NullFloat64{Float64: lng, Valid: true},
					ClockOutVerified: sql.NullBool{Bool: true, Valid: true},
				}, nil)
			},
			expectErr: false,
		},
		{
			name: "success unverified (fallback)",
			req: domain.ClockOutRequest{
				Latitude:  nil,
				Longitude: nil,
			},
			mock: func() {
				mockRepo.EXPECT().GetScheduleByID(gomock.Any(), scheduleID).Return(sqlc.Schedule{
					ID:        scheduleID,
					ClockInAt: sql.NullTime{Time: now.Add(-time.Hour), Valid: true}, // in progress
				}, nil)

				params := sqlc.ClockOutParams{
					ID:       scheduleID,
					Verified: sql.NullBool{Bool: false, Valid: true},
				}

				mockRepo.EXPECT().ClockOut(gomock.Any(), params).Return(sqlc.Schedule{
					ID:               scheduleID,
					ClockInAt:        sql.NullTime{Time: now.Add(-time.Hour), Valid: true},
					ClockOutAt:       sql.NullTime{Time: now, Valid: true},
					ClockOutVerified: sql.NullBool{Bool: false, Valid: true},
				}, nil)
			},
			expectErr: false,
		},
		{
			name: "failure - not in progress",
			req:  domain.ClockOutRequest{},
			mock: func() {
				mockRepo.EXPECT().GetScheduleByID(gomock.Any(), scheduleID).Return(sqlc.Schedule{
					ID:      scheduleID,
					StartAt: now.Add(time.Hour), // upcoming, not in progress
				}, nil)
			},
			expectErr:    true,
			expectErrMsg: "can only clock out of an in-progress schedule",
		},
		{
			name: "failure - not found",
			req:  domain.ClockOutRequest{},
			mock: func() {
				mockRepo.EXPECT().GetScheduleByID(gomock.Any(), scheduleID).Return(sqlc.Schedule{}, errors.New("not found"))
			},
			expectErr:    true,
			expectErrMsg: "not found",
		},
		{
			name: "failure - clock out failed",
			req: domain.ClockOutRequest{
				Latitude:  &lat,
				Longitude: &lng,
			},
			mock: func() {
				mockRepo.EXPECT().GetScheduleByID(gomock.Any(), scheduleID).Return(sqlc.Schedule{
					ID:        scheduleID,
					ClockInAt: sql.NullTime{Time: now.Add(-time.Hour), Valid: true}, // in progress
				}, nil)

				params := sqlc.ClockOutParams{
					ID:        scheduleID,
					Verified:  sql.NullBool{Bool: true, Valid: true},
					Latitude:  sql.NullFloat64{Float64: lat, Valid: true},
					Longitude: sql.NullFloat64{Float64: lng, Valid: true},
				}

				mockRepo.EXPECT().ClockOut(gomock.Any(), params).Return(sqlc.Schedule{}, errors.New("clock out failed"))
			},
			expectErr:    true,
			expectErrMsg: "clock out failed",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mock()

			res, err := u.ClockOut(ctx, scheduleID, tc.req)
			if tc.expectErr {
				assert.Error(t, err)
				if tc.expectErrMsg != "" {
					assert.Equal(t, tc.expectErrMsg, err.Error())
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
			}
		})
	}
}

func TestScheduleUsecase_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repoMock.NewMockIRepository(ctrl)
	u := NewScheduleUsecase(mockRepo, nil)

	ctx := context.Background()

	testCases := []struct {
		name         string
		mock         func()
		expectErr    bool
		expectErrMsg string
	}{
		{
			name: "success",
			mock: func() {
				mockRepo.EXPECT().ListSchedules(ctx).Return([]sqlc.Schedule{
					{ID: 1, ClientName: "Client A"},
					{ID: 2, ClientName: "Client B"},
				}, nil)
			},
			expectErr: false,
		},
		{
			name: "repository error",
			mock: func() {
				mockRepo.EXPECT().ListSchedules(ctx).Return(nil, errors.New("db error"))
			},
			expectErr: true,
		},
		{
			name: "empty list",
			mock: func() {
				mockRepo.EXPECT().ListSchedules(ctx).Return([]sqlc.Schedule{}, nil)
			},
			expectErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mock()
			res, err := u.List(ctx)
			if tc.expectErr {
				assert.Error(t, err)
				assert.Nil(t, res)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestScheduleUsecase_GetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repoMock.NewMockIRepository(ctrl)
	u := NewScheduleUsecase(mockRepo, nil)

	ctx := context.Background()
	scheduleID := int32(1)

	testCases := []struct {
		name         string
		id           int32
		mock         func()
		expectErr    bool
		expectErrMsg string
	}{
		{
			name: "success",
			id:   scheduleID,
			mock: func() {
				mockRepo.EXPECT().GetScheduleByID(gomock.Any(), scheduleID).Return(sqlc.Schedule{
					ID:         scheduleID,
					ClientName: "Client A",
				}, nil)

				mockRepo.EXPECT().GetTasksByScheduleID(ctx, scheduleID).Return([]sqlc.Task{
					{ID: 1, Title: "Task 1"},
				}, nil)
			},
			expectErr: false,
		},
		{
			name: "schedule not found",
			id:   scheduleID,
			mock: func() {
				mockRepo.EXPECT().GetScheduleByID(gomock.Any(), scheduleID).Return(sqlc.Schedule{}, errors.New("not found"))
			},
			expectErr:    true,
			expectErrMsg: "not found",
		},
		{
			name: "tasks error",
			id:   scheduleID,
			mock: func() {
				mockRepo.EXPECT().GetScheduleByID(gomock.Any(), scheduleID).Return(sqlc.Schedule{
					ID:         scheduleID,
					ClientName: "Client A",
				}, nil)

				mockRepo.EXPECT().GetTasksByScheduleID(ctx, scheduleID).Return(nil, errors.New("db error"))
			},
			expectErr:    true,
			expectErrMsg: "db error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mock()
			res, err := u.GetByID(ctx, tc.id)
			if tc.expectErr {
				assert.Error(t, err)
				if tc.expectErrMsg != "" {
					assert.Equal(t, tc.expectErrMsg, err.Error())
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
			}
		})
	}
}

func TestScheduleUsecase_GetToday(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repoMock.NewMockIRepository(ctrl)
	u := NewScheduleUsecase(mockRepo, nil)

	ctx := context.Background()

	testCases := []struct {
		name      string
		timezone  string
		mock      func()
		expectErr bool
	}{
		{
			name:     "success with specific timezone",
			timezone: "Asia/Tokyo",
			mock: func() {
				mockRepo.EXPECT().GetTodaySchedules(ctx, gomock.Any()).Return([]sqlc.Schedule{
					{ID: 1},
				}, nil)
			},
			expectErr: false,
		},
		{
			name:     "invalid timezone fallbacks to UTC",
			timezone: "Invalid/Zone",
			mock: func() {
				mockRepo.EXPECT().GetTodaySchedules(ctx, gomock.Any()).Return([]sqlc.Schedule{}, nil)
			},
			expectErr: false,
		},
		{
			name:     "repository error",
			timezone: "UTC",
			mock: func() {
				mockRepo.EXPECT().GetTodaySchedules(ctx, gomock.Any()).Return(nil, errors.New("db error"))
			},
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mock()
			res, err := u.GetToday(ctx, tc.timezone)
			if tc.expectErr {
				assert.Error(t, err)
				assert.Nil(t, res)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestScheduleUsecase_GetStats(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repoMock.NewMockIRepository(ctrl)
	u := NewScheduleUsecase(mockRepo, nil)

	ctx := context.Background()

	testCases := []struct {
		name      string
		timezone  string
		mock      func()
		expectErr bool
	}{
		{
			name:     "success",
			timezone: "UTC",
			mock: func() {
				mockRepo.EXPECT().GetScheduleStats(ctx, gomock.Any()).Return(sqlc.GetScheduleStatsRow{
					Total:     10,
					Missed:    2,
					Upcoming:  5,
					Completed: 3,
				}, nil)
			},
			expectErr: false,
		},
		{
			name:     "error map",
			timezone: "UTC",
			mock: func() {
				mockRepo.EXPECT().GetScheduleStats(ctx, gomock.Any()).Return(sqlc.GetScheduleStatsRow{}, errors.New("error"))
			},
			expectErr: true,
		},
		{
			name:     "invalid timezone fallbacks to UTC",
			timezone: "Invalid/Zone",
			mock: func() {
				mockRepo.EXPECT().GetScheduleStats(ctx, gomock.Any()).Return(sqlc.GetScheduleStatsRow{}, nil)
			},
			expectErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mock()
			res, err := u.GetStats(ctx, tc.timezone)
			if tc.expectErr {
				assert.Error(t, err)
				assert.Nil(t, res)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
			}
		})
	}
}

func TestScheduleUsecase_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repoMock.NewMockIRepository(ctrl)
	db, dbMock, _ := sqlmock.New()
	tx := sqlc.New(db)

	u := NewScheduleUsecase(mockRepo, db)

	ctx := context.Background()

	createScheduleSQL := regexp.QuoteMeta(`INSERT INTO schedules (
    client_name, start_at, end_at,
    location, latitude, longitude
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING id, client_name, start_at, end_at, location, latitude, longitude, clock_in_at, clock_in_lat, clock_in_lng, clock_out_at, clock_out_lat, clock_out_lng, created_at, updated_at, clock_in_verified, clock_out_verified`)

	createTaskSQL := regexp.QuoteMeta(`INSERT INTO tasks (
    schedule_id, title, status
) VALUES (
    $1, $2, 'pending'
)
RETURNING id, schedule_id, title, status, reason, created_at, updated_at`)

	testCases := []struct {
		name      string
		req       domain.CreateScheduleRequest
		mockFunc  func()
		expectErr bool
	}{
		{
			name: "success",
			req: domain.CreateScheduleRequest{
				ClientName: "Client A",
				StartAt:    "2026-03-06T10:00:00Z",
				EndAt:      "2026-03-06T12:00:00Z",
				Location:   "Location A",
				Latitude:   1.23,
				Longitude:  4.56,
				Tasks: []domain.CreateTaskRequest{
					{Title: "Task 1"},
				},
			},
			mockFunc: func() {
				dbMock.ExpectBegin()
				mockRepo.EXPECT().WithTx(gomock.Any()).Return(tx)

				startAt, _ := time.Parse(time.RFC3339, "2026-03-06T10:00:00Z")
				endAt, _ := time.Parse(time.RFC3339, "2026-03-06T12:00:00Z")

				dbMock.ExpectQuery(createScheduleSQL).
					WithArgs("Client A", startAt, endAt, "Location A", sql.NullFloat64{Float64: 1.23, Valid: true}, sql.NullFloat64{Float64: 4.56, Valid: true}).
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "client_name", "start_at", "end_at", "location", "latitude", "longitude",
						"clock_in_at", "clock_in_lat", "clock_in_lng", "clock_out_at", "clock_out_lat", "clock_out_lng",
						"created_at", "updated_at", "clock_in_verified", "clock_out_verified",
					}).AddRow(1, "Client A", startAt, endAt, "Location A", 1.23, 4.56, nil, nil, nil, nil, nil, nil, time.Now(), time.Now(), false, false))

				dbMock.ExpectQuery(createTaskSQL).
					WithArgs(1, "Task 1").
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "schedule_id", "title", "status", "reason", "created_at", "updated_at",
					}).AddRow(1, 1, "Task 1", "pending", nil, time.Now(), time.Now()))

				dbMock.ExpectCommit()
			},
			expectErr: false,
		},
		{
			name: "error creating schedule translates to rollback",
			req: domain.CreateScheduleRequest{
				ClientName: "Client A",
				StartAt:    "2026-03-06T10:00:00Z",
				EndAt:      "2026-03-06T12:00:00Z",
				Location:   "Location A",
				Latitude:   1.23,
				Longitude:  4.56,
				Tasks: []domain.CreateTaskRequest{
					{Title: "Task 1"},
				},
			},
			mockFunc: func() {
				dbMock.ExpectBegin()
				mockRepo.EXPECT().WithTx(gomock.Any()).Return(tx)

				startAt, _ := time.Parse(time.RFC3339, "2026-03-06T10:00:00Z")
				endAt, _ := time.Parse(time.RFC3339, "2026-03-06T12:00:00Z")

				dbMock.ExpectQuery(createScheduleSQL).
					WithArgs("Client A", startAt, endAt, "Location A", sql.NullFloat64{Float64: 1.23, Valid: true}, sql.NullFloat64{Float64: 4.56, Valid: true}).
					WillReturnError(errors.New("db error"))

				dbMock.ExpectRollback()
			},
			expectErr: true,
		},
		{
			name: "error creating task translates to rollback",
			req: domain.CreateScheduleRequest{
				ClientName: "Client A",
				StartAt:    "2026-03-06T10:00:00Z",
				EndAt:      "2026-03-06T12:00:00Z",
				Location:   "Location A",
				Latitude:   1.23,
				Longitude:  4.56,
				Tasks: []domain.CreateTaskRequest{
					{Title: "Task 1"},
				},
			},
			mockFunc: func() {
				dbMock.ExpectBegin()
				mockRepo.EXPECT().WithTx(gomock.Any()).Return(tx)

				startAt, _ := time.Parse(time.RFC3339, "2026-03-06T10:00:00Z")
				endAt, _ := time.Parse(time.RFC3339, "2026-03-06T12:00:00Z")

				dbMock.ExpectQuery(createScheduleSQL).
					WithArgs("Client A", startAt, endAt, "Location A", sql.NullFloat64{Float64: 1.23, Valid: true}, sql.NullFloat64{Float64: 4.56, Valid: true}).
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "client_name", "start_at", "end_at", "location", "latitude", "longitude",
						"clock_in_at", "clock_in_lat", "clock_in_lng", "clock_out_at", "clock_out_lat", "clock_out_lng",
						"created_at", "updated_at", "clock_in_verified", "clock_out_verified",
					}).AddRow(1, "Client A", startAt, endAt, "Location A", 1.23, 4.56, nil, nil, nil, nil, nil, nil, time.Now(), time.Now(), false, false))

				dbMock.ExpectQuery(createTaskSQL).
					WithArgs(1, "Task 1").
					WillReturnError(errors.New("db error"))

				dbMock.ExpectRollback()
			},
			expectErr: true,
		},
		{
			name:     "error parse start time",
			mockFunc: func() {},
			req: domain.CreateScheduleRequest{
				ClientName: "Client A",
				StartAt:    "invalid-time",
				EndAt:      "2026-03-06T12:00:00Z",
				Location:   "Location A",
				Latitude:   1.23,
				Longitude:  4.56,
				Tasks: []domain.CreateTaskRequest{
					{Title: "Task 1"},
				},
			},
			expectErr: true,
		},
		{
			name:     "error parse end time",
			mockFunc: func() {},
			req: domain.CreateScheduleRequest{
				ClientName: "Client A",
				StartAt:    "2026-03-06T10:00:00Z",
				EndAt:      "invalid-time",
				Location:   "Location A",
				Latitude:   1.23,
				Longitude:  4.56,
				Tasks: []domain.CreateTaskRequest{
					{Title: "Task 1"},
				},
			},
			expectErr: true,
		},
		{
			name: "error on begin",
			mockFunc: func() {
				dbMock.ExpectBegin().WillReturnError(errors.New("db error"))
			},
			req: domain.CreateScheduleRequest{
				ClientName: "Client A",
				StartAt:    "2026-03-06T10:00:00Z",
				EndAt:      "2026-03-06T12:00:00Z",
				Location:   "Location A",
				Latitude:   1.23,
				Longitude:  4.56,
				Tasks: []domain.CreateTaskRequest{
					{Title: "Task 1"},
				},
			},
			expectErr: true,
		},
		{
			name: "error on commit",
			mockFunc: func() {
				dbMock.ExpectBegin()
				mockRepo.EXPECT().WithTx(gomock.Any()).Return(tx)

				startAt, _ := time.Parse(time.RFC3339, "2026-03-06T10:00:00Z")
				endAt, _ := time.Parse(time.RFC3339, "2026-03-06T12:00:00Z")

				dbMock.ExpectQuery(createScheduleSQL).
					WithArgs("Client A", startAt, endAt, "Location A", sql.NullFloat64{Float64: 1.23, Valid: true}, sql.NullFloat64{Float64: 4.56, Valid: true}).
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "client_name", "start_at", "end_at", "location", "latitude", "longitude",
						"clock_in_at", "clock_in_lat", "clock_in_lng", "clock_out_at", "clock_out_lat", "clock_out_lng",
						"created_at", "updated_at", "clock_in_verified", "clock_out_verified",
					}).AddRow(1, "Client A", startAt, endAt, "Location A", 1.23, 4.56, nil, nil, nil, nil, nil, nil, time.Now(), time.Now(), false, false))

				dbMock.ExpectQuery(createTaskSQL).
					WithArgs(1, "Task 1").
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "schedule_id", "title", "status", "reason", "created_at", "updated_at",
					}).AddRow(1, 1, "Task 1", "pending", nil, time.Now(), time.Now()))

				dbMock.ExpectCommit().WillReturnError(errors.New("db error"))
			},
			req: domain.CreateScheduleRequest{
				ClientName: "Client A",
				StartAt:    "2026-03-06T10:00:00Z",
				EndAt:      "2026-03-06T12:00:00Z",
				Location:   "Location A",
				Latitude:   1.23,
				Longitude:  4.56,
				Tasks: []domain.CreateTaskRequest{
					{Title: "Task 1"},
				},
			},
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockFunc()

			res, err := u.Create(ctx, tc.req)
			if tc.expectErr {
				assert.Error(t, err)
				assert.Nil(t, res)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
			}
			assert.NoError(t, dbMock.ExpectationsWereMet())
		})
	}
}
