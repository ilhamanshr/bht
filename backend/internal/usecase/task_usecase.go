package usecase

import (
	"context"
	"database/sql"
	"errors"

	"bht-test/internal/domain"
	"bht-test/internal/repository"
	"bht-test/internal/repository/sqlc"
)

// TaskUsecase handles business logic for tasks
type TaskUsecase struct {
	repository repository.IRepository
}

// NewTaskUsecase creates a new task usecase
func NewTaskUsecase(repository repository.IRepository) domain.ITaskUsecase {
	return &TaskUsecase{repository: repository}
}

// Create adds a new task to a schedule
func (u *TaskUsecase) Create(ctx context.Context, req domain.CreateTaskRequest) (*domain.Task, error) {
	// Validate schedule exists
	_, err := u.repository.GetScheduleByID(ctx, req.ScheduleID)
	if err != nil {
		return nil, err
	}

	task, err := u.repository.CreateTask(ctx, sqlc.CreateTaskParams{
		ScheduleID: req.ScheduleID,
		Title:      req.Title,
	})
	if err != nil {
		return nil, err
	}

	return domain.NewTaskDomainFromQuery(task), nil
}

// UpdateStatus updates the status of a task
func (u *TaskUsecase) UpdateStatus(ctx context.Context, id int32, req domain.UpdateTaskRequest) (*domain.Task, error) {
	// Validate task exists
	_, err := u.repository.GetTaskByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// If marking as not_completed, reason is required
	if req.Status == "not_completed" && req.Reason == "" {
		return nil, errors.New("reason is required when marking task as not completed")
	}

	task, err := u.repository.UpdateTaskStatus(ctx, sqlc.UpdateTaskStatusParams{
		ID:     id,
		Status: req.Status,
		Reason: sql.NullString{String: req.Reason, Valid: req.Reason != ""},
	})
	if err != nil {
		return nil, err
	}

	return domain.NewTaskDomainFromQuery(task), nil
}
