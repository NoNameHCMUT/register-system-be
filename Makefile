.PHONY: help build run test test-unit test-e2e lint ci clean swagger seed docker-build docker-run

help:
	@echo "Available commands:"
	@echo "  make build         - Build the application"
	@echo "  make run           - Run the application"
	@echo "  make test          - Run all tests (unit + e2e)"
	@echo "  make test-unit     - Run unit tests only"
	@echo "  make test-e2e      - Run e2e tests (requires DB)"
	@echo "  make lint          - Run golangci-lint"
	@echo "  make ci            - Run full CI pipeline (build + vet + lint + test)"
	@echo "  make swagger       - Regenerate swagger docs"
	@echo "  make seed          - Migrate + seed database"
	@echo "  make clean         - Clean build artifacts"
	@echo "  make docker-build  - Build Docker image"
	@echo "  make docker-run    - Run with Docker Compose"

build:
	CGO_ENABLED=0 go build -o bin/server ./cmd/server

run: build
	./bin/server

test-unit:
	go clean -testcache
	go test -v -race ./business/ ./handler/ ./repo/ ./model/ ./config/ ./email/ ./upload/

test-e2e:
	go clean -testcache
	go test -v -race -timeout 120s ./tests/

test: test-unit test-e2e

lint:
	golangci-lint run ./...

ci: build
	go vet ./...
	golangci-lint run ./...
	go clean -testcache
	go test -v -race ./business/ ./handler/ ./repo/ ./model/ ./config/ ./email/ ./upload/
	go test -v -race -timeout 120s ./tests/

swagger:
	swag init -g cmd/server/main.go -o docs

seed:
	go run ./cmd/server &
	sleep 3 && kill $$!
	npm install && node seed.js

clean:
	rm -rf bin/
	go clean -testcache

docker-build:
	docker build -t register-system-be:latest .

docker-run:
	docker compose up -d

docker-down:
	docker compose down

docker-logs:
	docker compose logs -f app
