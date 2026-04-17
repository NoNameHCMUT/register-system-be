#!/bin/bash
set -e

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

info() { echo -e "${GREEN}[INFO]${NC} $1"; }
warn() { echo -e "${YELLOW}[WARN]${NC} $1"; }

setup_env() {
    if [ ! -f .env ]; then
        info "Copying .env.example -> .env"
        cp .env.example .env
    fi
}

start_db() {
    if ! docker compose ps | grep -q "db.*running\|db.*healthy"; then
        info "Starting PostgreSQL..."
        docker compose up -d
        sleep 3
    else
        info "PostgreSQL already running"
    fi
}

stop_db() {
    info "Stopping PostgreSQL..."
    docker compose down
}

setup() {
    info "Installing Go dependencies..."
    go mod tidy
    info "Generating swagger docs..."
    if ! command -v swag &> /dev/null && [ ! -f "$(go env GOPATH)/bin/swag" ]; then
        warn "swag not installed. Installing..."
        go install github.com/swaggo/swag/cmd/swag@latest
    fi
    $(go env GOPATH)/bin/swag init -g cmd/server/main.go -o docs
}

run() {
    setup_env
    start_db
    setup
    info "Running server..."
    go run ./cmd/server
}

dev() {
    if ! command -v air &> /dev/null; then
        warn "air not installed. Installing..."
        go install github.com/air-verse/air@latest
    fi
    setup_env
    start_db
    setup
    info "Starting dev server with hot reload..."
    $(go env GOPATH)/bin/air -c .air.toml
}

build() {
    info "Building binary..."
    go build -o bin/server ./cmd/server
    info "Binary saved to bin/server"
}

seed() {
    setup_env
    start_db
    info "Running migrations..."
    go run ./cmd/server &
    SERVER_PID=$!
    sleep 3
    kill $SERVER_PID 2>/dev/null
    wait $SERVER_PID 2>/dev/null
    info "Seeding database..."
    npm install
    node seed.js
}

swagger() {
    info "Generating swagger docs..."
    if ! command -v swag &> /dev/null; then
        warn "swag not installed. Installing..."
        go install github.com/swaggo/swag/cmd/swag@latest
    fi
    $(go env GOPATH)/bin/swag init -g cmd/server/main.go -o docs
    info "Done"
}

test_unit() {
    info "Running unit tests..."
    go clean -testcache
    go test -v -race ./business/ ./handler/ ./repo/ ./model/ ./config/ ./email/ ./upload/
}

test_e2e() {
    setup_env
    start_db
    info "Running e2e tests..."
    go clean -testcache
    go test -v -race -timeout 120s ./tests/
}

test() {
    test_unit
    test_e2e
}

lint() {
    info "Running golangci-lint..."
    if ! command -v golangci-lint &> /dev/null && [ ! -f "$(go env GOPATH)/bin/golangci-lint" ]; then
        warn "golangci-lint not installed. Installing..."
        go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
    fi
    golangci-lint run ./...
}

ci() {
    info "Running CI pipeline locally..."
    info "--- Step 1: Build ---"
    go build ./cmd/server
    info "--- Step 2: Vet ---"
    go vet ./...
    info "--- Step 3: Lint ---"
    lint
    info "--- Step 4: Unit tests ---"
    test_unit
    info "--- Step 5: E2E tests ---"
    test_e2e
    info "CI pipeline passed!"
}

case "${1:-}" in
    run)      run      ;;
    dev)      dev      ;;
    build)    build    ;;
    seed)     seed     ;;
    swagger)  swagger  ;;
    db:start) setup_env; start_db ;;
    db:stop)  stop_db  ;;
    test)     test     ;;
    test:unit) test_unit ;;
    test:e2e) test_e2e ;;
    lint)     lint     ;;
    ci)       ci       ;;
    *)
        echo "Usage: $0 {run|dev|build|seed|swagger|db:start|db:stop|test|test:unit|test:e2e|lint|ci}"
        echo ""
        echo "  run        Start DB + run the server"
        echo "  dev        Start DB + run with hot reload (air)"
        echo "  build      Compile binary to bin/server"
        echo "  seed       Migrate + seed database"
        echo "  swagger    Regenerate swagger docs"
        echo "  db:start   Start PostgreSQL only"
        echo "  db:stop    Stop PostgreSQL"
        echo "  test       Run all tests (unit + e2e)"
        echo "  test:unit  Run unit tests only"
        echo "  test:e2e   Run e2e tests only (requires DB)"
        echo "  lint       Run golangci-lint"
        echo "  ci         Run full CI pipeline locally"
        exit 1
        ;;
esac
