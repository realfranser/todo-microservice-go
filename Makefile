# Constants
SHELL := /bin/bash
SETUP_SCRIPT := ./scripts/setup.sh
ROOT_PATH :=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))


## Postgres
POSTGRES_USER := user
POSTGRES_PASSWORD := password
POSTGRES_PORT := 5432

## DB
DATABASE_MIGRATIONS_PATH := ./db/migrations/ 

install-tools:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.14.1 \
	go install github.com/kyleconroy/sqlc/cmd/sqlc@v1.6.0 \
	go install github.com/maxbrunsfeld/counterfeiter/v6@v6.3.0

setup:
	@bash $(SETUP_SCRIPT)

run:
	go run ./cmd/repository/main.go

migrate-tables-up:
	migrate -path $(DB_MIGRATIONS_PATH) -database "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:$(POSTGRES_PORT)/dbname?sslmode=disable" up

migrate-tables-down:
	migrate -path $(DB_MIGRATIONS_PATH) -database "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:$(POSTGRES_PORT)/dbname?sslmode=disable" down

docker-kill:
	docker kill todo_postgres_db

generate-sqlc:
	cd $(ROOT_PATH)/internal/postgresql \
	sqlc generate

generate-counterfeiter:
	cd $(ROOT_PATH)/internal/envvar/
	go generate counterfeiter -o envvartesting/provider.gen.go . Provider

generate: generate-sqlc generate-counterfeiter
