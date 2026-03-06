package controllers

import (
	"log/slog"
	"net/http"
	"strconv"

	"bht-test/internal/domain"

	"github.com/gin-gonic/gin"
)

// TaskController handles HTTP requests for tasks
type TaskController struct {
	usecase domain.ITaskUsecase
}

// NewTaskController creates a new task controller
func NewTaskController(u domain.ITaskUsecase) *TaskController {
	return &TaskController{usecase: u}
}

// Create adds a new task to a schedule
// @Summary      Add task to schedule
// @Description  Create a new care task associated with a specific schedule
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        body  body      domain.CreateTaskRequest  true  "Task details"
// @Success      201   {object}  domain.Response{data=domain.Task}
// @Failure      400   {object}  domain.ErrorResponse
// @Router       /tasks [post]
func (h *TaskController) Create(c *gin.Context) {
	var req domain.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.NewErrorResponse("title is required"))
		return
	}

	task, err := h.usecase.Create(c.Request.Context(), req)
	if err != nil {
		slog.Error("Failed to add task to schedule", slog.Int("schedule_id", int(req.ScheduleID)), slog.Any("error", err))
		c.JSON(http.StatusBadRequest, domain.NewErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, domain.NewSuccessResponse(task))
}

// Update updates a task's status
// @Summary      Update task status
// @Description  Mark a care activity as completed or not_completed. If not_completed, a reason is required.
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        taskId  path      int                      true  "Task ID"
// @Param        body    body      domain.UpdateTaskRequest  true  "Status and optional reason"
// @Success      200     {object}  domain.Response{data=domain.Task}
// @Failure      400     {object}  domain.ErrorResponse
// @Router       /tasks/{taskId}/update [post]
func (h *TaskController) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("taskId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.NewErrorResponse("invalid task ID"))
		return
	}

	var req domain.UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.NewErrorResponse(err.Error()))
		return
	}

	task, err := h.usecase.UpdateStatus(c.Request.Context(), int32(id), req)
	if err != nil {
		slog.Error("Failed to update task", slog.Int("task_id", id), slog.Any("error", err))
		c.JSON(http.StatusBadRequest, domain.NewErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, domain.NewSuccessResponse(task))
}
