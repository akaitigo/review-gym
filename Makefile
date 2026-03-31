.PHONY: build test lint format check quality clean migrate migrate-down

build:
	cd frontend && npm run build
	cd backend && make build

test:
	cd frontend && npm run test
	cd backend && make test

lint:
	cd frontend && npm run lint
	cd backend && make lint

format:
	cd frontend && npm run format
	cd backend && make format

check: lint test build
	@echo "All checks passed."

quality:
	@echo "=== Quality Gate ==="
	@test -f LICENSE || { echo "ERROR: LICENSE missing. Fix: add MIT LICENSE file"; exit 1; }
	@! grep -rn "TODO\|FIXME\|HACK\|console\.log\|println\|print(" frontend/src/ backend/ 2>/dev/null | grep -v "node_modules" | grep -v "_test\.go:.*\`" | grep -v "_test\.go:.*\"" || { echo "ERROR: debug output or TODO found. Fix: remove before ship"; exit 1; }
	@! grep -rn "password=\|secret=\|api_key=\|sk-\|ghp_" frontend/src/ backend/ 2>/dev/null | grep -v '\$${' | grep -v "node_modules" || { echo "ERROR: hardcoded secrets. Fix: use env vars with no default"; exit 1; }
	@test ! -f PRD.md || ! grep -q "\[ \]" PRD.md || { echo "ERROR: unchecked acceptance criteria in PRD.md"; exit 1; }
	@echo "OK: automated quality checks passed"
	@echo "Manual checks required: README quickstart, demo GIF, input validation, ADR >=1"

migrate:
	cd backend && make migrate

migrate-down:
	cd backend && make migrate-down

clean:
	cd frontend && rm -rf dist/ coverage/ node_modules/.cache/
	cd backend && make clean
