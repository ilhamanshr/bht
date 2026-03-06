package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"bht-test/internal/domain"
	domainMock "bht-test/internal/domain/mock"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTaskController_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := domainMock.NewMockITaskUsecase(ctrl)
	controller := NewTaskController(mockUsecase)

	router := setupTestRouter()
	router.POST("/tasks", controller.Create)

	testCases := []struct {
		name         string
		reqBody      interface{}
		mock         func()
		expectedCode int
	}{
		{
			name: "success",
			reqBody: domain.CreateTaskRequest{
				ScheduleID: 1,
				Title:      "New Task",
			},
			mock: func() {
				mockUsecase.EXPECT().Create(gomock.Any(), gomock.Any()).Return(&domain.Task{ID: 1, Title: "New Task"}, nil)
			},
			expectedCode: http.StatusCreated,
		},
		{
			name:         "bad request - invalid json",
			reqBody:      "invalid_json",
			mock:         func() {},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "usecase error",
			reqBody: domain.CreateTaskRequest{
				ScheduleID: 1,
				Title:      "New Task",
			},
			mock: func() {
				mockUsecase.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil, errors.New("schedule not found"))
			},
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mock()

			var reqBody *bytes.Buffer
			if strBody, ok := tc.reqBody.(string); ok {
				reqBody = bytes.NewBufferString(strBody)
			} else {
				jsonBody, _ := json.Marshal(tc.reqBody)
				reqBody = bytes.NewBuffer(jsonBody)
			}

			req, err := http.NewRequest(http.MethodPost, "/tasks", reqBody)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedCode, w.Code)

			if tc.expectedCode == http.StatusCreated {
				var response domain.Response
				err = json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.True(t, response.Success)
			}
		})
	}
}

func TestTaskController_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := domainMock.NewMockITaskUsecase(ctrl)
	controller := NewTaskController(mockUsecase)

	router := setupTestRouter()
	router.POST("/tasks/:taskId/update", controller.Update)

	testCases := []struct {
		name         string
		url          string
		reqBody      interface{}
		mock         func()
		expectedCode int
	}{
		{
			name: "success",
			url:  "/tasks/1/update",
			reqBody: domain.UpdateTaskRequest{
				Status: "completed",
			},
			mock: func() {
				mockUsecase.EXPECT().UpdateStatus(gomock.Any(), int32(1), gomock.Any()).Return(&domain.Task{ID: 1, Status: "completed"}, nil)
			},
			expectedCode: http.StatusOK,
		},
		{
			name:         "bad request - invalid id",
			url:          "/tasks/abc/update",
			reqBody:      nil,
			mock:         func() {},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "bad request - invalid json",
			url:          "/tasks/1/update",
			reqBody:      "invalid_json",
			mock:         func() {},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "usecase error",
			url:  "/tasks/1/update",
			reqBody: domain.UpdateTaskRequest{
				Status: "not_completed",
			}, // missing reason
			mock: func() {
				mockUsecase.EXPECT().UpdateStatus(gomock.Any(), int32(1), gomock.Any()).Return(nil, errors.New("reason is required"))
			},
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mock()

			var reqBody *bytes.Buffer
			if tc.reqBody == nil {
				reqBody = nil
			} else if strBody, ok := tc.reqBody.(string); ok {
				reqBody = bytes.NewBufferString(strBody)
			} else {
				jsonBody, _ := json.Marshal(tc.reqBody)
				reqBody = bytes.NewBuffer(jsonBody)
			}

			var req *http.Request
			var err error
			if reqBody != nil {
				req, err = http.NewRequest(http.MethodPost, tc.url, reqBody)
			} else {
				req, err = http.NewRequest(http.MethodPost, tc.url, nil)
			}
			require.NoError(t, err)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedCode, w.Code)

			if tc.expectedCode == http.StatusOK {
				var response domain.Response
				err = json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.True(t, response.Success)
			}
		})
	}
}
