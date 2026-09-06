.PHONY: migrate-up migrate-down migrate-status migrate-create migrate-reset db-up sqlc-generate test verify

# Load .env so MIGRATIONS_DIR etc. are available if needed locally
ifneq (,$(wildcard .env))
	include .env
	export
endif

## --- Application ---
up: ## Start the application and database
	docker compose up -d

down: ## Stop the application and database
	docker compose down

app-up: ## Start the application with Air hot reload
	docker compose up -d app 

app-down: ## Stop the application and database
	docker compose down

## --- Debugging ---

debug-up: ## Start the application with Delve debugging enabled
	docker compose --profile debug up app-debug

debug-down: ## Stop the debug application and database
	docker compose rm -sf app-debug

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

## --- SQLC (via the sqlc/sqlc image — no local install needed) ---

sqlc-generate: ## Regenerate sqlc code from queries/ and migrations/
	docker run --rm -v "$(CURDIR):/src" -w /src sqlc/sqlc generate

## --- Testing ---

test: ## Run all Go tests inside the running app container
	docker compose exec -T app go test ./... -v -count=1

verify: ## Run format check, build, vet and tests (verbose) inside the running app container
	docker compose exec -T app sh -c "set -e; echo '== gofmt =='; files=\$$(gofmt -l internal cmd); if [ -n \"\$$files\" ]; then echo 'Unformatted files:'; echo \"\$$files\"; exit 1; fi; echo '== build =='; go build -v ./...; echo '== vet =='; go vet -v ./...; echo '== test =='; go test -v ./... -count=1"