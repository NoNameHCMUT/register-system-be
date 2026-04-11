#!/bin/bash
set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

info()  { echo -e "${GREEN}[INFO]${NC}  $1"; }
warn()  { echo -e "${YELLOW}[WARN]${NC}  $1"; }
error() { echo -e "${RED}[ERROR]${NC} $1"; }

check_cmd() {
    if ! command -v "$1" &> /dev/null; then
        error "$1 is not installed. Please install it first."
        exit 1
    fi
}

setup_env() {
    if [ ! -f .env ]; then
        info "Copying .env.example -> .env"
        cp .env.example .env
    else
        info ".env already exists, skipping"
    fi
}

start_db() {
    if ! docker compose ps | grep -q "db.*running\|db.*healthy"; then
        info "Starting PostgreSQL..."
        docker compose up -d
        info "Waiting for PostgreSQL to be ready..."
        sleep 3
    else
        info "PostgreSQL is already running"
    fi
}

stop_db() {
    info "Stopping PostgreSQL..."
    docker compose down
}

run() {
    check_cmd go
    setup_env
    start_db
    info "Running server..."
    go run ./cmd/server
}

dev() {
    check_cmd go
    if ! command -v air &> /dev/null; then
        warn "air is not installed. Installing..."
        go install github.com/air-verse/air@latest
    fi
    setup_env
    start_db
    info "Starting dev server with hot reload..."
    air
}

build() {
    check_cmd go
    info "Building binary..."
    go build -o bin/server ./cmd/server
    info "Binary saved to bin/server"
}

migrate() {
    check_cmd go
    setup_env
    start_db
    info "Running migrations via server startup..."
    go run ./cmd/server
}

case "${1:-}" in
    run)      run      ;;
    dev)      dev      ;;
    build)    build    ;;
    db:start) setup_env; start_db ;;
    db:stop)  stop_db  ;;
    *)
        echo "Usage: $0 {run|dev|build|db:start|db:stop}"
        echo ""
        echo "Commands:"
        echo "  run        Start PostgreSQL + run the server"
        echo "  dev        Start PostgreSQL + run with hot reload (air)"
        echo "  build      Compile binary to bin/server"
        echo "  db:start   Start PostgreSQL only"
        echo "  db:stop    Stop PostgreSQL"
        exit 1
        ;;
esac
