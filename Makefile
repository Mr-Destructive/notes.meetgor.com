.PHONY: help init build run test test-verbose test-coverage test-unit clean deploy lint sqlc config

help:
	@echo "Blog CMS - Makefile Commands"
	@echo ""
	@echo "Development:"
	@echo "  make init           - Initialize database"
	@echo "  make build          - Build Go binary"
	@echo "  make run            - Run local server (default port 8080)"
	@echo "  make test           - Run all tests"
	@echo "  make test-verbose   - Run tests with verbose output"
	@echo "  make test-coverage  - Run tests with coverage report"
	@echo "  make test-unit      - Run only unit tests"
	@echo "  make lint           - Run Go linter (if available)"
	@echo "  make clean          - Remove build artifacts"
	@echo ""
	@echo "Database:"
	@echo "  make sqlc           - Generate sqlc code (requires sqlc installed)"
	@echo ""
	@echo "Deployment:"
	@echo "  make deploy         - Deploy to Netlify"
	@echo ""
	@echo "Configuration:"
	@echo "  make config         - Show environment setup instructions"

init:
	@echo "Initializing database..."
	@go run cmd/cms/main.go

build:
	@echo "Building CMS binary..."
	@go build -o cms ./cmd/functions/main.go
	@echo "✓ Built: ./cms (12MB)"

run: build
	@echo "Starting CMS server on http://localhost:8080"
	@./cms

test:
	@echo "Running all tests..."
	@go test -race -timeout 30s ./...

test-verbose:
	@echo "Running tests with verbose output..."
	@go test -v -race -timeout 30s ./...

test-coverage:
	@echo "Running tests with coverage..."
	@go test -race -timeout 30s -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

test-unit:
	@echo "Running unit tests only..."
	@go test -short -race ./...

lint:
	@echo "Running linter..."
	@which golangci-lint > /dev/null || (echo "Install: curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin" && exit 1)
	@golangci-lint run ./internal/...

sqlc:
	@echo "Generating sqlc code..."
	@which sqlc > /dev/null || (echo "Install sqlc: brew install sqlc  (or see https://github.com/sqlc-dev/sqlc/blob/main/docs/overview/install.md)" && exit 1)
	@sqlc generate

clean:
	@echo "Cleaning build artifacts..."
	@rm -f cms
	@rm -f blog.db
	@go clean

deploy:
	@echo "Deploying to Netlify..."
	@netlify deploy --prod

config:
	@echo ""
	@echo "Environment Configuration:"
	@echo ""
	@echo "  cp .env.example .env"
	@echo ""
	@echo "Then edit .env with:"
	@echo ""
	@echo "  DATABASE_URL=file:./blog.db"
	@echo "  ADMIN_PASSWORD=your-password"
	@echo "  JWT_SECRET=your-secret-key"
	@echo "  ENV=development"
	@echo ""
	@echo "For production (Turso):"
	@echo ""
	@echo "  DATABASE_URL=libsql://your-db-org.turso.io?authToken=token"
	@echo "  ADMIN_PASSWORD=your-password"
	@echo "  JWT_SECRET=your-secret-key"
	@echo "  ENV=production"
	@echo ""

.DEFAULT_GOAL := help
