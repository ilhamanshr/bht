# Mini EVV Logger - Backend

This is the Go-based backend REST API for the Mini EVV Logger application, providing robust schedule management, task tracking, and geolocation-verified clock-ins/outs.

## Tech Stack
- **Language:** Go 1.23+
- **Framework:** Gin (HTTP Web Framework)
- **Database:** PostgreSQL 16
- **Database Migrations:** golang-migrate
- **SQL Generation:** sqlc
- **Logging:** `log/slog` (Structured JSON Logging)

## Architecture & Key Decisions
- **Clean Architecture**: The codebase is strictly separated into `domain`, `repository`, `usecase`, and `controller` layers. This decoupling ensures the business logic is independent of the database and web framework, making it highly testable.
- **sqlc over ORM**: We use `sqlc` to generate type-safe Go code directly from raw SQL files. This provides the performance of raw SQL with compile-time safety, avoiding the silent reflection overheads commonly found in Go ORMs (like GORM).
- **Structured Logging**: We migrated from the standard `log` package to `log/slog` with a `JSONHandler`. This means all application logs, including unhandled panics, are emitted as structured JSON, making them ready for ingestion by modern observability platforms (e.g., Datadog).
- **Timezone Header**: The backend dynamically accepts an `X-Timezone` header to scope "Today's Schedules" accurately based on the caregiver's physical location, rather than the server's UTC time.

## Local Setup Instructions

### Prerequisites
- Go 1.23+
- PostgreSQL 16+ (or docker compose)
- `golang-migrate` CLI
- `sqlc` CLI (optional, for regenerating queries)

### 1. Database Setup
Ensure PostgreSQL is running. You can use the provided `docker compose` at the root of the project to spin up a local DB.

```bash
# Run migrations to create tables
make migrate-up

# (Optional) Seed the database with sample caregivers, schedules, and tasks
make seed
```

### 2. Environment Variables
Create a `.env` file in the `backend` directory (copy from `.env.example` if available) or rely on defaults:
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=bht_test
SERVER_PORT=8080
```

### 3. Run the Server
```bash
go run cmd/server/main.go
# Server starts on http://localhost:8080
```
- The JSON API is served at `http://localhost:8080/api`.
- 📖 **Swagger API Docs** are available at: `http://localhost:8080/docs/index.html`

### 4. Running Tests
The backend has 100% statement test coverage across the usecase and controller layers.
```bash
go test -v -cover ./...
```

## Assumptions
- Environment variables will default to local Docker credentials if `.env` is missing.
- Geolocation coordinates (Latitude/Longitude) are passed as floats and stored securely in the database to verify visits.

## Future Improvements
- Integrate standard OpenTelemetry tracing across the application.
- Implement JWT authentication middleware to secure the endpoints.
- Add Redis caching for frequently accessed dashboard stats.
