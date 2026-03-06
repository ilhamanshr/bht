package controllers

import (
	"log/slog"
	"net/http"
	"strconv"

	"bht-test/internal/domain"

	"github.com/gin-gonic/gin"
)

// ScheduleController handles HTTP requests for schedules
type ScheduleController struct {
	usecase domain.IScheduleUsecase
}

// NewScheduleController creates a new schedule controller
func NewScheduleController(u domain.IScheduleUsecase) *ScheduleController {
	return &ScheduleController{usecase: u}
}

// List returns all schedules
// @Summary      List all schedules
// @Description  Get a list of all caregiver schedules ordered by date
// @Tags         schedules
// @Produce      json
// @Success      200  {object}  domain.Response{data=[]domain.Schedule}
// @Failure      500  {object}  domain.ErrorResponse
// @Router       /schedules [get]
func (h *ScheduleController) List(c *gin.Context) {
	schedules, err := h.usecase.List(c.Request.Context())
	if err != nil {
		slog.Error("Failed to list schedules", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, domain.NewErrorResponse("failed to fetch schedules"))
		return
	}

	c.JSON(http.StatusOK, domain.NewSuccessResponse(schedules))
}

// GetToday returns today's schedules
// @Summary      Get today's schedules
// @Description  Get all schedules for the current date relative to user timezone
// @Tags         schedules
// @Produce      json
// @Param        X-Timezone  header    string  false  "User timezone (e.g. Asia/Tokyo)" default(UTC)
// @Success      200  {object}  domain.Response{data=[]domain.Schedule}
// @Failure      500  {object}  domain.ErrorResponse
// @Router       /schedules/today [get]
func (h *ScheduleController) GetToday(c *gin.Context) {
	timezone := c.GetHeader("X-Timezone")
	if timezone == "" {
		timezone = "UTC"
	}

	schedules, err := h.usecase.GetToday(c.Request.Context(), timezone)
	if err != nil {
		slog.Error("Failed to get today schedules", slog.String("timezone", timezone), slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, domain.NewErrorResponse("failed to fetch today's schedules"))
		return
	}

	c.JSON(http.StatusOK, domain.NewSuccessResponse(schedules))
}

// GetByID returns a single schedule by ID
// @Summary      Get schedule by ID
// @Description  Get a single schedule with its associated care tasks
// @Tags         schedules
// @Produce      json
// @Param        id   path      int  true  "Schedule ID"
// @Success      200  {object}  domain.Response{data=domain.Schedule}
// @Failure      400  {object}  domain.ErrorResponse
// @Failure      404  {object}  domain.ErrorResponse
// @Router       /schedules/{id} [get]
func (h *ScheduleController) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.NewErrorResponse("invalid schedule ID"))
		return
	}

	schedule, err := h.usecase.GetByID(c.Request.Context(), int32(id))
	if err != nil {
		slog.Error("Failed to get schedule", slog.Int("schedule_id", id), slog.Any("error", err))
		c.JSON(http.StatusNotFound, domain.NewErrorResponse("schedule not found"))
		return
	}

	c.JSON(http.StatusOK, domain.NewSuccessResponse(schedule))
}

// GetStats returns schedule statistics
// @Summary      Get schedule statistics
// @Description  Get dashboard statistics relative to user timezone
// @Tags         schedules
// @Produce      json
// @Param        X-Timezone  header    string  false  "User timezone (e.g. Asia/Tokyo)" default(UTC)
// @Success      200  {object}  domain.Response{data=domain.ScheduleStats}
// @Failure      500  {object}  domain.ErrorResponse
// @Router       /schedules/stats [get]
func (h *ScheduleController) GetStats(c *gin.Context) {
	timezone := c.GetHeader("X-Timezone")
	if timezone == "" {
		timezone = "UTC"
	}

	stats, err := h.usecase.GetStats(c.Request.Context(), timezone)
	if err != nil {
		slog.Error("Failed to get stats", slog.String("timezone", timezone), slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, domain.NewErrorResponse("failed to fetch stats"))
		return
	}

	c.JSON(http.StatusOK, domain.NewSuccessResponse(stats))
}

// Create creates a new schedule
// @Summary      Create a new schedule
// @Description  Create a new schedule along with its associated tasks
// @Tags         schedules
// @Accept       json
// @Produce      json
// @Param        body  body      domain.CreateScheduleRequest  true  "Schedule and tasks details"
// @Success      201   {object}  domain.Response{data=domain.Schedule}
// @Failure      400   {object}  domain.ErrorResponse
// @Failure      500   {object}  domain.ErrorResponse
// @Router       /schedules [post]
func (h *ScheduleController) Create(c *gin.Context) {
	var req domain.CreateScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.NewErrorResponse("invalid request body: "+err.Error()))
		return
	}

	schedule, err := h.usecase.Create(c.Request.Context(), req)
	if err != nil {
		slog.Error("Failed to create schedule", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, domain.NewErrorResponse("failed to create schedule: "+err.Error()))
		return
	}

	c.JSON(http.StatusCreated, domain.NewSuccessResponse(schedule))
}

// ClockIn starts a visit
// @Summary      Clock in (start visit)
// @Description  Start a visit by recording timestamp and geolocation. Schedule must be in 'upcoming' status.
// @Tags         schedules
// @Accept       json
// @Produce      json
// @Param        id    path      int                   true  "Schedule ID"
// @Param        body  body      domain.ClockInRequest  true  "Geolocation coordinates"
// @Success      200   {object}  domain.Response{data=domain.Schedule}
// @Failure      400   {object}  domain.ErrorResponse
// @Router       /schedules/{id}/clock-in [post]
func (h *ScheduleController) ClockIn(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.NewErrorResponse("invalid schedule ID"))
		return
	}

	var req domain.ClockInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.NewErrorResponse("latitude and longitude are required"))
		return
	}

	schedule, err := h.usecase.ClockIn(c.Request.Context(), int32(id), req)
	if err != nil {
		slog.Error("Failed to clock in schedule", slog.Int("schedule_id", id), slog.Any("error", err))
		c.JSON(http.StatusBadRequest, domain.NewErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, domain.NewSuccessResponse(schedule))
}

// ClockOut ends a visit
// @Summary      Clock out (end visit)
// @Description  End a visit by recording timestamp and geolocation. Schedule must be in 'in_progress' status.
// @Tags         schedules
// @Accept       json
// @Produce      json
// @Param        id    path      int                    true  "Schedule ID"
// @Param        body  body      domain.ClockOutRequest  true  "Geolocation coordinates"
// @Success      200   {object}  domain.Response{data=domain.Schedule}
// @Failure      400   {object}  domain.ErrorResponse
// @Router       /schedules/{id}/clock-out [post]
func (h *ScheduleController) ClockOut(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.NewErrorResponse("invalid schedule ID"))
		return
	}

	var req domain.ClockOutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.NewErrorResponse("latitude and longitude are required"))
		return
	}

	schedule, err := h.usecase.ClockOut(c.Request.Context(), int32(id), req)
	if err != nil {
		slog.Error("Failed to clock out schedule", slog.Int("schedule_id", id), slog.Any("error", err))
		c.JSON(http.StatusBadRequest, domain.NewErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, domain.NewSuccessResponse(schedule))
}
