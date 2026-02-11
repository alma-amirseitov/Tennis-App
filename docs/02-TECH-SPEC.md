# Technical Specification v2.0
## Теннисная платформа для Астаны

**Версия:** 2.0  
**Дата:** Февраль 2026  
**Статус:** Approved

---

## 1. Architecture Overview

### 1.1 High-Level Architecture

```
┌───────────────────────────────────────────────────────────┐
│                        Clients                             │
├──────────────┬───────────────┬────────────┬───────────────┤
│ iOS App      │ Android App   │ Web Admin  │ Web Superadmin│
│ (Expo/RN)    │ (Expo/RN)     │ (React)    │ (React)       │
└──────┬───────┴───────┬───────┴─────┬──────┴──────┬────────┘
       │               │             │             │
       └───────────────┼─────────────┼─────────────┘
                       │             │
                       ▼             ▼
              ┌─────────────┐ ┌──────────────┐
              │  REST API   │ │  WebSocket   │
              │  (JSON)     │ │  (Chat)      │
              └──────┬──────┘ └──────┬───────┘
                     │               │
                     ▼               ▼
              ┌────────────────────────────┐
              │       Go Backend           │
              │   (Monolith, Layered)      │
              │                            │
              │  Handler → Service → Repo  │
              └─────┬──────┬───────┬───────┘
                    │      │       │
         ┌──────────┘      │       └──────────┐
         ▼                 ▼                  ▼
┌──────────────┐  ┌──────────────┐  ┌──────────────┐
│ PostgreSQL   │  │    Redis     │  │ MinIO / S3   │
│ (Data)       │  │ (Cache,      │  │ (Files,      │
│              │  │  Sessions,   │  │  Photos)     │
│              │  │  WS pubsub)  │  │              │
└──────────────┘  └──────────────┘  └──────────────┘
                        │
                        ▼
              ┌──────────────────┐
              │  Firebase FCM    │
              │  (Push Notifs)   │
              └──────────────────┘
```

### 1.2 Design Principles
1. **Monolith First** — единый Go-сервис, разделение по слоям (не микросервисы)
2. **Layered Architecture** — Handler → Service → Repository (чистое разделение)
3. **API-First** — общий REST API для всех клиентов
4. **Cost-Effective** — минимум платных сервисов, open-source где возможно
5. **i18n from Day 1** — три языка в каждом компоненте
6. **Security by Default** — OTP, JWT, RBAC, rate limiting

---

## 2. Technology Stack

### 2.1 Backend

| Компонент | Технология | Версия | Обоснование |
|-----------|-----------|--------|-------------|
| Language | Go | 1.22+ | Производительность, простой деплой, concurrency |
| HTTP Router | Chi (go-chi/chi) | v5 | Лёгкий, middleware-friendly, idiomatic Go, stdlib-compatible |
| ORM | sqlc | latest | Type-safe SQL, генерация Go-кода из SQL |
| Migrations | golang-migrate | latest | SQL-файлы, версионирование |
| WebSocket | nhooyr.io/websocket | latest | Замена deprecated Gorilla |
| Auth (JWT) | golang-jwt/jwt | v5 | Стандарт, production-tested |
| Validation | go-playground/validator | v10 | Struct tag validation |
| Config | envconfig / viper | latest | Env-based config |
| Logging | slog (stdlib) | Go 1.22 | Structured logging, zero deps |
| Testing | stdlib + testify | latest | Table-driven tests |

**Почему не Gorilla WebSocket:** Archived / deprecated. nhooyr.io/websocket активно поддерживается.

**Почему sqlc вместо GORM:** Type-safe, нет magic, полный контроль над SQL, AI лучше генерирует plain SQL.

### 2.2 Mobile

| Компонент | Технология | Обоснование |
|-----------|-----------|-------------|
| Framework | React Native + Expo | iOS + Android, OTA updates, EAS Build |
| Language | TypeScript | Strict mode |
| Navigation | Expo Router | File-based routing, deep links |
| State (server) | TanStack Query (React Query) | Caching, refetching, optimistic updates |
| State (client) | Zustand | Простой, lightweight |
| Forms | React Hook Form + Zod | Validation, performance |
| HTTP | Axios | Interceptors, retry |
| WebSocket | Native WebSocket API | Достаточно для чата |
| Push | expo-notifications + FCM | Expo managed |
| Secure Storage | react-native-keychain | Tokens, PIN hash |
| i18n | i18next + react-i18next | Зрелое решение, pluralization |
| Charts | react-native-chart-kit | Графики рейтинга |
| Maps | react-native-maps | Карта кортов |

**Не используем:** AsyncStorage для токенов (небезопасно), Redux (избыточно), React Native Paper (будем строить свою дизайн-систему для уникального look & feel).

### 2.3 Web (Admin + Superadmin)

| Компонент | Технология | Обоснование |
|-----------|-----------|-------------|
| Framework | React + Vite | Fast build, HMR |
| Language | TypeScript | Strict mode |
| Routing | React Router v6+ | Standard |
| State (server) | TanStack Query | Shared patterns with mobile |
| State (client) | Zustand | Consistency с мобайлом |
| UI Kit | Shadcn/UI + Tailwind | Гибкость, красивые компоненты, не vendor lock-in |
| Tables | TanStack Table | Sorting, filtering, pagination |
| Charts | Recharts | Дашборды, графики |
| Forms | React Hook Form + Zod | Consistency с мобайлом |

### 2.4 Infrastructure

| Компонент | Сервис | Tier | Стоимость |
|-----------|--------|------|-----------|
| Backend hosting | Railway / Render | Hobby | $5–10/мес |
| PostgreSQL | Supabase / Railway | Free → Paid | $0–25/мес |
| Redis | Upstash | Free | $0 |
| File storage | Cloudflare R2 / MinIO | Free 10GB | $0 |
| CDN | Cloudflare | Free | $0 |
| Web hosting | Vercel | Free | $0 |
| Push | Firebase FCM | Free | $0 |
| SMS | smsc.kz / mobizon | Pay per SMS | ~$0.02/SMS |
| Monitoring | Sentry | Free tier | $0 |
| CI/CD | GitHub Actions | Free | $0 |
| **Total (MVP)** | | | **$5–35/мес + SMS** |

---

## 3. Backend Architecture

### 3.1 Layered Structure

```
cmd/server/main.go          ← Entry point
│
internal/
├── config/                  ← Environment config
├── handler/                 ← HTTP handlers (thin, validation + response)
│   ├── middleware/           ← Auth, RBAC, rate limit, CORS
│   └── dto/                 ← Request/Response DTOs
├── service/                 ← Business logic (core of the app)
├── repository/              ← Database queries (sqlc generated + custom)
├── model/                   ← Domain models
├── ws/                      ← WebSocket hub + clients
└── pkg/                     ← Shared: jwt, sms, firebase, elo, storage
```

### 3.2 API Design

**Base URL:** `https://api.tennisapp.kz/v1`

**Аутентификация:** Bearer JWT в Authorization header
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
```

**Формат ошибок:**
```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Номер телефона недействителен",
    "details": [
      { "field": "phone", "message": "Неверный формат" }
    ]
  }
}
```

**Пагинация:**
```json
{
  "data": [...],
  "pagination": {
    "page": 1,
    "per_page": 20,
    "total": 142,
    "total_pages": 8
  }
}
```

### 3.3 Ключевые API-модули

| Модуль | Endpoints | Описание |
|--------|-----------|----------|
| Auth | 5 | OTP send, verify, refresh, PIN set/verify |
| Users | 8 | CRUD, profile, public profile, search, friends |
| Events | 12 | CRUD, join, leave, lifecycle, results, confirm |
| Communities | 10 | CRUD, join, leave, members, roles, verify |
| Chat | 6 | List chats, messages, send, mark read |
| Notifications | 4 | List, mark read, settings, count |
| Rating | 4 | Global, community, history, recalculate |
| Courts | 4 | List, detail, map data, CRUD (admin) |
| Posts | 5 | CRUD, feed, like |
| Admin | 8 | Dashboard, members, events, results, export |
| Superadmin | 6 | Verify, ban, courts CRUD, stats |
| **Total** | **~72** | |

### 3.4 Authentication Flow

```
1. POST /auth/otp/send       { phone: "+77071234567" }
   → SMS отправлен, возвращает session_id

2. POST /auth/otp/verify      { session_id, code: "1234" }
   → Если новый пользователь: { is_new: true, temp_token }
   → Если существующий: { access_token, refresh_token, user }

3. POST /auth/profile/setup   { name, surname, gender, birth_year, city }
   (только для новых, Authorization: temp_token)
   → { access_token, refresh_token, user }

4. POST /auth/quiz/submit     { answers: [...] }
   → { level: "Любитель", ntrp: 3.0 }

5. POST /auth/pin/set         { pin: "1234" }
   → OK

6. POST /auth/pin/verify      { pin: "1234" }
   → { access_token, refresh_token }

7. POST /auth/refresh          { refresh_token }
   → { access_token, refresh_token }
```

### 3.5 WebSocket Protocol (Chat)

**Connection:** `wss://api.tennisapp.kz/ws?token={jwt}`

**Message format:**
```json
{
  "type": "message",
  "chat_id": "uuid",
  "content": "Привет!",
  "reply_to": null
}
```

**Server events:**
```json
{ "type": "message", "data": { ... } }
{ "type": "typing", "data": { "chat_id": "...", "user_id": "..." } }
{ "type": "read", "data": { "chat_id": "...", "last_read_at": "..." } }
{ "type": "notification", "data": { ... } }
```

**Reconnection:** Client implements exponential backoff (1s, 2s, 4s, 8s, max 30s).

### 3.6 Rating Algorithm (ELO)

```go
// Упрощённый ELO для тенниса
func CalculateNewRating(winnerRating, loserRating float64, kFactor int) (float64, float64) {
    expectedWinner := 1.0 / (1.0 + math.Pow(10, (loserRating-winnerRating)/400.0))
    expectedLoser := 1.0 - expectedWinner
    
    newWinnerRating := winnerRating + float64(kFactor)*(1.0-expectedWinner)
    newLoserRating := loserRating + float64(kFactor)*(0.0-expectedLoser)
    
    return newWinnerRating, newLoserRating
}

// K-фактор зависит от количества игр:
// < 10 игр: K=40 (быстрая калибровка)
// 10-30 игр: K=32
// > 30 игр: K=24 (стабильный)
```

---

## 4. Security

### 4.1 Authentication & Authorization
- **SMS OTP:** 4-значный код, TTL 5 мин, макс. 5 попыток
- **JWT Access Token:** TTL 15 мин, HS256
- **Refresh Token:** TTL 30 дней, stored server-side (Redis), one-time use (rotation)
- **PIN:** bcrypt hash в react-native-keychain, 3 попытки
- **RBAC:** Middleware проверяет роль пользователя в контексте сообщества

### 4.2 Rate Limiting
| Endpoint | Лимит |
|----------|-------|
| SMS отправка | 3/час, 10/день на номер |
| OTP проверка | 5 попыток за сессию |
| PIN проверка | 3 попытки |
| API (общий) | 100 req/min на пользователя |
| Регистрация | 3 аккаунта/час с одного устройства |
| Поиск | 30 req/min |
| Создание ивентов | 10/час |

### 4.3 Data Protection
- HTTPS everywhere (TLS 1.3)
- Prepared statements (SQL injection protection)
- Input validation на backend (не доверяем клиенту)
- Sanitization HTML в постах
- Secrets в environment variables
- Номера телефонов маскируются в логах

---

## 5. Performance

### 5.1 Backend Targets
| Метрика | Цель |
|---------|------|
| API response time (p95) | < 200ms |
| WebSocket message delivery | < 100ms |
| Database query time (p95) | < 50ms |
| Concurrent WebSocket connections | 500+ |
| Events feed load time | < 300ms |

### 5.2 Mobile Targets
| Метрика | Цель |
|---------|------|
| App launch → interactive | < 2s |
| Screen transition | < 300ms |
| Infinite scroll (load more) | < 500ms |
| Image load (cached) | < 100ms |
| Offline → online sync | < 3s |

### 5.3 Optimization Strategies

**Backend:**
- PostgreSQL indexes на все фильтруемые колонки
- Redis cache для hot data (рейтинги, профили)
- Connection pooling (pgx pool)
- Pagination everywhere (никогда не SELECT * без LIMIT)
- Background jobs для тяжёлых операций (рейтинг, push)

**Mobile:**
- React Query cache (staleTime 5min для справочников)
- FlatList с windowSize и maxToRenderPerBatch
- Image caching (expo-image)
- Lazy loading для тяжёлых экранов
- Debounce для поиска (300ms)
- Skeleton screens вместо spinners

**Web:**
- Code splitting per route
- Lazy loading для тяжёлых компонентов (charts, tables)
- React Query для серверных данных
- Virtualized tables для больших списков

---

## 6. Scalability Plan

| Фаза | MAU | Архитектура | Стоимость |
|------|-----|-------------|-----------|
| MVP | 0–500 | Single server + managed Postgres | $5–35/мес |
| Growth | 500–2,000 | Vertical scaling + Redis cache + CDN | $50–100/мес |
| Scale | 2,000–10,000 | Load balancer + 2-3 instances + read replicas | $150–300/мес |
| Expansion | 10,000+ | K8s + microservices extraction + sharding по городам | $500+/мес |

---

## 7. Monitoring & Observability

### MVP
- **Errors:** Sentry (free tier) — mobile + backend
- **Logs:** Structured logging (slog) → stdout → Railway logs
- **Uptime:** UptimeRobot (free) — API healthcheck
- **Analytics:** Firebase Analytics (mobile) + custom events

### Post-MVP
- Prometheus + Grafana для метрик
- Distributed tracing (OpenTelemetry)
- Custom dashboards для бизнес-метрик

---

## 8. Development Workflow

### Git Strategy
- `main` — production
- `develop` — integration
- `feature/*` — feature branches
- PR review required (если команда > 1)

### CI/CD
- **GitHub Actions:**
  - Backend: lint → test → build → deploy (Railway)
  - Web: lint → build → deploy (Vercel)
  - Mobile: lint → build → EAS Build (на PR в main)

### Environments
| Env | Backend | Database | Назначение |
|-----|---------|----------|-----------|
| local | localhost:8080 | Docker PostgreSQL | Разработка |
| staging | staging.api.tennisapp.kz | Staging DB | Тестирование |
| production | api.tennisapp.kz | Production DB | Продакшн |

---

## 9. Third-Party Integrations

| Сервис | Назначение | API | Fallback |
|--------|-----------|-----|----------|
| SMS Provider (TBD) | OTP отправка | HTTP API | Retry + queue |
| Firebase FCM | Push notifications | Firebase Admin SDK | In-app notifications |
| Google Maps API | Карта кортов | Maps SDK (mobile) | Static list |
| Cloudflare R2 | Файловое хранилище | S3-compatible | Local storage |
| Sentry | Error tracking | SDK | Логи |

---

## 10. Tech Debt (Conscious Compromises)

### MVP Compromises
1. **Чат в PostgreSQL** — не специализированная БД (достаточно для 500 MAU)
2. **Без image processing** — фото загружаются as-is (нет resize/compress)
3. **Минимальные тесты** — unit tests на service layer, без e2e
4. **Один сервер** — нет horizontal scaling
5. **Нет offline mode** — требуется интернет

### Planned Improvements (Phase 2)
- Image processing pipeline (compress + resize + thumbnails)
- E2E tests (Detox для mobile, Cypress для web)
- Background job queue (для рейтинга, push, email)
- Database read replicas
- Comprehensive monitoring

---

## 11. Decision Log

| # | Решение | Выбор | Альтернатива | Причина |
|---|---------|-------|-------------|---------|
| 1 | Backend language | Go | Node.js | Производительность, goroutines для WS |
| 2 | Mobile framework | React Native (Expo) | Flutter | React ecosystem, OTA updates |
| 3 | Database | PostgreSQL | MongoDB | Relational data, ACID, ELO calculations |
| 4 | ORM | sqlc | GORM | Type-safe, no magic, full SQL control |
| 5 | WebSocket | nhooyr.io/websocket | Gorilla | Gorilla deprecated |
| 6 | Mobile nav | Expo Router | React Navigation | File-based, deep links, newer |
| 7 | Server state | TanStack Query | SWR | Richer features, shared mobile+web |
| 8 | Web UI | Shadcn/UI | Ant Design / MUI | Flexibility, Tailwind, no vendor lock |
| 9 | Auth | SMS OTP + JWT | Email/password | KZ market, faster onboarding |
| 10 | Token storage | react-native-keychain | AsyncStorage | Security (hardware encryption) |
| 11 | Architecture | Monolith | Microservices | Solo/small team, simpler |
| 12 | Hosting | Railway | AWS/GCP | Simplicity, affordable, good DX |

---

**Статус:** APPROVED  
**Связанные документы:** PRD.md, MVP-SCOPE.md, DEV-FILE-STRUCTURE.md, database-schema.sql
