# Register System BE

Go + Gin backend skeleton with JWT auth and role-based access.

## Quick Start
docker compose up -d --build db api
docker compose run --rm seeder
go run ./cmd/server

```bash
# 1. Start PostgreSQL
./dev.sh db:start

# 2. Run the server
./dev.sh run
```

### Seed sample data (Docker)

The repo includes a Node-based seeder (`seed.js`) that can be run via Docker Compose.

```bash
docker compose up -d --build db api
docker compose run --rm seeder
```

#### (Optional) Log plaintext passwords (test only)

By default, `seed.js` does **not** print plaintext passwords. For testing, you can enable it:

```bash
SEED_LOG_PASSWORDS=true docker compose run --rm seeder
```

PowerShell (Windows):

```powershell
$env:SEED_LOG_PASSWORDS = "true"
docker compose run --rm seeder
```

> Only use this for local testing.

### Seed sample data (Local, without Docker)

1) Ensure PostgreSQL is running and `.env` is configured:

```bash
cp .env.example .env
# Edit DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME as needed
```

2) Run the API once to auto-migrate tables (GORM AutoMigrate):

```bash
./dev.sh run
```

3) Run the seeder:

```bash
npm install
npm run seed
```

To print plaintext passwords (test only):

```bash
SEED_LOG_PASSWORDS=true npm run seed
```

PowerShell (Windows):

```powershell
$env:SEED_LOG_PASSWORDS = "true"
npm run seed
```

Server starts at `http://localhost:8080`.

Swagger UI at `http://localhost:8080/swagger/index.html`.

## Swagger

Swagger UI is served at `/swagger/index.html` when the server is running.

Docs are auto-generated from annotations in handler files. After changing endpoints:

```bash
~/go/bin/swag init -g cmd/server/main.go -o docs
```

## All Commands

| Command | Description |
|---------|-------------|
| `./dev.sh run` | Start DB + run the server |
| `./dev.sh dev` | Start DB + run with hot reload (air) |
| `./dev.sh build` | Compile binary to `bin/server` |
| `./dev.sh db:start` | Start PostgreSQL only |
| `./dev.sh db:stop` | Stop PostgreSQL |

## Running without Docker (existing PostgreSQL)

```bash
cp .env.example .env
# Edit DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME as needed
./dev.sh run
```

## API Endpoints

| Method | Path                   | Auth  | Description       |
|--------|------------------------|-------|-------------------|
| POST   | /auth/register         | No    | Register user     |
| POST   | /auth/login            | No    | Login             |
| POST   | /auth/refresh          | No    | Refresh tokens    |
| GET    | /auth/me               | Yes   | Current user      |
| GET    | /admin/users/pending   | Admin | List pending users|
| POST   | /admin/users/:id/accept| Admin | Accept user       |
| POST   | /admin/users/:id/reject| Admin | Reject user      |
| GET    | /health                | No    | Health check      |
| GET    | /swagger/*any          | No    | Swagger UI        |

## Adding a New Endpoint

1. **Define** path in `constant/endpoint.go`
2. **Register** in `router/router.go`
3. **Code handler** in `handler/` (use `Parse[T]` for request body)
4. **Code business** in `business/` (define interface + impl)
5. **Code repo** in `repo/` (define interface + impl)

## Project Structure

```
├── cmd/server/       # Entry point
├── handler/          # HTTP handlers + Parse[T] helper
├── business/         # Business logic (interfaces + impl)
├── repo/             # Database operations (interfaces + impl)
├── middleware/        # Auth, Role, CORS
├── model/            # GORM models + DTOs
├── config/           # Env config
├── constant/         # API endpoint paths
├── router/           # Route registration
├── docs/             # Swagger spec
```
