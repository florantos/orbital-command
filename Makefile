ifneq (,$(wildcard ./.env))
  include .env
  export
endif

.PHONY: help dev down clean test lint

help:
	@echo "Available commands:"
	@echo "  make dev     Start the full stack"
	@echo "  make down    Stop containers, keep data"
	@echo "  make clean   Stop containers, wipe data"
	@echo "  make test    Run backend tests"
	@echo "  make lint    Run backend linter"

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