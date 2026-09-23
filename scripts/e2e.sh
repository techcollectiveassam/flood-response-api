#!/usr/bin/env bash
set -euo pipefail

E2E_PROJECT="${E2E_PROJECT:-flood-e2e}"
COMPOSE="docker compose -p ${E2E_PROJECT} -f docker-compose.yml -f docker-compose.e2e.yml"

mkdir -p test-results

teardown() {
  ${COMPOSE} down -v --remove-orphans >/dev/null 2>&1 || true
}
trap teardown EXIT

echo "==> Starting database and applying migrations (project: ${E2E_PROJECT})"
${COMPOSE} run --rm migrate up

echo "==> Starting API service"
${COMPOSE} up -d app

echo "==> Waiting for API to become healthy"
for i in $(seq 1 60); do
  if ${COMPOSE} exec -T app wget -qO- http://127.0.0.1:8080/health 2>/dev/null | grep -q '"status":"ok"'; then
    echo "API is healthy (attempt ${i})"
    break
  fi
  if [ "${i}" -eq 60 ]; then
    echo "ERROR: API did not become healthy in time" >&2
    exit 1
  fi
  sleep 1
done

echo "==> Running Bruno collection E2E tests"
if [ "$#" -gt 0 ]; then
  ${COMPOSE} run --rm bruno run -r --env ci --reporter-junit /results/bruno-junit.xml "$@"
else
  ${COMPOSE} run --rm bruno
fi
