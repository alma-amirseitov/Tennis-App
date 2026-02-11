# CLAUDE.md — Tennis Platform (Astana)

## Project

Tennis social platform for Astana, Kazakhstan. Mobile app (Expo/RN) + Web admin (React) + Go backend.

## Architecture

Monorepo. Layered monolith backend. REST API + WebSocket (chat).

```
apps/backend/     — Go 1.22+ (Chi router, sqlc, PostgreSQL, Redis)
apps/mobile/      — React Native + Expo (TypeScript, Expo Router)
apps/web-admin/   — React + Vite + Shadcn/UI + Tailwind
apps/web-superadmin/ — React + Vite + Shadcn/UI + Tailwind
packages/shared-types/ — Shared TypeScript types
packages/api-client/   — Shared API client
```

## Tech Stack

**Backend:** Go 1.22+, Chi router, sqlc (NOT GORM), golang-migrate, nhooyr.io/websocket (NOT gorilla), golang-jwt/jwt v5, go-playground/validator v10, slog (stdlib), pgx v5, go-redis/redis v9

**Mobile:** TypeScript strict, Expo SDK 52+, Expo Router, TanStack Query, Zustand, React Hook Form + Zod, Axios, react-native-keychain, i18next

**Web:** TypeScript strict, React 18+, Vite, Shadcn/UI, Tailwind CSS, TanStack Query + Table, Zustand, React Hook Form + Zod, Recharts

**Infra:** PostgreSQL 16, Redis 7, MinIO (S3-compatible), Firebase FCM, Docker Compose

## Backend Structure

```
apps/backend/
├── cmd/server/main.go           — Entry point, graceful shutdown
├── internal/
│   ├── config/config.go         — Env-based config (envconfig)
│   ├── handler/                 — HTTP handlers (thin: validate + respond)
│   │   ├── middleware/          — Auth, RBAC, rate limit, CORS, logging
│   │   ├── dto/                 — Request/Response structs
│   │   ├── auth.go
│   │   ├── user.go
│   │   ├── event.go
│   │   ├── community.go
│   │   ├── chat.go
│   │   ├── post.go
│   │   ├── notification.go
│   │   ├── rating.go
│   │   ├── court.go
│   │   ├── admin.go
│   │   └── superadmin.go
│   ├── service/                 — Business logic (core)
│   ├── repository/              — DB queries (sqlc generated + custom)
│   │   └── queries/             — .sql files for sqlc
│   ├── model/                   — Domain models
│   ├── ws/                      — WebSocket hub + client management
│   └── pkg/                     — Shared utilities
│       ├── jwt/
│       ├── sms/
│       ├── firebase/
│       ├── elo/
│       ├── storage/             — S3/MinIO wrapper
│       └── validator/
├── migrations/                  — SQL migration files
├── sqlc.yaml
├── go.mod
└── go.sum
```

## Coding Rules — Backend (Go)

### Naming
- Files: `snake_case.go`
- Packages: single lowercase word (`handler`, `service`, `repository`)
- Structs/Interfaces: `PascalCase` — `EventService`, `UserRepository`
- Functions: `PascalCase` for exported, `camelCase` for private
- Variables: `camelCase`
- Constants: `PascalCase` or `SCREAMING_SNAKE` for env keys
- DB columns: `snake_case`
- JSON fields: `snake_case`
- URL paths: `kebab-case` for multi-word (`/match-results`)

### Layer Rules

**Handler (thin):**
```go
func (h *EventHandler) Create(w http.ResponseWriter, r *http.Request) {
    var req dto.CreateEventRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        respondError(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request body")
        return
    }
    if err := h.validator.Struct(req); err != nil {
        respondValidationError(w, err)
        return
    }
    userID := middleware.GetUserID(r.Context())
    event, err := h.service.Create(r.Context(), userID, req.ToModel())
    if err != nil {
        handleServiceError(w, err)
        return
    }
    respondJSON(w, http.StatusCreated, dto.EventFromModel(event))
}
```

**Service (business logic):**
```go
func (s *EventService) Create(ctx context.Context, userID uuid.UUID, input model.CreateEventInput) (*model.Event, error) {
    // Business rules here
    if input.LevelMin > input.LevelMax {
        return nil, ErrInvalidLevelRange
    }
    if input.MaxParticipants < 2 {
        return nil, ErrTooFewParticipants
    }
    event, err := s.repo.Create(ctx, input)
    if err != nil {
        return nil, fmt.Errorf("create event: %w", err)
    }
    return event, nil
}
```

**Repository (data access, sqlc):**
```sql
-- queries/event.sql

-- name: CreateEvent :one
INSERT INTO events (community_id, creator_id, title, description, event_type, composition_type, level_min, level_max, max_participants, sets_count, start_time, end_time, location_name)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
RETURNING *;

-- name: GetEventByID :one
SELECT * FROM events WHERE id = $1 AND deleted_at IS NULL;

-- name: ListEvents :many
SELECT * FROM events
WHERE deleted_at IS NULL
  AND ($1::uuid IS NULL OR community_id = $1)
  AND ($2::text IS NULL OR event_type = $2)
  AND ($3::text IS NULL OR status = $3)
ORDER BY start_time DESC
LIMIT $4 OFFSET $5;
```

### Error Handling
```go
// Define typed errors in service layer
var (
    ErrNotFound          = errors.New("not found")
    ErrForbidden         = errors.New("forbidden")
    ErrAlreadyExists     = errors.New("already exists")
    ErrInvalidInput      = errors.New("invalid input")
    ErrInvalidLevelRange = fmt.Errorf("%w: level_min > level_max", ErrInvalidInput)
)

// Handler maps service errors to HTTP status codes
func handleServiceError(w http.ResponseWriter, err error) {
    switch {
    case errors.Is(err, ErrNotFound):
        respondError(w, 404, "NOT_FOUND", err.Error())
    case errors.Is(err, ErrForbidden):
        respondError(w, 403, "FORBIDDEN", err.Error())
    case errors.Is(err, ErrAlreadyExists):
        respondError(w, 409, "ALREADY_EXISTS", err.Error())
    case errors.Is(err, ErrInvalidInput):
        respondError(w, 400, "VALIDATION_ERROR", err.Error())
    default:
        slog.Error("unhandled error", "error", err)
        respondError(w, 500, "INTERNAL_ERROR", "Internal server error")
    }
}
```

### Response Format
```go
// Success
func respondJSON(w http.ResponseWriter, status int, data any) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(data)
}

// Error
func respondError(w http.ResponseWriter, status int, code, message string) {
    respondJSON(w, status, map[string]any{
        "error": map[string]any{
            "code":    code,
            "message": message,
        },
    })
}

// Paginated
type PaginatedResponse struct {
    Data       any        `json:"data"`
    Pagination Pagination `json:"pagination"`
}
type Pagination struct {
    Page      int `json:"page"`
    PerPage   int `json:"per_page"`
    Total     int `json:"total"`
    TotalPages int `json:"total_pages"`
}
```

### Database Rules
- ALWAYS use prepared statements (sqlc handles this)
- NEVER `SELECT *` without LIMIT in list queries
- ALWAYS add `deleted_at IS NULL` for soft-deleted tables
- ALWAYS use transactions for multi-table writes
- Index every column used in WHERE/ORDER BY
- UUID for all primary keys (gen_random_uuid())
- timestamps: `created_at`, `updated_at`, `deleted_at`

### Testing
```go
// Table-driven tests in service layer
func TestEventService_Create(t *testing.T) {
    tests := []struct {
        name    string
        input   model.CreateEventInput
        wantErr error
    }{
        {
            name:    "valid event",
            input:   validEventInput(),
            wantErr: nil,
        },
        {
            name:    "invalid level range",
            input:   model.CreateEventInput{LevelMin: 4.0, LevelMax: 2.0},
            wantErr: ErrInvalidLevelRange,
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // arrange, act, assert
        })
    }
}
```

## API Design

Base URL: `/v1`
Auth: `Authorization: Bearer <jwt>`
All responses: JSON, snake_case fields
Errors: `{ "error": { "code": "...", "message": "...", "details": [...] } }`
Pagination: `?page=1&per_page=20` → response includes `pagination` object
Full API spec: `docs/04-API-SPECIFICATION.md`

## Key Business Logic

- **ELO Rating:** See `docs/08-ELO-ALGORITHM.md` — K-factor varies by games played
- **Auth flow:** Phone → OTP → Profile Setup (new) → Quiz → PIN
- **RBAC:** Community roles: owner, admin, moderator, member
- **Event lifecycle:** draft → open → filling → full → in_progress → completed → cancelled
- **Match confirmation:** Both players must confirm result within 48h

## Environment Variables

```
DATABASE_URL=postgres://user:pass@localhost:5432/tennisapp?sslmode=disable
REDIS_URL=redis://localhost:6379
JWT_SECRET=your-secret-key
JWT_ACCESS_TTL=15m
JWT_REFRESH_TTL=720h
SMS_API_KEY=xxx
SMS_API_URL=https://smsc.kz/sys/send.php
FIREBASE_CREDENTIALS=path/to/firebase.json
S3_ENDPOINT=http://localhost:9000
S3_ACCESS_KEY=minioadmin
S3_SECRET_KEY=minioadmin
S3_BUCKET=tennisapp
ENVIRONMENT=development
PORT=8080
```

## Commands

```bash
make dev              # Run with hot-reload (air)
make build            # Build binary
make test             # Run all tests
make lint             # golangci-lint
make migrate-up       # Apply migrations
make migrate-down     # Rollback last migration
make migrate-create   # Create new migration (NAME=xxx)
make sqlc             # Generate Go code from SQL queries
make seed             # Seed database with test data
```

## Important: Do NOT

- Do NOT use GORM, use sqlc
- Do NOT use gorilla/websocket, use nhooyr.io/websocket
- Do NOT use global variables for state
- Do NOT put business logic in handlers
- Do NOT return 500 for user input errors
- Do NOT use `SELECT *` in production queries (list columns explicitly in sqlc)
- Do NOT store secrets in code
- Do NOT skip error wrapping (`fmt.Errorf("context: %w", err)`)
- Do NOT create endpoints not documented in `docs/04-API-SPECIFICATION.md`
