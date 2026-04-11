# Register System BE

Go + Gin backend with JWT auth and role-based access.

## Setup

```bash
cp .env.example .env
docker compose up -d
```

## Run

```bash
go run ./cmd/server
```

Dev mode with hot reload:

```bash
./dev.sh dev
```

## Seed Database

First run the server once to create tables, then seed:

```bash
go run ./cmd/server &   # creates tables via AutoMigrate
# Ctrl+C after "Server starting"
npm install
node seed.js
```

Or use the shortcut:

```bash
./dev.sh seed
```

Sample accounts:

| Username | Password | Role |
|---|---|---|
| admin_root | admin123 | admin |
| cb_bachkhoa | school456 | school |
| leader_binh_phuoc | local789 | community |
| student_nhu | student_nhu | student |

## Swagger

Available at `http://localhost:8080/swagger/index.html`.

Regenerate after endpoint changes:

```bash
./dev.sh swagger
```

## Commands

| Command | Description |
|---|---|
| `./dev.sh run` | Start DB + run server |
| `./dev.sh dev` | Start DB + hot reload |
| `./dev.sh build` | Compile to `bin/server` |
| `./dev.sh seed` | Migrate + seed DB |
| `./dev.sh swagger` | Regenerate swagger docs |
| `./dev.sh db:start` | Start PostgreSQL |
| `./dev.sh db:stop` | Stop PostgreSQL |

## API Endpoints

| Method | Path | Auth | Description |
|---|---|---|---|
| POST | /auth/register | No | Register user |
| POST | /auth/login | No | Login |
| POST | /auth/refresh | No | Refresh tokens |
| GET | /auth/me | Yes | Current user |
| GET | /admin/users/pending | Admin | List pending users |
| POST | /admin/users/:id/accept | Admin | Accept user |
| POST | /admin/users/:id/reject | Admin | Reject user |
| GET | /health | No | Health check |

## Adding a New Endpoint

1. **Define** path in `constant/endpoint.go`
2. **Register** in `router/router.go`
3. **Code handler** in `handler/` (use `Parse[T]`)
4. **Code business** in `business/` (interface + impl)
5. **Code repo** in `repo/` (interface + impl)

## Project Structure

```
cmd/server/     Entry point
handler/        HTTP handlers + Parse[T] helper
business/       Business logic (interfaces + impl)
repo/           Database operations (interfaces + impl)
middleware/      Auth, Role, CORS
model/          GORM models + DTOs
config/         Env config
constant/       API endpoint paths
router/         Route registration
docs/           Swagger spec (generated)
```
