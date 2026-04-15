# AGENTS.md

## Commands

- `./dev.sh dev` — start DB + hot reload (air)
- `./dev.sh run` — start DB + run server once
- `./dev.sh build` — compile to `bin/server`
- `./dev.sh seed` — run AutoMigrate + seed DB (needs `npm install` first run)
- `./dev.sh swagger` — regenerate swagger docs via `swag init`
- `./dev.sh db:start` / `db:stop` — manage PostgreSQL container

Plain run: `go run ./cmd/server` (requires DB already running)

## Prerequisites

- PostgreSQL running via `docker compose up -d` (or `./dev.sh db:start`)
- `.env` file — copy from `.env.example`
- `swag` CLI installed for swagger generation (`go install github.com/swaggo/swag/cmd/swag@latest`)
- `air` for hot reload (`go install github.com/air-verse/air@latest`)
- Node.js for seeding only (`npm install && node seed.js`)

## Architecture

Go module name: `register-system-be` (not a `github.com/` path — use this in imports)

Layered pattern — each feature has files in three packages:

1. **handler/** — HTTP handlers; use `Parse[T]` from `helper.go` for request binding
2. **business/** — business logic; define interface + implement in same file
3. **repo/** — database operations via GORM; define interface + implement in same file

Wiring: `main.go` constructs repo → business → handler, passes to `router.Setup()`.

- Route paths defined as constants in `constant/endpoint.go`
- Middleware in `middleware/`: `Auth` (JWT), `ActiveOnly`, `RequireRole`, `CORS`
- Models in `model/`: GORM structs, DTOs (`request.go`, `response.go`), JWT claims
- DB migrations are GORM `AutoMigrate` in `repo/db.go` — no manual migration files

## Adding an endpoint

1. Add path constant in `constant/endpoint.go`
2. Register route in `router/router.go`
3. Create handler in `handler/` (use `Parse[T]` for binding)
4. Create business logic in `business/` (interface + impl)
5. Create repo in `repo/` (interface + impl)
6. Wire in `cmd/server/main.go`

## CI

- `go build ./cmd/server`
- `go vet ./...`
- No test suite exists; no linter configured

## Notes

- `package.json` is for the seed script only — the app is pure Go
- Swagger docs in `docs/` are auto-generated; never edit by hand
- Seed script (`seed.js`) truncates `users` and `affiliations` before inserting
- User roles: `admin`, `student`, `school`, `community`
- CI workflow specifies Go 1.22 but `go.mod` says 1.25 — module may need alignment
