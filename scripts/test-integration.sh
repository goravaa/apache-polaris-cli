#!/usr/bin/env bash
set -e

DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null 2>&1 && pwd)"
ROOT_DIR="$(dirname "$DIR")"
COMPOSE_DIR="$ROOT_DIR/dev/docker/apache-polaris"

# Generate .env file
echo "Generating .env file..."
cat <<ENV > "$COMPOSE_DIR/.env"
ROOT_CLIENT_ID=root
ROOT_CLIENT_SECRET=secret
ENV

echo "Starting docker-compose environment..."
cd "$COMPOSE_DIR"
docker compose up -d

echo "Waiting for Polaris to be healthy..."
sleep 15 # give it a moment to boot
until curl -s -o /dev/null -w "%{http_code}" http://localhost:8181/api/catalog/v1/config | grep -qE '^(200|401)'; do
  echo "Waiting..."
  sleep 2
done
echo "Polaris is healthy!"

echo "Running integration tests..."
cd "$ROOT_DIR"
export POLARIS_TEST_HOST="http://localhost:8181"
export POLARIS_TEST_CLIENT_ID="root"
export POLARIS_TEST_CLIENT_SECRET="secret"

# Capture exit code so we can tear down even if tests fail
TEST_EXIT_CODE=0
go test -v -tags=integration ./test/integration/... || TEST_EXIT_CODE=$?

if [ $TEST_EXIT_CODE -ne 0 ]; then
  echo "Tests failed!"
  cd "$COMPOSE_DIR"
  docker compose logs
fi

echo "Tearing down environment..."
cd "$COMPOSE_DIR"
docker compose down -v

if [ $TEST_EXIT_CODE -ne 0 ]; then
  # fail the script if tests failed, but without exiting bash directly
  false
fi

echo "Integration tests passed!"
