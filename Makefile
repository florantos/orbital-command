ifneq (,$(wildcard ./.env))
  include .env
  export
endif

.PHONY: help db-setup db-setup-test db-truncate db-truncate-test dev build down clean test test-filter lint

help:
	@echo "Available commands:"
	@echo "  make db-setup           Create database schema"
	@echo "  make db-setup-test      Create test database schema"
	@echo "  make db-truncate        Wipe all data from dev database" 
	@echo "  make db-truncate-test   Wipe all data from test database"	
	@echo "  make dev                Start the full stack"
	@echo "  make build              Start the full stack, rebuild images"
	@echo "  make down               Stop containers, keep data"
	@echo "  make clean              Stop containers, wipe data"
	@echo "  make test               Run backend tests"
	@echo "  make test-filter        Run specific tests e.g. make test-filter FILTER=TestCrewRepo"
	@echo "  make lint               Run backend linter"

db-setup:
	docker exec -i orbital-command-postgres-1 psql -U $(POSTGRES_USER) -d $(POSTGRES_DB) < backend/db/schema.sql

db-setup-test:
	docker exec -i orbital-command-postgres-1 psql -U $(POSTGRES_USER) -d $(POSTGRES_DB)_test < backend/db/schema.sql

db-truncate:
	docker exec -i orbital-command-postgres-1 psql -U $(POSTGRES_USER) -d $(POSTGRES_DB) -c "TRUNCATE TABLE crew_capabilities, crew, audit_events, modules RESTART IDENTITY CASCADE;"

db-truncate-test:
	docker exec -i orbital-command-postgres-1 psql -U $(POSTGRES_USER) -d $(POSTGRES_DB)_test -c "TRUNCATE TABLE crew_capabilities, crew, audit_events, modules RESTART IDENTITY CASCADE;"

dev:
	docker compose up

build: 
	docker compose up --build

down:
	docker compose down

clean:
	docker compose down -v

test:
	cd backend && go test ./...

test-filter:
	cd backend && go test ./... -run $(FILTER) -v

lint:
	cd backend && golangci-lint run ./...