package domain

import (
	"bht-test/internal/repository/sqlc"
	"context"
	"time"
)

// Schedule represents a caregiver's shift schedule
type Schedule struct {
	ID               int32      `json:"id"`
	ClientName       string     `json:"client_name"`
	StartAt          time.Time  `json:"start_at"`
	EndAt            time.Time  `json:"end_at"`
	Location         string     `json:"location"`
	Status           string     `json:"status"`
	ClockInAt        *time.Time `json:"clock_in_at"`
	ClockInLat       *float64   `json:"clock_in_lat"`
	ClockInLng       *float64   `json:"clock_in_lng"`
	ClockOutAt       *time.Time `json:"clock_out_at"`
	ClockOutLat      *float64   `json:"clock_out_lat"`
	ClockOutLng      *float64   `json:"clock_out_lng"`
	Latitude         *float64   `json:"latitude"`
	Longitude        *float64   `json:"longitude"`
	Tasks            []*Task    `json:"tasks,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	ClockInVerified  bool       `json:"clock_in_verified"`
	ClockOutVerified bool       `json:"clock_out_verified"`
}

// ScheduleStats holds dashboard statistics
type ScheduleStats struct {
	Total     int `json:"total"`
	Missed    int `json:"missed"`
	Upcoming  int `json:"upcoming"`
	Completed int `json:"completed"`
}

// ClockInRequest represents the request body for clocking in
type ClockInRequest struct {
	Latitude  *float64 `json:"latitude" binding:"omitempty,gte=-90,lte=90"`
	Longitude *float64 `json:"longitude" binding:"omitempty,gte=-180,lte=180"`
}

// ClockOutRequest represents the request body for clocking out
type ClockOutRequest struct {
	Latitude  *float64 `json:"latitude" binding:"omitempty,gte=-90,lte=90"`
	Longitude *float64 `json:"longitude" binding:"omitempty,gte=-180,lte=180"`
}

type CreateScheduleRequest struct {
	ClientName string              `json:"client_name" binding:"required" example:"John Doe"`
	StartAt    string              `json:"start_at" binding:"required" example:"2026-03-06T10:00:00+09:00"`
	EndAt      string              `json:"end_at" binding:"required" example:"2026-03-06T18:00:00+09:00"`
	Location   string              `json:"location" binding:"required" example:"Tokyo Branch"`
	Latitude   float64             `json:"latitude" binding:"required,gte=-90,lte=90" example:"35.6895"`
	Longitude  float64             `json:"longitude" binding:"required,gte=-180,lte=180" example:"139.6917"`
	Tasks      []CreateTaskRequest `json:"tasks,omitempty"`
}

func NewScheduleDomainFromQuery(schedule sqlc.Schedule) *Schedule {
	now := time.Now()
	scheduleDomain := &Schedule{
		ID:         schedule.ID,
		ClientName: schedule.ClientName,
		StartAt:    schedule.StartAt,
		EndAt:      schedule.EndAt,
		Location:   schedule.Location,
		CreatedAt:  schedule.CreatedAt,
		UpdatedAt:  schedule.UpdatedAt,
	}

	if schedule.ClockInAt.Valid {
		scheduleDomain.ClockInAt = &schedule.ClockInAt.Time
	}
	if schedule.ClockInLat.Valid {
		scheduleDomain.ClockInLat = &schedule.ClockInLat.Float64
	}
	if schedule.ClockInLng.Valid {
		scheduleDomain.ClockInLng = &schedule.ClockInLng.Float64
	}
	if schedule.ClockOutAt.Valid {
		scheduleDomain.ClockOutAt = &schedule.ClockOutAt.Time
	}
	if schedule.ClockOutLat.Valid {
		scheduleDomain.ClockOutLat = &schedule.ClockOutLat.Float64
	}
	if schedule.ClockOutLng.Valid {
		scheduleDomain.ClockOutLng = &schedule.ClockOutLng.Float64
	}
	if schedule.Latitude.Valid {
		scheduleDomain.Latitude = &schedule.Latitude.Float64
	}
	if schedule.Longitude.Valid {
		scheduleDomain.Longitude = &schedule.Longitude.Float64
	}
	if schedule.ClockInVerified.Valid {
		scheduleDomain.ClockInVerified = schedule.ClockInVerified.Bool
	}
	if schedule.ClockOutVerified.Valid {
		scheduleDomain.ClockOutVerified = schedule.ClockOutVerified.Bool
	}

	// Dynamic Status Calculation
	if scheduleDomain.ClockOutAt != nil {
		scheduleDomain.Status = "completed"
	} else if scheduleDomain.ClockInAt != nil {
		scheduleDomain.Status = "in_progress"
	} else if scheduleDomain.EndAt.Before(now) {
		scheduleDomain.Status = "missed"
	} else {
		scheduleDomain.Status = "upcoming"
	}

	return scheduleDomain
}

func NewScheduleDomainsFromQueries(schedules []sqlc.Schedule) []*Schedule {
	var result []*Schedule
	for _, schedule := range schedules {
		result = append(result, NewScheduleDomainFromQuery(schedule))
	}
	return result
}

//go:generate mockgen -source=schedule.go -destination=mock/schedule.go -package=mock
type IScheduleUsecase interface {
	List(ctx context.Context) ([]*Schedule, error)
	GetByID(ctx context.Context, id int32) (*Schedule, error)
	GetToday(ctx context.Context, timezone string) ([]*Schedule, error)
	GetStats(ctx context.Context, timezone string) (*ScheduleStats, error)
	ClockIn(ctx context.Context, id int32, req ClockInRequest) (*Schedule, error)
	ClockOut(ctx context.Context, id int32, req ClockOutRequest) (*Schedule, error)
	Create(ctx context.Context, req CreateScheduleRequest) (*Schedule, error)
}
