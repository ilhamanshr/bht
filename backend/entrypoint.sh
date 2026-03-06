#!/bin/sh
set -e

DB_URL="postgresql://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSLMODE}"
echo "Running database migrations..."
/usr/local/bin/migrate -path /app/db/migrations -database "${DB_URL}" up

echo "Checking if database needs seeding..."
SEED_COUNT=$(psql "${DB_URL}" -tAc "SELECT count(*) FROM schedules")

if [ "$SEED_COUNT" -eq 0 ]; then
    echo "Seeding database..."
    psql "${DB_URL}" -f /app/db/seeds/000001_init_seed.sql
else
    echo "Database already contains data ($SEED_COUNT schedules), skipping seeding."
fi

echo "Starting server..."
exec ./server
