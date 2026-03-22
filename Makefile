ifneq (,$(wildcard ./.env))
  include .env
  export
endif

.PHONY: help db-setup dev down clean test lint

help:
	@echo "Available commands:"
	@echo "  make db-setup       Create database schema"
	@echo "  make db-setup-test  Create test database schema"
	@echo "  make dev            Start the full stack"
	@echo "  make down           Stop containers, keep data"
	@echo "  make clean          Stop containers, wipe data"
	@echo "  make test           Run backend tests"
	@echo "  make lint           Run backend linter"

db-setup:
	docker exec -i orbital-command-postgres-1 psql -U $(POSTGRES_USER) -d $(POSTGRES_DB) < backend/db/schema.sql

db-setup-test:
	docker exec -i orbital-command-postgres-1 psql -U $(POSTGRES_USER) -d $(POSTGRES_DB)_test < backend/db/schema.sql

dev:
	docker compose up

down:
	docker compose down

clean:
	docker compose down -v

test:
	cd backend && go test ./...

lint:
	cd backend && golangci-lint run ./...