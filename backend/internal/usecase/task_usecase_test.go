package usecase

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"bht-test/internal/domain"
	repoMock "bht-test/internal/repository/mock"
	"bht-test/internal/repository/sqlc"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestTaskUsecase_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repoMock.NewMockIRepository(ctrl)
	u := NewTaskUsecase(mockRepo)

	ctx := context.Background()

	testCases := []struct {
		name      string
		req       domain.CreateTaskRequest
		mockFunc  func()
		expectErr bool
	}{
		{
			name: "success",
			req: domain.CreateTaskRequest{
				ScheduleID: int32(1),
				Title:      "New Task",
			},
			mockFunc: func() {
				mockRepo.EXPECT().GetScheduleByID(ctx, int32(1)).Return(sqlc.Schedule{ID: 1}, nil)
				mockRepo.EXPECT().CreateTask(ctx, sqlc.CreateTaskParams{
					ScheduleID: int32(1),
					Title:      "New Task",
				}).Return(sqlc.Task{
					ID:         1,
					ScheduleID: 1,
					Title:      "New Task",
					Status:     "pending",
				}, nil)
			},
			expectErr: false,
		},
		{
			name: "schedule not found",
			req: domain.CreateTaskRequest{
				ScheduleID: int32(1),
				Title:      "New Task",
			},
			mockFunc: func() {
				mockRepo.EXPECT().GetScheduleByID(ctx, int32(1)).Return(sqlc.Schedule{}, errors.New("not found"))
			},
			expectErr: true,
		},
		{
			name: "error creating task",
			req: domain.CreateTaskRequest{
				ScheduleID: int32(1),
				Title:      "New Task",
			},
			mockFunc: func() {
				mockRepo.EXPECT().GetScheduleByID(ctx, int32(1)).Return(sqlc.Schedule{ID: 1}, nil)
				mockRepo.EXPECT().CreateTask(ctx, gomock.Any()).Return(sqlc.Task{}, errors.New("db error"))
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
				assert.Equal(t, tc.req.Title, res.Title)
			}
		})
	}
}

func TestTaskUsecase_UpdateStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repoMock.NewMockIRepository(ctrl)
	u := NewTaskUsecase(mockRepo)

	ctx := context.Background()
	taskID := int32(1)

	testCases := []struct {
		name      string
		req       domain.UpdateTaskRequest
		mockFunc  func()
		expectErr bool
		errMsg    string
	}{
		{
			name: "success marked completed",
			req: domain.UpdateTaskRequest{
				Status: "completed",
			},
			mockFunc: func() {
				mockRepo.EXPECT().GetTaskByID(ctx, taskID).Return(sqlc.Task{ID: taskID}, nil)
				mockRepo.EXPECT().UpdateTaskStatus(ctx, sqlc.UpdateTaskStatusParams{
					ID:     taskID,
					Status: "completed",
					Reason: sql.NullString{String: "", Valid: false},
				}).Return(sqlc.Task{
					ID:     taskID,
					Status: "completed",
				}, nil)
			},
			expectErr: false,
		},
		{
			name: "success marked not_completed with reason",
			req: domain.UpdateTaskRequest{
				Status: "not_completed",
				Reason: "patient absent",
			},
			mockFunc: func() {
				mockRepo.EXPECT().GetTaskByID(ctx, taskID).Return(sqlc.Task{ID: taskID}, nil)
				mockRepo.EXPECT().UpdateTaskStatus(ctx, sqlc.UpdateTaskStatusParams{
					ID:     taskID,
					Status: "not_completed",
					Reason: sql.NullString{String: "patient absent", Valid: true},
				}).Return(sqlc.Task{
					ID:     taskID,
					Status: "not_completed",
					Reason: sql.NullString{String: "patient absent", Valid: true},
				}, nil)
			},
			expectErr: false,
		},
		{
			name: "fail marked not_completed without reason",
			req: domain.UpdateTaskRequest{
				Status: "not_completed",
				Reason: "", // missing reason
			},
			mockFunc: func() {
				mockRepo.EXPECT().GetTaskByID(ctx, taskID).Return(sqlc.Task{ID: taskID}, nil)
			},
			expectErr: true,
			errMsg:    "reason is required when marking task as not completed",
		},
		{
			name: "task not found",
			req: domain.UpdateTaskRequest{
				Status: "completed",
			},
			mockFunc: func() {
				mockRepo.EXPECT().GetTaskByID(ctx, taskID).Return(sqlc.Task{}, errors.New("not found"))
			},
			expectErr: true,
			errMsg:    "not found",
		},
		{
			name: "error updating status",
			req: domain.UpdateTaskRequest{
				Status: "completed",
			},
			mockFunc: func() {
				mockRepo.EXPECT().GetTaskByID(ctx, taskID).Return(sqlc.Task{ID: taskID}, nil)
				mockRepo.EXPECT().UpdateTaskStatus(ctx, gomock.Any()).Return(sqlc.Task{}, errors.New("db error"))
			},
			expectErr: true,
			errMsg:    "db error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockFunc()

			res, err := u.UpdateStatus(ctx, taskID, tc.req)
			if tc.expectErr {
				assert.Error(t, err)
				assert.Nil(t, res)
				if tc.errMsg != "" {
					assert.Equal(t, tc.errMsg, err.Error())
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, tc.req.Status, res.Status)
			}
		})
	}
}
