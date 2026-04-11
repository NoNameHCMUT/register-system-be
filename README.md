# Register System BE

Go + Gin backend skeleton with JWT auth and role-based access.

## Quick Start

```bash
# Copy env and edit as needed
cp .env.example .env

# Run with Docker (spins up app + PostgreSQL)
docker compose up --build

# Or run locally (requires running PostgreSQL)
go run ./cmd/server
```

Server starts at `http://localhost:8080`.

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
