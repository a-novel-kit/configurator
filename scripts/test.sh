#!/bin/bash

TEST_TOOL_PKG="gotest.tools/gotestsum@latest"

# First, we set up a temporary directory to receive the coverage (binary)files...
GOCOVERTMPDIR="$(mktemp -d)"
trap 'rm -rf -- "$GOCOVERTMPDIR"' EXIT

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

