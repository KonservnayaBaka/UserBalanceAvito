GO_BINARY := ./cmd/app/main
GO_PACKAGES := $(shell go list ./... | grep -v /vendor/)

all: build

build:
	@echo "Building the project..."
	@go build -o $(GO_BINARY) ./cmd/app

run: build
	@echo "Running the project..."
	@./$(GO_BINARY)

test:
	@echo "Running tests..."
	@go test $(GO_PACKAGES)

lint:
	@echo "Running linter..."
	@golangci-lint run

clean:
	@echo "Cleaning up..."
	@rm -f $(GO_BINARY)

docker-build:
	@echo "Building Docker image..."
	@docker build -t user-balance-avito .

docker-run: docker-build
	@echo "Running Docker container..."
	@docker-compose up --build

docker-down:
	@echo "Stopping and removing Docker containers..."
	@docker-compose down

help:
	@echo "Makefile commands:"
	@echo "  all         - Build the project"
	@echo "  build       - Build the project"
	@echo "  run         - Run the project"
	@echo "  test        - Run tests"
	@echo "  lint        - Run linter"
	@echo "  clean       - Clean up build artifacts"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-run  - Run Docker container"
	@echo "  docker-down - Stop and remove Docker containers"
	@echo "  help        - Show this help message"
