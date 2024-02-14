# Makefile

# Variables
BIN_DIR=bin
BINARY=$(BIN_DIR)/app
DB_MIGRATIONS_DIR=migrations

.PHONY: build run test migrate clean docker-up docker-down

build:
	@echo "Building the application..."
	@go build -o $(BINARY)

run: build
	@echo "Running the application..."
	@./$(BINARY)

test:
	@echo "Running tests..."
	@go test -v ./... -count=1

migrate-up:
	@echo "Running database migrations up..."
	@goose -dir $(DB_MIGRATIONS_DIR) postgres "$(DSN)" up

migrate-down:
	@echo "Reverting database migrations..."
	@goose -dir $(DB_MIGRATIONS_DIR) postgres "$(DSN)" down

clean:
	@echo "Cleaning up..."
	@rm -rf $(BIN_DIR)

docker-up:
	@echo "Starting Docker containers..."
	@docker-compose up -d

docker-down:
	@echo "Stopping Docker containers..."
	@docker-compose down

