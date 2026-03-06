package routes

import (
	"bht-test/internal/api/controllers"
	"bht-test/internal/domain"
	"bht-test/internal/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// NewRouter creates and configures the Gin router
func NewRouter(sh *controllers.ScheduleController, th *controllers.TaskController) *gin.Engine {
	r := gin.New()

	// Global middleware
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())

	// CORS configuration
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, domain.NewSuccessResponse(gin.H{"status": "ok"}))
	})

	// Swagger UI
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API routes
	api := r.Group("/api")
	{
		// Schedule routes
		schedules := api.Group("/schedules")
		{
			schedules.POST("", sh.Create)
			schedules.GET("", sh.List)
			schedules.GET("/today", sh.GetToday)
			schedules.GET("/stats", sh.GetStats)
			schedules.GET("/:id", sh.GetByID)
			schedules.POST("/:id/clock-in", sh.ClockIn)
			schedules.POST("/:id/clock-out", sh.ClockOut)
			schedules.POST("/:id/tasks", th.Create)
		}

		// Task routes
		tasks := api.Group("/tasks")
		{
			tasks.POST("", th.Create)
			tasks.POST("/:taskId/update", th.Update)
		}
	}

	return r
}
