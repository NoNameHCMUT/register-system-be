# Register System BE

Go + Gin backend with JWT auth, role-based access, and multi-stage project approval flow.

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

## Commands

| Command | Description |
|---|---|
| `./dev.sh run` | Start DB + run server |
| `./dev.sh dev` | Start DB + hot reload (air) |
| `./dev.sh build` | Compile to `bin/server` |
| `./dev.sh seed` | Migrate + seed DB |
| `./dev.sh swagger` | Regenerate swagger docs |
| `./dev.sh db:start` | Start PostgreSQL |
| `./dev.sh db:stop` | Stop PostgreSQL |
| `./dev.sh test` | Run all tests (unit + e2e) |
| `./dev.sh test:unit` | Run unit tests only |
| `./dev.sh test:e2e` | Run e2e tests (requires DB) |
| `./dev.sh lint` | Run golangci-lint |
| `./dev.sh ci` | Full CI pipeline (build + vet + lint + test) |

Or with Make:

```bash
make test         # all tests
make test-unit    # unit only
make test-e2e     # e2e only (requires DB)
make lint         # golangci-lint
make ci           # full CI pipeline
```

## Seed Database

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

## API Endpoints

### Public

| Method | Path | Description |
|---|---|---|
| POST | /auth/register | Register user |
| POST | /auth/login | Login (active users only) |
| POST | /auth/refresh | Refresh tokens |
| GET | /affiliations | List affiliations |
| GET | /health | Health check |

### Auth (any active user)

| Method | Path | Description |
|---|---|---|
| GET | /auth/me | Current user |
| PATCH | /users/me | Update profile (phone, full_name) |
| POST | /users/me/avatar | Upload avatar |

### Admin

| Method | Path | Description |
|---|---|---|
| GET | /admin/users/pending | List pending users |
| POST | /admin/users/:id/accept | Accept user |
| POST | /admin/users/:id/reject | Reject user |
| GET | /admin/projects | All projects |
| GET | /admin/applications | All applications |
| POST | /admin/affiliations | Create affiliation |
| PATCH | /admin/affiliations/:id | Update affiliation |
| DELETE | /admin/affiliations/:id | Delete affiliation |

### School

| Method | Path | Description |
|---|---|---|
| GET | /schools/projects | Projects in school's affiliation |
| GET | /schools/projects/pending | Projects awaiting school approval |
| POST | /schools/projects/:id/approve | Approve a project |
| GET | /schools/projects/:id/applicants | View applicants for a project |
| POST | /schools/applicants/action | Batch approve/reject applicants |

### Community

| Method | Path | Description |
|---|---|---|
| GET | /community/projects | Own projects |
| GET | /community/projects/:id/applicants | View applicants for own project |
| POST | /community/applicants/action | Batch approve/reject applicants |

### Student

| Method | Path | Description |
|---|---|---|
| GET | /students/projects | Approved projects for student's affiliation |
| POST | /students/projects/:id/apply | Apply to a project |
| GET | /students/applications | View own applications |

### Project (community/admin)

| Method | Path | Description |
|---|---|---|
| POST | /projects | Create project |
| PATCH | /projects/:id | Update project |
| POST | /projects/:id/banner | Upload project banner |

## Application Approval Flow

```
Student applies → SCHOOL_PENDING
    ↓ School approves → COMMUNITY_PENDING
        ↓ Community approves → APPROVED
    ↓ School rejects → SCHOOL_REJECT
        ↓ Community rejects → COMMUNITY_REJECT
```

- Students can only apply to school-approved projects within their affiliation
- Applications are only accepted during the project's form registration period
- Batch approve/reject supported (multiple application IDs per request)

## Adding a New Endpoint

1. **Define** path in `constant/endpoint.go`
2. **Register** in `router/router.go`
3. **Code handler** in `handler/` (use `Parse[T]`)
4. **Code business** in `business/` (interface + impl)
5. **Code repo** in `repo/` (interface + impl)
6. **Wire** in `cmd/server/main.go`

## Project Structure

```
cmd/server/     Entry point
handler/        HTTP handlers + Parse[T] helper
business/       Business logic (interfaces + impl)
repo/           Database operations (interfaces + impl)
middleware/      Auth, ActiveOnly, RequireRole, CORS
model/          GORM models + DTOs + enums
config/         Env config
constant/       API endpoint paths
router/         Route registration
email/          SMTP email sender
upload/         File upload utility
docs/           Swagger spec (generated)
tests/          E2E tests
```
