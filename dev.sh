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

run() {
    setup_env
    start_db
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
    info "Starting dev server with hot reload..."
    air
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
    npm install --silent
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

case "${1:-}" in
    run)      run      ;;
    dev)      dev      ;;
    build)    build    ;;
    seed)     seed     ;;
    swagger)  swagger  ;;
    db:start) setup_env; start_db ;;
    db:stop)  stop_db  ;;
    *)
        echo "Usage: $0 {run|dev|build|seed|swagger|db:start|db:stop}"
        echo ""
        echo "  run       Start DB + run the server"
        echo "  dev       Start DB + run with hot reload (air)"
        echo "  build     Compile binary to bin/server"
        echo "  seed      Migrate + seed database"
        echo "  swagger   Regenerate swagger docs"
        echo "  db:start  Start PostgreSQL only"
        echo "  db:stop   Stop PostgreSQL"
        exit 1
        ;;
esac
