.PHONY: build test lint format check clean tidy migrate migrate-down seed

# Go build settings
GOFLAGS ?= -trimpath
LDFLAGS ?= -s -w

build:
	go build $(GOFLAGS) -ldflags "$(LDFLAGS)" ./...

test:
	go test -v -race -count=1 -coverprofile=coverage.out ./...

lint:
	golangci-lint run ./...

format:
	gofumpt -w .
	goimports -w .

tidy:
	go mod tidy

check: format tidy lint test build
	@echo "All checks passed."

migrate:
	go run ./cmd/migrate -direction up

migrate-down:
	go run ./cmd/migrate -direction down

seed:
	go run ./cmd/seed

clean:
	go clean -cache -testcache
	rm -f coverage.out
