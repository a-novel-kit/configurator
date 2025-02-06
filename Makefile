# Run tests.
test:
	bash -c "set -m; bash '$(CURDIR)/scripts/test.sh'"

# Check code quality.
lint:
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@latest run

# Reformat code so it passes the code style lint checks.
format:
	go mod tidy
	go fmt ./...
	go run github.com/daixiang0/gci@latest write \
		--skip-generated \
		-s standard -s default \
		-s "prefix(github.com/a-novel-kit)" \
		-s "prefix(github.com/a-novel-kit/configurator)" \
		.
	go run mvdan.cc/gofumpt@latest -l -w .
	go run golang.org/x/tools/cmd/goimports@latest -w -local github.com/a-novel-kit .
