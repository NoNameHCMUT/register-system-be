# Implementation Plan — NEW_UPDATE Requirements

---

## Phase 1: Schema Changes

**`model/project.go`**
- Add `DateApproved *time.Time` (nullable — nil = not yet approved by school)
- Add `BannerURL string`

**`model/project.go` (StudentProject)**
- Replace `Status string` with a typed `ApplicationStatus` enum:
  - `SCHOOL_PENDING`, `SCHOOL_REJECT`, `COMMUNITY_PENDING`, `COMMUNITY_REJECT`, `APPROVED`
- Default status: `SCHOOL_PENDING`

**`model/user.go`**
- Add `AvatarURL string`
- Add `Phone string` (contact info for settings page)

**`model/affiliation.go`**
- Add `Description string` (optional)

**`repo/db.go`**
- Update `AutoMigrate` call — GORM handles new columns automatically

**`model/response.go`**
- Add `DateApproved`, `BannerURL` to `ProjectResponse`
- Add `AvatarURL`, `Phone` to `UserResponse`
- Add `Description` to `AffiliationResponse`
- Update `StudentProjectResponse` to use `ApplicationStatus`
- Update `ToUserResponse` helper

**`model/request.go`**
- Add `UserUpdateRequest` (phone, full_name)
- Add `ApplicationActionRequest` (batch: list of IDs + action)
- Add `AffiliationCreateRequest` (std_name, description)
- Add `AffiliationUpdateRequest` (std_name, description)

---

## Phase 2: Access Control

**`business/auth.go` — Login restriction**
- Check `user.IsActive` before generating tokens; return "account not approved" error if false

**Middleware — restrict inactive users**
- Extend `ActiveOnly` or create a new middleware that only allows `GET /auth/me` for inactive users. All other endpoints return 403.

---

## Phase 3: Project Approval Flow (School approves project)

New endpoints:

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| `GET` | `/schools/projects` | school | All projects matching school's affiliation |
| `GET` | `/schools/projects/pending` | school | Projects pending school approval (`date_approved` is nil) |
| `POST` | `/schools/projects/:id/approve` | school | Approve a project (sets `date_approved`) |

**Business logic:**
- Only school-role users whose `AffiliationID` matches the project's `AffiliationID` can approve
- Project must have `DateApproved == nil` to be approvable
- Setting `DateApproved` makes the project visible to students

**Visibility rules:**
- **Students**: only see projects where `date_approved != nil` AND `affiliation_id` matches their own
- **School**: see all projects matching their `affiliation_id`
- **Community**: see only projects they created (`community_user_id` matches)
- **Admin**: see all projects

---

## Phase 4: Student Application Flow

New endpoints:

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| `POST` | `/students/projects/:id/apply` | student | Apply to a project |
| `GET` | `/students/applications` | student | View own applications with status |

**Apply validation:**
- Only `student` role can apply
- Project must have `date_approved != nil`
- Student's `AffiliationID` must match project's `AffiliationID`
- Must be within form registration period (`form_start_day <= now <= form_end_day`)
- Student must not already have an application for this project (enforced by unique index)
- Creates `StudentProject` with status `SCHOOL_PENDING`

---

## Phase 5: Application Approval Flow (Batch)

**School endpoints:**

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| `GET` | `/schools/projects/:id/applicants` | school | View applicants for a project |
| `POST` | `/schools/applicants/action` | school | Batch approve/reject applicants |

**Community endpoints:**

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| `GET` | `/community/projects/:id/applicants` | community | View applicants for own project |
| `POST` | `/community/applicants/action` | community | Batch approve/reject applicants |

**Batch action request shape:**
```json
{
  "application_ids": [1, 2, 3],
  "action": "approve" | "reject"
}
```

**State transitions:**
- School `approve` → `SCHOOL_PENDING` → `COMMUNITY_PENDING`
- School `reject` → `SCHOOL_PENDING` → `SCHOOL_REJECT`
- Community `approve` → `COMMUNITY_PENDING` → `APPROVED`
- Community `reject` → `COMMUNITY_PENDING` → `COMMUNITY_REJECT`
- Validate that the acting role can only touch applications in the correct current state
- Validate ownership: school must share affiliation with project; community must own the project

---

## Phase 6: Admin Enhancements

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| `GET` | `/admin/projects` | admin | Master view of all projects |
| `GET` | `/admin/applications` | admin | View all applications |
| `POST` | `/admin/affiliations` | admin | Create affiliation |
| `PATCH` | `/admin/affiliations/:id` | admin | Update affiliation |
| `DELETE` | `/admin/affiliations/:id` | admin | Remove affiliation |

---

## Phase 7: User Settings & Profile

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| `PATCH` | `/users/me` | any active | Update contact info (phone, full_name) |
| `POST` | `/users/me/avatar` | any active | Upload avatar image |

---

## Phase 8: File Upload & Static Serving

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| `POST` | `/projects/:id/banner` | community/admin | Upload project banner |
| `GET` | `/uploads/avatars/:filename` | public | Serve avatar file |
| `GET` | `/uploads/banners/:filename` | public | Serve banner file |

**Implementation:**
- Store in `./uploads/avatars/` and `./uploads/banners/`
- Gin static file serving via `r.Static("/uploads", "./uploads")`
- Accept `multipart/form-data`, validate file type (jpg, png, gif) and size
- Save with UUID filename, store relative path in DB
- `ProjectCreateRequest` — add optional `banner` file upload
- `ProjectUpdateRequest` — add optional `banner` file upload

**Config additions** (`.env.example`):
```
UPLOAD_DIR=uploads
MAX_UPLOAD_SIZE=5242880
```

Add `uploads/` to `.gitignore`.

---

## Phase 9: Email Notifications

**New package: `email/`**

**Config additions** (`.env.example`):
```
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASSWORD=your-app-password
SMTP_FROM=Register System <your-email@gmail.com>
```

**`config/config.go`** — add SMTP fields

**`email/email.go`** — interface + SMTP implementation:
```go
type EmailSender interface {
    Send(to, subject, body string) error
}
```

**Send emails when:**
- Admin approves user account → notify user
- Admin rejects user account → notify user
- School approves/rejects student application → notify student
- Community approves/rejects student application → notify student

---

## Phase 10: Route Registration & Wiring

**`constant/endpoint.go`** — add all new path constants

**`router/router.go`** — register new route groups:
- `/schools/...` — Auth + ActiveOnly + RequireRole(school)
- `/community/...` — Auth + ActiveOnly + RequireRole(community)
- `/students/...` — Auth + ActiveOnly + RequireRole(student)
- `/admin/projects`, `/admin/applications`, `/admin/affiliations` — existing admin group
- `/users/me` — Auth + ActiveOnly
- `/uploads/...` — public static

**`cmd/server/main.go`** — wire new handlers with injected `email.EmailSender`

---

## Phase 11: Repo Layer Additions

**`repo/project.go`** — add:
- `FindByAffiliation(affiliationID uint)` — projects for a school's affiliation
- `FindByAffiliationApproved(affiliationID uint)` — approved projects only
- `FindByCreator(userID uint)` — community's own projects
- `FindAll()` — admin view
- `UpdateDateApproved(id uint, t time.Time)` — school approves project

**`repo/user.go`** — add:
- `UpdateProfile(userID uint, phone, avatarURL string)` — settings

**`repo/affiliation.go`** — add:
- `Create(aff *model.Affiliation)` — admin creates
- `Update(aff *model.Affiliation)` — admin updates
- `Delete(id uint)` — admin removes

New **`repo/application.go`** (or extend project repo):
- `FindByProjectID(projectID uint)` — list applications for a project
- `FindByStudentID(userID uint)` — student's own applications
- `FindByID(id uint)` — single application
- `BatchUpdateStatus(ids []uint, status model.ApplicationStatus)` — batch approve/reject
- `FindAll()` — admin view

---

## Phase 12: Testing

**Unit tests** (`*_test.go` alongside business files):
- Auth: login blocked for inactive users
- Project: approval flow, visibility rules
- Application: state transitions, date validation, duplicate prevention
- Application: batch approve/reject validation

**E2E test** (`tests/e2e_test.go`):
1. Admin pre-created
2. Admin creates 2 affiliations (1 school, 1 community)
3. School/community/student register → admin approves
4. Community creates project → sees it's unapproved
5. Student sees nothing (project not approved)
6. School approves project
7. Community sees project approved
8. Student sees the project
9. School views empty applicant list
10. Student applies (within form dates)
11. School sees applicant with SCHOOL_PENDING
12. Community sees no applicant (wrong state)
13. School approves → status becomes COMMUNITY_PENDING
14. Student sees COMMUNITY_PENDING
15. Community sees applicant
16. Community approves → status becomes APPROVED
17. Student sees APPROVED, school sees APPROVED

---

## Implementation Order

1. **Phase 1** (schema) — foundation for everything
2. **Phase 2** (access control) — small, blocking
3. **Phase 8** (file upload) — standalone, needed by phases 7 & project banner
4. **Phase 3** (project approval) — unblocks student view
5. **Phase 4** (student application) — depends on phase 3
6. **Phase 5** (application approval) — depends on phase 4
7. **Phase 6** (admin enhancements) — can be done in parallel with 3–5
8. **Phase 7** (user settings) — standalone
9. **Phase 9** (email) — cross-cutting, integrate after core flows
10. **Phase 10** (routes & wiring) — as each feature is built
11. **Phase 11** (repo additions) — alongside each feature
12. **Phase 12** (testing) — after all features
