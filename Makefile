.PHONY: migrateup migratedown migrateforce migratefix migrateverify migratereset create-migration setup swagger run-with-docs install-swagger-deps all build test clean generate-graphql generate-grpc install-deps dev lint format test-unit test-e2e test-all deploy-minikube delete-minikube

DATABASE_DSN=postgresql://metabase:ffdd1f0a568f407daaa2e176b5fd5481@localhost:5432/sil_backend_assessment_db?sslmode=disable

migrate_init:
	go run cmd/migrate/main.go init

migrate_create:
	@read -p "Enter migration name: " name; \
	go run cmd/migrate/main.go create $$name

migrate_up:
	go run cmd/migrate/main.go up

migrate_down:
	go run cmd/migrate/main.go down

migrate_reset:
	@echo "Dropping schema_migrations table..."
	go run cmd/migrate/main.go reset

migrate_verify:
	go run cmd/migrate/main.go status

migrate_mark_applied:
	@read -p "Enter migration name: " name; \
	go run cmd/migrate/main.go mark_applied $$name

migrate_all_applied:
	go run cmd/migrate/main.go mark_all_applied

migrate_unlock:
	go run cmd/migrate/main.go unlock

migrate_help:
	go run cmd/migrate/main.go help

# Development setup
dev-setup:
	cp .env.example .env
	@echo "Please update .env file with your configuration"
	go mod download

all: build

build:
	go build -o bin/server cmd/server/main.go

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Run unit tests only
test-unit:
	./scripts/test.sh --unit-only

# Run E2E tests only
test-e2e:
	./scripts/test.sh --e2e-only

# Run all tests (unit + e2e)
test-all:
	./scripts/test.sh

# Run tests with coverage and linting
test-full:
	./scripts/test.sh --no-security

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

# Install dependencies
deps:
	go mod download
	go mod tidy

lint:
	golangci-lint run

fmt:
	go fmt ./...

generate-graphql:
	GOFLAGS=-mod=mod go run github.com/99designs/gqlgen@v0.17.78 generate

# Docker build
docker-build:
	docker build -t sil-backend-assessment .

# Docker run
docker-run:
	docker run -p 8080:8080 --env-file .env sil-backend-assessment

# Docker compose up
docker-up:
	docker-compose up -d

# Docker compose down
docker-down:
	docker-compose down -v

# Deploy to Minikube
deploy-minikube:
	./scripts/deploy-minikube.sh

# Delete from Minikube
delete-minikube:
	./scripts/delete-minikube.sh

# Install testing tools
install-test-tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest

# Security scan
security-scan:
	gosec ./...

# Lint with golangci-lint
lint-full:
	golangci-lint run --timeout=5m

# Format code
format:
	go fmt ./...
	goimports -w .
