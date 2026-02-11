# Deployment & Infrastructure

---

## 1. Environments

| Env | Backend | Database | Redis | Purpose |
|-----|---------|----------|-------|---------|
| local | localhost:8080 | Docker Postgres:5432 | Docker Redis:6379 | Development |
| staging | staging.api.tennisapp.kz | Railway Postgres | Upstash Redis | Testing |
| production | api.tennisapp.kz | Railway/Supabase Postgres | Upstash Redis | Production |

---

## 2. Local Development

### docker-compose.yml
```yaml
version: '3.8'

services:
  postgres:
    image: postgres:16-alpine
    environment:
      POSTGRES_DB: tennisapp
      POSTGRES_USER: tennisapp
      POSTGRES_PASSWORD: tennisapp
    ports:
      - "5432:5432"
    volumes:
      - pgdata:/var/lib/postgresql/data

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"

  minio:
    image: minio/minio:latest
    command: server /data --console-address ":9001"
    environment:
      MINIO_ROOT_USER: minioadmin
      MINIO_ROOT_PASSWORD: minioadmin
    ports:
      - "9000:9000"
      - "9001:9001"
    volumes:
      - miniodata:/data

volumes:
  pgdata:
  miniodata:
```

### .env.example
```env
# Server
PORT=8080
ENVIRONMENT=development

# Database
DATABASE_URL=postgres://tennisapp:tennisapp@localhost:5432/tennisapp?sslmode=disable

# Redis
REDIS_URL=redis://localhost:6379

# JWT
JWT_SECRET=dev-secret-change-in-production
JWT_ACCESS_TTL=15m
JWT_REFRESH_TTL=720h

# SMS
SMS_PROVIDER=mock
SMS_API_URL=https://smsc.kz/sys/send.php
SMS_API_LOGIN=
SMS_API_PASSWORD=

# Storage
S3_ENDPOINT=http://localhost:9000
S3_ACCESS_KEY=minioadmin
S3_SECRET_KEY=minioadmin
S3_BUCKET=tennisapp
S3_PUBLIC_URL=http://localhost:9000/tennisapp

# Firebase
FIREBASE_CREDENTIALS=

# Sentry
SENTRY_DSN=

# Dev settings
DEV_OTP_CODE=1234
```

### Makefile
```makefile
.PHONY: dev build test lint migrate-up migrate-down sqlc seed

# Development
dev:
	cd apps/backend && air

build:
	cd apps/backend && go build -o bin/server cmd/server/main.go

# Testing
test:
	cd apps/backend && go test ./... -v -count=1

test-coverage:
	cd apps/backend && go test ./... -coverprofile=coverage.out
	cd apps/backend && go tool cover -html=coverage.out

lint:
	cd apps/backend && golangci-lint run ./...

# Database
migrate-up:
	cd apps/backend && migrate -path migrations -database "$(DATABASE_URL)" up

migrate-down:
	cd apps/backend && migrate -path migrations -database "$(DATABASE_URL)" down 1

migrate-create:
	cd apps/backend && migrate create -ext sql -dir migrations -seq $(NAME)

# Code generation
sqlc:
	cd apps/backend && sqlc generate

# Data
seed:
	cd apps/backend && go run scripts/seed/main.go

# Docker
infra-up:
	docker-compose up -d

infra-down:
	docker-compose down

infra-logs:
	docker-compose logs -f

# Full setup
setup: infra-up migrate-up seed
	@echo "Ready! Run 'make dev' to start server"
```

---

## 3. Production — Backend (Railway)

### Deployment
```
GitHub push to main → Railway auto-deploy
```

### railway.toml
```toml
[build]
builder = "dockerfile"
dockerfilePath = "apps/backend/Dockerfile"

[deploy]
startCommand = "/app/server"
healthcheckPath = "/health"
healthcheckTimeout = 5
restartPolicyType = "on_failure"
```

### Dockerfile (apps/backend/Dockerfile)
```dockerfile
FROM golang:1.22-alpine AS builder
WORKDIR /build
COPY apps/backend/go.mod apps/backend/go.sum ./
RUN go mod download
COPY apps/backend/ .
RUN CGO_ENABLED=0 GOOS=linux go build -o server cmd/server/main.go

FROM alpine:3.19
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app
COPY --from=builder /build/server .
COPY apps/backend/migrations ./migrations
EXPOSE 8080
CMD ["./server"]
```

### Railway Environment Variables
```
DATABASE_URL        → Railway Postgres (auto-provisioned)
REDIS_URL           → Upstash Redis URL
JWT_SECRET          → Generate: openssl rand -hex 32
JWT_ACCESS_TTL      → 15m
JWT_REFRESH_TTL     → 720h
SMS_PROVIDER        → smsc
SMS_API_LOGIN       → xxx
SMS_API_PASSWORD    → xxx
S3_ENDPOINT         → Cloudflare R2 endpoint
S3_ACCESS_KEY       → R2 key
S3_SECRET_KEY       → R2 secret
S3_BUCKET           → tennisapp
S3_PUBLIC_URL       → https://cdn.tennisapp.kz
FIREBASE_CREDENTIALS → Base64-encoded service account JSON
SENTRY_DSN          → Sentry DSN
ENVIRONMENT         → production
PORT                → 8080
```

---

## 4. Production — Web (Vercel)

### vercel.json (apps/web-admin/)
```json
{
  "buildCommand": "npm run build",
  "outputDirectory": "dist",
  "framework": "vite",
  "rewrites": [{ "source": "/(.*)", "destination": "/index.html" }]
}
```

Deploy: GitHub push → Vercel auto-deploy.

---

## 5. Production — Mobile (EAS Build)

### eas.json (apps/mobile/)
```json
{
  "cli": { "version": ">= 5.0.0" },
  "build": {
    "development": {
      "developmentClient": true,
      "distribution": "internal"
    },
    "preview": {
      "distribution": "internal",
      "channel": "preview"
    },
    "production": {
      "channel": "production",
      "autoIncrement": true
    }
  },
  "submit": {
    "production": {
      "ios": { "appleId": "xxx", "ascAppId": "xxx" },
      "android": { "serviceAccountKeyPath": "./google-play-key.json" }
    }
  }
}
```

Commands:
```bash
eas build --platform ios --profile production
eas build --platform android --profile production
eas submit --platform ios
eas submit --platform android
eas update --channel production --message "Bug fixes"  # OTA update
```

---

## 6. CI/CD — GitHub Actions

### .github/workflows/backend.yml
```yaml
name: Backend CI

on:
  push:
    branches: [main, develop]
    paths: ['apps/backend/**']
  pull_request:
    branches: [main]
    paths: ['apps/backend/**']

jobs:
  test:
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgres:16-alpine
        env:
          POSTGRES_DB: tennisapp_test
          POSTGRES_USER: test
          POSTGRES_PASSWORD: test
        ports: ['5432:5432']
      redis:
        image: redis:7-alpine
        ports: ['6379:6379']

    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with: { go-version: '1.22' }
      - run: cd apps/backend && go mod download
      - run: cd apps/backend && go vet ./...
      - run: cd apps/backend && go test ./... -v -race
      - run: cd apps/backend && go build ./...
```

### .github/workflows/mobile.yml
```yaml
name: Mobile CI

on:
  push:
    branches: [main]
    paths: ['apps/mobile/**']

jobs:
  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with: { node-version: 20 }
      - run: cd apps/mobile && npm ci
      - run: cd apps/mobile && npm run lint
      - run: cd apps/mobile && npx tsc --noEmit
```

---

## 7. Database Migrations

### Migration Files
```
apps/backend/migrations/
├── 000001_init_schema.up.sql
├── 000001_init_schema.down.sql
├── 000002_add_badges.up.sql
├── 000002_add_badges.down.sql
└── ...
```

### Rules
- NEVER edit existing applied migrations
- Create new migration for schema changes
- Up migration must be reversible by down migration
- Test both up AND down before committing
- Run on staging before production

### Production Migration
```bash
# Railway runs migrations automatically via start command:
# /app/server --migrate-up && /app/server
# OR manually:
railway run migrate -path migrations -database "$DATABASE_URL" up
```

---

## 8. Monitoring

| Tool | Purpose | Tier |
|------|---------|------|
| Sentry | Error tracking (backend + mobile) | Free |
| UptimeRobot | API health monitoring | Free |
| Railway Logs | Backend logs (stdout) | Included |
| Firebase Analytics | Mobile usage analytics | Free |
| Custom `/health` | Liveness + readiness probe | Built-in |

### Health Check
```go
// GET /health
{
  "status": "ok",
  "version": "1.0.0",
  "uptime": "24h15m",
  "database": "connected",
  "redis": "connected"
}
```

---

## 9. Backup Strategy

| Data | Method | Frequency |
|------|--------|-----------|
| PostgreSQL | Railway/Supabase auto-backup | Daily |
| Redis | Not backed up (ephemeral cache) | — |
| S3/R2 files | Cloudflare R2 built-in redundancy | Automatic |
| Code | GitHub | Every push |
| Secrets | 1Password / env vars in Railway | Manual |

---

## 10. Domain & SSL

```
tennisapp.kz              → Vercel (web landing, future)
api.tennisapp.kz           → Railway (backend)
admin.tennisapp.kz         → Vercel (web-admin)
superadmin.tennisapp.kz    → Vercel (web-superadmin)
cdn.tennisapp.kz           → Cloudflare R2 (files)
```

SSL: automatic via Railway (API) and Vercel (web). All HTTPS.
