#!/bin/bash
set -e

load_env() {
    if [ -f .env ]; then
        set -a; source .env; set +a
    fi
}

DB_USER="${DB_USER:-postgres}"
DB_NAME="${DB_NAME:-register_system}"

info()  { echo -e "\033[0;32m[SEED]\033[0m $1"; }
error() { echo -e "\033[0;31m[ERROR]\033[0m $1"; exit 1; }

psql_exec() {
    local container
    container=$(docker ps -q --filter "ancestor=postgres:16-alpine" | head -1)
    if [ -n "$container" ]; then
        docker exec "$container" psql -U "$DB_USER" -d "$DB_NAME" -c "$1"
    elif command -v psql &> /dev/null; then
        psql -h "${DB_HOST:-localhost}" -p "${DB_PORT:-5432}" -U "$DB_USER" -d "$DB_NAME" -c "$1"
    else
        error "No PostgreSQL client found (neither docker container nor psql)"
    fi
}

load_env

info "Seeding affiliations..."
psql_exec "INSERT INTO affiliations (id, std_name) VALUES (0, 'Admin'), (1, 'School of Engineering') ON CONFLICT (id) DO NOTHING;"

info "Seeding admin user..."
HASH=$(go run ./cmd/hashpw)
psql_exec "INSERT INTO users (username, full_name, password_hash, email, affiliation_id, role, is_active) VALUES ('admin', 'System Admin', '$HASH', 'admin@register-system.local', 0, 'admin', true) ON CONFLICT (username) DO NOTHING;"

info "Done! admin / admin123"
