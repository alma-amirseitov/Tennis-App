# User Stories & Sprint Plan v2.1
## Теннисная платформа для Астаны

**Версия:** 2.1 (добавлены Acceptance Criteria для Sprint 1-2)  
**Дата:** Февраль 2026  
**Методология:** 2-недельные спринты  

---

## Легенда

**Приоритет:**
- 🔴 **P0** — Блокер. Без этого продукт не запустится
- 🟠 **P1** — Критично для MVP. Нужно до релиза
- 🟡 **P2** — Важно. Улучшает UX, но можно выпустить без этого
- 🟢 **P3** — Nice-to-have. Можно в Phase 2

**Story Points (SP):** Fibonacci (1, 2, 3, 5, 8, 13)
- 1 SP = ~0.5 дня
- 2 SP = ~1 день
- 3 SP = ~2 дня
- 5 SP = ~3-4 дня
- 8 SP = ~1 неделя
- 13 SP = ~1.5-2 недели (нужно разбить)

---

## SPRINT 1: Foundation (Недели 1-2)
**Цель:** Backend skeleton + инфра + mobile shell

---

### INF-1: Monorepo scaffold
**SP:** 3 | **P:** 🔴 | **Компонент:** 🔧

**Задача:** Настроить monorepo структуру проекта

**Acceptance Criteria:**
- [ ] Структура `apps/backend/`, `apps/mobile/`, `apps/web-admin/`, `apps/web-superadmin/`, `packages/` создана
- [ ] `apps/backend/go.mod` инициализирован с модулем `github.com/{user}/tennisapp/apps/backend`
- [ ] `README.md` в корне с описанием структуры и командами запуска
- [ ] `.gitignore` покрывает: Go binaries, node_modules, .env, .DS_Store, IDE files
- [ ] `.editorconfig` — 2 spaces для TS/JSON, tabs для Go

---

### INF-2: Docker Compose для local dev
**SP:** 2 | **P:** 🔴 | **Компонент:** 🔧

**Задача:** Docker Compose с PostgreSQL, Redis, MinIO

**Acceptance Criteria:**
- [ ] `docker-compose.yml` в корне проекта
- [ ] PostgreSQL 16-alpine: порт 5432, db=tennisapp, user=tennisapp, password=tennisapp, volume для persistence
- [ ] Redis 7-alpine: порт 6379
- [ ] MinIO: порты 9000 (API) + 9001 (console), default credentials minioadmin/minioadmin, volume
- [ ] `docker-compose up -d` запускается без ошибок
- [ ] `docker-compose down && docker-compose up -d` — данные PostgreSQL сохраняются (volume)

---

### INF-3: Database migrations + начальная схема
**SP:** 5 | **P:** 🔴 | **Компонент:** 🔧

**Задача:** golang-migrate setup + миграция со всеми таблицами из `docs/03-DATABASE-SCHEMA.sql`

**Acceptance Criteria:**
- [ ] golang-migrate установлен как зависимость
- [ ] `apps/backend/migrations/000001_init_schema.up.sql` содержит все 21 таблицу из docs/03-DATABASE-SCHEMA.sql
- [ ] `apps/backend/migrations/000001_init_schema.down.sql` — DROP TABLE IF EXISTS CASCADE для всех таблиц (обратный порядок)
- [ ] `make migrate-up` применяет миграцию, все таблицы создаются
- [ ] `make migrate-down` откатывает миграцию, все таблицы удаляются
- [ ] Повторный `make migrate-up` после down — работает без ошибок (идемпотентность)
- [ ] Все constraints, indexes, triggers, views из schema включены

---

### INF-4: Go backend skeleton
**SP:** 5 | **P:** 🔴 | **Компонент:** 🔧

**Задача:** HTTP server с chi router, middleware chain, config, logging

**Acceptance Criteria:**
- [ ] `cmd/server/main.go` — entry point с graceful shutdown (SIGINT/SIGTERM)
- [ ] `internal/config/config.go` — envconfig struct с полями из `.env.example`
- [ ] Chi router с middleware chain: Logger → Recovery → CORS → RequestID
- [ ] `GET /health` возвращает `{"status":"ok","version":"0.1.0"}`
- [ ] `GET /health` проверяет DB и Redis connection, возвращает `"database":"connected"` / `"disconnected"`
- [ ] slog structured logging: каждый request логируется с method, path, status, duration
- [ ] `.env.example` со всеми переменными (DB, Redis, JWT, SMS, S3, Firebase, Sentry)
- [ ] `make dev` запускает сервер с hot-reload (air)
- [ ] `make build` компилирует бинарник без ошибок
- [ ] Сервер слушает на порту из `PORT` env variable (default 8080)

---

### INF-5: sqlc setup + users queries
**SP:** 3 | **P:** 🔴 | **Компонент:** 🔧

**Задача:** Настроить sqlc, создать первые queries для users

**Acceptance Criteria:**
- [ ] `apps/backend/sqlc.yaml` сконфигурирован (postgres, Go output path = `internal/repository`)
- [ ] `apps/backend/internal/repository/queries/users.sql` — минимум: CreateUser, GetUserByID, GetUserByPhone, UpdateUser
- [ ] `make sqlc` генерирует Go код без ошибок
- [ ] Сгенерированные типы используют `uuid.UUID` для ID, `sql.NullString` для nullable fields
- [ ] `internal/repository/db.go` — connection pool (pgx) с настраиваемыми параметрами

---

### INF-6: CI/CD: GitHub Actions
**SP:** 3 | **P:** 🟠 | **Компонент:** 🔧

**Задача:** GitHub Actions pipeline для backend

**Acceptance Criteria:**
- [ ] `.github/workflows/backend.yml` — триггер на push/PR в main/develop, paths: apps/backend/**
- [ ] Steps: checkout → setup-go 1.22 → go mod download → go vet → go test → go build
- [ ] Postgres и Redis services запускаются как containers для тестов
- [ ] Pipeline проходит на пустом проекте (хотя бы compilation check)

---

### INF-7: Expo project + navigation shell
**SP:** 3 | **P:** 🔴 | **Компонент:** 📱

**Задача:** Инициализировать Expo проект с Expo Router и 5 табами

**Acceptance Criteria:**
- [ ] `apps/mobile/` — Expo SDK 52+, TypeScript strict, ESLint
- [ ] Expo Router file-based routing настроен
- [ ] Tab navigator с 5 табами: Home, Players, Events, Communities, Profile
- [ ] Каждый таб показывает placeholder экран с названием
- [ ] Tab bar соответствует дизайн-системе: высота 80px, иконки, цвета из `docs/13-DESIGN-SYSTEM.md`
- [ ] Stack navigator внутри каждого таба (для sub-screens)
- [ ] `npx expo start` запускается без ошибок
- [ ] Работает на iOS Simulator и Android Emulator

---

### INF-8: i18n setup
**SP:** 3 | **P:** 🔴 | **Компонент:** 📱

**Задача:** Настроить i18next с 3 языками

**Acceptance Criteria:**
- [ ] i18next + react-i18next настроены
- [ ] 3 файла переводов: `ru.json`, `kk.json`, `en.json` (минимум ключи для auth + tab names)
- [ ] Язык по умолчанию: ru
- [ ] Детектор языка устройства (если kk → kk, если en → en, иначе → ru)
- [ ] `useTranslation()` hook работает в компонентах
- [ ] Переключение языка в runtime (для экрана настроек)
- [ ] Ни один user-visible string не захардкожен — всё через t()
- [ ] Структура ключей: `{screen}.{element}` — например `auth.phone_title`, `tabs.events`

---

### INF-9: Базовая дизайн-система
**SP:** 5 | **P:** 🔴 | **Компонент:** 📱

**Задача:** Реализовать токены и базовые компоненты из `docs/13-DESIGN-SYSTEM.md`

**Acceptance Criteria:**
- [ ] `src/shared/theme/colors.ts` — все цвета из Design System
- [ ] `src/shared/theme/typography.ts` — fontSize, fontWeight, textStyles
- [ ] `src/shared/theme/spacing.ts` — spacing scale (4, 8, 12, 16, 20, 24, 32)
- [ ] `src/shared/theme/radius.ts` — border radius tokens
- [ ] `src/shared/ui/Button.tsx` — Primary, Secondary, Outline, Small variants; disabled + loading states
- [ ] `src/shared/ui/Input.tsx` — default, focused, error states; label + error message
- [ ] `src/shared/ui/Card.tsx` — bg, radius, padding, border as per design system
- [ ] `src/shared/ui/Avatar.tsx` — sizes 24-80, circle, initials fallback
- [ ] `src/shared/ui/Badge.tsx` — primary, success, warning, danger, info, muted variants
- [ ] Все компоненты используют tokens, не raw hex values
- [ ] TypeScript strict — все props typed, no `any`

---

### INF-10: Shared types package
**SP:** 3 | **P:** 🟠 | **Компонент:** 🔧📱

**Задача:** TypeScript типы = API контракты

**Acceptance Criteria:**
- [ ] `packages/shared-types/src/index.ts` экспортирует все типы
- [ ] Типы для всех основных entities: User, Community, Event, Match, Chat, Message, Notification
- [ ] Enum типы: EventType, EventStatus, CommunityType, UserRole, CommunityRole
- [ ] API response типы: ApiResponse<T>, PaginatedResponse<T>, ApiError
- [ ] Request типы для Sprint 1-2: SendOTPRequest, VerifyOTPRequest, SetupProfileRequest, QuizRequest
- [ ] Mobile и web-admin могут импортировать: `import { User, Event } from '@tennisapp/shared-types'`

---

## SPRINT 2: Auth (Недели 3-4)
**Цель:** Полный auth flow — backend + mobile

---

### AUTH-1: POST /auth/otp/send
**SP:** 5 | **P:** 🔴 | **Компонент:** 🔧

**Задача:** Endpoint отправки SMS OTP

**Acceptance Criteria:**
- [ ] `POST /v1/auth/otp/send` принимает `{"phone": "+77071234567"}`
- [ ] Валидация формата телефона: `^\+7[0-9]{10}$` — иначе 400 `VALIDATION_ERROR`
- [ ] Генерация 4-значного случайного OTP кода
- [ ] Session создаётся в Redis: ключ `otp:{session_id}`, значение `{phone, code, attempts: 0}`, TTL 5 мин
- [ ] SMS отправляется через SMS provider (mock в dev: код всегда 1234, SMS не отправляется)
- [ ] Rate limit: 3 SMS/час на номер — иначе 429 `RATE_LIMITED`
- [ ] Rate limit: 10 SMS/день на номер — иначе 429 `RATE_LIMITED`
- [ ] Response: `{"session_id": "uuid", "expires_in": 300}`
- [ ] Телефон маскируется в логах: +7707***4567
- [ ] Unit тест: валидация телефона (valid, invalid, empty)
- [ ] Unit тест: rate limit (3-й запрос за час → ошибка)

**Endpoint spec:** `docs/04-API-SPECIFICATION.md`, секция 1, endpoint 1.1

---

### AUTH-2: POST /auth/otp/verify
**SP:** 5 | **P:** 🔴 | **Компонент:** 🔧

**Задача:** Проверка OTP, создание/авторизация пользователя

**Acceptance Criteria:**
- [ ] `POST /v1/auth/otp/verify` принимает `{"session_id": "uuid", "code": "1234"}`
- [ ] Загружает сессию из Redis по session_id — если не найдена: 400 `OTP_SESSION_EXPIRED`
- [ ] Сравнивает код — если неверный: инкремент attempts, если attempts >= 5: удалить сессию, 400 `OTP_MAX_ATTEMPTS`
- [ ] Если код верный — удалить сессию из Redis
- [ ] Поиск пользователя по телефону в БД
- [ ] **Новый пользователь:** создать запись в users (phone, status=active), вернуть `{"is_new_user": true, "temp_token": "jwt", "user": {...}}`
- [ ] **Существующий пользователь:** вернуть `{"is_new_user": false, "access_token": "jwt", "refresh_token": "jwt", "user": {...}}`
- [ ] Access token: JWT HS256, TTL 15 мин, claims: sub=user_id, role=user
- [ ] Refresh token: JWT, TTL 30 дней, jti=unique_id, сохранён в Redis: `refresh:{jti}` → user_id
- [ ] Unit тест: верный код → tokens
- [ ] Unit тест: неверный код → ошибка, attempts++
- [ ] Unit тест: 5 попыток → session deleted

**Endpoint spec:** `docs/04-API-SPECIFICATION.md`, секция 1, endpoint 1.2

---

### AUTH-3: POST /auth/refresh
**SP:** 3 | **P:** 🔴 | **Компонент:** 🔧

**Задача:** Обновление access token через refresh token с rotation

**Acceptance Criteria:**
- [ ] `POST /v1/auth/refresh` принимает `{"refresh_token": "jwt"}`
- [ ] Валидация JWT подписи и expiration
- [ ] Проверка jti в Redis — если не найден: 401 `TOKEN_REVOKED` (возможно compromise)
- [ ] Удаление старого jti из Redis (one-time use)
- [ ] Генерация нового access_token + нового refresh_token с новым jti
- [ ] Сохранение нового jti в Redis
- [ ] Response: `{"access_token": "...", "refresh_token": "..."}`
- [ ] Если refresh_token reuse detected (jti уже удалён) → удалить ВСЕ refresh tokens пользователя
- [ ] Unit тест: valid refresh → new tokens
- [ ] Unit тест: expired refresh → 401
- [ ] Unit тест: reused refresh → revoke all

---

### AUTH-4: Auth middleware
**SP:** 3 | **P:** 🔴 | **Компонент:** 🔧

**Задача:** JWT validation middleware, user context injection

**Acceptance Criteria:**
- [ ] Middleware извлекает token из `Authorization: Bearer {token}` header
- [ ] Валидирует JWT подпись (HS256, JWT_SECRET)
- [ ] Проверяет expiration — если expired: 401 `TOKEN_EXPIRED`
- [ ] Инжектит `user_id` и `role` в request context
- [ ] Helper функции: `middleware.GetUserID(ctx)`, `middleware.GetUserRole(ctx)`
- [ ] Отсутствие header → 401 `UNAUTHORIZED`
- [ ] Malformed token → 401 `UNAUTHORIZED`
- [ ] Requests к `/health`, `/v1/auth/otp/*`, `/v1/auth/refresh` — без middleware (public)
- [ ] Unit тест: valid token → user_id in context
- [ ] Unit тест: expired token → 401
- [ ] Unit тест: no header → 401

---

### AUTH-5: Rate limiting middleware
**SP:** 3 | **P:** 🔴 | **Компонент:** 🔧

**Задача:** Redis-based rate limiting middleware

**Acceptance Criteria:**
- [ ] Generic rate limiter: `RateLimit(key, limit, window)`
- [ ] Redis sliding window counter (INCR + EXPIRE)
- [ ] Configurable per-route: SMS = 3/hour, API general = 100/min
- [ ] Key extraction: user_id (authenticated) или IP (unauthenticated)
- [ ] Response при превышении: 429, header `Retry-After: {seconds}`
- [ ] Rate limit headers: `X-RateLimit-Limit`, `X-RateLimit-Remaining`, `X-RateLimit-Reset`
- [ ] Unit тест: requests within limit → pass
- [ ] Unit тест: requests over limit → 429

---

### AUTH-6: Экран ввода телефона
**SP:** 3 | **P:** 🔴 | **Компонент:** 📱

**Задача:** Mobile экран ввода телефона для OTP

**Acceptance Criteria:**
- [ ] Маска ввода: +7 (XXX) XXX-XX-XX
- [ ] Автоматический +7 prefix, нельзя удалить
- [ ] Клавиатура: numeric keypad
- [ ] Кнопка "Получить код" — disabled пока номер не полный (10 цифр после +7)
- [ ] При нажатии: вызов `POST /v1/auth/otp/send`
- [ ] Loading state на кнопке при отправке
- [ ] При успехе: навигация на экран OTP, передать session_id
- [ ] При ошибке rate limit: показать toast "Подождите X минут"
- [ ] При ошибке сети: показать toast "Проверьте подключение"
- [ ] i18n: все строки через t()
- [ ] Спецификация: `docs/mobile-app/02-auth-onboarding.md`

---

### AUTH-7: Экран ввода OTP
**SP:** 3 | **P:** 🔴 | **Компонент:** 📱

**Задача:** 4 ячейки для ввода OTP кода

**Acceptance Criteria:**
- [ ] 4 отдельных ячейки ввода, автофокус на первую
- [ ] При вводе цифры — автопереход на следующую ячейку
- [ ] При удалении — автопереход на предыдущую
- [ ] Paste поддержка: вставка 4 цифр из clipboard
- [ ] Таймер повторной отправки: 60 секунд → "Отправить повторно"
- [ ] При заполнении всех 4 — автоматический вызов `POST /v1/auth/otp/verify`
- [ ] Loading state при проверке
- [ ] Неверный код: анимация shake на ячейках, toast "Неверный код"
- [ ] При `is_new_user: true` → навигация на Profile Setup
- [ ] При `is_new_user: false` → навигация на Home (tabs)
- [ ] Сохранение tokens в react-native-keychain
- [ ] i18n: все строки через t()

---

### AUTH-8: Token management + Axios interceptor
**SP:** 3 | **P:** 🔴 | **Компонент:** 📱

**Задача:** Безопасное хранение токенов + auto-refresh

**Acceptance Criteria:**
- [ ] Tokens хранятся в react-native-keychain (не AsyncStorage!)
- [ ] Axios instance с baseURL = API_BASE_URL
- [ ] Request interceptor: добавляет `Authorization: Bearer {access_token}`
- [ ] Response interceptor: при 401 → пробует refresh token → повторяет request
- [ ] Если refresh тоже 401 → clear tokens → redirect на auth screen
- [ ] Queue: параллельные запросы при refresh — ждут один refresh, потом все retry
- [ ] Auth state в Zustand: `{isAuthenticated, user, isLoading}`
- [ ] При app launch: проверка наличия token → try refresh → set state
- [ ] Auth guard: если не authenticated → redirect на auth flow

---

### AUTH-9: POST /auth/profile/setup
**SP:** 3 | **P:** 🔴 | **Компонент:** 🔧

**Задача:** Заполнение профиля нового пользователя

**Acceptance Criteria:**
- [ ] `POST /v1/auth/profile/setup` — protected (temp_token)
- [ ] Body: `{"first_name": "Алмас", "last_name": "Б.", "gender": "male", "birth_year": 1995, "city": "Astana", "district": "Есильский"}`
- [ ] Валидация: first_name 2-50 chars, last_name 2-50, gender in [male, female], birth_year 1940-2012, city required
- [ ] Обновляет users запись
- [ ] Возвращает полный access_token + refresh_token (upgrade from temp_token)
- [ ] Только для `is_new_user` — если профиль уже заполнен → 400 `PROFILE_ALREADY_SET`
- [ ] Unit тест: валидация полей
- [ ] Unit тест: дублирующий setup → ошибка

---

### AUTH-10: Экран заполнения профиля
**SP:** 3 | **P:** 🔴 | **Компонент:** 📱

**Задача:** Форма профиля для нового пользователя

**Acceptance Criteria:**
- [ ] Поля: Имя, Фамилия, Пол (toggle male/female), Год рождения (picker), Город (Астана), Район (dropdown)
- [ ] React Hook Form + Zod validation
- [ ] Inline validation errors под каждым полем
- [ ] Кнопка "Продолжить" — disabled пока форма невалидна
- [ ] При submit: `POST /v1/auth/profile/setup`
- [ ] Loading state
- [ ] При успехе: сохранение новых tokens → навигация на Quiz
- [ ] i18n: все строки через t()
- [ ] Районы Астаны: Есильский, Алматинский, Сарыаркинский, Байконурский, Нуринский

---

### AUTH-11: Quiz endpoints
**SP:** 3 | **P:** 🟠 | **Компонент:** 🔧

**Задача:** Skill quiz для определения начального уровня

**Acceptance Criteria:**
- [ ] `GET /v1/quiz` — возвращает список вопросов (3-5 вопросов с вариантами ответов)
- [ ] Вопросы hardcoded (не из БД): опыт игры, частота, уровень соперников, результаты
- [ ] `POST /v1/quiz` принимает `{"answers": [{"question_id": 1, "answer_id": 2}, ...]}`
- [ ] Алгоритм расчёта: каждый ответ имеет weight, sum → initial NTRP level
- [ ] Обновляет `users.level` и `users.rating_score` (initial rating based on level)
- [ ] Response: `{"level": "Любитель", "ntrp": 3.0, "initial_rating": 1150}`
- [ ] Unit тест: разные ответы → разные уровни

---

### AUTH-12: Экран Quiz
**SP:** 3 | **P:** 🟠 | **Компонент:** 📱

**Задача:** Skill quiz на мобайле

**Acceptance Criteria:**
- [ ] Загрузка вопросов: `GET /v1/quiz`
- [ ] По 1 вопросу на экране, swipe/button для перехода
- [ ] Карточки вариантов ответов (tap для выбора)
- [ ] Progress bar сверху (1/5, 2/5...)
- [ ] При завершении: `POST /v1/quiz` → показать результат
- [ ] Экран результата: "Ваш уровень: Любитель (NTRP 3.0)" + кнопка "Начать"
- [ ] При нажатии "Начать" → навигация на Home (tabs)
- [ ] Можно пропустить (Skip) → default level "Новичок" (NTRP 2.5)

---

### AUTH-13: Онбординг
**SP:** 2 | **P:** 🟡 | **Компонент:** 📱

**Задача:** Swipe-экраны при первом запуске

**Acceptance Criteria:**
- [ ] 3-4 экрана с иллюстрациями (placeholder images) и текстом
- [ ] Swipe между экранами, pagination dots
- [ ] Кнопка "Пропустить" на каждом экране
- [ ] На последнем: "Начать" → экран ввода телефона
- [ ] Показывается только при первом запуске (AsyncStorage flag)

---

### AUTH-14: PIN-код endpoints
**SP:** 3 | **P:** 🟡 | **Компонент:** 🔧

**Задача:** Установка и проверка PIN-кода

**Acceptance Criteria:**
- [ ] `POST /v1/auth/pin/set` — `{"pin": "1234"}` → bcrypt hash сохранён в users.pin_hash
- [ ] `POST /v1/auth/pin/verify` — `{"pin": "1234"}` → bcrypt compare → tokens
- [ ] Валидация: exactly 4 digits
- [ ] Max 3 попытки verify — потом require OTP re-auth
- [ ] Attempts counter в Redis: `pin_attempts:{user_id}`, TTL 1 hour
- [ ] Unit тест: set + verify → success
- [ ] Unit тест: wrong pin 3 times → locked

---

### AUTH-15: PIN-код экраны
**SP:** 3 | **P:** 🟡 | **Компонент:** 📱

**Задача:** Установка и ввод PIN на мобайле

**Acceptance Criteria:**
- [ ] Экран установки: 4 ячейки + подтверждение (ввести дважды)
- [ ] Если не совпадают → shake animation + "PIN-коды не совпадают"
- [ ] При совпадении: `POST /v1/auth/pin/set` → success toast
- [ ] Экран ввода (при повторном входе): 4 ячейки, biometric prompt (если available)
- [ ] 3 неверных попытки → "Войдите через SMS" → redirect на phone screen
- [ ] Кнопка "Забыли PIN?" → redirect на phone screen

---

## SPRINT 3-4: Core — Events & Communities (Недели 5-8)
**Цель:** Основная ценность продукта — создание и поиск игр

### EPIC 3: Пользователи и профиль

| ID | Story | SP | P | Компонент |
|----|-------|----|---|-----------|
| USR-1 | **Backend:** GET /users/me + PATCH /users/me — профиль пользователя | 3 | 🔴 | 🔧 |
| USR-2 | **Backend:** POST /users/me/avatar — загрузка аватара в S3/MinIO | 3 | 🔴 | 🔧 |
| USR-3 | **Backend:** GET /users/:id — публичный профиль (с badges, communities, stats) | 3 | 🔴 | 🔧 |
| USR-4 | **Backend:** GET /users/search — поиск с фильтрами (trgm, level, district, gender) | 5 | 🔴 | 🔧 |
| USR-5 | **Mobile:** Таб «Профиль» — мой профиль (7 секций) | 5 | 🔴 | 📱 |
| USR-6 | **Mobile:** Экран редактирования профиля | 3 | 🔴 | 📱 |
| USR-7 | **Mobile:** Экран публичного профиля (другого игрока) + кнопки действий | 3 | 🔴 | 📱 |
| USR-8 | **Mobile:** Экран настроек (язык, уведомления, PIN, приватность, о приложении) | 3 | 🟠 | 📱 |
| USR-9 | **Mobile:** Таб «Игроки» — каталог с поиском и фильтрами | 5 | 🟠 | 📱 |
| USR-10 | **Backend:** PATCH /users/me/notifications + /users/me/privacy | 2 | 🟠 | 🔧 |

---

### EPIC 4: Сообщества

| ID | Story | SP | P | Компонент |
|----|-------|----|---|-----------|
| COM-1 | **Backend:** POST /communities — создание сообщества (4 типа, авто verification_status) | 5 | 🔴 | 🔧 |
| COM-2 | **Backend:** GET /communities + GET /communities/:id — список и детали | 3 | 🔴 | 🔧 |
| COM-3 | **Backend:** POST /communities/:id/join + DELETE — вступление / выход | 3 | 🔴 | 🔧 |
| COM-4 | **Backend:** GET /communities/:id/members + PATCH role/status — управление участниками | 5 | 🔴 | 🔧 |
| COM-5 | **Backend:** RBAC middleware — проверка роли в сообществе (owner/admin/moderator/member) | 5 | 🔴 | 🔧 |
| COM-6 | **Backend:** POST /communities/:id/members/review — одобрение/отклонение заявок | 3 | 🔴 | 🔧 |
| COM-7 | **Mobile:** Таб «Сообщества» — список с поиском и фильтрами | 5 | 🔴 | 📱 |
| COM-8 | **Mobile:** Экран сообщества (6 внутренних табов: лента, ивенты, рейтинг, участники, чат, фото) | 8 | 🔴 | 📱 |
| COM-9 | **Mobile:** Экран создания сообщества (для групп — любой; для клубов/лиг — заявка) | 3 | 🟠 | 📱 |
| COM-10 | **Mobile:** Экран «Мои сообщества» | 2 | 🟠 | 📱 |
| COM-11 | **Backend:** GET /communities/:id/leaderboard — рейтинг участников | 3 | 🟠 | 🔧 |
| COM-12 | **Backend:** GET /communities/:id/feed — лента постов сообщества | 3 | 🟠 | 🔧 |

---

### EPIC 5: Ивенты (ключевой модуль)

| ID | Story | SP | P | Компонент |
|----|-------|----|---|-----------|
| EVT-1 | **Backend:** POST /events — создание ивента (все 8 шагов конструктора в 1 запросе) | 5 | 🔴 | 🔧 |
| EVT-2 | **Backend:** GET /events — лента с фильтрами (тип, статус, уровень, дата, район, сообщество) | 5 | 🔴 | 🔧 |
| EVT-3 | **Backend:** GET /events/:id — детали ивента (+ participants, my_status, can_join) | 3 | 🔴 | 🔧 |
| EVT-4 | **Backend:** POST /events/:id/join + DELETE — запись / отписка | 3 | 🔴 | 🔧 |
| EVT-5 | **Backend:** PATCH /events/:id/status — lifecycle transitions | 3 | 🔴 | 🔧 |
| EVT-6 | **Backend:** GET /events/calendar — ивенты по месяцу (grouped by day) | 3 | 🟠 | 🔧 |
| EVT-7 | **Backend:** GET /events/my — мои ивенты (created / joined / past) | 3 | 🟠 | 🔧 |
| EVT-8 | **Mobile:** Таб «Ивенты» — 3 внутренних таба (Лента / Календарь / Мои) | 5 | 🔴 | 📱 |
| EVT-9 | **Mobile:** Лента ивентов с карточками и фильтрами | 5 | 🔴 | 📱 |
| EVT-10 | **Mobile:** Конструктор ивента — wizard (8 шагов с анимациями) | 8 | 🔴 | 📱 |
| EVT-11 | **Mobile:** Экран деталей ивента (инфо, участники, кнопка записи) | 5 | 🔴 | 📱 |
| EVT-12 | **Mobile:** Календарь (месячный вид + дневной вид с ивентами) | 5 | 🟠 | 📱 |
| EVT-13 | **Mobile:** Мои ивенты (3 подтаба) | 3 | 🟠 | 📱 |

---

## SPRINT 5-6: Matches, Rating, Chat (Недели 9-12)
**Цель:** Результаты матчей, рейтинговая система, чат

### EPIC 6: Матчи и результаты

| ID | Story | SP | P | Компонент |
|----|-------|----|---|-----------|
| MTH-1 | **Backend:** POST /events/:id/matches — создание матча внутри ивента | 3 | 🔴 | 🔧 |
| MTH-2 | **Backend:** POST /matches/:id/result — внесение результата (score JSONB) | 5 | 🔴 | 🔧 |
| MTH-3 | **Backend:** POST /matches/:id/confirm — подтверждение / оспаривание результата | 5 | 🔴 | 🔧 |
| MTH-4 | **Backend:** ELO calculation service — расчёт нового рейтинга после подтверждения | 5 | 🔴 | 🔧 |
| MTH-5 | **Backend:** Обновление player_stats_global + community_members rating после матча | 3 | 🔴 | 🔧 |
| MTH-6 | **Backend:** rating_history — запись всех изменений рейтинга | 2 | 🔴 | 🔧 |
| MTH-7 | **Backend:** GET /matches/my — история матчей пользователя | 3 | 🟠 | 🔧 |
| MTH-8 | **Backend:** GET /rating/global + /rating/me + /rating/history | 3 | 🟠 | 🔧 |
| MTH-9 | **Mobile:** Экран ввода результата (выбор сетов, тай-брейк, winner) | 5 | 🔴 | 📱 |
| MTH-10 | **Mobile:** Экран подтверждения результата (push → экран с деталями → Confirm/Dispute) | 3 | 🔴 | 📱 |
| MTH-11 | **Mobile:** История матчей (в профиле) | 3 | 🟠 | 📱 |
| MTH-12 | **Mobile:** График динамики рейтинга (в профиле, 6 мес.) | 3 | 🟠 | 📱 |

### EPIC 7: Чат

| ID | Story | SP | P | Компонент |
|----|-------|----|---|-----------|
| CHT-1 | **Backend:** WebSocket server — connection, auth, hub, rooms | 8 | 🔴 | 🔧 |
| CHT-2 | **Backend:** Chat CRUD — создание personal/community/event чатов | 3 | 🔴 | 🔧 |
| CHT-3 | **Backend:** Messages — отправка, получение, пагинация (cursor-based) | 5 | 🔴 | 🔧 |
| CHT-4 | **Backend:** Read status — отметка прочтения, unread count | 3 | 🔴 | 🔧 |
| CHT-5 | **Backend:** Автоматическое создание чата сообщества (при создании community) | 2 | 🔴 | 🔧 |
| CHT-6 | **Backend:** Автоматическое создание чата ивента (при создании event с участниками) | 2 | 🔴 | 🔧 |
| CHT-7 | **Mobile:** Экран списка чатов (с preview последнего сообщения, unread badge) | 5 | 🔴 | 📱 |
| CHT-8 | **Mobile:** Экран чата (сообщения, reply, typing indicator, auto-scroll) | 8 | 🔴 | 📱 |
| CHT-9 | **Mobile:** WebSocket connection manager (connect, reconnect, exponential backoff) | 5 | 🔴 | 📱 |
| CHT-10 | **Mobile:** Header badge (💬 с кол-вом непрочитанных) | 2 | 🟠 | 📱 |
| CHT-11 | **Mobile:** Mute/unmute чата | 1 | 🟡 | 📱 |

### EPIC 8: Уведомления

| ID | Story | SP | P | Компонент |
|----|-------|----|---|-----------|
| NTF-1 | **Backend:** Notification service — создание, хранение, FCM push | 5 | 🔴 | 🔧 |
| NTF-2 | **Backend:** Firebase FCM integration — отправка push (по token, по topic) | 5 | 🔴 | 🔧 |
| NTF-3 | **Backend:** Триггеры уведомлений: новый отклик, напоминание, счёт, сообщение, рейтинг | 5 | 🔴 | 🔧 |
| NTF-4 | **Backend:** GET /notifications + POST /notifications/read + unread-count | 3 | 🔴 | 🔧 |
| NTF-5 | **Backend:** Scheduler: напоминания за 24ч и 1ч до игры (cron job) | 3 | 🟠 | 🔧 |
| NTF-6 | **Backend:** Quiet hours (не слать push между 23:00-07:00) | 2 | 🟡 | 🔧 |
| NTF-7 | **Mobile:** expo-notifications setup + permission request + token registration | 3 | 🔴 | 📱 |
| NTF-8 | **Mobile:** Экран уведомлений (группировка: Сегодня / Вчера / Ранее) | 3 | 🔴 | 📱 |
| NTF-9 | **Mobile:** Deep linking: tap push → конкретный экран (матч, чат, ивент) | 5 | 🟠 | 📱 |
| NTF-10 | **Mobile:** Header badge (🔔 с кол-вом непрочитанных) | 2 | 🟠 | 📱 |

---

## SPRINT 7-8: Главная, Посты, Доп. функции (Недели 13-16)
**Цель:** Home tab, feed, badges, friends, courts map

### EPIC 9: Главная страница

| ID | Story | SP | P | Компонент |
|----|-------|----|---|-----------|
| HOM-1 | **Backend:** GET /feed — глобальная лента (посты сообществ + друзей + результаты матчей) | 5 | 🔴 | 🔧 |
| HOM-2 | **Backend:** Home dashboard data — рейтинг, ближайшие игры, quick stats | 3 | 🔴 | 🔧 |
| HOM-3 | **Mobile:** Таб «Главная» — виджет рейтинга + quick actions | 5 | 🔴 | 📱 |
| HOM-4 | **Mobile:** Секция «Ближайшие игры» (до 3 из моих сообществ) | 3 | 🟠 | 📱 |
| HOM-5 | **Mobile:** Лента с табами (Новости / Лента) + бесконечный скролл | 5 | 🟠 | 📱 |
| HOM-6 | **Mobile:** Карточка поста (текст + фото + лайки) | 3 | 🟠 | 📱 |
| HOM-7 | **Mobile:** Автоматическая карточка результата матча | 3 | 🟡 | 📱 |

### EPIC 10: Посты

| ID | Story | SP | P | Компонент |
|----|-------|----|---|-----------|
| PST-1 | **Backend:** POST /posts — создание поста (текст + до 5 фото) | 3 | 🟠 | 🔧 |
| PST-2 | **Backend:** POST/DELETE /posts/:id/like | 2 | 🟠 | 🔧 |
| PST-3 | **Backend:** Автоматический пост-результат при подтверждении матча | 3 | 🟡 | 🔧 |
| PST-4 | **Mobile:** Экран создания поста (текст + выбор фото) | 3 | 🟠 | 📱 |
| PST-5 | **Mobile:** Like animation + count | 2 | 🟡 | 📱 |

### EPIC 11: Достижения / Бейджи

| ID | Story | SP | P | Компонент |
|----|-------|----|---|-----------|
| BDG-1 | **Backend:** Badge check service — проверка условий после каждого матча | 5 | 🟠 | 🔧 |
| BDG-2 | **Backend:** GET /rating/badges — earned + in_progress badges | 2 | 🟠 | 🔧 |
| BDG-3 | **Backend:** Push notification при получении нового бейджа | 2 | 🟡 | 🔧 |
| BDG-4 | **Mobile:** Секция бейджей в профиле (earned = цветные, in_progress = серые с прогрессом) | 3 | 🟠 | 📱 |
| BDG-5 | **Mobile:** Celebration animation при получении нового бейджа | 2 | 🟡 | 📱 |

### EPIC 12: Друзья и избранные

| ID | Story | SP | P | Компонент |
|----|-------|----|---|-----------|
| FRD-1 | **Backend:** POST/DELETE /friends/:user_id — добавить/удалить | 2 | 🟠 | 🔧 |
| FRD-2 | **Backend:** GET /friends — список друзей | 2 | 🟠 | 🔧 |
| FRD-3 | **Mobile:** Кнопка «В друзья» на публичном профиле | 1 | 🟠 | 📱 |
| FRD-4 | **Mobile:** Экран списка друзей (из профиля) | 2 | 🟡 | 📱 |

### EPIC 13: Карта кортов

| ID | Story | SP | P | Компонент |
|----|-------|----|---|-----------|
| CRT-1 | **Backend:** GET /courts + GET /courts/map — список и данные для карты | 3 | 🟠 | 🔧 |
| CRT-2 | **Backend:** GET /courts/:id — детали корта | 1 | 🟠 | 🔧 |
| CRT-3 | **Mobile:** Экран карты кортов (react-native-maps + маркеры) | 5 | 🟠 | 📱 |
| CRT-4 | **Mobile:** Bottomsheet с деталями корта при нажатии на маркер | 3 | 🟠 | 📱 |
| CRT-5 | **Mobile:** Интеграция карты в конструктор ивента (шаг 5: выбор корта) | 3 | 🟠 | 📱 |

---

## SPRINT 9-10: Web Panels (Недели 17-20)
**Цель:** Веб-панель админа + минимальный суперадмин

### EPIC 14: Web Admin Panel

| ID | Story | SP | P | Компонент |
|----|-------|----|---|-----------|
| ADM-1 | **Web:** Project setup (Vite + React + Shadcn/UI + React Router + TanStack Query) | 3 | 🔴 | 🖥 |
| ADM-2 | **Web:** Auth flow (телефон + OTP, тот же backend) | 3 | 🔴 | 🖥 |
| ADM-3 | **Web:** Layout (sidebar navigation, community switcher) | 3 | 🔴 | 🖥 |
| ADM-4 | **Web:** Dashboard — метрики, графики роста, последняя активность | 5 | 🔴 | 🖥 |
| ADM-5 | **Web:** Members page — таблица, фильтры, сортировка, массовые действия | 5 | 🔴 | 🖥 |
| ADM-6 | **Web:** Join requests — очередь заявок, approve/reject | 3 | 🔴 | 🖥 |
| ADM-7 | **Web:** Events page — таблица ивентов, создание, статусы | 5 | 🔴 | 🖥 |
| ADM-8 | **Web:** Event detail — участники, ввод результатов матчей | 5 | 🔴 | 🖥 |
| ADM-9 | **Web:** Posts — создание от имени сообщества, модерация | 3 | 🟠 | 🖥 |
| ADM-10 | **Web:** Leaderboard — рейтинг участников, настройка параметров | 3 | 🟠 | 🖥 |
| ADM-11 | **Web:** Settings — настройки сообщества (info, access, rating, danger zone) | 3 | 🟠 | 🖥 |
| ADM-12 | **Web:** Export data (CSV — members, matches, ratings) | 3 | 🟡 | 🖥 |
| ADM-13 | **Backend:** GET /admin/communities/:id/dashboard — dashboard API | 3 | 🔴 | 🔧 |
| ADM-14 | **Backend:** GET /admin/communities/:id/export — CSV export | 3 | 🟡 | 🔧 |

### EPIC 15: Superadmin Panel (минимальная)

| ID | Story | SP | P | Компонент |
|----|-------|----|---|-----------|
| SUP-1 | **Web:** Project setup (отдельное React приложение, Shadcn/UI) | 2 | 🔴 | 🖥 |
| SUP-2 | **Web:** Auth (только superadmin role) | 2 | 🔴 | 🖥 |
| SUP-3 | **Web:** Dashboard — общая статистика платформы | 3 | 🔴 | 🖥 |
| SUP-4 | **Web:** Verification queue — заявки на верификацию сообществ | 3 | 🔴 | 🖥 |
| SUP-5 | **Web:** User management — поиск, бан/разбан | 3 | 🔴 | 🖥 |
| SUP-6 | **Web:** Courts CRUD — добавление, редактирование, удаление кортов | 3 | 🔴 | 🖥 |
| SUP-7 | **Backend:** Superadmin endpoints (6 штук из API spec) | 5 | 🔴 | 🖥🔧 |

---

## SPRINT 11: Polish & Integration (Недели 21-22)
**Цель:** Интеграционное тестирование, полировка

### EPIC 16: Polish

| ID | Story | SP | P | Компонент |
|----|-------|----|---|-----------|
| POL-1 | Empty states для всех списков (ивенты, чаты, уведомления, участники) | 3 | 🟠 | 📱 |
| POL-2 | Error states + retry buttons | 3 | 🟠 | 📱 |
| POL-3 | Skeleton screens (loading placeholders) | 3 | 🟠 | 📱 |
| POL-4 | Pull-to-refresh на всех списках | 2 | 🟠 | 📱 |
| POL-5 | Haptic feedback на кнопках и действиях | 1 | 🟡 | 📱 |
| POL-6 | Splash screen + app icon | 2 | 🔴 | 📱 |
| POL-7 | App Store screenshots + описание | 3 | 🔴 | 📱 |
| POL-8 | Backend: seed data для демо (тестовые пользователи, сообщества, ивенты) | 3 | 🟠 | 🔧 |
| POL-9 | Performance testing — API response times, WebSocket load | 3 | 🟠 | 🔧 |
| POL-10 | Security audit — rate limits, input validation, SQL injection tests | 3 | 🔴 | 🔧 |

---

## SPRINT 12: QA & Beta (Недели 23-24)
**Цель:** Тестирование, исправление багов, beta launch

### EPIC 17: QA & Launch

| ID | Story | SP | P | Компонент |
|----|-------|----|---|-----------|
| QA-1 | Manual testing всех user flows (happy path + edge cases) | 8 | 🔴 | All |
| QA-2 | Bug fixes — critical и high priority | 8 | 🔴 | All |
| QA-3 | Deploy production (backend → Railway, web → Vercel) | 3 | 🔴 | 🔧🖥 |
| QA-4 | EAS Build — iOS (TestFlight) + Android (Google Play Beta) | 3 | 🔴 | 📱 |
| QA-5 | Онбординг 2 партнёрских сообществ (перенос данных, обучение админов) | 5 | 🔴 | All |
| QA-6 | Monitoring setup — Sentry + UptimeRobot + Firebase Analytics | 3 | 🟠 | 🔧📱 |
| QA-7 | Seed: корты Астаны (полный список) | 2 | 🟠 | 🔧 |

---

## Итого

### По спринтам
| Sprint | Недели | Фокус | SP (approx) |
|--------|--------|-------|-------------|
| 1-2 | 1-4 | Foundation + Auth | ~82 |
| 3-4 | 5-8 | Users + Communities + Events | ~130 |
| 5-6 | 9-12 | Matches + Rating + Chat + Notifications | ~125 |
| 7-8 | 13-16 | Home + Posts + Badges + Friends + Courts | ~80 |
| 9-10 | 17-20 | Web Admin + Superadmin | ~70 |
| 11 | 21-22 | Polish + Integration | ~26 |
| 12 | 23-24 | QA + Beta Launch | ~32 |
| **Total** | **24 недели** | | **~545 SP** |

### По компоненту
| Компонент | Stories | SP |
|-----------|---------|-----|
| 🔧 Backend | ~55 | ~215 |
| 📱 Mobile | ~55 | ~210 |
| 🖥 Web (Admin + Superadmin) | ~20 | ~75 |
| All (QA, Polish) | ~10 | ~45 |
| **Total** | **~140** | **~545** |

### Critical Path
```
INF-1 → INF-4 → AUTH-1 → AUTH-2 → AUTH-4
    ↓
USR-1 → COM-1 → EVT-1 → MTH-1 → MTH-2 → MTH-4
    ↓
CHT-1 → CHT-2 → NTF-1 → HOM-1
    ↓
ADM-1 → ADM-4 → SUP-1
    ↓
QA-1 → QA-3 → QA-4 → QA-5 (LAUNCH)
```

---

**Статус:** APPROVED v2.1  
**Использование:** Промпты для AI-агентов, планирование спринтов  
**Связанные документы:** MVP-SCOPE.md, PRD.md, API-SPECIFICATION.md
