-include .env
-include .envrc

DB_ADDR         ?= postgres://supervillager:adminpassword@localhost:5433/social?sslmode=disable
MIGRATIONS_PATH ?= ./cmd/migrate/migrations

.PHONY: migrate-up migrate-down migration seed

print-env:
	@echo "DB_ADDR=$(DB_ADDR)"
	@echo "MIGRATIONS_PATH=$(MIGRATIONS_PATH)"

migrate-up:
	@migrate -path="$(MIGRATIONS_PATH)" -database="$(DB_ADDR)" -verbose up

migrate-down:
	@migrate -path="$(MIGRATIONS_PATH)" -database="$(DB_ADDR)" -verbose down 1

migration:
	@if [ -z "$(name)" ]; then echo "Usage: make migration name=<snake_case_name>"; exit 1; fi
	@migrate create -seq -ext sql -dir "$(MIGRATIONS_PATH)" "$(name)"

seed:
	@go run cmd/migrate/migrations/seed/main.go
