#!/bin/bash

COMPOSE_FILE="docker-compose.test.yaml"
TEST_TOOL_PKG="gotest.tools/gotestsum@latest"

# First, we set up a temporary directory to receive the coverage (binary)files...
GOCOVERTMPDIR="$(mktemp -d)"
trap 'rm -rf -- "$GOCOVERTMPDIR"' EXIT

# Ensure containers are properly shut down when the program exits abnormally.
int_handler()
{
    docker compose -f ${COMPOSE_FILE} down
}
trap int_handler INT

# Setup test containers.
docker compose -f ${COMPOSE_FILE} up -d

export PORT=8080
export DSN="postgres://test:test@localhost:5432/test?sslmode=disable"

# Clear old coverage files.
if [ -d "$GOCOVERTMPDIR" ]; then rm -Rf $GOCOVERTMPDIR; fi
mkdir $GOCOVERTMPDIR

# Execute tests.
go run ${TEST_TOOL_PKG} --format pkgname -- \
  -cover -v -count=1 \
  $(go list ./... | grep -v /mocks) \
  -args -test.gocoverdir=$GOCOVERTMPDIR

# Collect test coverage.
go tool covdata textfmt -i="$GOCOVERTMPDIR" -o=cover.out
go tool cover -html=cover.out -o=cover.html
go tool cover -func=cover.out -o=cover.out

# Normal execution: containers are shut down.
docker compose -f ${COMPOSE_FILE} down
