.PHONY: migrate-up migrate-down migrate-status migrate-create migrate-reset db-up

# Load .env so MIGRATIONS_DIR etc. are available if needed locally
ifneq (,$(wildcard .env))
	include .env
	export
endif

## --- Migrations (Goose, via the `migrate` service in docker-compose.yml) ---
## Uses the ghcr.io/pressly/goose image — no local Goose install needed.
## `db` must be healthy first; `docker compose run` will wait for it automatically.

db-up: ## Start just the db service in the background
	docker compose up -d db

migrate-up: db-up ## Apply all pending migrations
	docker compose run --rm migrate up

migrate-down: db-up ## Roll back the last applied migration
	docker compose run --rm migrate down

migrate-status: db-up ## Show which migrations are applied / pending
	docker compose run --rm migrate status

migrate-create: ## Create a new migration file: make migrate-create name=create_incidents
	docker compose run --rm migrate create $(name) sql

migrate-reset: db-up ## Roll back ALL migrations (destructive — dev only)
	docker compose run --rm migrate reset