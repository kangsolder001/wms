.PHONY: build run dev test lint clean docker-up docker-down

# Variables
APP_NAME=wms
BUILD_DIR=./bin
MAIN_PATH=./cmd/wms

# Build
build:
	go build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PATH)

# Run
run: build
	$(BUILD_DIR)/$(APP_NAME)

# Development
dev:
	air -c .air.toml

# Test
test:
	go test -v ./...

# Lint
lint:
	golangci-lint run ./...

# Clean
clean:
	rm -rf $(BUILD_DIR)

# Docker
docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

# Database
db-reset:
	docker-compose down -v
	docker-compose up -d postgres
	sleep 3
	go run $(MAIN_PATH)/main.go

# Tidy dependencies
tidy:
	go mod tidy

# Format code
fmt:
	go fmt ./...
