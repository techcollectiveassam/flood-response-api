# Flood Response API

A REST API for managing flood response, built with Go and PostGIS for geospatial capabilities.

## Tech Stack

- **Language:** Go 1.26
- **Framework:** [Gin](https://github.com/gin-gonic/gin)
- **Database:** PostgreSQL 16 with [PostGIS](https://postgis.net/) extension
- **Database access:** [pgx](https://github.com/jackc/pgx) v5 (`pgxpool`) + [sqlc](https://sqlc.dev) for type-safe queries
- **Migrations:** [Goose](https://github.com/pressly/goose)
- **Testing:** Go testing + [Testify](https://github.com/stretchr/testify)
- **Containerization:** Docker + Docker Compose

## Requirements

- Docker Desktop
- Docker Compose
- `make`

You do not need Go installed on your machine — Docker runs the API, tests, PostgreSQL, and development tools.

## Getting started

Follow these steps in order to run the API locally and start using it.

### 1. Configure your environment

Docker Compose reads local settings from `.env`, which is gitignored. Create it from the template:

```bash
cp .env-example .env
```

The defaults work for local development; edit the database credentials or ports in `.env` only if you need to.

### 2. Start the stack

```bash
make up
# or
docker compose up -d
```

This starts two services:

- `db` — PostgreSQL 16 with PostGIS
- `app` — the API `http://localhost:8080`, watched by [Air](https://github.com/air-verse/air) (hot reloads on Go file changes)

`make up` does **not** apply database migrations — that's step 3.

### 3. Apply database migrations

On a fresh database the schema is empty, so apply all pending migrations:

```bash
make migrate-up
```

`make migrate-up` starts the `db` service if it isn't already running, then applies every pending migration with [Goose](https://github.com/pressly/goose). Check what has and hasn't been applied with `make migrate-status`.

### 4. Verify the API is healthy

```bash
curl http://localhost:8080/health
```

Expected response:

```json
{
  "status": "ok"
}
```

### 5. Explore the API with Bruno

The repository tracks a versioned [Bruno](https://www.usebruno.com/) collection covering every endpoint — no external API client setup needed.

1. Generate your local Bruno environment (gitignored):

   ```bash
   make bruno-local-env
   ```

2. Install the [Bruno desktop client](https://www.usebruno.com/download) (if you haven't already) and open it.
3. **File → Open Collection** and select the `tools/bruno/collections/flood-response-api` folder (the one containing `bruno.json`).
4. Select the **local** environment (bottom-right of the Bruno window).
5. Run the **Disaster** folder first — it captures the reference `disasterId` (creating a disaster if the database is empty). Later requests use it automatically.
6. Run a single request, a folder, or the whole collection.

### 6. Browse the API reference

A Redocly page is served by the app:

```text
http://localhost:8080/docs
```

It is generated from `docs/openapi.yaml` — regenerate with `make docs-build` after editing the spec.

### 7. Stop the stack

```bash
make down
```

To also delete local database data (destructive):

```bash
docker compose down -v
```

## Reference: make targets

Every command below can be run two ways — the `make` target or the plain `docker compose` command. They are equivalent; `make` just saves you typing. Both lines are shown for each command.

### Application

- **Start the application and database** — `make up` *or* `docker compose up -d`
- **Stop the application and database** — `make down` *or* `docker compose down`
- **Start only the application (Air hot reload)** — `make app-up` *or* `docker compose up -d app`
- **Stop the application** — `make app-down` *or* `docker compose down`

The `app` container runs [Air](https://github.com/air-verse/air), which watches Go files and restarts the API after changes.

### Testing

- **Run all Go tests (verbose)** — `make test` *or* `docker compose exec -T app go test ./... -v -count=1`

- **Format check, build, vet, and tests in one go (verbose)** — either of:

  ```bash
  make verify
  # or
  docker compose exec -T app sh -c "test -z \"\$(gofmt -l internal cmd)\" && go build ./... && go vet ./... && go test ./... -count=1"
  ```

Both require the app container to be running (`make up` / `docker compose up -d` first). `verify` fails if any file in `internal/` or `cmd/` is not `gofmt`-formatted.

- **Run the Bruno collection E2E tests** — `make e2e`

  `e2e` boots the API and a **fresh** PostgreSQL database in an isolated `flood-e2e` compose project (own network and volumes), applies migrations, then runs the full Bruno collection in `tools/bruno/collections/flood-response-api` via the official `usebruno/cli` image against `http://app:8080` over the container network. The **Disaster** folder runs first (seq 1) and establishes the reference disaster id — `create/valid` creates one (capturing its id) if the database is empty, otherwise `list/all` captures the first existing id — so no external seeding or manual IDs are needed. It publishes **no host ports** (see `docker-compose.e2e.yml`), so it runs side-by-side with `make up` without touching your running dev stack, and it tears the e2e containers and volumes down afterward (even on failure). Exit code is non-zero if any test fails.

  Pass extra args to `bru run` with `BRUNO_ARGS`, e.g. `make e2e BRUNO_ARGS=Disaster`.

- **Run the Bruno collection against the running dev stack (fast inner loop)** — `make test-bruno`

  Runs the collection through the same `usebruno/cli` container on the dev network. Requires the dev stack to be up (`make up`). Unlike `e2e` it reuses existing dev data — the **Disaster** folder runs first and captures IDs from existing data (creating a new disaster only when the API is empty) — so tests that assume fresh IDs may fail; use `e2e` for an authoritative run. JUnit output is written to `test-results/bruno-junit.xml`.

### Bruno API collection

The API collection is tracked in this repository and works with [Bruno](https://www.usebruno.com/) (GUI and CLI) — everything lives in `tools/bruno/collections/flood-response-api/` and is versioned alongside the code, with no external API client needed to explore the endpoints.

**Using the collection in the local Bruno GUI:**

1. Install the [Bruno desktop client](https://www.usebruno.com/download) and open it.
2. **File → Open Collection**, then select the `tools/bruno/collections/flood-response-api` folder (the one containing `bruno.json`).
3. Create your local environment (gitignored so you can't accidentally commit it): `make bruno-local-env` copies `environments/local.bru.example` → `environments/local.bru` with `baseUrl: http://localhost:8080`.
4. Start the dev stack with `make up`, then select the **local** environment (bottom-right of the Bruno window).
5. Run a single request, a folder, or the whole collection. Run the **Disaster** folder first — it captures the reference `disasterId` (creating a disaster if the database is empty), and later requests use it automatically.

**Environments:**

| Environment | `baseUrl` | Used for |
| --- | --- | --- |
| `ci` | `http://app:8080` | CLI/CI runs (`make e2e`, `make test-bruno`) |
| `local` | `http://localhost:8080` | Local GUI runs (gitignored; generate with `make bruno-local-env`) |
| `staging` | *(empty — fill in the GUI)* | GUI runs against the staging deployment |
| `local.bru.example` | `http://localhost:8080` | Committed template for `local.bru` |

### Migrations (Goose)

Migrations live in `internal/pkg/database/migrations`. Goose runs via the `migrate` service — no local Goose install needed; the `db` service is started automatically.

- **Apply all pending migrations** — `make migrate-up` *or* `docker compose run --rm migrate up`
- **Roll back the last applied migration** — `make migrate-down` *or* `docker compose run --rm migrate down`
- **Show applied / pending migrations** — `make migrate-status` *or* `docker compose run --rm migrate status`
- **Create a new migration file** — `make migrate-create name=create_incidents` *or* `docker compose run --rm migrate create create_incidents sql`
- **Roll back ALL migrations (destructive — dev only)** — `make migrate-reset` *or* `docker compose run --rm migrate reset`

### sqlc

Query definitions live in `internal/pkg/database/queries`; generated code in `internal/pkg/database/sqlcgen` (regenerated, never hand-edited).

- **Regenerate sqlc code** — `make sqlc-generate` *or* `docker run --rm -v "$PWD:/src" -w /src sqlc/sqlc generate`

### Debugging

- **Start the application with Delve debugging enabled** — `make debug-up` *or* `docker compose --profile debug up app-debug`
- **Stop the debug application** — `make debug-down` *or* `docker compose rm -sf app-debug`

## Environment configuration

Docker Compose provides development environment variables automatically from `.env`.

| Variable | Description | Default |
| --- | --- | --- |
| `APP_ENV` | Application environment | `development` |
| `APP_PORT` | HTTP server port | `8080` |
| `DATABASE_URL` | PostgreSQL connection URL | Set by Compose |
| `DATABASE_MAX_OPEN_CONNS` | Maximum open database connections | `10` |
| `DATABASE_MAX_IDLE_CONNS` | Maximum idle database connections | `5` |
| `DATABASE_CONN_MAX_LIFETIME` | Maximum connection lifetime | `30m` |
| `DATABASE_CONN_MAX_IDLE_TIME` | Maximum idle connection lifetime | `5m` |
| `DATABASE_CONNECT_TIMEOUT` | Startup connection timeout | `5s` |
| `POSTGRES_DB` | PostgreSQL database name | `flood_response` |
| `POSTGRES_USER` | PostgreSQL user | `techcollectiveassam` |
| `POSTGRES_PASSWORD` | PostgreSQL password | `techcollectiveassam` |
| `POSTGRES_PORT` | Exposed PostgreSQL port | `5432` |
| `MIGRATIONS_DIR` | Migrations directory | `internal/pkg/database/migrations` |

See `.env-example` for the full development configuration.