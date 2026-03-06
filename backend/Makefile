.PHONY: migrate-up migrate-down sqlc seed run test build swagger

# Database URL
DB_URL ?= postgres://postgres:postgres@localhost:5433/evv_logger?sslmode=disable

migrate-up:
	migrate -path db/migrations -database "$(DB_URL)" -verbose up

migrate-down:
	migrate -path db/migrations -database "$(DB_URL)" -verbose down

migrate-create:
	migrate create -ext sql -dir db/migrations -seq $(name)

sqlc:
	sqlc generate

swagger:
	swag init -g cmd/server/main.go -o docs --parseDependency --parseInternal

seed:
	psql "$(DB_URL)" -f db/seeds/000001_init_seed.sql

run:
	go run cmd/server/main.go

test:
	go test ./... -v -count=1

build:
	go build -o bin/server cmd/server/main.go

db-setup: migrate-up seed
	@echo "Database setup complete"

db-reset: migrate-down migrate-up seed
	@echo "Database reset complete"
