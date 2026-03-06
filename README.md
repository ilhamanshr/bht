# Mini EVV Logger – Caregiver Shift Tracker

A full-stack web application designed for **Electronic Visit Verification (EVV)** compliance. Caregivers can view their schedules, clock in/out with real-time geolocation, and track daily care activity progress seamlessly.

## Project Structure
- `frontend/` - React 19 SPA powered by Vite and TypeScript. [Frontend README](./frontend/README.md)
- `backend/` - Go REST API powered by Gin and PostgreSQL. [Backend README](./backend/README.md)

## Overall Tech Stack

| Layer | Technology |
|-------|-----------|
| **Frontend** | React 19 + TypeScript, Vite |
| **Styling** | Custom CSS Variables (Dark Theme) |
| **Backend** | Go (Gin framework) |
| **Database** | PostgreSQL 16 |
| **Architecture** | Clean Architecture |
| **DB Migrations** | golang-migrate |
| **SQL Queries** | sqlc (type-safe SQL generation) |
| **Logging** | `log/slog` (Structured JSON Logging) |
| **Containerization** | Docker Compose |

## Architecture & Key Decisions

1. **Clean Architecture (Backend)** — Strict separation of concerns between `domain`, `usecase`, `repository`, and HTTP `controller`. The dependency rule points inward, ensuring business logic (`usecase`) never depends on external frameworks.
2. **sqlc over ORMs** — Instead of relying on reflection-heavy ORMs (like GORM), we define pure `.sql` queries. `sqlc` generates type-safe Go code for data access, marrying performance with compile-time safety.
3. **Structured Logging (`slog`)** — The backend strictly uses the native Go 1.21+ `slog.JSONHandler` to emit machine-readable log context suitable for log aggregation platforms (e.g., Datadog).
4. **Geolocation API** — EVV compliance mandates geolocation on clock-in and clock-out. The frontend integrates directly with the native browser `navigator.geolocation` API.
5. **Timezone Accuracy** — To display an accurate "Today" schedule, the frontend passes the caregiver's local `X-Timezone` via a React Axios interceptor, so the backend calculates midnight relative to the physical user.

## Geolocation & Compliance Fallbacks

- **Primary Method**: Browser Geolocation API.
- **Fallback Handling**: If GPS is unavailable (permission denied or device failure), caregivers are prompted to proceed anyway. The entry will be accepted but permanently flagged as **"Unverified" ⚠️** for administrative audit. This ensures caregivers aren't completely blocked in signal-dead zones.

## Live Environment

The application is deployed and available at the following locations:

- **Frontend**: [https://bht-frontend-production.up.railway.app/](https://bht-frontend-production.up.railway.app/)
- **Backend API**: [https://bht-backend-production.up.railway.app/](https://bht-backend-production.up.railway.app/)
- **API Documentation (Swagger)**: [https://bht-backend-production.up.railway.app/docs/index.html](https://bht-backend-production.up.railway.app/docs/index.html)

---

## Quick Start (Docker - Recommended)

You can spin up the entire application using Docker Compose with zero manual configuration.

```bash
# Start PostgreSQL, Backend, and Frontend (builds images)
docker compose up --build
```
- The **Frontend** will be available at: `http://localhost:5173`
- The **Backend API** will be available at: `http://localhost:8080/api`

*(Note: Ensure you seed the test data if running for the first time by running `make seed` from the `backend/` directory).*

## Local Development Setup

To run the application manually on your host machine with minimal setup:

### 1. Start the Database
```bash
# Start PostgreSQL via Docker
docker compose up postgres -d

# Run migrations to build the tables
cd backend
make migrate-up

# Seed sample data (schedules & tasks)
make seed
```

### 2. Start the Backend (API)
```bash
cd backend
go run cmd/server/main.go
# Server starts on http://localhost:8080
```

### 3. Start the Frontend (Client)
In a new terminal wrapper:
```bash
cd frontend
npm install
npm run dev
# App opens at http://localhost:5173
```

## API Endpoints Overview

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/schedules` | List all schedules |
| `GET` | `/api/schedules/today` | Today's schedules (scoped to `X-Timezone`) |
| `GET` | `/api/schedules/stats` | Dashboard statistics |
| `GET` | `/api/schedules/:id` | Schedule details and associated tasks |
| `POST` | `/api/schedules` | Create a new schedule |
| `POST` | `/api/schedules/:id/clock-in` | Clock in (requires lat/lng) |
| `POST` | `/api/schedules/:id/clock-out` | Clock out (requires lat/lng) |
| `POST` | `/api/tasks` | Add a new task to a schedule |
| `POST` | `/api/tasks/:taskId/update` | Update task status (completed/not_completed) |

### 📖 Swagger API Documentation
Interactive OpenAPI (Swagger) documentation is built into the backend. 
Once the backend is running, you can view and test all API endpoints directly from your browser by navigating to:
👉 **[http://localhost:8080/docs/index.html](http://localhost:8080/docs/index.html)**

## Core Assumptions
1. **Single Caregiver Mode**: The application runs assuming a single, authenticated caregiver context for demonstration purposes. There is no active login screen or JSON Web Token (JWT) handling in place.
2. **Visit State Machine**: Schedules linearly transition from `upcoming` → `in_progress` → `completed`.
3. **Task Rules**: Tasks can only be checked off if the schedule is currently `in_progress`.
4. **Local DB Credentials**: The `.env` configurations fall back gracefully to the default `postgres` Docker credentials if they are absent.

## Optional Improvements (Given More Time)
If I had additional time to iterate on this architecture, I would implement:
- **Authentication**: JWT-based login for caregivers, complete with Role-Based Access Control (RBAC) separating administrative users from caregivers.
- **Offline Mode**: A Progressive Web App (PWA) ServiceWorker caching strategy, allowing caregivers in rural areas without reception to complete tasks and sync state when connectivity is restored.
- **Real-time Synchronization**: WebSockets (using socket.io or raw Go WebSockets) to live-update the schedule stats dashboard as caregivers clock in/out.
- **Observability**: OpenTelemetry tracing decorators surrounding the sqlc database transactions and Gin HTTP controllers.
