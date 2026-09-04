# Flood Response API

A REST API for managing flood response, built with Go and PostGIS for geospatial capabilities.

## Tech Stack

- **Language:** Go 1.26
- **Framework:** [Gin](https://github.com/gin-gonic/gin)
- **Database:** PostgreSQL 16 with [PostGIS](https://postgis.net/) extension
- **Testing:** Go testing + [Testify](https://github.com/stretchr/testify)
- **Containerization:** Docker + Docker Compose

## Requirements

- Docker Desktop
- Docker Compose

## Run locally

Start the API, PostgreSQL, and Air live reload:

```bash
docker compose up --build
```

The API is available at:

```text
http://localhost:8080
```

## Check application health

After Docker Compose starts, confirm the API is running:

```bash
curl http://localhost:8080/health
```

Expected response:

```json
{
  "status": "ok"
}
```

Stop the development stack:

```bash
docker compose down
```

To also remove local database data:

```bash
docker compose down -v
```

## Development workflow

The `app` container runs [Air](https://github.com/air-verse/air), which watches Go files and restarts the API after changes.

You do not need Go installed on your machine. Docker runs the API, tests, PostgreSQL, and development tools.

## Run tests

Run all tests inside the application container:

```bash
docker compose exec app go test ./...
```

Run tests with coverage:

```bash
docker compose exec app go test -cover ./...
```

Run one package:

```bash
docker compose exec app go test ./internal/app/affectedarea
```

## Environment configuration

Docker Compose provides development environment variables automatically.

| Variable | Description | Default |
| --- | --- | --- |
| `APP_ENV` | Application environment | `development` |
| `PORT` | HTTP server port | `8080` |
| `DATABASE_URL` | PostgreSQL connection URL | Set by Compose |
| `DATABASE_MAX_OPEN_CONNS` | Maximum open database connections | `10` |
| `DATABASE_MAX_IDLE_CONNS` | Maximum idle database connections | `5` |
| `DATABASE_CONN_MAX_LIFETIME` | Maximum connection lifetime | `30m` |
| `DATABASE_CONN_MAX_IDLE_TIME` | Maximum idle connection lifetime | `5m` |
| `DATABASE_CONNECT_TIMEOUT` | Startup connection timeout | `5s` |

See `.env-example` for the full development configuration.
