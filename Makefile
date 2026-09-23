.PHONY: migrate-up migrate-down migrate-status migrate-create migrate-reset db-up sqlc-generate docs-build bruno-local-env test test-bruno verify e2e

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

## --- API Docs (Redoc, via the node:alpine image — no local install needed) ---

docs-build: ## Regenerate the embedded Redoc page (internal/api/docs/static) from docs/openapi.yaml
	docker run --rm -v "$(CURDIR):/src" -w /src node:22-alpine sh -c "node /src/scripts/docs-build.mjs"

## --- Bruno ---

bruno-local-env: ## Generate the Bruno local environment from local.bru.example
	cp tools/bruno/collections/flood-response-api/environments/local.bru.example tools/bruno/collections/flood-response-api/environments/local.bru
	@echo "Created tools/bruno/collections/flood-response-api/environments/local.bru — edit it with your machine's values."

## --- Testing ---

test: ## Run all Go tests inside the running app container
	docker compose exec -T app go test ./... -v -count=1

test-bruno: ## Run Bruno collection tests against the running dev stack (fast inner loop)
	@mkdir -p test-results
	@docker compose up -d app >/dev/null 2>&1
	docker compose run --rm --remove-orphans bruno run -r --env ci --reporter-junit /results/bruno-junit.xml $(BRUNO_ARGS)

e2e: ## Run the Bruno collection E2E tests against a fresh, isolated stack (tears down afterward)
	./scripts/e2e.sh $(BRUNO_ARGS)

verify: ## Run format check, build, vet and tests (verbose) inside the running app container
	docker compose exec -T app sh -c "set -e; echo '== gofmt =='; files=\$$(gofmt -l internal cmd); if [ -n \"\$$files\" ]; then echo 'Unformatted files:'; echo \"\$$files\"; exit 1; fi; echo '== build =='; go build -v ./...; echo '== vet =='; go vet -v ./...; echo '== test =='; go test -v ./... -count=1"