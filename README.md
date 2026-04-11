# Register System BE

Go + Gin backend skeleton with JWT auth and role-based access.

## Quick Start

```bash
# 1. Start PostgreSQL
./dev.sh db:start

# 2. Run the server
./dev.sh run
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
