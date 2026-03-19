.PHONY: dev down clean test lint-backend lint-frontend lint

ifneq (,$(wildcard ./.env))
  include .env
  export
endif

dev:
	docker compose up

down:
	docker compose down

clean:
	docker compose down -v

test-backend:
	@if [ -f .env ]; then source .env; fi && cd backend && go test ./...

lint-backend:
	cd backend && golangci-lint run ./...

lint-frontend:
	cd frontend && pnpm lint

lint: lint-backend lint-frontend