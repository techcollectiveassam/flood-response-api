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

## Quick start

Start the API, PostgreSQL, and Air live reload — either line works:

```bash
make up
# or
docker compose up -d
```

The API is available at:

```text
http://localhost:8080
```

## Check application health

After the stack starts, confirm the API is running:

```bash
curl http://localhost:8080/health
```

Expected response:

```json
{
  "status": "ok"
}
```

Stop the development stack — either line works:

```bash
make down
# or
docker compose down
```

To also remove local database data (destructive):

```bash
docker compose down -v
```

## Make targets

Every command below can be run two ways — the `make` target or the plain `docker compose` command. They are equivalent; `make` just saves you typing. Both lines are shown for each command.

### Application

- **Start the application and database** — `make up` *or* `docker compose up -d`
- **Stop the application and database** — `make down` *or* `docker compose down`
- **Start only the application (Air hot reload)** — `make app-up` *or* `docker compose up -d app`
- **Stop the application and database** — `make app-down` *or* `docker compose down`

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