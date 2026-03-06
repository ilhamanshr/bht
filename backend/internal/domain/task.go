package domain

import (
	"bht-test/internal/repository/sqlc"
	"context"
	"time"
)

// Task represents a care activity within a schedule
type Task struct {
	ID         int32     `json:"id"`
	ScheduleID int32     `json:"schedule_id"`
	Title      string    `json:"title"`
	Status     string    `json:"status"`
	Reason     *string   `json:"reason"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// UpdateTaskRequest represents the request body for updating a task
type UpdateTaskRequest struct {
	Status string `json:"status" binding:"required,oneof=completed not_completed"`
	Reason string `json:"reason" binding:"omitempty"`
}

// CreateTaskRequest represents the request body for creating a task
type CreateTaskRequest struct {
	ScheduleID int32  `json:"schedule_id" binding:"required"`
	Title      string `json:"title" binding:"required"`
}

func NewTaskDomainFromQuery(task sqlc.Task) *Task {
	taskDomain := &Task{
		ID:         task.ID,
		ScheduleID: task.ScheduleID,
		Title:      task.Title,
		Status:     task.Status,
		CreatedAt:  task.CreatedAt,
		UpdatedAt:  task.UpdatedAt,
	}

	if task.Reason.Valid {
		taskDomain.Reason = &task.Reason.String
	}

	return taskDomain
}

func NewTaskDomainsFromQueries(tasks []sqlc.Task) []*Task {
	var result []*Task
	for _, task := range tasks {
		result = append(result, NewTaskDomainFromQuery(task))
	}
	return result
}

//go:generate mockgen -source=task.go -destination=mock/task.go -package=mock
type ITaskUsecase interface {
	Create(ctx context.Context, req CreateTaskRequest) (*Task, error)
	UpdateStatus(ctx context.Context, id int32, req UpdateTaskRequest) (*Task, error)
}
