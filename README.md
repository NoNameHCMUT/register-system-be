# Register System BE

Go + Gin backend skeleton with JWT auth and role-based access.

## Quick Start

### 1. Start PostgreSQL + API (Docker)

```bash
docker compose up -d --build db api
```

### 2. Seed sample data (Docker)

```bash
docker compose run --rm seeder
```

### 3. Configure environment

```bash
cp .env.example .env
# Edit .env if needed (defaults match the Docker PostgreSQL)
```

### 4. Run the server (local Go)

```bash
go run ./cmd/server
```

Server starts at `http://localhost:8080`.

## Running without Docker (existing PostgreSQL)

Set the env vars in `.env` to point to your local PostgreSQL instance:

```bash
cp .env.example .env
# Edit DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME as needed
go run ./cmd/server
```

## API Endpoints

| Method | Path           | Auth | Description   |
|--------|----------------|------|---------------|
| POST   | /auth/register | No   | Register user |
| POST   | /auth/login    | No   | Login         |
| POST   | /auth/refresh  | No   | Refresh tokens|
| GET    | /auth/me       | Yes  | Current user  |
| GET    | /health        | No   | Health check  |

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
```
