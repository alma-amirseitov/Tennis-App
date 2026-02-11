# 🎾 Tennis App — Полный гайд по разработке
## Пошаговая инструкция с промптами для AI-агентов

---

# ⚡ ПЕРЕД СТАРТОМ: Setup (30 минут)

## Шаг 0.1 — Создать GitHub репозиторий

```bash
# На GitHub: Create new repository → tennisapp → Private
git clone https://github.com/{username}/tennisapp.git
cd tennisapp
```

## Шаг 0.2 — Скопировать документацию

```bash
mkdir -p docs/mobile-app docs/web-platform

# Скопировать все файлы из project-v2/ в docs/
# Корневые файлы:
cp 00-PROJECT-OVERVIEW.md docs/
cp 01-PRD.md docs/
cp 02-TECH-SPEC.md docs/
cp 03-DATABASE-SCHEMA.sql docs/
cp 04-API-SPECIFICATION.md docs/
cp 05-USER-STORIES.md docs/
cp 07-CODING-CONVENTIONS.md docs/
cp 08-ELO-ALGORITHM.md docs/
cp 09-INTEGRATIONS.md docs/
cp 10-DEPLOYMENT.md docs/
cp 11-SECURITY.md docs/
cp 12-TESTING-STRATEGY.md docs/
cp 13-DESIGN-SYSTEM.md docs/
cp 14-ERROR-CODES.md docs/
cp 15-I18N-GUIDE.md docs/
cp 16-ANALYTICS.md docs/
cp 17-SEED-DATA.md docs/
cp 18-COMPETITIVE-ANALYSIS.md docs/
cp AI-AGENT-WORKFLOW.md docs/
cp MVP-SCOPE.md docs/
cp DEV-FILE-STRUCTURE.md docs/

# Мобайл и веб спеки:
cp mobile-app/*.md docs/mobile-app/
cp web-platform/*.md docs/web-platform/

# Агентские файлы — В КОРЕНЬ репо (не в docs/):
cp CLAUDE.md ./
cp .cursorrules ./

# Прототип:
cp 06-PROTOTYPE.jsx docs/
```

## Шаг 0.3 — Первый коммит

```bash
git add .
git commit -m "docs: full project documentation v2.1"
git push origin main
```

## Шаг 0.4 — Настроить инструменты

1. **Claude Code** — установить: `npm install -g @anthropic-ai/claude-code`
2. **Cursor** — скачать с cursor.com, открыть папку tennisapp/
3. **Docker Desktop** — установить и запустить

---

# 🏗️ SPRINT 1: Foundation (Неделя 1-2)

> **Цель:** Backend skeleton + инфра + mobile shell
> **Агент:** Claude Code (backend), Cursor (mobile)

---

## Задача 1.1 — Monorepo + Docker + Makefile

**Агент:** Claude Code
**Где:** В терминале, в папке tennisapp/

### Промпт:

```
Прочитай CLAUDE.md и docs/10-DEPLOYMENT.md, docs/DEV-FILE-STRUCTURE.md.

Создай скелет monorepo проекта:

1. Структура папок:
   apps/backend/
   apps/mobile/
   apps/web-admin/
   apps/web-superadmin/
   packages/shared-types/

2. docker-compose.yml в корне:
   - PostgreSQL 16 alpine (порт 5432, db=tennisapp, user=tennisapp, pass=tennisapp, volume)
   - Redis 7 alpine (порт 6379)
   - MinIO (порты 9000+9001, user=minioadmin, pass=minioadmin, volume)

3. apps/backend/:
   - go.mod (module github.com/{username}/tennisapp/apps/backend, go 1.22)
   - .env.example со всеми переменными из docs/10-DEPLOYMENT.md
   - Makefile со всеми командами: dev, build, test, lint, migrate-up, migrate-down, migrate-create, sqlc, seed

4. .gitignore для Go + Node + .env + IDE + OS

5. .editorconfig (2 spaces для TS/JSON, tabs для Go)

6. README.md с описанием проекта и командами запуска

НЕ инициализируй Expo и React проекты — только структуру папок.
```

### Проверка:

```bash
docker-compose up -d    # PostgreSQL, Redis, MinIO запускаются
docker-compose ps       # Все 3 сервиса running
```

---

## Задача 1.2 — Database migrations

**Агент:** Claude Code

### Промпт:

```
Прочитай docs/03-DATABASE-SCHEMA.sql полностью.

Создай database migrations:

1. Установи golang-migrate в go.mod
2. Создай apps/backend/migrations/000001_init_schema.up.sql
   — Скопируй ВСЮ схему из docs/03-DATABASE-SCHEMA.sql: все 21 таблицу, views, triggers, indexes, constraints
3. Создай apps/backend/migrations/000001_init_schema.down.sql
   — DROP TABLE IF EXISTS CASCADE для всех таблиц в обратном порядке зависимостей
   — DROP VIEW, DROP FUNCTION, DROP TYPE

4. Проверь что make migrate-up создаёт все таблицы без ошибок
5. Проверь что make migrate-down удаляет всё
6. Проверь что повторный make migrate-up после down работает
```

### Проверка:

```bash
make migrate-up
# Подключись к БД и проверь:
docker exec -it tennisapp-postgres psql -U tennisapp -c "\dt"
# Должно показать 21 таблицу
```

---

## Задача 1.3 — Go backend skeleton

**Агент:** Claude Code

### Промпт:

```
Прочитай CLAUDE.md (секции Architecture, Backend Structure, Commands).
Прочитай docs/07-CODING-CONVENTIONS.md (Go секция).
Прочитай docs/11-SECURITY.md (Middleware chain).

Создай Go backend скелет:

1. cmd/server/main.go
   - Загрузка конфига из ENV
   - Подключение к PostgreSQL (pgx pool)
   - Подключение к Redis (go-redis)
   - Chi router с middleware chain: Logger → Recovery → CORS → RequestID
   - Graceful shutdown (SIGINT/SIGTERM)
   - Логирование через slog (structured, JSON в production)

2. internal/config/config.go
   - Struct с envconfig тегами
   - Все поля из .env.example
   - Validation (required поля)

3. internal/handler/health.go
   - GET /health → {"status":"ok","version":"0.1.0","database":"connected","redis":"connected"}
   - Реально проверяет pg.Ping() и redis.Ping()

4. internal/handler/router.go
   - Инициализация всех routes
   - Middleware chain
   - CORS: allow origins *, methods GET/POST/PUT/PATCH/DELETE, headers Authorization/Content-Type

5. internal/handler/middleware/
   - logger.go — slog request logging (method, path, status, duration)
   - recovery.go — panic recovery → 500 + log
   - requestid.go — X-Request-ID header
   - cors.go — CORS middleware

Убедись что make dev запускает сервер и curl http://localhost:8080/health отвечает 200.
```

### Проверка:

```bash
make dev
# В другом терминале:
curl http://localhost:8080/health
# {"status":"ok","version":"0.1.0","database":"connected","redis":"connected"}
```

---

## Задача 1.4 — sqlc setup + users queries

**Агент:** Claude Code

### Промпт:

```
Прочитай CLAUDE.md (секция Database Rules, sqlc).
Прочитай docs/03-DATABASE-SCHEMA.sql (таблица users).
Прочитай docs/04-API-SPECIFICATION.md (секция 3 — Users).

Настрой sqlc:

1. apps/backend/sqlc.yaml
   - version: "2"
   - engine: postgresql
   - queries: internal/repository/queries/
   - schema: migrations/
   - gen go: package: repository, out: internal/repository

2. internal/repository/queries/users.sql:
   - CreateUser (phone, всё остальное опционально)
   - GetUserByID
   - GetUserByPhone
   - UpdateUser (все поля кроме id и phone)
   - SearchUsers (с фильтрами: level, district, gender, name search через pg_trgm)

3. Запусти make sqlc — сгенерируй Go код
4. Убедись что код компилируется без ошибок

ВАЖНО: в sqlc queries всегда перечисляй колонки явно, НИКОГДА не используй SELECT *.
ВАЖНО: всегда добавляй WHERE deleted_at IS NULL.
```

---

## Задача 1.5 — Expo project + navigation

**Агент:** Cursor
**Где:** Открой apps/mobile/ в Cursor

### Промпт:

```
Прочитай .cursorrules в корне проекта.
Прочитай docs/mobile-app/01-app-structure.md.
Прочитай docs/13-DESIGN-SYSTEM.md.

Инициализируй Expo проект в apps/mobile/:

1. npx create-expo-app@latest . --template blank-typescript
2. Установи зависимости:
   - expo-router, expo-linking, expo-constants
   - @tanstack/react-query
   - zustand
   - i18next, react-i18next, expo-localization
   - react-native-keychain
   - axios
   - react-hook-form, @hookform/resolvers, zod
   - expo-haptics

3. Настрой Expo Router (file-based routing):
   src/app/
   ├── _layout.tsx          — Root layout (QueryClientProvider, i18n)
   ├── (auth)/
   │   ├── _layout.tsx      — Auth stack layout
   │   ├── index.tsx         — Phone input screen (placeholder)
   │   ├── otp.tsx           — OTP screen (placeholder)
   │   ├── profile-setup.tsx — Profile setup (placeholder)
   │   └── quiz.tsx          — Quiz (placeholder)
   ├── (tabs)/
   │   ├── _layout.tsx      — Tab navigator (5 табов)
   │   ├── index.tsx         — Home tab (placeholder)
   │   ├── players.tsx       — Players tab (placeholder)
   │   ├── events/
   │   │   └── index.tsx     — Events tab (placeholder)
   │   ├── communities/
   │   │   └── index.tsx     — Communities tab (placeholder)
   │   └── profile/
   │       └── index.tsx     — Profile tab (placeholder)

4. Tab bar:
   - 5 табов: 🏠 Главная, 👥 Игроки, 🎾 Ивенты, 🏛 Сообщества, 👤 Профиль
   - Высота 80px, цвета из design system (primary green, neutral gray)

5. Каждый placeholder экран: название таба по центру

Убедись что npx expo start запускается и показывает 5 табов.
```

---

## Задача 1.6 — i18n setup

**Агент:** Cursor

### Промпт:

```
Прочитай docs/15-I18N-GUIDE.md полностью.

Настрой i18n в apps/mobile/:

1. src/shared/i18n/index.ts — инициализация i18next
2. src/shared/i18n/locales/ru.json — скопируй ВСЕ ключи из docs/15-I18N-GUIDE.md секция 5
3. src/shared/i18n/locales/kk.json — скопируй русские значения (переводы позже)
4. src/shared/i18n/locales/en.json — english values

5. Язык по умолчанию: ru
6. Определение по locale устройства
7. Обнови все placeholder экраны — используй t() вместо хардкода

ПРАВИЛО: Ни одного user-visible string напрямую. Всё через t().
```

---

## Задача 1.7 — Design system components

**Агент:** Cursor

### Промпт:

```
Прочитай docs/13-DESIGN-SYSTEM.md полностью.
Прочитай .cursorrules (секция Styling Rules).

Создай базовую дизайн-систему:

1. src/shared/theme/
   - colors.ts — ВСЕ цвета из Design System (primary, accent, neutrals, semantic, status)
   - typography.ts — fontSize (11-36), fontWeight (400-800), textStyles (h1-h4, body, caption, rating, badge)
   - spacing.ts — scale: xs=4, sm=8, md=12, base=16, lg=20, xl=24, 2xl=32, 3xl=40, 4xl=48, 5xl=56, 6xl=64
   - radius.ts — sm=8, md=12, lg=16, xl=20, pill=100, full=9999
   - shadows.ts — sm, md, lg (с shadowColor, shadowOffset, elevation)
   - index.ts — реэкспорт всего

2. src/shared/ui/
   - Button.tsx
     Variants: primary (green bg), secondary (white bg + border), outline (transparent), small
     States: default, pressed (scale 0.97), disabled (opacity 0.5), loading (ActivityIndicator)
     Height: 52 (default), 36 (small)
     Props: variant, title, onPress, disabled, loading, icon

   - Input.tsx
     States: default, focused (primary border), error (danger border), disabled
     Height: 52, radius: 12
     Props: label, placeholder, value, onChangeText, error, secureTextEntry

   - Card.tsx
     Flat card: bg white, radius 16, padding 16, border 1px neutral-200
     Props: children, style, onPress?

   - Avatar.tsx
     Sizes: xs=24, sm=32, md=40, lg=56, xl=80
     Circle shape, initials fallback (first letter of name), online indicator (green dot)
     Props: uri, name, size, showOnline

   - Badge.tsx
     Variants: primary, success, warning, danger, info, muted
     Background: color + 15% opacity, text: full color
     Props: variant, text

   - Skeleton.tsx
     Pulse animation (1.5s), customizable width/height/radius
     Props: width, height, radius

ВСЕ компоненты используют только tokens из theme/, НИКАКИХ raw hex values.
TypeScript strict — все props typed, никакого any.
```

---

## Задача 1.8 — CI/CD

**Агент:** Claude Code

### Промпт:

```
Прочитай docs/10-DEPLOYMENT.md (секция CI/CD).

Создай GitHub Actions:

1. .github/workflows/backend.yml
   - Trigger: push/PR на main и develop, paths: apps/backend/**
   - Services: postgres:16-alpine (port 5432), redis:7-alpine (port 6379)
   - Steps: checkout → setup-go 1.22 → go mod download → go vet → go test -race ./... → go build ./cmd/server

2. .github/workflows/mobile.yml
   - Trigger: push/PR на main и develop, paths: apps/mobile/**
   - Steps: checkout → setup-node 20 → npm ci → npx tsc --noEmit → npx eslint .
```

---

## Коммит Sprint 1:

```bash
git add .
git commit -m "feat(sprint-1): foundation - backend skeleton, mobile shell, infra"
git push origin main
```

---

# 🔐 SPRINT 2: Auth (Неделя 3-4)

> **Цель:** Полный auth flow: OTP → JWT → Profile → Quiz
> **Агент:** Claude Code (backend), Cursor (mobile)

---

## Задача 2.1 — OTP send + verify endpoints

**Агент:** Claude Code

### Промпт:

```
Прочитай docs/04-API-SPECIFICATION.md (секция 1 — Auth, endpoints 1.1 и 1.2).
Прочитай docs/11-SECURITY.md (секции Authentication Flow, OTP).
Прочитай docs/14-ERROR-CODES.md (секции OTP Errors, Auth Errors).
Прочитай docs/05-USER-STORIES.md (AUTH-1 и AUTH-2 — полные acceptance criteria).

Реализуй OTP auth flow:

1. internal/service/auth.go — AuthService
   - SendOTP(ctx, phone string) → (sessionID string, err error)
     · Валидация телефона: ^\+7[0-9]{10}$
     · Rate limit проверка: Redis INCR sms_rate:{phone}:hour (limit 3, TTL 3600)
     · Rate limit проверка: Redis INCR sms_rate:{phone}:day (limit 10, TTL 86400)
     · Генерация 4-digit кода (crypto/rand)
     · Сохранение в Redis: otp:{session_id} → {phone, code, attempts:0}, TTL 300
     · Отправка SMS (mock в dev: логируем код, не отправляем)
   
   - VerifyOTP(ctx, sessionID, code string) → (result *AuthResult, err error)
     · Загрузка сессии из Redis: otp:{session_id}
     · Если нет → ErrOTPSessionExpired
     · Если attempts >= 5 → delete session → ErrOTPMaxAttempts
     · Если код неверный → INCR attempts → ErrOTPInvalidCode
     · Если код верный → delete session
     · Ищем user по phone в БД
     · Если новый: CreateUser → return {is_new_user: true, temp_token}
     · Если существующий: return {is_new_user: false, access_token, refresh_token}

2. internal/service/token.go — TokenService
   - GenerateAccessToken(userID, role) → JWT HS256, TTL 15 min, claims: sub, role, iat, exp
   - GenerateRefreshToken(userID) → JWT, TTL 30 days, jti=uuid, save jti in Redis
   - ValidateToken(tokenString) → claims, error

3. internal/handler/auth.go — HTTP handlers
   - POST /v1/auth/otp/send → {session_id, expires_in: 300}
   - POST /v1/auth/otp/verify → {is_new_user, access_token?, refresh_token?, temp_token?, user}

4. internal/handler/dto/auth.go — request/response structs с validator tags

5. Тесты:
   - service/auth_test.go: valid phone, invalid phone, rate limit, correct code, wrong code, max attempts

Используй коды ошибок ТОЛЬКО из docs/14-ERROR-CODES.md. В dev-режиме OTP всегда = 1234.
```

---

## Задача 2.2 — Refresh token + Auth middleware

**Агент:** Claude Code

### Промпт:

```
Прочитай docs/05-USER-STORIES.md (AUTH-3 и AUTH-4 — acceptance criteria).
Прочитай docs/11-SECURITY.md (секции JWT Tokens, Refresh Token Rotation, Middleware Chain).

Реализуй:

1. POST /v1/auth/refresh
   - Принимает {refresh_token}
   - Валидирует JWT подпись и expiration
   - Проверяет jti в Redis (refresh:{jti})
   - Если jti не найден → TOKEN_REVOKED (possible compromise) + revoke ALL user refresh tokens
   - Удаляет старый jti (one-time use)
   - Генерирует новую пару access + refresh tokens
   - Сохраняет новый jti в Redis

2. Auth middleware (internal/handler/middleware/auth.go)
   - Извлекает Bearer token из Authorization header
   - Валидирует JWT
   - Инжектит user_id и role в context
   - Helpers: GetUserID(ctx), GetUserRole(ctx)
   - Public routes (без middleware): /health, /v1/auth/otp/*, /v1/auth/refresh

3. Rate limit middleware (internal/handler/middleware/ratelimit.go)
   - Redis sliding window
   - Configurable per-route
   - Headers: X-RateLimit-Limit, X-RateLimit-Remaining, X-RateLimit-Reset
   - 429 + Retry-After при превышении

4. Обнови router.go:
   - Public group (без auth): health, auth endpoints
   - Protected group (с auth middleware): всё остальное

5. Тесты: valid refresh, expired refresh, reused refresh, valid auth header, expired token, no header
```

---

## Задача 2.3 — Profile setup + Quiz endpoints

**Агент:** Claude Code

### Промпт:

```
Прочитай docs/05-USER-STORIES.md (AUTH-9 и AUTH-11 — acceptance criteria).
Прочитай docs/04-API-SPECIFICATION.md (секция 1, endpoints 1.4 и секция 2 Quiz).

Реализуй:

1. POST /v1/auth/profile/setup
   - Protected (temp_token)
   - Body: {first_name, last_name, gender, birth_year, city, district}
   - Валидация: имя 2-50, фамилия 2-50, gender in [male, female], birth_year 1940-2012
   - Обновляет users запись
   - Если профиль уже заполнен → PROFILE_ALREADY_SET
   - Возвращает access_token + refresh_token (upgrade from temp_token)

2. GET /v1/quiz
   - Возвращает hardcoded вопросы (3-5 штук):
     · "Как давно вы играете в теннис?" → Никогда / Меньше года / 1-3 года / 3+ лет
     · "Как часто вы играете?" → Редко / 1-2 раза в месяц / 1-2 раза в неделю / 3+ раз в неделю
     · "Ваш средний уровень соперников?" → Не знаю / Начинающие / Средние / Продвинутые
   - Каждый вариант имеет weight (0-4)

3. POST /v1/quiz
   - Body: {answers: [{question_id, answer_id}]}
   - Алгоритм: sum(weights) → NTRP level → initial rating
     · 0-3 → 2.0 (800) / 4-6 → 2.5 (950) / 7-9 → 3.0 (1100) / 10-12 → 3.5 (1250) / 13+ → 4.0 (1400)
   - Обновляет users.level и users.rating_score
   - Response: {level, ntrp, initial_rating}
```

---

## Задача 2.4 — Mobile: Auth screens

**Агент:** Cursor

### Промпт:

```
Прочитай docs/mobile-app/02-auth-onboarding.md полностью.
Прочитай docs/05-USER-STORIES.md (AUTH-6, AUTH-7, AUTH-8 — acceptance criteria).
Прочитай docs/13-DESIGN-SYSTEM.md (компоненты Button, Input).
Прочитай docs/14-ERROR-CODES.md (секция Frontend Handling).

Создай auth flow экраны:

1. src/app/(auth)/index.tsx — Экран ввода телефона
   - Маска +7 (XXX) XXX-XX-XX
   - Numeric keypad
   - Кнопка "Получить код" (disabled пока < 10 цифр)
   - При нажатии: POST /v1/auth/otp/send
   - Loading state на кнопке
   - Error: toast при rate limit / network error
   - При успехе: router.push('/otp', { session_id })
   - Использует компоненты из shared/ui (Button, Input)
   - Все строки через t()

2. src/app/(auth)/otp.tsx — Экран OTP
   - 4 ячейки, автофокус, автопереход
   - Paste support (clipboard 4 digits)
   - Таймер 60 сек → "Отправить повторно"
   - При заполнении: автоматический POST /v1/auth/otp/verify
   - Shake animation при неверном коде
   - is_new_user: true → router.replace('/profile-setup')
   - is_new_user: false → router.replace('/(tabs)')

3. src/shared/api/client.ts — Axios instance
   - baseURL из ENV (default http://localhost:8080)
   - Request interceptor: добавляет Bearer token
   - Response interceptor: при 401 → refresh → retry
   - Queue для параллельных запросов при refresh

4. src/shared/stores/auth.ts — Zustand store
   - State: {isAuthenticated, user, isLoading}
   - Actions: login(tokens), logout(), loadFromKeychain()
   - Tokens: react-native-keychain (set/get/reset)

5. src/shared/api/auth.ts — API функции
   - sendOTP(phone) → {session_id}
   - verifyOTP(session_id, code) → {is_new_user, tokens}
   - refreshToken(refresh_token) → {tokens}
```

---

## Задача 2.5 — Mobile: Profile setup + Quiz

**Агент:** Cursor

### Промпт:

```
Прочитай docs/05-USER-STORIES.md (AUTH-10 и AUTH-12 — acceptance criteria).
Прочитай docs/mobile-app/02-auth-onboarding.md (секции Profile Setup и Quiz).

Создай:

1. src/app/(auth)/profile-setup.tsx
   - Поля: Имя, Фамилия (Input), Пол (2 toggle buttons), Год рождения (picker), Район (dropdown)
   - React Hook Form + Zod schema validation
   - Inline errors под каждым полем
   - Кнопка "Продолжить" (disabled пока невалидно)
   - POST /v1/auth/profile/setup
   - При успехе: сохранить tokens → router.replace('/quiz')
   - Районы Астаны: Есильский, Алматинский, Сарыаркинский, Байконурский, Нуринский

2. src/app/(auth)/quiz.tsx
   - GET /v1/quiz → загрузить вопросы
   - По 1 вопросу на экране
   - Карточки вариантов (tap для выбора, зелёная обводка selected)
   - Progress bar сверху (1/5, 2/5...)
   - Кнопка "Далее" → следующий вопрос
   - На последнем → POST /v1/quiz
   - Экран результата: "Ваш уровень: Любитель (NTRP 3.0)" + "Начать"
   - Skip → router.replace('/(tabs)')

Все строки через t(). Все компоненты из shared/ui/.
```

---

## Коммит Sprint 2:

```bash
git add .
git commit -m "feat(sprint-2): auth flow - OTP, JWT, profile setup, quiz"
git push origin main
```

---

# 👥 SPRINT 3-4: Core — Users, Communities, Events (Неделя 5-8)

> **Цель:** Главная ценность — поиск и создание игр
> **Агент:** Claude Code (backend), Cursor (mobile)

---

## Задача 3.1 — Users backend (CRUD + search + avatar)

**Агент:** Claude Code

### Промпт:

```
Прочитай docs/04-API-SPECIFICATION.md (секция 3 — Users, все 8 endpoints).
Прочитай docs/09-INTEGRATIONS.md (секция Object Storage — для аватаров).
Прочитай docs/05-USER-STORIES.md (USR-1 через USR-4).

Реализуй Users module:

1. sqlc queries: GetUserByID, UpdateUser, SearchUsers (pg_trgm), GetUserStats

2. internal/service/user.go:
   - GetProfile(ctx, userID) — мой профиль (с communities, stats, badges)
   - UpdateProfile(ctx, userID, input) — обновить поля
   - GetPublicProfile(ctx, userID) — чужой профиль
   - SearchUsers(ctx, filters) — поиск с фильтрами (level, district, gender, name)
   - UploadAvatar(ctx, userID, file) — загрузка в MinIO/S3

3. internal/service/storage.go:
   - Upload(ctx, bucket, key, reader, contentType) → url
   - Delete(ctx, bucket, key)
   - MinIO client initialization
   - Bucket: avatars
   - File validation: max 5MB, jpeg/png/webp only, check magic bytes

4. Handlers:
   - GET /v1/users/me
   - PATCH /v1/users/me
   - POST /v1/users/me/avatar (multipart/form-data)
   - GET /v1/users/:id
   - GET /v1/users/search?name=&level_min=&level_max=&district=&gender=&page=&per_page=

Все endpoint'ы protected (auth middleware).
Пагинация: page + per_page (default 20, max 100).
```

---

## Задача 3.2 — Communities backend

**Агент:** Claude Code

### Промпт:

```
Прочитай docs/04-API-SPECIFICATION.md (секция 7 — Communities, все 12 endpoints).
Прочитай docs/11-SECURITY.md (секция Authorization RBAC — community roles).
Прочитай docs/05-USER-STORIES.md (COM-1 через COM-6).

Реализуй Communities module:

1. sqlc queries для communities, community_members

2. internal/service/community.go:
   - Create(ctx, userID, input) — создатель = owner, group/organizer = verified, club/league = pending
   - List(ctx, filters) — с пагинацией и фильтрами (type, search)
   - GetByID(ctx, communityID) — детали + member count + my_role
   - Join(ctx, userID, communityID) — open: сразу member, closed: pending
   - Leave(ctx, userID, communityID) — owner не может уйти
   - ListMembers(ctx, communityID, filters)
   - UpdateMemberRole(ctx, communityID, targetUserID, role) — только owner/admin
   - ReviewRequest(ctx, communityID, targetUserID, approve bool) — owner/admin/moderator

3. internal/handler/middleware/community_role.go:
   - RequireCommunityRole(roles ...string) — middleware
   - Проверяет роль в community_members для текущего community_id
   - Роли: owner > admin > moderator > member

4. Handlers: все 12 endpoints из API spec

RBAC правила строго по docs/11-SECURITY.md.
```

---

## Задача 3.3 — Events backend

**Агент:** Claude Code

### Промпт:

```
Прочитай docs/04-API-SPECIFICATION.md (секция 5 — Events, все 14 endpoints).
Прочитай docs/mobile-app/05-events-tab.md (для понимания бизнес-логики).
Прочитай docs/05-USER-STORIES.md (EVT-1 через EVT-7).

Реализуй Events module:

1. sqlc queries для events, event_participants

2. internal/service/event.go:
   - Create(ctx, userID, communityID, input) — все поля из wizard (8 шагов в 1 запросе)
   - List(ctx, filters) — фильтры: type, status, level, date_from, date_to, district, community_id
   - GetByID(ctx, eventID) — детали + participants list + my_status + can_join
   - Join(ctx, userID, eventID) — проверки: not full, right level, right status, not already joined
   - Leave(ctx, userID, eventID) — only if status open/filling
   - UpdateStatus(ctx, eventID, status) — lifecycle: draft→open→filling→full→in_progress→completed/cancelled
   - GetCalendar(ctx, userID, year, month) — ивенты grouped by day
   - GetMyEvents(ctx, userID, filter) — created / joined / past

3. Lifecycle validation:
   - open→filling: когда current_participants > 0
   - filling→full: когда current_participants == max_participants
   - Только создатель или admin может менять статус
   - cancelled: возможно из любого статуса кроме completed

4. Handlers: все 14 endpoints

Бизнес-правила join:
- EVENT_FULL: current >= max
- EVENT_CLOSED: status не open/filling
- EVENT_WRONG_LEVEL: user.ntrp < event.level_min || user.ntrp > event.level_max
- ALREADY_JOINED_EVENT
```

---

## Задача 3.4 — Mobile: Profile tab

**Агент:** Cursor

### Промпт:

```
Прочитай docs/mobile-app/07-profile-tab.md полностью.
Прочитай docs/13-DESIGN-SYSTEM.md.

Создай Profile tab:

1. src/app/(tabs)/profile/index.tsx — Мой профиль
   - Header: аватар (Avatar xl), имя, уровень (Badge), рейтинг
   - Quick stats: матчи, победы, win rate (3 карточки в ряд)
   - Секции: Мои сообщества, Достижения (badges placeholder), Друзья (placeholder)
   - Кнопка "Редактировать"

2. src/app/(tabs)/profile/edit.tsx — Редактирование
   - React Hook Form, те же поля что при setup + аватар (image picker)

3. src/app/(tabs)/profile/settings.tsx — Настройки
   - Язык (3 варианта), Уведомления, PIN, О приложении, Выход
   - "Выйти" → AlertDialog подтверждение → clear tokens → auth screen

4. src/features/profile/ — компоненты:
   - ProfileHeader.tsx
   - StatsCard.tsx
   - CommunitiesList.tsx

API hooks через TanStack Query:
- useProfile() → GET /v1/users/me
- useUpdateProfile() → PATCH /v1/users/me

Все строки через t(). Skeleton loading. Error state с retry.
```

---

## Задача 3.5 — Mobile: Players tab

**Агент:** Cursor

### Промпт:

```
Прочитай docs/mobile-app/04-players-tab.md полностью.

Создай Players tab:

1. src/app/(tabs)/players.tsx
   - Search bar сверху
   - Фильтры: Уровень (chip select), Район (dropdown), Пол (toggle)
   - FlatList с карточками игроков
   - Infinite scroll (TanStack Query useInfiniteQuery)
   - Pull-to-refresh
   - Empty state когда нет результатов

2. src/app/player/[id].tsx — Публичный профиль
   - Аналог своего профиля, но с кнопками: "В друзья", "Написать", "Позвать играть"

3. src/features/players/
   - PlayerCard.tsx — аватар, имя, NTRP badge, рейтинг, win rate
   - PlayerFilters.tsx — фильтры

API: useSearchPlayers(filters) → GET /v1/users/search
```

---

## Задача 3.6 — Mobile: Communities tab

**Агент:** Cursor

### Промпт:

```
Прочитай docs/mobile-app/06-communities-tab.md полностью.

Создай Communities:

1. src/app/(tabs)/communities/index.tsx — Список
   - Поиск + фильтр по типу
   - Карточки сообществ (лого, название, тип badge, members count)
   - "Мои сообщества" секция сверху
   - FlatList + infinite scroll

2. src/app/community/[id].tsx — Экран сообщества
   - Header: лого, название, тип, member count, кнопка Join/Joined
   - 6 внутренних табов (Material Top Tabs): Лента, Ивенты, Рейтинг, Участники, Чат, Фото
   - MVP: Ивенты и Участники табы полные, остальные placeholder

3. src/app/community/create.tsx — Создание
   - Форма: название, тип (4 карточки), описание, доступ (open/closed)

4. src/features/communities/
   - CommunityCard.tsx
   - CommunityHeader.tsx
   - MembersList.tsx

API: useCommunities(), useCommunity(id), useJoinCommunity(), useLeaveCommunity()
```

---

## Задача 3.7 — Mobile: Events tab (ключевой экран)

**Агент:** Cursor

### Промпт:

```
Прочитай docs/mobile-app/05-events-tab.md ПОЛНОСТЬЮ (это самый детальный spec — 402 строки).
Прочитай docs/13-DESIGN-SYSTEM.md (статусные цвета).

Создай Events:

1. src/app/(tabs)/events/index.tsx — Таб ивентов
   - 3 внутренних таба: Лента / Календарь / Мои
   - FAB кнопка "+" (создать ивент)

2. Лента ивентов:
   - Фильтры toolbar (тип, уровень, дата, район)
   - Карточки ивентов: тип icon, название, дата/время, место, spots (3/4), level badge, status badge
   - Статусные цвета: open=green, filling=blue, full=orange, completed=gray, cancelled=red
   - Infinite scroll

3. src/app/event/[id].tsx — Детали ивента
   - Header с status badge
   - Инфо: тип, формат, уровень, сеты, место (карта link), дата/время
   - Участники (аватары в ряд + имена)
   - Кнопка "Записаться" / "Вы записаны ✓" / "Мест нет"
   - Pull-to-refresh

4. src/app/event/create.tsx — Wizard создания (8 шагов)
   - Step 1: Тип (4 карточки с иконками)
   - Step 2: Формат (singles/doubles/mixed/team)
   - Step 3: Уровень (slider min-max)
   - Step 4: Кол-во участников, кол-во сетов
   - Step 5: Место (выбор корта — placeholder, потом карта)
   - Step 6: Дата + время (date/time picker)
   - Step 7: Описание + правила
   - Step 8: Review → "Создать"
   - Progress bar сверху
   - Анимация перехода между шагами

5. src/features/events/
   - EventCard.tsx
   - EventFilters.tsx
   - EventWizard/ (каждый шаг = отдельный компонент)

API: useEvents(filters), useEvent(id), useJoinEvent(), useLeaveEvent(), useCreateEvent()
```

---

## Коммит Sprint 3-4:

```bash
git add .
git commit -m "feat(sprint-3-4): users, communities, events - full CRUD + mobile screens"
git push origin main
```

---

# 🏆 SPRINT 5-6: Matches, Rating, Chat (Неделя 9-12)

> **Цель:** Результаты матчей, ELO рейтинг, чат

---

## Задача 5.1 — Matches + ELO backend

**Агент:** Claude Code

### Промпт:

```
Прочитай docs/04-API-SPECIFICATION.md (секция 6 — Matches).
Прочитай docs/08-ELO-ALGORITHM.md ПОЛНОСТЬЮ — там полный Go код.
Прочитай docs/05-USER-STORIES.md (MTH-1 через MTH-8).

Реализуй Matches + Rating:

1. internal/service/elo/elo.go
   — Скопируй реализацию из docs/08-ELO-ALGORITHM.md
   — Calculate() для singles, CalculateDoubles() для doubles
   — K-factors: 40 (<10 games), 32 (10-30), 24 (>30)

2. internal/service/match.go:
   - CreateMatch(ctx, eventID, player1ID, player2ID, matchType)
   - SubmitResult(ctx, matchID, submitterID, score) — score как JSONB [{p1: 6, p2: 4}, ...]
   - ConfirmResult(ctx, matchID, confirmerID) → trigger ELO calculation
   - DisputeResult(ctx, matchID, disputerID, reason)
   - GetMyMatches(ctx, userID, filters)

3. Confirm flow:
   - Player A submits score → status: pending_confirmation
   - Player B confirms → status: confirmed → ELO calculation
   - Player B disputes → status: disputed → admin resolves

4. After confirmation:
   - elo.Calculate(winner, loser) → new ratings
   - Update users.rating_score for both
   - Update player_stats_global (games, wins, losses, win_rate)
   - Insert into rating_history
   - Update community_members rating (if community event)

5. internal/service/rating.go:
   - GetGlobalLeaderboard(ctx, page)
   - GetCommunityLeaderboard(ctx, communityID, page)
   - GetMyRatingHistory(ctx, userID)
   - GetMyStats(ctx, userID)

6. Тесты для ELO: equal players, upset, expected, clamp, doubles

7. Handlers для всех match + rating endpoints
```

---

## Задача 5.2 — Chat + WebSocket backend

**Агент:** Claude Code

### Промпт:

```
Прочитай docs/04-API-SPECIFICATION.md (секция 9 — Chat, секция 15 — WebSocket).
Прочитай docs/09-INTEGRATIONS.md (секция Redis — pub/sub).
Прочитай docs/11-SECURITY.md (секция WebSocket Security).

Реализуй Chat:

1. internal/service/chat.go:
   - CreatePersonalChat(ctx, user1ID, user2ID)
   - GetOrCreateCommunityChat(ctx, communityID) — auto-create
   - GetOrCreateEventChat(ctx, eventID) — auto-create
   - ListMyChats(ctx, userID) — с last_message, unread_count
   - GetMessages(ctx, chatID, cursor, limit) — cursor-based pagination
   - SendMessage(ctx, chatID, userID, content)
   - MarkAsRead(ctx, chatID, userID, messageID)
   - GetUnreadCount(ctx, userID) — total across all chats

2. internal/ws/hub.go — WebSocket Hub
   - Connections map: userID → []*websocket.Conn
   - Rooms map: chatID → []userID
   - Register/Unregister clients
   - Broadcast to room

3. internal/ws/handler.go — WebSocket handler
   - GET /v1/ws?token={jwt} — upgrade connection
   - JWT validation from query param
   - Message types: message, typing, read, join_room, leave_room
   - JSON format: {"type": "message", "chat_id": "...", "content": "..."}

4. Redis pub/sub для multi-instance:
   - Publish: ws:chat:{chat_id} → message JSON
   - Subscribe: каждый instance слушает все активные чаты

5. Rate limit: 60 messages/min per user

6. Handlers:
   - POST /v1/chats/personal
   - GET /v1/chats
   - GET /v1/chats/:id/messages
   - POST /v1/chats/:id/messages (REST fallback)
   - POST /v1/chats/:id/read
   - GET /v1/chats/unread-count
```

---

## Задача 5.3 — Notifications backend

**Агент:** Claude Code

### Промпт:

```
Прочитай docs/04-API-SPECIFICATION.md (секция 10 — Notifications).
Прочитай docs/09-INTEGRATIONS.md (секция Firebase FCM).

Реализуй Notifications:

1. internal/service/notification.go:
   - Create(ctx, userID, type, title, body, data) — сохранить в БД + отправить push
   - List(ctx, userID, page) — с пагинацией, grouped by date
   - MarkAsRead(ctx, userID, notificationID)
   - MarkAllAsRead(ctx, userID)
   - GetUnreadCount(ctx, userID)

2. internal/service/firebase.go:
   - SendPush(ctx, deviceToken, title, body, data)
   - SendToTopic(ctx, topic, title, body, data)
   - Init from service account JSON (env var: FIREBASE_SERVICE_ACCOUNT)

3. Notification triggers (вызываются из других services):
   - match_result_pending → когда Player A submit result
   - match_result_confirmed → когда Player B confirms
   - rating_changed → после ELO calculation
   - event_joined → когда кто-то записался на мой ивент
   - event_reminder → за 24ч и 1ч (cron job — следующий спринт)
   - new_message → при новом сообщении в чате (если не online в WS)
   - community_request → заявка на вступление

4. Handlers:
   - GET /v1/notifications
   - POST /v1/notifications/:id/read
   - POST /v1/notifications/read-all
   - GET /v1/notifications/unread-count

В dev-режиме: FCM mock (логируем push, не отправляем).
```

---

## Задача 5.4 — Mobile: Match screens

**Агент:** Cursor

### Промпт:

```
Прочитай docs/mobile-app/05-events-tab.md (секция результаты матчей).

Создай Match screens:

1. src/app/match/submit-result.tsx — Ввод результата
   - Выбор количества сетов (2 или 3)
   - Для каждого сета: 2 инпута (score P1, score P2)
   - Тай-брейк toggle (если сет 7:6)
   - Автоматическое определение winner
   - Preview перед отправкой
   - POST /v1/matches/:id/result

2. src/app/match/confirm-result.tsx — Подтверждение
   - Показывает: оппонент, счёт по сетам, кто победил
   - 2 кнопки: "Подтвердить ✓" / "Оспорить ✕"
   - POST /v1/matches/:id/confirm

3. Rating change animation: +16 (зелёный) / -16 (красный) после confirmation
```

---

## Задача 5.5 — Mobile: Chat

**Агент:** Cursor

### Промпт:

```
Прочитай docs/mobile-app/08-chat.md полностью.

Создай Chat:

1. src/app/chat/index.tsx — Список чатов
   - Карточки: аватар, название, last message preview, unread badge, time
   - Сортировка по последнему сообщению

2. src/app/chat/[id].tsx — Экран чата
   - FlatList inverted (новые внизу)
   - Мои сообщения справа (зелёный), чужие слева (серый)
   - Input bar внизу: текст + кнопка отправки
   - Typing indicator
   - Auto-scroll при новом сообщении
   - Cursor-based pagination при scroll вверх (load more)

3. src/shared/lib/websocket.ts — WS connection manager
   - Connect с JWT token
   - Auto-reconnect с exponential backoff (1s, 2s, 4s, 8s, max 30s)
   - Message handler: dispatch to correct chat
   - Typing events

4. Zustand store для chat state:
   - messages by chatID
   - unread counts
   - typing indicators
```

---

## Коммит Sprint 5-6:

```bash
git add .
git commit -m "feat(sprint-5-6): matches, ELO rating, chat, notifications"
git push origin main
```

---

# 🏠 SPRINT 7-8: Home, Posts, Badges, Courts (Неделя 13-16)

> **Цель:** Home tab, feed, gamification, карта кортов

---

## Задача 7.1 — Home + Feed + Posts backend

**Агент:** Claude Code

### Промпт:

```
Прочитай docs/04-API-SPECIFICATION.md (секции 8 Posts, секция для feed/home).
Прочитай docs/05-USER-STORIES.md (HOM-1, HOM-2, PST-1, PST-2).

Реализуй:

1. Home dashboard: GET /v1/home
   - my_rating, rating_change_week, upcoming_events (3 max), quick_stats
2. Feed: GET /v1/feed?page=&per_page=
   - Посты из моих сообществ + match results
3. Posts CRUD:
   - POST /v1/communities/:id/posts (text + до 5 images)
   - GET /v1/communities/:id/posts
   - POST /v1/posts/:id/like + DELETE
4. Auto-post при подтверждении матча
5. Badges service:
   - CheckBadges(ctx, userID) — после каждого матча
   - Badges: first_game (1 game), streak_3 (3 wins), streak_5 (5 wins), veteran (50 games), top_rating (1500+)
   - GET /v1/rating/badges
6. Friends: POST/DELETE /v1/friends/:user_id, GET /v1/friends
7. Courts: GET /v1/courts, GET /v1/courts/:id, GET /v1/courts/map
```

---

## Задача 7.2 — Mobile: Home tab + remaining screens

**Агент:** Cursor

### Промпт:

```
Прочитай docs/mobile-app/03-home-tab.md.

Создай:

1. Home tab: рейтинг виджет, quick actions, ближайшие игры, feed
2. Post creation screen
3. Post card component (text + images + like button)
4. Badges section в Profile (earned = цветные, in-progress = серые)
5. Friends list screen
6. Courts map screen (react-native-maps + markers)
7. Court details bottomsheet
```

---

## Коммит Sprint 7-8:

```bash
git add .
git commit -m "feat(sprint-7-8): home tab, posts, badges, friends, courts map"
git push origin main
```

---

# 🖥️ SPRINT 9-10: Web Panels (Неделя 17-20)

> **Цель:** Admin panel + Superadmin
> **Агент:** Cursor (web), Claude Code (backend endpoints)

---

## Задача 9.1 — Web Admin setup + auth

**Агент:** Cursor (открой apps/web-admin/)

### Промпт:

```
Прочитай .cursorrules (секция Web Admin).
Прочитай docs/web-platform/01-platform-overview.md.

Создай web admin проект:

1. npx create-vite apps/web-admin --template react-ts
2. Установи: shadcn/ui, tailwindcss, @tanstack/react-query, @tanstack/react-table, zustand, react-router-dom, recharts, react-hook-form, zod, axios
3. Настрой shadcn/ui (npx shadcn-ui@latest init)
4. Layout:
   - Sidebar (256px): лого, navigation (Dashboard, Members, Events, Content, Rating, Settings), community switcher
   - Main content area (max-width 1200px)
5. Auth: телефон + OTP (те же endpoints что мобайл)
6. Protected routes
```

---

## Задача 9.2 — Dashboard + Members + Events pages

**Агент:** Cursor

### Промпт:

```
Прочитай docs/web-platform/02-admin-dashboard.md ПОЛНОСТЬЮ (155 строк с wireframes).
Прочитай docs/web-platform/03-member-management.md ПОЛНОСТЬЮ (198 строк с wireframes).
Прочитай docs/web-platform/04-event-management.md.

Создай:

1. Dashboard page:
   - 4 metric cards (members, active, events, avg NTRP)
   - 3 charts (Recharts): growth line, level pie, activity bar
   - Quick actions (3 buttons)
   - Recent activity feed
   - Period filter dropdown

2. Members page:
   - TanStack Table: columns (avatar+name, NTRP, rating, role, joined, activity, actions)
   - Search, filters (NTRP, role, activity)
   - Bulk actions toolbar
   - Individual actions menu (view, change role, ban, remove)
   - Side panel: member detail view
   - Tab "Заявки" для closed communities

3. Events page:
   - TanStack Table: events list with status badges
   - Create event dialog
   - Event detail page (participants, match results)

API: тот же backend, те же endpoints. Новый endpoint:
GET /v1/admin/communities/:id/dashboard (попроси Claude Code создать).
```

---

## Задача 9.3 — Admin backend endpoints

**Агент:** Claude Code

### Промпт:

```
Прочитай docs/04-API-SPECIFICATION.md (секция 13 — Admin).

Создай admin-specific endpoints:

1. GET /v1/admin/communities/:id/dashboard — metrics, charts data, recent activity
2. GET /v1/admin/communities/:id/export?type=members|matches|ratings — CSV export
3. Все endpoints require admin/owner role в community
```

---

## Задача 9.4 — Superadmin panel

**Агент:** Cursor

### Промпт:

```
Прочитай docs/web-platform/06-superadmin.md.
Прочитай docs/04-API-SPECIFICATION.md (секция 14 — Superadmin).

Создай superadmin panel (apps/web-superadmin/):

1. Отдельное Vite + React приложение (копия структуры web-admin)
2. Auth: только role=superadmin
3. Dashboard: total users, communities, matches, growth chart
4. Verification queue: pending communities, approve/reject
5. User management: search, ban/unban
6. Courts CRUD: add, edit, delete (название, тип, координаты, адрес, фото)
```

---

## Коммит Sprint 9-10:

```bash
git add .
git commit -m "feat(sprint-9-10): web admin panel + superadmin panel"
git push origin main
```

---

# ✨ SPRINT 11: Polish (Неделя 21-22)

> **Цель:** UX polish, edge cases, performance

---

## Задача 11.1 — Mobile polish

**Агент:** Cursor

### Промпт:

```
Прочитай docs/13-DESIGN-SYSTEM.md (компоненты Skeleton, Empty State, Error State).

Добавь во ВСЕ экраны:

1. Skeleton loading для каждого списка/экрана
2. Empty states с иконкой + текстом + action button
3. Error states с "Повторить"
4. Pull-to-refresh на всех FlatList
5. Haptic feedback (expo-haptics) на кнопках
6. Splash screen + app icon
7. Проверь ВСЕ строки через t() — ни одного хардкода

Проверь все экраны на iOS и Android.
```

---

## Задача 11.2 — Backend polish

**Агент:** Claude Code

### Промпт:

```
1. Создай scripts/seed/main.go по docs/17-SEED-DATA.md — 20 users, 3 communities, 10 events, 15 matches, 12 courts
2. Event reminder cron: напоминания за 24ч и 1ч (простой ticker в goroutine)
3. Проверь все rate limits из docs/11-SECURITY.md
4. Проверь все input validations
5. Добавь Sentry integration (docs/09-INTEGRATIONS.md)
6. Оптимизация: добавь cache (Redis) для leaderboards (TTL 5 min)
```

---

# 🚀 SPRINT 12: QA & Launch (Неделя 23-24)

---

## Задача 12.1 — Deploy

**Агент:** Claude Code

### Промпт:

```
Прочитай docs/10-DEPLOYMENT.md полностью.

1. Backend → Railway:
   - Dockerfile (multi-stage build)
   - railway.toml
   - Set all env variables

2. Web Admin → Vercel:
   - vercel.json
   - Connect GitHub repo

3. Mobile → EAS Build:
   - eas.json (development, preview, production profiles)
   - eas build --platform all --profile preview
   - eas submit (TestFlight + Google Play Beta)

4. DNS: api.tennisapp.kz → Railway, admin.tennisapp.kz → Vercel
```

---

## Задача 12.2 — Manual testing

**Ты сам** (не агент):

### Чеклист:

```
[ ] Auth: новый пользователь (phone → OTP → profile → quiz → home)
[ ] Auth: существующий пользователь (phone → OTP → home)
[ ] Auth: PIN (set → use → forgot)
[ ] Profile: view → edit → avatar upload
[ ] Players: search → filter → view profile → add friend
[ ] Events: browse → filter → view → join → leave
[ ] Events: create (all 8 steps wizard)
[ ] Communities: browse → join → view → leave
[ ] Communities: create → manage members → change roles
[ ] Match: submit result → opponent confirms → ELO updates
[ ] Match: dispute flow
[ ] Chat: personal → send message → receive
[ ] Chat: community chat → send → receive
[ ] Notifications: receive push → tap → deep link
[ ] Rating: check leaderboard → rating history graph
[ ] Web Admin: login → dashboard → members → events
[ ] Superadmin: login → verify community → manage courts
[ ] i18n: switch ru → kk → en (all screens)
[ ] Edge: no internet → retry
[ ] Edge: expired token → auto refresh
[ ] Edge: full event → proper error
[ ] iOS + Android: both platforms
```

---

# 📋 ШПАРГАЛКА: Промпт для начала любой задачи

Перед каждой задачей давай агенту этот шаблон:

```
Прочитай следующие документы:
1. {CLAUDE.md или .cursorrules}
2. docs/{relevant-spec}.md
3. docs/05-USER-STORIES.md (секция {ID} — acceptance criteria)

Задача: {описание}

Требования:
- Следуй coding conventions из docs/07-CODING-CONVENTIONS.md
- Коды ошибок ТОЛЬКО из docs/14-ERROR-CODES.md
- Все строки через i18n (docs/15-I18N-GUIDE.md)
- Компоненты из shared/ui/ (docs/13-DESIGN-SYSTEM.md)

После реализации:
- Проверь что код компилируется
- Проверь что тесты проходят
- Покажи какие файлы были созданы/изменены
```

---

# ⏱️ РЕАЛИСТИЧНЫЕ СРОКИ

| Sprint | Недели | Оптимистично | Реалистично | С буфером |
|--------|--------|-------------|-------------|-----------|
| 1 Foundation | 1-2 | 2 нед. | 2 нед. | 2 нед. |
| 2 Auth | 3-4 | 2 нед. | 3 нед. | 3 нед. |
| 3-4 Core | 5-8 | 4 нед. | 5 нед. | 6 нед. |
| 5-6 Match+Chat | 9-12 | 4 нед. | 5 нед. | 6 нед. |
| 7-8 Home+Posts | 13-16 | 4 нед. | 4 нед. | 5 нед. |
| 9-10 Web | 17-20 | 4 нед. | 5 нед. | 6 нед. |
| 11 Polish | 21-22 | 2 нед. | 3 нед. | 3 нед. |
| 12 Launch | 23-24 | 2 нед. | 2 нед. | 3 нед. |
| **TOTAL** | | **24 нед.** | **29 нед.** | **34 нед.** |

---

# 🔑 ЗОЛОТЫЕ ПРАВИЛА

1. **Docs First** — Сначала обнови документацию, потом давай задачу агенту
2. **1 задача = 1 промпт** — Не смешивай backend и frontend в одном промпте
3. **Указывай конкретные файлы** — "Прочитай docs/04-API-SPECIFICATION.md секция 5" > "посмотри API"
4. **Проверяй после каждой задачи** — Не копи 5 задач без проверки
5. **Коммить после каждого блока** — git commit часто, не теряй прогресс
6. **Не спорь с агентом** — Если делает не так, обнови документацию и дай заново
7. **Claude Code = backend, Cursor = frontend** — Не смешивай роли
