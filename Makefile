.PHONY: dev run build test clean sqlc migrate-up migrate-down migrate-create docker-up docker-down swagger help

# Default target
.DEFAULT_GOAL := help

# Variables
BINARY_NAME=blog-api
MAIN_PATH=./cmd/server
MIGRATIONS_PATH=./migrations
DB_URL=postgres://postgres:postgres@localhost:5432/blog?sslmode=disable

## help: Show this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'

## dev: Run with hot reload (air)
dev:
	air

## run: Run the application
run:
	go run $(MAIN_PATH)/main.go

## build: Build the application
build:
	go build -o bin/$(BINARY_NAME) $(MAIN_PATH)/main.go

## test: Run tests
test:
	go test -v ./...

## clean: Clean build files
clean:
	rm -rf bin/ tmp/

## sqlc: Generate SQL code
sqlc:
	sqlc generate

## migrate-up: Run database migrations
migrate-up:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" up

## migrate-down: Rollback database migrations
migrate-down:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" down 1

## migrate-create: Create a new migration (usage: make migrate-create name=migration_name)
migrate-create:
	migrate create -ext sql -dir $(MIGRATIONS_PATH) -seq $(name)

## docker-up: Start development containers
docker-up:
	docker-compose -f docker-compose.dev.yml up -d

## docker-down: Stop development containers
docker-down:
	docker-compose -f docker-compose.dev.yml down

## docker-logs: View container logs
docker-logs:
	docker-compose -f docker-compose.dev.yml logs -f

## swagger: Generate Swagger documentation
swagger:
	swag init -g $(MAIN_PATH)/main.go -o ./docs/swagger

## tidy: Tidy and verify go modules
tidy:
	go mod tidy
	go mod verify

## lint: Run linter
lint:
	golangci-lint run ./...

## install-tools: Install development tools
install-tools:
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	go install github.com/air-verse/air@latest
	go install github.com/swaggo/swag/cmd/swag@latest
