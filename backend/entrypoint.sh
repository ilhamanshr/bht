#!/bin/sh
set -e

echo "Running database migrations..."
migrate -path /app/db/migrations -database "${DATABASE_URL:-postgres://postgres:postgres@postgres:5432/evv_logger?sslmode=disable}" up

echo "Seeding database..."
PGPASSWORD=${DB_PASSWORD:-postgres} psql -h ${DB_HOST:-postgres} -U ${DB_USER:-postgres} -d ${DB_NAME:-evv_logger} -f /app/db/seeds/000001_init_seed.sql 2>/dev/null || true

echo "Starting server..."
exec ./server
