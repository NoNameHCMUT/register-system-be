.PHONY: help build run test lint clean docker-build docker-run

help:
	@echo "Available commands:"
	@echo "  make build         - Build the application"
	@echo "  make run           - Run the application"
	@echo "  make test          - Run tests"
	@echo "  make lint          - Run linter"
	@echo "  make clean         - Clean build artifacts"
	@echo "  make docker-build  - Build Docker image"
	@echo "  make docker-run    - Run with Docker Compose"

build:
	CGO_ENABLED=0 go build -o bin/server ./cmd/server

run: build
	./bin/server

test:
	go test -v -race -coverprofile=coverage.out ./...

lint:
	golangci-lint run

clean:
	rm -rf bin/
	go clean

docker-build:
	docker build -t register-system-be:latest .

docker-run:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f app
