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

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestScheduleController_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := domainMock.NewMockIScheduleUsecase(ctrl)
	controller := NewScheduleController(mockUsecase)

	router := setupTestRouter()
	router.GET("/schedules", controller.List)

	testCases := []struct {
		name         string
		mock         func()
		expectedCode int
	}{
		{
			name: "success",
			mock: func() {
				mockUsecase.EXPECT().List(gomock.Any()).Return([]*domain.Schedule{
					{ID: 1, ClientName: "Client A"},
					{ID: 2, ClientName: "Client B"},
				}, nil)
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "internal server error",
			mock: func() {
				mockUsecase.EXPECT().List(gomock.Any()).Return(nil, errors.New("usecase error"))
			},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mock()

			req, err := http.NewRequest(http.MethodGet, "/schedules", nil)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedCode, w.Code)

			if tc.expectedCode == http.StatusOK {
				var response domain.Response
				err = json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.True(t, response.Success)
				assert.NotNil(t, response.Data)
			} else {
				var response domain.ErrorResponse
				err = json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.False(t, response.Success)
				assert.Contains(t, response.Error, "failed to fetch schedules")
			}
		})
	}
}

func TestScheduleController_GetToday(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := domainMock.NewMockIScheduleUsecase(ctrl)
	controller := NewScheduleController(mockUsecase)

	router := setupTestRouter()
	router.GET("/schedules/today", controller.GetToday)

	testCases := []struct {
		name         string
		setupReq     func(req *http.Request)
		mock         func()
		expectedCode int
	}{
		{
			name:     "success with default timezone",
			setupReq: func(req *http.Request) {},
			mock: func() {
				mockUsecase.EXPECT().GetToday(gomock.Any(), "UTC").Return([]*domain.Schedule{
					{ID: 1, ClientName: "Client A"},
				}, nil)
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "success with custom timezone",
			setupReq: func(req *http.Request) {
				req.Header.Set("X-Timezone", "Asia/Tokyo")
			},
			mock: func() {
				mockUsecase.EXPECT().GetToday(gomock.Any(), "Asia/Tokyo").Return([]*domain.Schedule{
					{ID: 1, ClientName: "Client A"},
				}, nil)
			},
			expectedCode: http.StatusOK,
		},
		{
			name:     "internal server error",
			setupReq: func(req *http.Request) {},
			mock: func() {
				mockUsecase.EXPECT().GetToday(gomock.Any(), "UTC").Return(nil, errors.New("usecase error"))
			},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mock()

			req, err := http.NewRequest(http.MethodGet, "/schedules/today", nil)
			require.NoError(t, err)
			tc.setupReq(req)

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

func TestScheduleController_GetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := domainMock.NewMockIScheduleUsecase(ctrl)
	controller := NewScheduleController(mockUsecase)

	router := setupTestRouter()
	router.GET("/schedules/:id", controller.GetByID)

	testCases := []struct {
		name         string
		url          string
		mock         func()
		expectedCode int
	}{
		{
			name: "success",
			url:  "/schedules/1",
			mock: func() {
				mockUsecase.EXPECT().GetByID(gomock.Any(), int32(1)).Return(&domain.Schedule{
					ID: 1, ClientName: "Client A",
				}, nil)
			},
			expectedCode: http.StatusOK,
		},
		{
			name:         "invalid id",
			url:          "/schedules/abc",
			mock:         func() {},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "not found",
			url:  "/schedules/99",
			mock: func() {
				mockUsecase.EXPECT().GetByID(gomock.Any(), int32(99)).Return(nil, errors.New("not found"))
			},
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mock()

			req, err := http.NewRequest(http.MethodGet, tc.url, nil)
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

func TestScheduleController_GetStats(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := domainMock.NewMockIScheduleUsecase(ctrl)
	controller := NewScheduleController(mockUsecase)

	router := setupTestRouter()
	router.GET("/schedules/stats", controller.GetStats)

	testCases := []struct {
		name         string
		mock         func()
		expectedCode int
	}{
		{
			name: "success",
			mock: func() {
				mockUsecase.EXPECT().GetStats(gomock.Any(), "UTC").Return(&domain.ScheduleStats{
					Total:     10,
					Missed:    2,
					Upcoming:  5,
					Completed: 3,
				}, nil)
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "error",
			mock: func() {
				mockUsecase.EXPECT().GetStats(gomock.Any(), "UTC").Return(nil, errors.New("usecase error"))
			},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mock()

			req, err := http.NewRequest(http.MethodGet, "/schedules/stats", nil)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedCode, w.Code)
		})
	}
}

func TestScheduleController_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := domainMock.NewMockIScheduleUsecase(ctrl)
	controller := NewScheduleController(mockUsecase)

	router := setupTestRouter()
	router.POST("/schedules", controller.Create)

	testCases := []struct {
		name         string
		reqBody      interface{}
		mock         func()
		expectedCode int
	}{
		{
			name: "success",
			reqBody: domain.CreateScheduleRequest{
				ClientName: "Client A",
				StartAt:    "2026-03-06T10:00:00Z",
				EndAt:      "2026-03-06T12:00:00Z",
				Location:   "Location",
				Latitude:   35.888,
				Longitude:  108.0999,
			},
			mock: func() {
				mockUsecase.EXPECT().Create(gomock.Any(), gomock.Any()).Return(&domain.Schedule{ID: 1}, nil)
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
			name: "internal server error",
			reqBody: domain.CreateScheduleRequest{
				ClientName: "Client A",
				StartAt:    "2026-03-06T10:00:00Z",
				EndAt:      "2026-03-06T12:00:00Z",
				Location:   "Location",
				Latitude:   1.23,
				Longitude:  4.56,
			},
			mock: func() {
				mockUsecase.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil, errors.New("create error"))
			},
			expectedCode: http.StatusInternalServerError,
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

			req, err := http.NewRequest(http.MethodPost, "/schedules", reqBody)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedCode, w.Code)
		})
	}
}

func TestScheduleController_ClockIn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := domainMock.NewMockIScheduleUsecase(ctrl)
	controller := NewScheduleController(mockUsecase)

	router := setupTestRouter()
	router.POST("/schedules/:id/clock-in", controller.ClockIn)

	lat := 1.23
	lng := 4.56

	testCases := []struct {
		name         string
		url          string
		reqBody      interface{}
		mock         func()
		expectedCode int
	}{
		{
			name: "success",
			url:  "/schedules/1/clock-in",
			reqBody: domain.ClockInRequest{
				Latitude:  &lat,
				Longitude: &lng,
			},
			mock: func() {
				mockUsecase.EXPECT().ClockIn(gomock.Any(), int32(1), gomock.Any()).Return(&domain.Schedule{ID: 1, ClockInVerified: true}, nil)
			},
			expectedCode: http.StatusOK,
		},
		{
			name:         "bad request - invalid id",
			url:          "/schedules/abc/clock-in",
			reqBody:      nil,
			mock:         func() {},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "bad request - invalid json",
			url:          "/schedules/1/clock-in",
			reqBody:      "invalid_json",
			mock:         func() {},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "usecase error",
			url:  "/schedules/1/clock-in",
			reqBody: domain.ClockInRequest{
				Latitude:  &lat,
				Longitude: &lng,
			},
			mock: func() {
				mockUsecase.EXPECT().ClockIn(gomock.Any(), int32(1), gomock.Any()).Return(nil, errors.New("already in progress"))
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

func TestScheduleController_ClockOut(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := domainMock.NewMockIScheduleUsecase(ctrl)
	controller := NewScheduleController(mockUsecase)

	router := setupTestRouter()
	router.POST("/schedules/:id/clock-out", controller.ClockOut)

	lat := 1.23
	lng := 4.56

	testCases := []struct {
		name         string
		url          string
		reqBody      interface{}
		mock         func()
		expectedCode int
	}{
		{
			name: "success",
			url:  "/schedules/1/clock-out",
			reqBody: domain.ClockOutRequest{
				Latitude:  &lat,
				Longitude: &lng,
			},
			mock: func() {
				mockUsecase.EXPECT().ClockOut(gomock.Any(), int32(1), gomock.Any()).Return(&domain.Schedule{ID: 1, ClockOutVerified: true}, nil)
			},
			expectedCode: http.StatusOK,
		},
		{
			name:         "bad request - invalid id",
			url:          "/schedules/abc/clock-out",
			reqBody:      nil,
			mock:         func() {},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "bad request - invalid json",
			url:          "/schedules/1/clock-out",
			reqBody:      "invalid_json",
			mock:         func() {},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "usecase error",
			url:  "/schedules/1/clock-out",
			reqBody: domain.ClockOutRequest{
				Latitude:  &lat,
				Longitude: &lng,
			},
			mock: func() {
				mockUsecase.EXPECT().ClockOut(gomock.Any(), int32(1), gomock.Any()).Return(nil, errors.New("not in progress"))
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
