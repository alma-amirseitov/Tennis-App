# Tennis App — Astana

Tennis social platform for Astana, Kazakhstan. Mobile app (Expo/RN) + Web admin (React) + Go backend.

## Architecture

Monorepo with layered monolith backend, REST API + WebSocket (chat).

```
apps/backend/          — Go 1.22+ (Chi, sqlc, PostgreSQL, Redis)
apps/mobile/           — React Native + Expo (TypeScript)
apps/web-admin/        — React + Vite + Shadcn/UI
apps/web-superadmin/   — React + Vite + Shadcn/UI
packages/shared-types/ — Shared TypeScript types
packages/api-client/   — Shared API client
```

## Quick Start

### Prerequisites

- Go 1.22+
- Docker & Docker Compose
- Node.js 20+
- [golang-migrate](https://github.com/golang-migrate/migrate)
- [sqlc](https://sqlc.dev/)
- [air](https://github.com/air-verse/air) (hot-reload)
- [golangci-lint](https://golangci-lint.run/)

### 1. Start infrastructure

```bash
docker compose up -d
```

This starts PostgreSQL 16, Redis 7, and MinIO (S3-compatible storage).

### 2. Setup backend

```bash
cp apps/backend/.env.example apps/backend/.env
cd apps/backend
make migrate-up
make seed
make dev
```

API will be available at `http://localhost:8080`.

### Services

| Service    | Port  | URL                           |
|------------|-------|-------------------------------|
| Backend    | 8080  | http://localhost:8080          |
| PostgreSQL | 5432  | postgres://localhost:5432      |
| Redis      | 6379  | redis://localhost:6379         |
| MinIO API  | 9000  | http://localhost:9000          |
| MinIO UI   | 9001  | http://localhost:9001          |

### Backend Commands

Run from `apps/backend/`:

```bash
make dev              # Run with hot-reload (air)
make build            # Build binary
make test             # Run all tests
make test-coverage    # Run tests with coverage report
make lint             # golangci-lint
make migrate-up       # Apply migrations
make migrate-down     # Rollback last migration
make migrate-create NAME=xxx  # Create new migration
make sqlc             # Generate Go code from SQL queries
make seed             # Seed database with test data
```

### Infrastructure Commands

Run from project root:

```bash
docker compose up -d     # Start services
docker compose down      # Stop services
docker compose logs -f   # View logs
```

## Documentation

See `docs/` directory for full project documentation.
