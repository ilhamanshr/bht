package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"time"

	"bht-test/internal/api/controllers"
	"bht-test/internal/api/routes"
	"bht-test/internal/config"
	"bht-test/internal/repository/sqlc"
	"bht-test/internal/usecase"

	_ "bht-test/docs" // Swagger generated docs

	_ "github.com/lib/pq" // PostgreSQL driver
)

// @title           Mini EVV Logger API
// @version         1.0
// @description     Caregiver Shift Tracker API for Electronic Visit Verification (EVV) compliance.
// @description     Caregivers can view schedules, clock in/out with geolocation, and track care activity progress.

// @host      localhost:8080
// @BasePath  /api

// @contact.name   Blue Horn Tech
// @contact.email  hr@bluehorntek.com

func main() {
	// Setup structured JSON logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Load configuration
	cfg := config.Load()

	// Connect to PostgreSQL
	dbConn, errDB := sql.Open("postgres", cfg.DatabaseURL())
	dbConn.SetMaxIdleConns(5)
	dbConn.SetConnMaxIdleTime(10 * time.Second)
	dbConn.SetMaxOpenConns(90)
	if errDB != nil {
		slog.Error("Unable to connect to database", "error", errDB)
		os.Exit(1)
	}
	defer dbConn.Close()

	// Verify connection
	if err := dbConn.Ping(); err != nil {
		slog.Error("Failed to ping database", slog.Any("error", err))
		panic(err)
	}
	slog.Info("Database connection successful")

	// Repository layer
	queries := sqlc.New(dbConn)

	// Usecase layer
	scheduleUsecase := usecase.NewScheduleUsecase(queries, dbConn)
	taskUsecase := usecase.NewTaskUsecase(queries)

	// Controller layer
	scheduleController := controllers.NewScheduleController(scheduleUsecase)
	taskController := controllers.NewTaskController(taskUsecase)

	// Setup router
	router := routes.NewRouter(scheduleController, taskController)

	// Start server on all interfaces (required for Railway/Docker)
	addr := fmt.Sprintf("0.0.0.0:%s", cfg.ServerPort)
	slog.Info("Server started", slog.String("addr", addr))
	if err := router.Run(addr); err != nil {
		slog.Error("Failed to start server", slog.Any("error", err))
	}
}
