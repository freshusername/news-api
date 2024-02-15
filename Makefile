# Makefile

# Variables
BIN_DIR=bin
BINARY=$(BIN_DIR)/app
DB_MIGRATIONS_DIR=database/migrations

DB_USER=postgres
DB_PASSWORD=admin
GOOSE_DRIVER=postgres
GOOSE_DBSTRING=postgresql://$(DB_USER):$(DB_PASSWORD)@127.0.0.1:5432/news?sslmode=disable

.PHONY: build run test migrate clean docker-up docker-down docker-restart

goose-install:
	go get -u github.com/pressly/goose/v3/cmd/goose

build:
	echo "Building the application..."
	go build -o $(BINARY) ./api

run: build
	echo "Running the application..."
	./$(BINARY)

test:
	echo "Running tests..."
	go test -v ./... -count=1
	
coverage:
	echo "Running tests with coverage..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	echo "Coverage report generated: coverage.html"

migrate-up:
	echo "Running database migrations up..."
	goose -dir $(DB_MIGRATIONS_DIR) up

migrate-down:
	echo "Reverting database migrations..."
	goose -dir $(DB_MIGRATIONS_DIR) down

clean:
	echo "Cleaning up..."
	rm -rf $(BIN_DIR)

docker-up:
	echo "Starting Docker containers..."
	docker-compose up -d

docker-down:
	echo "Stopping Docker containers..."
	docker-compose down

docker-restart:
	echo "Restarting Docker containers..."
	docker-compose down && docker-compose up -d