# Constants
SHELL := /bin/bash
SETUP_SCRIPT := ./scripts/setup.sh

## Postgres
POSTGRES_USER := user
POSTGRES_PASSWORD := password
POSTGRES_PORT := 5432

## DB
DATABASE_MIGRATIONS_PATH := ./db/migrations/ 

install-migrate:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.14.1

migrate-tables-up:
	migrate -path $(DB_MIGRATIONS_PATH) -database "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:$(POSTGRES_PORT)/dbname?sslmode=disable" up

migrate-tables-down:
	migrate -path $(DB_MIGRATIONS_PATH) -database "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:$(POSTGRES_PORT)/dbname?sslmode=disable" down

setup:
	@bash $(SETUP_SCRIPT)

direnv-allow:
	direnv allow .

docker-kill:
	docker kill todo_postgres_db

run:
	go run ./cmd/repository/main.go

generate-sqlc:
	cd ./internal/postgresql \
	sqlc generate
